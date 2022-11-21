package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
)

const discoveryTag = "discoveryTest"

func main() {

	ctx := context.Background()
	// To construct a simple host with all the default settings, just use `New`
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello World, my hosts ID is %s\n", host.ID().ShortString())

	peerChan := initMDNS(host, discoveryTag)
	for { // allows multiple peers to join
		peer := <-peerChan // will block untill we discover a peer
		fmt.Println("Found peer:", peer.ID.ShortString(), ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}

}
