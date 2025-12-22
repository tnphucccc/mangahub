package websocket

import "github.com/gin-gonic/gin"

func Handler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	}
}
