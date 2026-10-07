package dto

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// AgentChatRequest Agent 对话请求
type AgentChatRequest struct {
	Message string `json:"message" binding:"required,max=4000" example:"列出最近发布的文章"`
}

func AgentChatReqToDTO(c *gin.Context) (*AgentChatReqDTO, error) {
	var req AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Message) == "" {
		return nil, errors.New("消息不能为空")
	}

	return &AgentChatReqDTO{
		Message: req.Message,
	}, nil
}

type AgentChatReqDTO struct {
	Message string
}
