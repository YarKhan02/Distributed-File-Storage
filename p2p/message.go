package p2p

import "net"

// message holds any arbitary data that is being sent over the
// each transport between two nodes in the network
type RPC struct {
	From net.Addr
	Payload []byte
} 