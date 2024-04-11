package socks5

import (
	"bytes"
	"encoding/binary"
	"net"
	"net/netip"
	"time"

	"github.com/juju/errors"
)

type socksConnection struct {
}

func (c *socksConnection) negotiate(addr string) (*net.TCPConn, error) {
	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	conn, err := net.DialTCP("tcp", nil, net.TCPAddrFromAddrPort(addrPort))
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		conn.Close()
		return nil, errors.Trace(err)
	}
	err = conn.SetKeepAlivePeriod(10 * time.Second)
	if err != nil {
		conn.Close()
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
