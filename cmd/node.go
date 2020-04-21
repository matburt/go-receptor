package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/project-receptor/go-receptor/connection"
	"github.com/project-receptor/go-receptor/message"
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
	fmt.Println("Running as node id", viper.GetString("node_id"))
	if viper.GetBool("debug") {
		fmt.Println("Listening on", listenAddresses)
		fmt.Println("Connecting to peers", viper.GetStringSlice("peer"))
	}
	for _, listenInterface := range listenAddresses {
		go connection.Listen(listenInterface, onboardConnection)
	}
	for _, peer := range viper.GetStringSlice("peer") {
		go connection.Open(peer, onboardConnection)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c // wait for interrupt
}

func onboardConnection(rw *bufio.ReadWriter, conn net.Conn) {
	fmt.Println("Connection established", conn.RemoteAddr().String())
	helloMessage, err := message.MakeHiMessage()
	if err != nil {
		fmt.Println("Error building Hi message", err)
		return
	}
	fmt.Println("Built Hi message", helloMessage)
	bbuffer := new(bytes.Buffer)
	helloMessage.Serialize(bbuffer)
	rw.Write(bbuffer.Next(bbuffer.Len()))
	ferror := rw.Flush()
	if ferror != nil {
		fmt.Println("Flush error", ferror)
	}
	fmt.Println("Wrote bytes")
	fmt.Println("Received bytes")
	for {
		lr := io.LimitReader(rw.Reader, 16482)
		if _, err := io.Copy(os.Stdout, lr); err != nil {
			time.Sleep(1)
		}
	}
}

func init() {
	nodeCmd.Flags().StringArrayVarP(&listenAddresses, "listen", "l", nil, "Address to listen for peer connections")
	rootCmd.AddCommand(nodeCmd)
}
