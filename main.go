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
package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/FMotalleb/udp_listener/cmd"
	hs "github.com/FMotalleb/udp_listener/http_server"
	udp "github.com/FMotalleb/udp_listener/udp_server"
	vh "github.com/FMotalleb/udp_listener/value_holder"
)

func main() {
	cmd.Execute()
	if cmd.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cmd.ListenAddr, cmd.UDPListenPort))
	if err != nil {
		slog.Error(fmt.Sprintf("cannot resolve udp listen address: %s", err))
		return
	}
	wg := new(sync.WaitGroup)
	buf := vh.NewValueHolder(cmd.Zero)
	wg.Add(1)
	go creteUdpServer(addr, wg, buf)
	go wg.Add(1)
	go createHttpServer(wg, buf)

	wg.Wait()
}

func createHttpServer(wg *sync.WaitGroup, value *vh.ValueHolder) {
	defer func() {
		wg.Done()
	}()
	for {
		_ = hs.StartHttpServer(value)
		slog.Warn("an error occurred during tcp server's starting/listening phase restarting the server after one second")
		time.Sleep(time.Second)
	}
}
func creteUdpServer(addr *net.UDPAddr, wg *sync.WaitGroup, buf io.Writer) {
	defer func() {
		wg.Done()
	}()
	for {
		err := udp.StartUdpServer(addr, buf)
		if err == nil {
			slog.Warn("an error occurred during udp server's starting/listening phase restarting the server after one second")
			time.Sleep(time.Second)
			continue
		} else {
			slog.Warn(
				fmt.Sprintf(
					"a fatal error occurred during udp server's starting phase panicing out: %s",
					err,
				),
			)
			return
		}
	}
}
