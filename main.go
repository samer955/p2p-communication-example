package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

const discoveryTag = "discoveryTest"
const topic = "greet"

func main() {

	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	flag.Parse()

	ctx := context.Background()
	// To construct a simple host with all the default settings, just use `New`
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	nick := *nickFlag
	if len(nick) == 0 {
		nick = host.ID().ShortString()
	}

	fmt.Printf("Hello LAN, my hosts ID is %s, my Nickname: %s \n", host.ID().ShortString(), nick)

	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	testTopic, _ := pubsub.Join(topic)
	topicSubscr, _ := testTopic.Subscribe()

	// allows multiple peers to join
	// will block untill we discover a peer

	go findPeers(host, ctx)
	go receive(host, ctx, topicSubscr)
	go greet(host, nick, ctx, testTopic)

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

}

func findPeers(host host.Host, ctx context.Context) {
	peerChan := initMDNS(host, discoveryTag)
	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer.ID.ShortString(), ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}
}

func greet(host host.Host, nick string, ctx context.Context, topic *pubsub.Topic) {

	for {
		if host.Peerstore().Peers().Len() == 0 {
			continue
		}
		topic.Publish(ctx, []byte("Hello peer from "+nick))
		time.Sleep(5 * time.Second)
	}
}

func receive(host host.Host, ctx context.Context, subscr *pubsub.Subscription) {

	for {
		m, _ := subscr.Next(ctx)
		if m.ReceivedFrom.ShortString() == host.ID().ShortString() {
			continue
		}
		fmt.Printf("Received message from <%s: %s> \n", m.ReceivedFrom.ShortString(), m.Data)
	}
}
