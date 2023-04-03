package apis

import (
	"github.com/spf13/cast"
	"time"
)

type ApiNormal struct {
	Domain     string
	RequestUri string
	TimeOut    time.Duration
	Token      map[int]string
	Cid        int
}

func (c *ApiNormal) Uri() string {
	return c.Domain + c.RequestUri
}

func (c *ApiNormal) GetTimeOut() time.Duration {
	return c.TimeOut
}

func (c *ApiNormal) SetUri(uri string) {
	c.RequestUri = uri
}

// BandData 返回默认token，适合所有渠道一样的token
func (c *ApiNormal) BandData(args ...interface{}) *map[string]interface{} {
	var d = map[string]interface{}{}
	if len(args) == 0 {
		return &d
	}
	if val, ok := args[0].(map[string]interface{}); ok {
		d = val
	}
	d["token"] = c.SetToken(cast.ToInt(d["cid"]))
	return &d
}

// SetToken 设置token
func (c *ApiNormal) SetToken(cid int) string {
	if len(c.Token) == 0 {
		return ""
	}

	return c.Token[cid]
}

func NewApiNormal(apiType string) *ApiNormal {
	switch apiType {
	//测试用的
	case "baikebcs":
		return &ApiNormal{
			Domain:  "https://baikebcs.cdn.bcebos.com",
			TimeOut: time.Second * 3,
		}
	default:
		panic("api interface undefined")
	}
}
