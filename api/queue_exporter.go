/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package api

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/core"
	"strings"
	"time"

	jsonrpc "github.com/gorilla/rpc/v2/json2"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func send(result *core.StorageExportResult) error {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "topic-A",
		Balancer: &kafka.LeastBytes{},
	})

	msgs := []kafka.Message{}
	for k, v := range result.Data {
		msgs = append(msgs,
			kafka.Message{
				Key:   []byte(k),
				Value: v.([]byte),
			})
	}

	err := w.WriteMessages(context.Background(), msgs...)

	w.Close()

	return err
}

type StorageWrapper struct {
	StorageExporter core.StorageExporter `inject:""`
}

func QueueExporter() error {
	expWrapper := StorageWrapper{}
	exp := expWrapper.StorageExporter
	ctx := context.TODO()

	currentPulse, err := exp.GetCurrentPulse(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current pulse data")
	}
	previousPulse := currentPulse

	for true {
		currentPulse, err := exp.GetCurrentPulse(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to get current pulse data")
		}

		dif := int(currentPulse.PulseNumber - previousPulse.PulseNumber)
		if dif > 0 {
			result, err := exp.Export(ctx, previousPulse.PulseNumber, dif)
			if err != nil {
				if strings.Contains(err.Error(), "failed to fetch pulse data") {
					return &jsonrpc.Error{
						Code:    errCodePulseNotFound,
						Message: "[ Export ]: " + err.Error(),
						Data:    nil,
					}
				}
				return errors.Wrap(err, "[ Export ]")
			} else {
				err := send(result)
				if err != nil {
					return err
				}
			}
			previousPulse = currentPulse
		} else {
			time.Sleep(1000)
		}
	}

	return nil
}

// Start runs api server
func (sw *StorageWrapper) Start(ctx context.Context) error {
	fmt.Println("StorageWrapper started")
	return nil
}

// Stop stops api server
func (sw *StorageWrapper) Stop(ctx context.Context) error {
	fmt.Println("StorageWrapper stopped")

	return nil
}
