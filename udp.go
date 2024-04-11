package socks5

import (
	"bytes"
	"encoding/binary"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/juju/errors"
)

var _ net.PacketConn = (*UDPAssociateConn)(nil)

type UDPAssociateConn struct {
	socksConnection
	addr      string
	localAddr *net.UDPAddr
	socksConn net.Conn
	udpConn   *net.UDPConn
	udpAddr   *net.UDPAddr

	lock    sync.Mutex
	ticker  *time.Ticker
	counter uint64
}

func NewUDPAssociateConn(proxyAddr string, localAddr *net.UDPAddr) *UDPAssociateConn {
	return &UDPAssociateConn{
		addr:      proxyAddr,
		localAddr: localAddr,
	}
}

func (u *UDPAssociateConn) disconnect() {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.ticker != nil {
		u.ticker.Stop()
	}
	u.counter = 0
	if u.udpConn != nil {
		u.udpConn.Close()
		u.udpConn = nil
	}
	if u.socksConn != nil {
		u.socksConn.Close()
		u.socksConn = nil
	}
}

func (u *UDPAssociateConn) connect() error {
	var err error
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.udpConn != nil {
		return nil
	}
	u.socksConn, err = u.negotiate(u.addr)
	if err != nil {
		return errors.Trace(err)
	}
	ip := [4]byte{}
	if u.localAddr.IP != nil {
		ip = u.localAddr.AddrPort().Addr().As4()
	}
	req := Request{
		Version:  Version5,
		Command:  CommandUDPAssociate,
		AddrType: AddrTypeIPv4,
		DstAddr:  ip,
		DstPort:  uint16(u.localAddr.Port),
	}
	err = binary.Write(u.socksConn, binary.BigEndian, req)
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
		u.udpConn, err = net.ListenUDP("udp", u.localAddr)
		if err != nil {
			return errors.Trace(err)
		}
		u.udpAddr = &net.UDPAddr{
			IP:   res.BindAddr[:],
			Port: int(res.BindPort),
		}
		u.ticker = time.NewTicker(time.Second * 10)
		go func() {
			for _ = range u.ticker.C {
				if u.counter == 0 {
					u.disconnect()
				} else {
					u.counter = 0
				}
			}
		}()
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
	if u.udpConn == nil {
		err = u.connect()
		if err != nil {
			err = errors.Trace(err)
			u.disconnect()
			return
		}
	}
	buf := make([]byte, len(p))
	n, _, err = u.udpConn.ReadFrom(buf)
	if err != nil {
		err = errors.Trace(err)
		u.disconnect()
		return
	}
	u.counter += 1

	hdr := UDPRequestHeader{}
	err = binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &hdr)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	addr = &net.UDPAddr{
		IP:   hdr.DstAddr[:],
		Port: int(hdr.DstPort),
	}
	actualLen := n - headerSize
	copy(p, buf[headerSize:n])
	n = actualLen
	return
}

// WriteTo implements net.PacketConn.
func (u *UDPAssociateConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if u.udpConn == nil {
		err = u.connect()
		if err != nil {
			err = errors.Trace(err)
			u.disconnect()
			return
		}
	}
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
	_, err = u.udpConn.WriteTo(buf.Bytes(), u.udpAddr)
	if err != nil {
		err = errors.Trace(err)
		u.disconnect()
		return
	}
	u.counter += 1
	n = len(p)
	return
}

// Close implements net.Conn.
func (u *UDPAssociateConn) Close() error {
	u.disconnect()
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
