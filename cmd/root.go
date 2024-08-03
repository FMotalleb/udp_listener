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
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "udp_listener",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}
var AllowedUDPClients []string = make([]string, 0)
var BasicAuth string = ""
var ListenAddr string = ""
var UDPListenPort uint16 = 0
var TCPListenPort uint16 = 0
var UDPBufferSize uint16 = 0

var Zero string = ""
var Verbose bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&ListenAddr, "listen-addr", "l", "0.0.0.0", "listen ip address")
	rootCmd.Flags().Uint16VarP(&TCPListenPort, "http-port", "p", 8080, "listen http port")
	rootCmd.Flags().Uint16VarP(&UDPListenPort, "udp-port", "u", 7982, "listen udp port")
	rootCmd.Flags().Uint16VarP(&UDPBufferSize, "udp-buffer-size", "b", 2048, "udp packet buffer size")

	rootCmd.Flags().StringVar(&BasicAuth, "user", "", "http basic authentication")
	rootCmd.Flags().StringVar(&Zero, "zero", "", "value to accept as zero")
	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "sets logger to be more verbose")

	rootCmd.Flags().StringSliceVar(&AllowedUDPClients, "allowd-udp-clients", make([]string, 0), "allow udp connection from these addressess")

}
