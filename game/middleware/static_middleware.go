package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

var staticFileType = []string{
	"js", "css", "jpg", "png", "gif", "html", "ico", "swf", "vue", "json", "map", "ico",
}

// 调试的时候用go处理静态文件，线上用nginx处理
func StaticMdl(staticRoot string) gin.HandlerFunc {
	return func(c *gin.Context) {
		UrlPath := c.Request.URL.Path
		items := strings.Split(UrlPath, ".")
		if len(items) > 1 {
			isStaticFile := false
			for _, t := range staticFileType {
				if t == items[len(items)-1] {
					isStaticFile = true
					break
				}
			}
			if isStaticFile {
				// 处理静态文件
				c.File(staticRoot + UrlPath)
				c.Abort()
			}
		}
	}
}
