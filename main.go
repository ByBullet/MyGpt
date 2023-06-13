package main

import (
	"ai/messages"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/russross/blackfriday"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var ms = new(messages.MessageStorage)

func init() {
	ms = new(messages.MessageStorage)
}

func authHandler(c *gin.Context) {
	session := sessions.Default(c)
	key := session.Get("key")
	if key == nil {
		key = uuid.NewString()
		session.Set("key", key)
		session.Save()
	}
	c.Next()
}

func loadSession(r *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store), authHandler)
}

func InitRouter(r *gin.Engine) {
	loadSession(r)
	var files []string
	filepath.Walk("./www/template", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	r.LoadHTMLFiles(files...)
	r.Static("/static/", "./www/static/")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})
	r.POST("/chat", func(c *gin.Context) {
		json := make(map[string]string) //注意该结构接受的内容
		c.ShouldBind(&json)

		req := json["req"]
		session := sessions.Default(c)
		key := session.Get("key")
		msg, err := ms.GetMessage(key.(string), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
			return
		}

		ans, err := msg.Chat(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"mgs": err})
			return
		}
		ans = string(blackfriday.MarkdownCommon([]byte(ans)))

		c.JSON(http.StatusOK, gin.H{"msg": ans})
	})

}

func main() {
	router := gin.Default()
	InitRouter(router)
	router.Run(":80")
}
