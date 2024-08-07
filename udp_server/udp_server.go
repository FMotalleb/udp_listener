package udp

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

// StartUdpServer starts a UDP server listening on the specified address and writes received data to the provided io.Writer.
// It continuously listens for incoming UDP packets and checks if the client's IP address is authorized.
// If the client is authorized, the received data is written to the provided io.Writer.
// If the client is not authorized, a warning message is logged and the server continues listening for incoming packets.
// If an error occurs while starting the server or reading UDP data, it is logged and returned as an error.
func StartUdpServer(addr *net.UDPAddr, allowList []string, buf io.Writer) error {

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
		hasAccess := checkAccess(allowList, addr)
		if !hasAccess {
			slog.Warn(fmt.Sprintf("uanuthorized access to udp server from ip: %s", addr))
			continue
		}
		if _, err = buf.Write(ans[0:n]); err != nil {
			slog.Error(fmt.Sprintf("error reading udp data after %d data was received: %s", n, err))
		}
	}
}

// checkAccess checks if the client's IP address is authorized by comparing it with the list of allowed IP addresses.
// If the list of allowed IP addresses is empty, the function returns true, indicating that the client is authorized.
// If the client's IP address is found in the list of allowed IP addresses, the function returns true.
// Otherwise, the function returns false, indicating that the client is not authorized.
func checkAccess(allowList []string, addr net.Addr) bool {
	if len(allowList) == 0 {
		return true
	}
	for _, i := range allowList {
		ip := strings.Split(addr.String(), ":")[0]
		if ip == i {
			return true
		}
	}
	return false
}
