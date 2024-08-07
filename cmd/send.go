/*
Copyright © 2024 Motalleb Fallahnezhad

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

var (
	sendIp           = ""
	sendPort  uint16 = 7982
	localPort uint16 = 7983
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send <data>",
	Short: "sends a single message to the server provided",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("sending to given address")
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", sendIp, sendPort))
		if err != nil {
			slog.Error(fmt.Sprintf("cannot resolve udp listen address: %s", err))
			return
		}
		laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "127.0.0.1", localPort))
		if err != nil {
			slog.Error(fmt.Sprintf("cannot resolve udp local address: %s", err))
			return
		}

		conn, err := net.DialUDP("udp", laddr, addr)
		if err != nil {
			slog.Error(fmt.Sprintf("cannot connect to given server address(%s) from local(%s): %s", addr.String(), laddr.String(), err))
			return
		}
		strMsg := strings.Join(args, " ")
		n, err := conn.Write([]byte(strMsg))
		if err != nil {
			slog.Error(fmt.Sprintf("cannot send data to given server address(%s) from local(%s): %s", addr.String(), laddr.String(), err))
			return
		}
		slog.Debug(fmt.Sprintf("%d bytes written to %s", n, addr))
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&sendIp, "ip", "i", "127.0.0.1", "specify ip to connect")
	sendCmd.Flags().Uint16VarP(&sendPort, "port", "p", 7982, "specify port to connect")
	sendCmd.Flags().Uint16VarP(&localPort, "local-port", "l", 7983, "specify the port used to connect")
}
