package newspush

import (
	"assistant_game_server/client"
	"assistant_game_server/gameredis"
	"assistant_game_server/message"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func newspushHandler(c *gin.Context) {
	var err error
	defer func() {
		if err == nil {
			c.String(http.StatusOK, "1")
			glog.Info("push succeed")
		} else {
			c.String(http.StatusOK, "0")
			glog.Info("push failed, ", err)
		}
	}()
	glog.Info(c.Request.URL)
	articleID := c.Query("articleId")
	title := c.Query("title")
	tags := c.Query("tags")
	picurl := c.Query("picurl")

	var articleInfo = &message.ArticleInfo{}
	articleInfo.ArticleID = articleID
	articleInfo.Title = title
	articleInfo.Tags = tags
	articleInfo.PicURL = picurl

	go pushArticle(articleInfo)
}

func pushArticle(articleInfo *message.ArticleInfo) {
	gameredis.PushArticle(articleInfo)
	notify := &message.PushNewsNotify{}
	notify.Meta.MessageType = "PushNewsNotify"
	notify.Data.Articles = append(notify.Data.Articles, *articleInfo)
	client.Mgr.Broadcast(notify)
}

// RunHTTPServer ...
func RunHTTPServer() {
	router := gin.Default()
	router.POST("/newspush", newspushHandler)
	router.GET("/newspush", newspushHandler)

	router.Run(":7001")
}
