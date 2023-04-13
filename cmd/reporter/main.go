package main

import (
	"fmt"
	"net/http"
	"twitch_chat_analysis/internal/cache"
	"twitch_chat_analysis/internal/environment"

	"github.com/gin-gonic/gin"
)

func main() {

	env := environment.MustGetReporterEnv()
	cacheClient := cache.New(env.RedisHost, env.RedisPort)

	r := gin.Default()

	r.GET("/message/list", func(c *gin.Context) {

		senderParam, ok := c.GetQuery("sender")
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		receiverParam, ok := c.GetQuery("receiver")
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		res, err := cacheClient.List(
			c.Request.Context(),
			senderParam,
			receiverParam,
		)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)

	})
	r.Run(fmt.Sprintf(":%d", env.ApiPort))
}
