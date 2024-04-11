package socks5

import (
	"net"
)

type ClientOptions struct {
	Addr string
}

type Client struct {
	options ClientOptions
}

func NewClient(options ClientOptions) *Client {
	return &Client{
		options: options,
	}
}

func (c *Client) UDPAssociate(localAddr *net.UDPAddr) (net.PacketConn, error) {
	udp := NewUDPAssociateConn(c.options.Addr, localAddr)
	return udp, nil
}
