package httpserver

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/FMotalleb/udp_listener/cmd"
	vh "github.com/FMotalleb/udp_listener/value_holder"
	"github.com/gin-gonic/gin"
)

func StartHttpServer(value *vh.ValueHolder) error {
	server := gin.New()
	if auth := strings.Split(cmd.BasicAuth, ":"); len(auth) == 2 {
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
		"api/v1/wight/current",
		func(g *gin.Context) {
			g.JSON(200, map[string]any{
				"code":    1,
				"message": "Weight fetched Successfully",
				"data":    value.ToMap(),
			})
		},
	)
	err := server.Run(fmt.Sprintf("%s:%d", cmd.ListenAddr, cmd.TCPListenPort))
	if err != nil {
		slog.Error(fmt.Sprintf("http server fatal error: %s", err))
		return err
	}
	return nil
}
