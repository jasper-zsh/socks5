package socks5

import (
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	localAddr = "100.64.0.11:22344"
)

func TestKeepalive(t *testing.T) {
	proxy := NewClient(ClientOptions{
		Addr: "100.64.0.9:41080",
	})
	sAddr, _ := netip.ParseAddrPort(localAddr)
	cAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 22345,
	}

	s, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(sAddr))
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, s) {
		return
	}

	c, err := proxy.UDPAssociate(cAddr)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	send := func() {
		payload := []byte("foo")
		n, err := c.WriteTo(payload, net.UDPAddrFromAddrPort(sAddr))
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, len(payload), n) {
			return
		}
		buf := make([]byte, 4096)
		n, _, err = s.ReadFrom(buf)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Equal(t, len(payload), n) {
			return
		}
		if !assert.Equal(t, payload[:n], buf[:n]) {
			return
		}
	}

	send()

	time.Sleep(20 * time.Second)

	send()
}

func TestUDPAssociate(t *testing.T) {
	proxy := NewClient(ClientOptions{
		Addr: "100.64.0.9:41080",
	})
	// sAddr := &net.UDPAddr{
	// 	IP:   net.IPv4zero,
	// 	Port: 22344,
	// }
	sAddr, _ := netip.ParseAddrPort(localAddr)
	cAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 22345,
	}

	s, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(sAddr))
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, s) {
		return
	}

	c, err := proxy.UDPAssociate(cAddr)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	payload := []byte("foo")
	n, err := c.WriteTo(payload, net.UDPAddrFromAddrPort(sAddr))
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, len(payload), n) {
		return
	}
	buf := make([]byte, 4096)
	n, addr, err := s.ReadFrom(buf)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, len(payload), n) {
		return
	}
	if !assert.Equal(t, payload[:n], buf[:n]) {
		return
	}

	payload = []byte("bar")
	n, err = s.WriteTo(payload, addr)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, len(payload), n) {
		return
	}

	n, _, err = c.ReadFrom(buf)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, len(payload), n) {
		return
	}
	if !assert.Equal(t, payload[:n], buf[:n]) {
		return
	}
}
