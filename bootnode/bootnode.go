package bootnode

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	mrand "math/rand"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

var prvKey crypto.PrivKey

func Start(listenHost string, port int, path string) {

	flag.Parse()

	fmt.Printf("[*] Listening on: %s with port: %d\n", listenHost, port)

	ctx := context.Background()
	r := mrand.New(mrand.NewSource(int64(port)))

	filePath := path + "/bootnode.bin"

	exists, err := common.FileExists(filePath)
	if err != nil {
		panic(err)
	}

	if exists {
		buff, err := ioutil.ReadFile(filePath)

		if err != nil {
			panic(err)
		}

		prvKey, err = crypto.UnmarshalPrivateKey(buff)

		if err != nil {
			panic(err)
		}
	} else {
		// Creates a new RSA key pair for this host.
		prvKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			panic(err)
		}

		buff, err := crypto.MarshalPrivateKey(prvKey)

		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(filePath, buff, 0644)
		if err != nil {
			panic(err)
		}
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", listenHost, port))

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

	_, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", listenHost, port, host.ID().Pretty())
	fmt.Println("")
	select {}
}
