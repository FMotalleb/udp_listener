package udp

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/FMotalleb/udp_listener/cmd"
)

func StartUdpServer(addr *net.UDPAddr, buf io.Writer) error {

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		slog.Error(fmt.Sprintf("cannot start udp server: %s", err))
		return err
	}

	defer conn.Close()
	slog.Info(fmt.Sprintf("listening on %s", conn.LocalAddr()))
	for {
		ans := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(ans)

		if err != nil {
			slog.Error(fmt.Sprintf("error reading udp data after %d data was received: %s", n, err))
		}
		hasAccess := checkAccess(addr)
		if !hasAccess {
			slog.Warn(fmt.Sprintf("uanuthorized access to udp server from ip: %s", addr))
			continue
		}
		if _, err = buf.Write(ans[0:n]); err != nil {
			slog.Error(fmt.Sprintf("error reading udp data after %d data was received: %s", n, err))
		}
	}
}

func checkAccess(addr net.Addr) bool {
	if len(cmd.AllowedUDPClients) == 0 {
		return true
	}
	for _, i := range cmd.AllowedUDPClients {
		ip := strings.Split(addr.String(), ":")[0]
		if ip == i {
			return true
		}
	}

	return false
}
