package socks5

import (
	"bytes"
	"encoding/binary"
	"net"
	"net/netip"
	"time"

	"github.com/juju/errors"
)

var _ net.PacketConn = (*UDPAssociateConn)(nil)

type UDPAssociateConn struct {
	socksConn net.Conn
	udpConn   *net.UDPConn
}

func NewUDPAssociateConn(socksConn net.Conn) *UDPAssociateConn {
	return &UDPAssociateConn{
		socksConn: socksConn,
	}
}

func (u *UDPAssociateConn) connect(localAddr *net.UDPAddr) error {
	req := Request{
		Version:  Version5,
		Command:  CommandUDPAssociate,
		AddrType: AddrTypeIPv4,
		DstAddr:  localAddr.AddrPort().Addr().As4(),
		DstPort:  uint16(localAddr.Port),
	}
	err := binary.Write(u.socksConn, binary.BigEndian, req)
	if err != nil {
		u.socksConn.Close()
		return errors.Trace(err)
	}
	res := Reply{}
	err = binary.Read(u.socksConn, binary.BigEndian, &res)
	if err != nil {
		return errors.Trace(err)
	}
	switch res.Reply {
	case ReplySucceed:
		u.udpConn, err = net.ListenUDP("udp", localAddr)
		if err != nil {
			return errors.Trace(err)
		}
		return nil
	case ReplyCommandNotSupported:
		return errors.New("udp associate not supported")
	default:
		return errors.Errorf("socks error %x", res.Reply)
	}
}

var (
	headerSize = binary.Size(UDPRequestHeader{})
)

// ReadFrom implements net.PacketConn.
func (u *UDPAssociateConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	buf := make([]byte, len(p))
	n, addr, err = u.udpConn.ReadFrom(buf)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	actualLen := n - headerSize
	copy(p, buf[headerSize:n])
	n = actualLen
	return
}

// WriteTo implements net.PacketConn.
func (u *UDPAssociateConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	var addrPort netip.AddrPort
	addrPort, err = netip.ParseAddrPort(addr.String())
	if err != nil {
		err = errors.Trace(err)
		return
	}
	ip := addrPort.Addr().As4()
	req := UDPRequestHeader{
		AddrType: AddrTypeIPv4,
		DstAddr:  ip,
		DstPort:  addrPort.Port(),
	}
	buf := bytes.NewBuffer(make([]byte, 0, binary.Size(req)+len(p)))
	err = binary.Write(buf, binary.BigEndian, req)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	buf.Write(p)
	_, err = u.udpConn.WriteTo(buf.Bytes(), addr)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	n = len(p)
	return
}

// Close implements net.Conn.
func (u *UDPAssociateConn) Close() error {
	err := u.udpConn.Close()
	err = u.socksConn.Close()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// LocalAddr implements net.Conn.
func (u *UDPAssociateConn) LocalAddr() net.Addr {
	return u.udpConn.LocalAddr()
}

// RemoteAddr implements net.Conn.
func (u *UDPAssociateConn) RemoteAddr() net.Addr {
	return u.udpConn.RemoteAddr()
}

// SetDeadline implements net.Conn.
func (u *UDPAssociateConn) SetDeadline(t time.Time) error {
	return u.udpConn.SetDeadline(t)
}

// SetReadDeadline implements net.Conn.
func (u *UDPAssociateConn) SetReadDeadline(t time.Time) error {
	return u.udpConn.SetReadDeadline(t)
}

// SetWriteDeadline implements net.Conn.
func (u *UDPAssociateConn) SetWriteDeadline(t time.Time) error {
	return u.udpConn.SetWriteDeadline(t)
}
