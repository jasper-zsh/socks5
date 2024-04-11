package socks5

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/juju/errors"
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
	conn, err := c.negotiate(c.options.Addr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	udp := NewUDPAssociateConn(conn)
	err = udp.connect(localAddr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return udp, nil
}

func (c *Client) negotiate(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	hdr := MethodSelectionRequestHeader{
		Version:  Version5,
		NMethods: 1,
	}
	req := bytes.NewBuffer(make([]byte, 0, binary.Size(hdr)+1))
	err = binary.Write(req, binary.BigEndian, hdr)
	if err != nil {
		conn.Close()
		return nil, errors.Trace(err)
	}
	req.WriteByte(MethodNoAuth)
	_, err = conn.Write(req.Bytes())
	if err != nil {
		conn.Close()
		return nil, errors.Trace(err)
	}
	res := MethodSelectionResponse{}
	err = binary.Read(conn, binary.BigEndian, &res)
	if err != nil {
		conn.Close()
		return nil, errors.Trace(err)
	}
	if res.Method == MethodNoAcceptableMethods {
		conn.Close()
		return nil, errors.Errorf("no acceptable methods")
	}
	return conn, nil
}
