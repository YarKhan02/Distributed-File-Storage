package p2p

import (
	// "bytes"
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents a remote node over a established TCP connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a connection => outbound = true
	// if we accept and retrieve a connection => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer {
		conn: conn,
		outbound: outbound,
	}
}

// close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	HandshakeFunc HandshakeFunc
	Decoder Decoder
	ListenAddr string
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	rpcch chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct {}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	rpc := RPC{}

	// Read for loop
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
		fmt.Printf("message: %+v\n", rpc)
	}
}