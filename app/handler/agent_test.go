package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gin-web/dto"
	"gin-web/internal/agent"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type testArticleLister struct{ called bool }

func (l *testArticleLister) List(_ context.Context, req *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	l.called = req.Keyword == "Go" && req.Pager.PageSize == 2
	return &dto.ListArticleRespDTO{}, nil
}

type testRoundTripper func(*http.Request) (*http.Response, error)

func (f testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestAgentChatStreamsToolCallAndFinalAnswer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	requestCount := 0
	transport := testRoundTripper(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || !strings.HasSuffix(req.URL.Path, "/v1/responses") {
			t.Errorf("unexpected OpenAI request: %s %s", req.Method, req.URL)
		}
		if req.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("missing bearer authorization")
		}
		var body struct {
			Input []json.RawMessage `json:"input"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, err
		}
		requestCount++
		var payload string
		if requestCount == 1 {
			payload = `{"output":[{"type":"function_call","name":"list_articles","call_id":"call_1","arguments":"{\"keyword\":\"Go\",\"page_size\":2}"}]}`
		} else {
			if len(body.Input) < 3 || !strings.Contains(string(body.Input[len(body.Input)-1]), "function_call_output") {
				t.Errorf("follow-up request did not contain tool result: %s", body.Input)
			}
			payload = `{"output_text":"找到 Go 文章。","output":[{"type":"message","role":"assistant"}]}`
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
			Request:    req,
		}, nil
	})

	articles := &testArticleLister{}
	agentService := agent.NewServiceWithClient(articles, "test-key", "test-model", "https://relay.test/codex", &http.Client{Transport: transport})
	services := &service.Services{Agent: agentService}
	router := gin.New()
	router.POST("/admin/agent/chat", NewAgentHandler(services).Chat)
	req := httptest.NewRequest(http.MethodPost, "/admin/agent/chat", strings.NewReader(`{"message":"搜索 Go 文章"}`))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	if response.Code != http.StatusOK || !strings.Contains(response.Header().Get("Content-Type"), "text/event-stream") {
		t.Fatalf("expected SSE 200 response, got status=%d content-type=%q body=%s", response.Code, response.Header().Get("Content-Type"), response.Body.String())
	}
	if !articles.called || requestCount != 2 {
		t.Fatalf("tool was not executed and followed up: called=%t requests=%d", articles.called, requestCount)
	}
	body := response.Body.String()
	orderedEvents := []string{
		"event: start",
		"event: tool_call",
		`"name":"list_articles"`,
		`"keyword":"Go"`,
		"event: message",
		"找到 Go 文章。",
		"event: done",
	}
	lastIndex := -1
	for _, part := range orderedEvents {
		index := strings.Index(body, part)
		if index <= lastIndex {
			t.Fatalf("SSE part %q out of order or missing in response: %s", part, body)
		}
		lastIndex = index
	}
}
