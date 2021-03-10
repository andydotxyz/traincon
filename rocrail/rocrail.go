package rocrail

import (
	"fmt"
	"net"
)

// Connection is the main type for communicating with a Rocrail server
type Connection struct {
	host string
	port int
	conn net.Conn
}

// Connect establishes a connection to a Rocrail server and returns a Connection.
// If host or port are zero it will default to localhost and 8051.
// Should an error occur it will return an error and a nil Connection instead
func Connect(host string, port int) (*Connection, error) {
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 8051
	}
	ret := &Connection{host: host, port: port}
	err := ret.connect()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// SendXML is the low-level call that allows XML based content to be sent to the server.
// Normally this will not be used as the client formats messages, but it may be useful for custom actions.
func (c *Connection) SendXML(cmd, xml string) error {
	str := "<xmlh><xml size=\"%d\" name=\"%s\"/></xmlh>%s"
	_, err := fmt.Fprintf(c.conn, str, len(xml), cmd, xml)
	return err
}

func (c *Connection) connect() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		c.conn = conn
	}
	return err
}
