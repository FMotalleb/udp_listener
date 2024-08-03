package udp

import (
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/FMotalleb/udp_listener/cmd"
)

func StartUdpServer(buf io.Writer) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cmd.ListenAddr, cmd.UDPListenPort))
	if err != nil {
		slog.Error(fmt.Sprintf("cannot resolve udp listen address: %s", err))
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error(fmt.Sprintf("cannot start udp server: %s", err))
		return err
	}
	defer conn.Close()
	slog.Info(fmt.Sprintf("listening on %s", conn.LocalAddr()))
	for {
		n, err := io.Copy(buf, conn)
		if err != nil {
			slog.Error(fmt.Sprintf("error reading udp data after %d data was received: %s", n, err))
		}
	}
}
