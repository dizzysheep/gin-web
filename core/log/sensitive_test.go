package log

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	// 本包init会按配置在 ./runtime/logs/ 生成日志文件，测试后清理
	_ = os.RemoveAll("runtime")
	os.Exit(code)
}

func TestHideSensitiveInfo(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		leaks []string // 不应再出现的敏感明文
	}{
		{
			name:  "json格式密码",
			in:    `{"username":"admin","password":"secret123"}`,
			leaks: []string{"secret123"},
		},
		{
			name:  "表单格式密码",
			in:    "username=admin&password=secret123&age=18",
			leaks: []string{"secret123"},
		},
		{
			name:  "手机号",
			in:    "contact=13812345678&ok=1",
			leaks: []string{"13812345678"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := HideSensitiveInfo(c.in)
			for _, leak := range c.leaks {
				if strings.Contains(got, leak) {
					t.Fatalf("sensitive value %q not hidden, got: %s", leak, got)
				}
			}
		})
	}
}

func TestHideSensitiveInfoKeepsNormalContent(t *testing.T) {
	in := `{"username":"admin","title":"hello world"}`
	if got := HideSensitiveInfo(in); got != in {
		t.Fatalf("normal content should not change, got: %s", got)
	}
}
