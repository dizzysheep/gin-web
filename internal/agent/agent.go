package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gin-web/core/config"
	"gin-web/core/xtime"
	"gin-web/dto"
)

const (
	defaultBaseURL = "https://new.sharedchat.cc/codex"
	defaultModel   = "5.6-sol"
	maxToolRounds  = 4
)

type ArticleLister interface {
	List(ctx context.Context, req *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error)
}

type ToolDefinition struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	Strict      bool           `json:"strict"`
}

type ToolCallLog struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
	Output    any             `json:"output,omitempty"`
}

type ChatResult struct {
	Answer    string        `json:"answer"`
	ToolCalls []ToolCallLog `json:"tool_calls"`
}

type Service struct {
	apiKey   string
	baseURL  string
	model    string
	articles ArticleLister
	client   *http.Client
	tools    []ToolDefinition
	handlers map[string]func(context.Context, json.RawMessage) (any, error)
}

// NewServiceFromConfig reads the [agent] settings from config/app.toml.
// OPENAI_API_KEY, OPENAI_BASE_URL, and OPENAI_MODEL take precedence when set.
func NewServiceFromConfig(articles ArticleLister) *Service {
	apiKey := config.GetString("agent.key")
	if value, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		apiKey = value
	}
	baseURL := config.GetString("agent.domain")
	if value, ok := os.LookupEnv("OPENAI_BASE_URL"); ok {
		baseURL = value
	}
	model := config.GetString("agent.model")
	if value, ok := os.LookupEnv("OPENAI_MODEL"); ok {
		model = value
	}
	return newService(
		articles,
		apiKey,
		model,
		baseURL,
		&http.Client{Timeout: 60 * time.Second},
	)
}

// NewServiceWithClient creates a service with explicit connection settings.
// The injected client also makes API integration tests deterministic.
func NewServiceWithClient(articles ArticleLister, apiKey, model, baseURL string, client *http.Client) *Service {
	return newService(articles, apiKey, model, baseURL, client)
}

type responseRequest struct {
	Model        string           `json:"model"`
	Instructions string           `json:"instructions"`
	Input        []any            `json:"input"`
	Tools        []ToolDefinition `json:"tools"`
	ToolChoice   string           `json:"tool_choice"`
}

type responseBody struct {
	OutputText string            `json:"output_text"`
	Output     []json.RawMessage `json:"output"`
}

type responseItem struct {
	Type      string               `json:"type"`
	Name      string               `json:"name"`
	CallID    string               `json:"call_id"`
	Arguments string               `json:"arguments"`
	Role      string               `json:"role"`
	Content   []responseItemOutput `json:"content"`
}

type responseItemOutput struct {
	Type string `json:"type"` // output_text
	Text string `json:"text"`
}

func newService(articles ArticleLister, apiKey, model, baseURL string, client *http.Client) *Service {
	if strings.TrimSpace(model) == "" {
		model = defaultModel
	}
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultBaseURL
	}
	s := &Service{
		apiKey:   strings.TrimSpace(apiKey),
		baseURL:  strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		model:    model,
		articles: articles,
		client:   client,
		tools:    registeredTools(),
		handlers: make(map[string]func(context.Context, json.RawMessage) (any, error)),
	}
	s.handlers["list_articles"] = s.listArticles
	return s
}

func registeredTools() []ToolDefinition {
	return []ToolDefinition{{
		Type:        "function",
		Name:        "list_articles",
		Description: "Search published blog articles. Look up a single article by its exact id, or search by keyword (use an empty keyword to list recent published articles). This tool is read-only.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":        map[string]any{"type": "integer", "description": "Exact article ID to look up; use 0 when searching by keyword or listing recent articles."},
				"keyword":   map[string]any{"type": "string", "description": "Words to search in the title, description, or body; use an empty string for recent articles. Ignored when id is greater than 0."},
				"page_size": map[string]any{"type": "integer", "minimum": 1, "maximum": 10, "description": "Maximum number of articles to return."},
			},
			"required":             []string{"id", "keyword", "page_size"},
			"additionalProperties": false,
		},
		Strict: true,
	}}
}

func (s *Service) Tools() []ToolDefinition {
	tools := make([]ToolDefinition, len(s.tools))
	copy(tools, s.tools)
	return tools
}

func (s *Service) Chat(ctx context.Context, message string) (*ChatResult, error) {
	return s.ChatWithEvents(ctx, message, nil)
}

// ChatWithEvents runs the tool loop and reports each completed tool call as it
// happens. The final answer is returned after the model finishes the turn.
func (s *Service) ChatWithEvents(ctx context.Context, message string, onToolCall func(ToolCallLog)) (*ChatResult, error) {
	if s.apiKey == "" {
		return nil, errors.New("Agent 未配置 OPENAI_API_KEY")
	}
	if s.articles == nil {
		return nil, errors.New("Agent 文章查询服务未配置")
	}
	if strings.TrimSpace(message) == "" {
		return nil, errors.New("消息不能为空")
	}

	input := []any{map[string]any{"role": "user", "content": message}}
	result := &ChatResult{ToolCalls: make([]ToolCallLog, 0)}
	for round := 0; round < maxToolRounds; round++ {
		resp, err := s.createResponse(ctx, input)
		if err != nil {
			return nil, err
		}

		calls := make([]responseItem, 0)
		for _, raw := range resp.Output {
			var item responseItem
			if err := json.Unmarshal(raw, &item); err != nil {
				return nil, fmt.Errorf("解析模型响应失败: %w", err)
			}
			if item.Type == "function_call" {
				calls = append(calls, item)
			}
		}
		if len(calls) == 0 {
			result.Answer = resp.OutputText
			return result, nil
		}

		// Keep all response items so reasoning models can continue the same turn.
		for _, raw := range resp.Output {
			input = append(input, json.RawMessage(raw))
		}
		for _, call := range calls {
			toolOutput, err := s.executeTool(ctx, call)
			if err != nil {
				return nil, err
			}
			toolCallLog := ToolCallLog{
				Name: call.Name, Arguments: json.RawMessage(call.Arguments), Output: toolOutput,
			}
			result.ToolCalls = append(result.ToolCalls, toolCallLog)
			if onToolCall != nil {
				onToolCall(toolCallLog)
			}
			encodedOutput, err := json.Marshal(toolOutput)
			if err != nil {
				return nil, fmt.Errorf("编码工具结果失败: %w", err)
			}
			input = append(input, map[string]any{
				"type":    "function_call_output",
				"call_id": call.CallID,
				"output":  string(encodedOutput),
			})
		}
	}
	return nil, fmt.Errorf("模型连续调用工具超过 %d 轮", maxToolRounds)
}

func (s *Service) createResponse(ctx context.Context, input []any) (*responseBody, error) {
	body, err := json.Marshal(responseRequest{
		Model:        s.model,
		Instructions: "你是博客后台助手。需要查询文章时调用 list_articles；只能根据工具返回的数据回答，不得猜测文章信息。回答中列举文章时，必须逐条展示文章标题，并可附带发布日期、分类与标签；若结果为空需明确说明没有找到。",
		Input:        input,
		Tools:        s.tools,
		ToolChoice:   "auto",
	})
	if err != nil {
		return nil, fmt.Errorf("编码 OpenAI 请求失败: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.responsesURL(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建 OpenAI 请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用 OpenAI API 失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("OpenAI API 返回 HTTP %d", resp.StatusCode)
	}
	var decoded responseBody
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("解析 OpenAI API 响应失败: %w", err)
	}
	// 部分中转站不返回顶层 output_text 便捷字段，回退到 output[] 中提取 message 文本
	if decoded.OutputText == "" {
		decoded.OutputText = extractOutputText(decoded.Output)
	}
	return &decoded, nil
}

func extractOutputText(items []json.RawMessage) string {
	var b strings.Builder
	for _, raw := range items {
		var item responseItem
		if err := json.Unmarshal(raw, &item); err != nil || item.Type != "message" {
			continue
		}
		for _, c := range item.Content {
			if c.Type == "output_text" {
				b.WriteString(c.Text)
			}
		}
	}
	return b.String()
}

// agentArticleSummary 回传给模型的文章摘要：只保留回答所需的字段，
// 避免完整 VO（desc/cover/嵌套对象等）挤占上下文、稀释模型对标题的注意力。
type agentArticleSummary struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	PublishedOn string   `json:"published_on,omitempty"`
}

func (s *Service) responsesURL() string {
	if strings.HasSuffix(s.baseURL, "/v1") {
		return s.baseURL + "/responses"
	}
	return s.baseURL + "/v1/responses"
}

func (s *Service) executeTool(ctx context.Context, call responseItem) (any, error) {
	handler, ok := s.handlers[call.Name]
	if !ok {
		return nil, fmt.Errorf("模型请求了未注册的工具 %q", call.Name)
	}
	arguments := json.RawMessage(call.Arguments)
	if len(arguments) == 0 || !json.Valid(arguments) {
		return nil, fmt.Errorf("工具 %q 的参数不是有效 JSON", call.Name)
	}
	return handler(ctx, arguments)
}

func (s *Service) listArticles(ctx context.Context, raw json.RawMessage) (any, error) {
	var args struct {
		ID       int64  `json:"id"`
		Keyword  string `json:"keyword"`
		PageSize int    `json:"page_size"`
	}
	if err := json.Unmarshal(raw, &args); err != nil {
		return nil, fmt.Errorf("解析 list_articles 参数失败: %w", err)
	}
	if args.PageSize < 1 || args.PageSize > 10 {
		return nil, errors.New("list_articles 的 page_size 必须在 1 到 10 之间")
	}
	if args.ID < 0 {
		return nil, errors.New("list_articles 的 id 不能为负数")
	}
	articles, err := s.articles.List(ctx, &dto.ListArticleReqDTO{
		ID:      args.ID,
		Keyword: args.Keyword,
		Pager:   dto.PagerReqToDTO(1, args.PageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}
	vo := articles.ToVO()
	list := make([]agentArticleSummary, 0, len(vo.List))
	for _, item := range vo.List {
		summary := agentArticleSummary{
			ID:          item.ID,
			Title:       item.Title,
			PublishedOn: xtime.GetTimeFormat(int64(item.PublishedOn), "2006-01-02"),
		}
		if item.Category != nil {
			summary.Category = item.Category.Name
		}
		if len(item.Tags) > 0 {
			summary.Tags = make([]string, 0, len(item.Tags))
			for _, tag := range item.Tags {
				summary.Tags = append(summary.Tags, tag.Name)
			}
		}
		list = append(list, summary)
	}
	return map[string]any{
		"pager": vo.Pager,
		"list":  list,
	}, nil
}
