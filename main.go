package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

const (
	PROCESS_BATCH_TOPIC   = "process_batch"
	BATCH_PROCESSED_TOPIC = "batch_processed_batch"
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

	go processRecords(channel)("./database.json")

	router.AddHandler(
		"batch_processor",
		PROCESS_BATCH_TOPIC,
		channel,
		BATCH_PROCESSED_TOPIC,
		channel,
		func(msg *message.Message) ([]*message.Message, error) {
			var records []Vehicle
			if err := json.Unmarshal(msg.Payload, &records); err != nil {
				return nil, err
			}

			for _, v := range records {
				fmt.Printf("%+v", v)
			}

			msg.Ack()
			return nil, nil
		},
	)

	ctx := context.Background()
	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}

func processRecords(publisher message.Publisher) func(path string) {
	return func(path string) {
		var data []Vehicle

		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, &data)
		if err != nil {
			panic(err)
		}

		var batch []Vehicle

		for i, v := range data {
			batch = append(batch, v)
			records, err := json.Marshal(batch)
			if err != nil {
				panic(err)
			}
			if i%10 == 0 {
				publisher.Publish(PROCESS_BATCH_TOPIC, message.NewMessage(watermill.NewUUID(), records))
			}
		}
	}
}
