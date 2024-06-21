package p2p

// message holds any arbitary data that is being sent over the
// each transport between two nodes in the network
type Message struct {
	Payload []byte
} 