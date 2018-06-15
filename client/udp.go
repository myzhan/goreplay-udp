package client

import (
	"log"
	"net"
	"time"
)

type UDPClient struct {
	address        string
	timeout        time.Duration
	ignoreResponse bool

	conn *net.UDPConn
}

func NewUDPClient(address string, timeout time.Duration, ignoreResponse bool) (c *UDPClient) {
	c = new(UDPClient)
	c.address = address
	c.timeout = timeout
	c.ignoreResponse = ignoreResponse

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Error initialize UDP Client %s\n", address)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Error dialing %s\n", address)
	}

	c.conn = conn
	return
}

func (c *UDPClient) Send(data []byte) (resp []byte, err error) {

	c.conn.SetReadDeadline(time.Now().Add(c.timeout))

	_, err = c.conn.Write(data)
	if err != nil {
		log.Printf("UDP Write Error: %v\n", err)
	}

	if c.ignoreResponse {
		return nil, nil
	}

	resp = make([]byte, 4096)
	respLength, err := c.conn.Read(resp)
	if err != nil {
		log.Printf("UDP Read Error: %v\n", err)
	}
	if len(resp) <= respLength {
		log.Printf("UDP Response may be truncated, length of response is %d\n", respLength)
	}

	return resp, err
}
