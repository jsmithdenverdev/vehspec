package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type batchDataInput struct {
	path      string
	batchSize int
}

func batchData(publisher message.Publisher) func(input batchDataInput) {
	return func(input batchDataInput) {
		var data []Vehicle

		f, err := os.Open(input.path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		b, err := io.ReadAll(f)
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
			if i%input.batchSize == 0 {
				publisher.Publish(PROCESS_BATCH_TOPIC, message.NewMessage(watermill.NewUUID(), records))
			} else {
				if len(data)-i < input.batchSize {
					publisher.Publish(PROCESS_BATCH_TOPIC, message.NewMessage(watermill.NewUUID(), records))
					return
				}
			}
		}
	}
}
