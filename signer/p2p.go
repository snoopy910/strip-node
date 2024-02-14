package signer

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

var topic *pubsub.Topic
var topicNameFlag = "renode"

func initDHT(ctx context.Context, h host.Host, bootnode []multiaddr.Multiaddr) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range bootnode {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func discoverPeers(h host.Host, bootnode []multiaddr.Multiaddr) {
	ctx := context.Background()
	kademliaDHT := initDHT(ctx, h, bootnode)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, topicNameFlag)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		//fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, topicNameFlag)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peer)
			if err != nil {
				fmt.Println("Failed connecting to ", peer.ID.Pretty(), ", error:", err)
			} else {
				fmt.Println("Connected to:", peer.ID.Pretty())
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
}

func subscribe(h host.Host) {
	ctx := context.Background()

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	topic, err = ps.Join(topicNameFlag)
	if err != nil {
		panic(err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	handleMessageFromSub(ctx, sub)
}

func handleMessageFromSub(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}

		go handleIncomingMessage(m.Message.Data)
	}
}

func createHost(listenHost string, port int, bootnodeURL string) (host.Host, multiaddr.Multiaddr) {
	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", listenHost, port)))
	if err != nil {
		panic(err)
	}

	addr, err := multiaddr.NewMultiaddr(bootnodeURL)
	if err != nil {
		panic(err)
	}

	return h, addr
}
