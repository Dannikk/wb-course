package nats

import (
	"fmt"
	"log"
	"ordermngmt/internal/config"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

type Subscription struct {
	NatsConn *nats.Conn
	StanConn stan.Conn
	Sub      stan.Subscription
}

func NewSubscription(cfg config.StanCfg, data chan<- []byte) (*Subscription, error) {
	URL := fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Port)

	opts := []nats.Option{nats.Name("NATS Streaming Example Subscriber")}

	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("while nats.Connect() %w", err)
	}

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost, reason: %v", reason)
		}))

	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("can't connect: %w.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}

	startOpt, err := getStartOpt(cfg)
	if err != nil {
		sc.Close()
		nc.Close()
		return nil, fmt.Errorf("can't getStartOpt: %w", err)
	}

	log.Printf("before sub: Subject: %v, Qgroup: %v, Durable: %v", cfg.Subject, cfg.Qgroup, cfg.Durable)
	sub, err := sc.QueueSubscribe(
		cfg.Subject,
		cfg.Qgroup,
		func(msg *stan.Msg) {
			data <- msg.Data
		},
		startOpt,
		stan.DurableName(cfg.Durable),
	)

	if err != nil {
		sc.Close()
		nc.Close()
		return nil, fmt.Errorf("can't subscribe: %w", err)
	}

	return &Subscription{nc, sc, sub}, nil
}

func (sub *Subscription) Close() {
	sub.Sub.Unsubscribe()
	sub.StanConn.Close()
	sub.NatsConn.Close()
}

func getStartOpt(cfg config.StanCfg) (stan.SubscriptionOption, error) {
	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	if cfg.StartSeq != 0 {
		startOpt = stan.StartAtSequence(cfg.StartSeq)
		log.Println("start at sequence:", cfg.StartSeq)
	} else if cfg.DeliverLast {
		startOpt = stan.StartWithLastReceived()
		log.Println("start with last received")
	} else if cfg.DeliverAll && !cfg.NewOnly {
		startOpt = stan.DeliverAllAvailable()
		log.Println("start with all available")
	} else if cfg.StartDelta != "" {
		ago, err := time.ParseDuration(cfg.StartDelta)
		if err != nil {
			log.Println("Error parsing duration:", err)
			return nil, err
		}
		startOpt = stan.StartAtTimeDelta(ago)
		log.Println("start at delta:", cfg.StartDelta)
	}

	return startOpt, nil
}
