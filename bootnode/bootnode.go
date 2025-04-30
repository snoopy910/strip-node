package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	var listenHost, portStr string
	var ok bool

	if listenHost, ok = os.LookupEnv("LISTEN_HOST"); !ok {
		listenHost = "0.0.0.0"
	}

	if portStr, ok = os.LookupEnv("PORT"); !ok {
		portStr = "4001"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[*] Listening on: %s with port: %d\n", listenHost, port)

	filePath := "bootnode.bin"

	var prvKey crypto.PrivKey

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("Loading private key from file")
		buff, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		prvKey, err = crypto.UnmarshalPrivateKey(buff)
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Generating new private key")
		priv, _, err := crypto.GenerateKeyPair(
			crypto.RSA,
			2048,
		)
		if err != nil {
			panic(err)
		}

		buff, err := crypto.MarshalPrivateKey(priv)

		if err != nil {
			panic(err)
		}

		err = os.WriteFile(filePath, buff, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", listenHost, port))
	if err != nil {
		panic(err)
	}

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
		libp2p.ForceReachabilityPublic(),
	)

	if err != nil {
		panic(err)
	}

	_, err = dht.New(context.Background(), host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", listenHost, port, host.ID())
	select {}
}
