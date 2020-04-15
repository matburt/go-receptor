package connection

import (
	"bufio"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Open returns a connection to a Receptor Peer
// NOTE: This doesn't handle errors gracefully or attempt to reconnect
func Open(address string, handler AcceptFunc) error {
	if viper.GetBool("debug") {
		fmt.Println("Connecting to peer at address", address)
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return errors.Wrap(err, "Connecting "+address+" failed")
	}
	go handler(
		bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		conn,
	)
	return nil
}
