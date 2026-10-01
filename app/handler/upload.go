package handler

import (
	"io"

	"gin-web/app/response"
	"gin-web/core/config"
	"gin-web/core/storage"
	"gin-web/internal/errcode"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	driver storage.Driver
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{driver: storage.NewDriver()}
}

type UploadResponse struct {
	URL string `json:"url" example:"/uploads/20261001/6f96c9c9-1e07-4a1e-8f7a-3b6a9d4c2e01.png"`
}

// UploadImage godoc
// @Summary 管理-上传图片
// @Description 支持jpg/jpeg/png/gif/webp，默认上限5MB，返回可访问URL
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	// 扩展名与大小校验
	if err := storage.Validate(fileHeader); err != nil {
		response.FailErr(c, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.FailErr(c, errcode.NewCustomError(errcode.ErrFileInvalid))
		return
	}
	defer func() { _ = file.Close() }()

	maxSizeMB := config.GetInt("upload.maxSize")
	if maxSizeMB <= 0 {
		maxSizeMB = 5
	}
	maxBytes := (int64(maxSizeMB) << 20) + 1
	data, err := io.ReadAll(io.LimitReader(file, maxBytes))
	if err != nil {
		response.FailErr(c, errcode.NewCustomError(errcode.ErrFileInvalid))
		return
	}
	if len(data) > maxSizeMB<<20 {
		response.FailErr(c, errcode.NewCustomError(errcode.ErrFileInvalid))
		return
	}
	if err := storage.ValidateContent(data); err != nil {
		response.FailErr(c, err)
		return
	}

	url, err := h.driver.Save(fileHeader.Filename, data)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, &UploadResponse{URL: url})
}
