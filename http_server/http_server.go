package httpserver

import (
	"fmt"
	"log/slog"
	"strings"

	vh "github.com/FMotalleb/udp_listener/value_holder"
	"github.com/gin-gonic/gin"
)

// StartHttpServer starts an HTTP server that serves the current state of the given valueHolder.
// It uses the Gin web framework and listens on the specified address and port.
// If basic authentication is enabled (providing a user:pass to cli), it requires a valid username and password.
// The server responds to GET requests at the "/api/v1/state/current" endpoint with the current state of the system.
func StartHttpServer(serveAddr string, auth string, valueHolder *vh.ValueHolder) error {
	server := gin.New()
	if auth := strings.Split(auth, ":"); len(auth) == 2 {
		user := auth[0]
		pass := auth[1]
		server.Use(
			gin.BasicAuth(
				gin.Accounts{
					user: pass,
				},
			),
		)
	}
	server.GET(
		"api/v1/state/current",
		func(g *gin.Context) {
			g.JSON(200, map[string]any{
				"code":    1,
				"message": "",
				"data":    valueHolder.ToMap(),
			})
		},
	)
	err := server.Run(serveAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("http server fatal error: %s", err))
		return err
	}
	return nil
}
