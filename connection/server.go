package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Listen on an interface:port
func Listen(address string, handler AcceptFunc) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", address)
	}
	if viper.GetBool("debug") {
		fmt.Println("Listening on", listener.Addr().String())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Accepted connection from", conn)
		go handler(
			bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
			conn,
		)
	}
}
