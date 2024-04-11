package socks5

const (
	CommandConnect      = 0x01
	CommandBind         = 0x02
	CommandUDPAssociate = 0x03

	AddrTypeIPv4       = 0x01
	AddrTypeDomainName = 0x03
	AddrTypeIPv6       = 0x04
	AddrLenIPv4        = 4
	AddrLenIPv6        = 16
)

type Request struct {
	Version  byte
	Command  byte
	Reserved byte
	AddrType byte
	DstAddr  []byte
	DstPort  [2]byte
}

type Reply struct {
	Version  byte
	Reply    byte
	Reserved byte
	AddrType byte
	BindAddr []byte
	BindPort [2]byte
}

type UDPRequest struct {
	Reserved [2]byte
	Fragment byte
	AddrType byte
	DstAddr  []byte
	DstPort  [2]byte
	Data     []byte
}
