package cmd

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/project-receptor/go-receptor/connection"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Run a Receptor Node",
		Run:   nodeRun,
	}
	listenAddresses []string
)

func nodeRun(cmd *cobra.Command, args []string) {
	if viper.GetBool("debug") {
		fmt.Println("Running the run command and listening on", listenAddresses)
		fmt.Println("Connecting to peers", viper.GetStringSlice("peer"))
	}
	for _, listenInterface := range listenAddresses {
		go connection.Listen(listenInterface, acceptConnection)
	}
	for _, peer := range viper.GetStringSlice("peer") {
		go connection.Open(peer, acceptConnection)
	}
	for {
		time.Sleep(1)
	}
}

func acceptConnection(rw *bufio.ReadWriter, conn net.Conn) {
	fmt.Println("Connection established", conn)
}

func init() {
	nodeCmd.Flags().StringArrayVarP(&listenAddresses, "listen", "l", nil, "Address to listen for peer connections")
	rootCmd.AddCommand(nodeCmd)
}
