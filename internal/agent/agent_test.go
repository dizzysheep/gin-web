package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"gin-web/dto"
)

type articleListerFunc func(context.Context, *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error)

func (f articleListerFunc) List(ctx context.Context, req *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	return f(ctx, req)
}

func TestChatExecutesRegisteredToolAndReturnsAnswer(t *testing.T) {
	requests := 0
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.String() != "https://relay.example/codex/v1/responses" {
			return nil, fmt.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			return nil, fmt.Errorf("authorization = %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode request: %w", err)
		}
		if len(body["tools"].([]any)) != 1 {
			return nil, fmt.Errorf("expected registered tool in request, got %#v", body["tools"])
		}
		requests++
		var payload string
		if requests == 1 {
			payload = `{"output":[{"type":"function_call","name":"list_articles","call_id":"call_1","arguments":"{\"keyword\":\"Go\",\"page_size\":3}"}]}`
		} else {
			input, ok := body["input"].([]any)
			if !ok || len(input) < 3 {
				return nil, fmt.Errorf("expected tool output in follow-up input: %#v", body["input"])
			}
			last, ok := input[len(input)-1].(map[string]any)
			if !ok || last["type"] != "function_call_output" || last["call_id"] != "call_1" {
				return nil, fmt.Errorf("unexpected tool output item: %#v", input[len(input)-1])
			}
			payload = `{"output_text":"找到 1 篇文章。","output":[{"type":"message","role":"assistant"}]}`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
			Request:    r,
		}, nil
	})

	listed := false
	articles := articleListerFunc(func(_ context.Context, req *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
		listed = true
		if req.Keyword != "Go" || req.Pager.PageSize != 3 {
			return nil, fmt.Errorf("unexpected list request: keyword=%q pager=%+v", req.Keyword, req.Pager)
		}
		return &dto.ListArticleRespDTO{}, nil
	})
	service := newService(articles, "test-key", "test-model", "https://relay.example/codex", &http.Client{Transport: transport})

	var streamedToolCalls []ToolCallLog
	result, err := service.ChatWithEvents(context.Background(), "查一下 Go 文章", func(call ToolCallLog) {
		streamedToolCalls = append(streamedToolCalls, call)
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if !listed || requests != 2 {
		t.Fatalf("listed=%t requests=%d, want true and 2", listed, requests)
	}
	if result.Answer != "找到 1 篇文章。" || len(result.ToolCalls) != 1 || result.ToolCalls[0].Name != "list_articles" {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(streamedToolCalls) != 1 || streamedToolCalls[0].Name != "list_articles" {
		t.Fatalf("tool call was not streamed as it completed: %+v", streamedToolCalls)
	}
}

func TestNewServiceFromEnvReadsConfigAndAppliesOverrides(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_MODEL", "")
	service := NewServiceFromConfig(nil)
	if service.baseURL != "https://new.sharedchat.cc/codex" || service.model != "5.6-sol" {
		t.Fatalf("config values not loaded: baseURL=%q model=%q", service.baseURL, service.model)
	}

	t.Setenv("OPENAI_API_KEY", "override-key")
	t.Setenv("OPENAI_BASE_URL", "https://override.example/api")
	t.Setenv("OPENAI_MODEL", "override-model")
	service = NewServiceFromConfig(nil)
	if service.apiKey != "override-key" || service.baseURL != "https://override.example/api" || service.model != "override-model" {
		t.Fatalf("environment overrides not applied: keySet=%t baseURL=%q model=%q", service.apiKey != "", service.baseURL, service.model)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
