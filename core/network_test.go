package core

import (
	"flag"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/glog"
)

func TestConnect(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	// Set up and connect two peers
	_, _, err := setupAndConnectPeers(46)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDisonnect(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	// Set up and connect two peers
	p1, p2, err := setupAndConnectPeers(46)
	if err != nil {
		t.Fatal(err)
	}

	// Disconnect the two peers
	err = p1.n.Disconnect(p2.ID())
	if err != nil {
		t.Fatal(err)
	}
	err = p2.n.Disconnect(p1.ID())
	if err != nil {
		t.Fatal(err)
	}

	// Try to send a datagram
	c := ChanID(5)
	start := ChunkID(3)
	end := ChunkID(5)
	msend, err := messagize(HaveMsg{Start: start, End: end})
	if err != nil {
		t.Fatal(err)
	}
	dsend := datagramize(c, msend)
	err = p1.n.SendDatagram(*dsend, p2.ID())
	if err == nil {
		t.Fatal("SendDatagram should fail")
	} else {
		glog.Infof("SendDatagram returned error as expected: %v", err)
	}
}

func TestSendDatagram(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	// Set up and connect two peers
	p1, p2, err := setupAndConnectPeers(36)
	if err != nil {
		t.Fatal(err)
	}

	// SendDatagram p1 -> p2
	c := ChanID(5)
	start := ChunkID(3)
	end := ChunkID(5)
	msend, err := messagize(HaveMsg{Start: start, End: end})
	if err != nil {
		t.Fatal(err)
	}
	dsend := datagramize(c, msend)
	p1.n.SendDatagram(*dsend, p2.ID())

	// Sleep and then check that p2 received the datagram
	time.Sleep(1 * time.Second)
	prot2, ok := p2.P.(*StubProtocol)
	if !ok {
		t.Fatal("type assertion failed")
	}
	if num := prot2.NumHandledDatagrams(); num != 1 {
		t.Fatalf("should have handled 1 datagram, got %d", num)
	}
	drecv, idrecv, err := prot2.ReadHandledDatagram()
	if err != nil {
		t.Fatal(err)
	}
	if idrecv != p1.ID() {
		t.Errorf("should have received from id=%v, got %v", p1.ID(), idrecv)
	}
	if !reflect.DeepEqual(drecv, dsend) {
		t.Errorf("drecv != dsend: drecv=%v, dsend=%v", drecv, dsend)
	}
}

// setupAndConnectPeers creates two libp2p peers with StubProtocol and connects them
func setupAndConnectPeers(seed int64) (*Peer, *Peer, error) {
	// Set up and connect two peers
	rand.Seed(seed)
	port1 := rand.Intn(100) + 10000
	prot1 := newStubProtocol()
	p1, err := NewLibp2pPeer(port1, prot1)
	if err != nil {
		return nil, nil, err
	}
	port2 := port1 + 1
	prot2 := newStubProtocol()
	p2, err := NewLibp2pPeer(port2, prot2)
	if err != nil {
		return nil, nil, err
	}
	peerExchangeIDAddr(p1, p2)
	err = p1.n.Connect(p2.ID())
	if err != nil {
		return nil, nil, err
	}
	err = p2.n.Connect(p1.ID())
	if err != nil {
		return nil, nil, err
	}

	return p1, p2, nil
}

// peerExchangeIDAddr magic exchange of peer IDs and addrs
func peerExchangeIDAddr(p1 *Peer, p2 *Peer) {
	addrs1 := p1.Addrs()
	addrs2 := p2.Addrs()
	p1.AddAddrs(p2.ID(), addrs2)
	p2.AddAddrs(p1.ID(), addrs1)
}
