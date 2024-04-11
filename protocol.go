package socks5

const (
	Version5 = 0x05

	MethodNoAuth              = 0x00
	MethodGSSAPI              = 0x01
	MethodUserPass            = 0x02
	MethodNoAcceptableMethods = 0xFF

	CommandConnect      = 0x01
	CommandBind         = 0x02
	CommandUDPAssociate = 0x03

	AddrTypeIPv4       = 0x01
	AddrTypeDomainName = 0x03
	AddrTypeIPv6       = 0x04
	AddrLenIPv4        = 4
	AddrLenIPv6        = 16

	ReplySucceed                   = 0x00
	ReplyGeneralSocksServerFailure = 0x01
	ReplyConnectionNowAllowed      = 0x02
	ReplyNetworkUnreachable        = 0x03
	ReplyHostUnreachable           = 0x04
	ReplyConnectionRefused         = 0x05
	ReplyTTLExpired                = 0x06
	ReplyCommandNotSupported       = 0x07
	ReplyAddressTypeNotSupported   = 0x08
)

type MethodSelectionRequestHeader struct {
	Version  byte
	NMethods byte
}

type MethodSelectionResponse struct {
	Version byte
	Method  byte
}

type Request struct {
	Version  byte
	Command  byte
	Reserved byte
	AddrType byte
	DstAddr  [4]byte
	DstPort  uint16
}

type Reply struct {
	Version  byte
	Reply    byte
	Reserved byte
	AddrType byte
	BindAddr [4]byte
	BindPort uint16
}

type UDPRequestHeader struct {
	Reserved [2]byte
	Fragment byte
	AddrType byte
	DstAddr  [4]byte
	DstPort  uint16
}
