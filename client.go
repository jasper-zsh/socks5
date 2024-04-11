package socks5

import "net"

type ClientOptions struct {
	Addr string
}

type Client struct {
	options ClientOptions
}

func (c *Client) UDPAssociate(localAddr net.UDPAddr) (net.Conn, error) {
	return nil, nil
}
