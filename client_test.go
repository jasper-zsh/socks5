package socks5

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDPAssociate(t *testing.T) {
	proxy := NewClient(ClientOptions{
		Addr: "100.64.0.9:41080",
	})
	sAddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 22344,
	}
	cAddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 22345,
	}
	s, err := proxy.UDPAssociate(sAddr)
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
	n, err := c.WriteTo(payload, sAddr)
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
}
