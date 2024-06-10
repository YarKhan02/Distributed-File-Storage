package p2p


// Peer is an interface that represents remote node.
type Peer interface {

}


// transport is anything that handles communication 
//between nodes in the network. (TCP, UDP, websockets ...)
type Transport interface {
	ListenAndAccept() error
}