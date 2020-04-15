package connection

import (
	"bufio"
	"net"
)

// AcceptFunc accepts a connection and handles it
type AcceptFunc func(*bufio.ReadWriter, net.Conn)
