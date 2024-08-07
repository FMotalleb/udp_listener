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
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	hs "github.com/FMotalleb/udp_listener/http_server"
	udp "github.com/FMotalleb/udp_listener/udp_server"
	vh "github.com/FMotalleb/udp_listener/value_holder"
	"github.com/spf13/cobra"
)

type serverConfig struct {
	allowedUDPClients []string
	basicAuth         string
	listenAddr        string
	UDPListenPort     uint16
	TCPListenPort     uint16
	UDPBufferSize     uint16
	zero              string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "udp_listener",
	Short: "Create a HTTP api to get last state emitted by the client of udp server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}
var rootCFG = new(serverConfig)
var verbose bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rootCFG.listenAddr, "listen-addr", "l", "0.0.0.0", "listen ip address")
	rootCmd.Flags().Uint16VarP(&rootCFG.TCPListenPort, "http-port", "p", 8080, "listen http port")
	rootCmd.Flags().Uint16VarP(&rootCFG.UDPListenPort, "udp-port", "u", 7982, "listen udp port")
	rootCmd.Flags().Uint16VarP(&rootCFG.UDPBufferSize, "udp-buffer-size", "b", 2048, "udp packet buffer size")

	rootCmd.Flags().StringVar(&rootCFG.basicAuth, "user", "", "http basic authentication")
	rootCmd.Flags().StringVar(&rootCFG.zero, "zero", "", "value to accept as zero")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "sets logger to be more verbose")

	rootCmd.Flags().StringSliceVar(&rootCFG.allowedUDPClients, "allowed-udp-clients", make([]string, 0), "allow udp connection from these addresses")
}

func runServer() {

	addr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", rootCFG.listenAddr, rootCFG.UDPListenPort),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("cannot resolve udp listen address: %s", err))
		return
	}

	wg := new(sync.WaitGroup)
	buf := vh.NewValueHolder(rootCFG.zero)

	// Only single wait group so if any of udp server or http server is down the process should shutdown immediately
	// and restart itself
	wg.Add(1)

	go creteUdpServer(wg, addr, buf)
	go createHttpServer(wg, buf)

	wg.Wait()
}

/*
 * createHttpServer starts and continuously restarts the HTTP server.
 * It takes a sync.WaitGroup and a pointer to a ValueHolder instance as parameters.
 * The sync.WaitGroup is used to manage the goroutines and ensure they finish before the main function exits.
 * The ValueHolder instance is used to store and retrieve data from the HTTP server.
 * The function will keep restarting the HTTP server in a loop, with a one-second delay after each restart.
 * If an error occurs during the starting/listening phase of the HTTP server, it will be logged as a warning and the server will be restarted.
 * The function does not return any value.
 */
func createHttpServer(wg *sync.WaitGroup, buffer *vh.ValueHolder) {
	// Unlocks the wait group if the process panics in order to shutdown the whole process
	defer func() {
		wg.Done()
	}()

	for {
		_ = hs.StartHttpServer(
			fmt.Sprintf(
				"%s:%d",
				rootCFG.listenAddr,
				rootCFG.TCPListenPort,
			),
			rootCFG.basicAuth,
			buffer,
		)
		slog.Warn("an error occurred during tcp server's starting/listening phase restarting the server after one second")
		time.Sleep(time.Second)
	}
}

/*
 * createUdpServer starts and continuously restarts the UDP server.
 * It takes a sync.WaitGroup, a pointer to a net.UDPAddr instance, and an io.Writer instance as parameters.
 * The sync.WaitGroup is used to manage the goroutines and ensure they finish before the main function exits.
 * The net.UDPAddr instance is used to specify the address and port for the UDP server to listen on.
 * The io.Writer instance is used to write data received from the UDP server.
 * The function will keep restarting the UDP server in a loop, with a one-second delay after each restart.
 * If an error occurs during the starting/listening phase of the UDP server, it will be logged as a warning and the server will be restarted.
 * If a fatal error occurs during the starting phase of the UDP server, the function will panic and exit.
 * The function does not return any value.
 */
func creteUdpServer(wg *sync.WaitGroup, addr *net.UDPAddr, buf io.Writer) {
	defer func() {
		wg.Done()
	}()
	for {
		err := udp.StartUdpServer(addr, rootCFG.allowedUDPClients, buf)
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
			panic(err)
		}
	}
}
