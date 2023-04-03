package v1

import (
	"gin-web/core/apis"
	"gin-web/core/ginc"
	"github.com/gin-gonic/gin"
)

type BaikeRes struct {
	Title     string `json:"title"`
	ItemId    string `json:"itemId"`
	ItemOrder string `json:"itemOrder"`
}

func GetTest(c *gin.Context) {
	res := []BaikeRes{}
	uri := "baikebcs:/cms/pc_index/knowledge_topic_menu.json"
	_ = apis.Get(uri).ToJSON(&res)
	//log.Get(c).Info("test info")
	ginc.Ok(c, res)
}
