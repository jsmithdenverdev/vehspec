package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

const (
	PROCESS_BATCH_TOPIC   = "process_batch"
	BATCH_PROCESSED_TOPIC = "batch_processed"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

func main() {
	router, err := message.NewRouter(message.RouterConfig{}, logger)

	if err != nil {
		panic(err)
	}

	channel := gochannel.NewGoChannel(gochannel.Config{}, logger)

	go batchData(channel)(batchDataInput{
		path:      "./database.json",
		batchSize: 100,
	})

	router.AddHandler(
		"batch_processor",
		PROCESS_BATCH_TOPIC,
		channel,
		BATCH_PROCESSED_TOPIC,
		channel,
		func(msg *message.Message) ([]*message.Message, error) {
			msg.Ack()
			fmt.Printf("processing batch %s\n", msg.UUID)
			var records []Vehicle
			if err := json.Unmarshal(msg.Payload, &records); err != nil {
				return nil, err
			}

			for i := 0; i <= len(records); i++ {
				time.Sleep(100 * time.Millisecond)
			}

			fmt.Printf("processed batch %s\n", msg.UUID)
			return nil, nil
		},
	)

	ctx := context.Background()
	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
