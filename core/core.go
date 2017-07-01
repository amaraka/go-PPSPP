package core

import (
	"github.com/golang/glog"
	ma "github.com/multiformats/go-multiaddr"
)

// FIXME: where should this go?
const proto = "/goppspp/0.0.1"

// PeerID identifies a peer
type PeerID interface {
	String() string
}

// ChanID identifies a channel
type ChanID uint32

// Datagram holds a protocol datagram
type Datagram struct {
	ChanID ChanID
	Msgs   []Msg
}

// Peer implements protocol logic and underlying network
type Peer struct {

	// P handles protocol logic.
	P Protocol

	// n handles the underlying network.
	// private because no one should touch this except P
	n Network
}

// NewPeer makes and initializes a new peer
func NewPeer(port int) *Peer {

	prot := newPpspp()

	// This determines the network implementation (libp2p)
	n := newLibp2pNetwork(port)

	p := Peer{n: n, P: prot}

	// set the network's datagram handler
	p.n.SetDatagramHandler(prot.HandleDatagram)

	p.P.SetDatagramSender(n.SendDatagram)

	return &p
}

// ID returns the peer ID
func (p *Peer) ID() PeerID {
	return p.n.ID()
}

// AddAddrs adds multiaddresses for the remote peer to this peer's store
func (p *Peer) AddAddrs(remote PeerID, addrs []ma.Multiaddr) {
	p.n.AddAddrs(remote, addrs)
}

// Addrs returns multiaddresses for this peer
func (p *Peer) Addrs() []ma.Multiaddr {
	return p.n.Addrs()
}

// Connect creates a stream from p to the peer at id and sets a stream handler
func (p *Peer) Connect(id PeerID) error {
	glog.Infof("%s: Connecting to %s", p.ID(), id)
	return p.n.Connect(id)
}

// Disconnect closes the stream that p is using to connect to the peer at id
func (p *Peer) Disconnect(id PeerID) error {
	glog.Infof("%s: Disconnecting from %s", p.ID(), id)
	return p.n.Connect(id)
}
