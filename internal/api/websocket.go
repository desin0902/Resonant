package api

import (
	ws "Resonant/internal/websocket"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		userID, err := ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &ws.Client{
			Hub:    hub,
			Conn:   wsConn,
			Send:   make(chan ws.Message),
			UserID: userID,
		}

		hub.register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}

func ValidateJWT(token string) (uint, error) {
	// TODO: Replace this with real JWT validation later
	if token == "" {
		return 0, fmt.Errorf("empty token")
	}
	return 1, nil // fake user ID for now
}
