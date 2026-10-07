package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gin-web/app/response"
	"gin-web/dto"
	"gin-web/internal/agent"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	service *service.Services
}

func NewAgentHandler(service *service.Services) *AgentHandler {
	return &AgentHandler{service: service}
}

func (h *AgentHandler) Tools(c *gin.Context) {
	response.Ok(c, h.service.Agent.Tools())
}

// Chat 以 SSE 流式返回 Agent 对话过程：
// event: tool_call（单个工具调用完成即推送）、message（最终回答）、error、done。
// 参数校验失败仍返回普通 JSON，便于复用统一错误结构。
func (h *AgentHandler) Chat(c *gin.Context) {
	reqDTO, err := dto.AgentChatReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-store")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	flusher, _ := c.Writer.(http.Flusher)

	writeEvent := func(name string, payload any) {
		data, err := json.Marshal(payload)
		if err != nil {
			return
		}
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", name, data)
		if flusher != nil {
			flusher.Flush()
		}
	}

	// 连接建立即推送，前端可据此显示"已连接"状态
	writeEvent("start", map[string]string{"status": "started"})

	result, err := h.service.Agent.ChatWithEvents(c.Request.Context(), reqDTO.Message, func(call agent.ToolCallLog) {
		writeEvent("tool_call", call)
	})
	if err != nil {
		writeEvent("error", map[string]string{"message": err.Error()})
		return
	}
	writeEvent("message", result)
	writeEvent("done", map[string]string{"status": "ok"})
}
