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
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/exporter"
)

// StorageExporterService is a service that provides API for exporting storage data.
type StorageQueueExporterService struct {
	runner  *Runner
	brokers []string
	topic   string

	kafkaProducer *kafka.Writer
	parent        context.Context
}

// NewStorageExporterService creates new StorageExporter service instance.
func NewStorageQueueExporterService(runner *Runner, brokers []string, topic string) (*StorageQueueExporterService, error) {
	// connect to kafka
	config := kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer: &kafka.Dialer{
			Timeout:  10 * time.Second,
			ClientID: "my-kafka-client",
		},
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	kafkaProducer := kafka.NewWriter(config)
	parent := context.Background()

	return &StorageQueueExporterService{runner: runner, brokers: brokers, topic: topic, kafkaProducer: kafkaProducer, parent: parent}, nil
}

func (s *StorageQueueExporterService) Close() error {
	err := s.kafkaProducer.Close()
	if err != nil {
		return err
	}

	s.parent.Done()

	return nil
}

func (s *StorageQueueExporterService) send(result *core.StorageExportResult, inslog core.Logger) error {
	msgs := []kafka.Message{}
	for k, pd := range result.Data {
		value, err := json.Marshal(pd.([]*exporter.PulseData))
		if err != nil {
			return errors.Wrap(err, "Json marshal error")
		} else {
		}
		msgs = append(msgs,
			kafka.Message{
				Key:   []byte(k),
				Value: value,
			})
	}

	err := s.kafkaProducer.WriteMessages(s.parent, msgs...) //io: read/write on closed pipeQueueExporterError

	return err
}

func (s *StorageQueueExporterService) Exporter(r *http.Request, args *StorageExporterArgs, reply *StorageExporterReply) error {
	ctx := context.TODO()
	inslog := inslogger.FromContext(ctx)
	s.QueueExporter(ctx, inslog)
	return nil
}

func (s *StorageQueueExporterService) QueueExporter(ctx context.Context, inslog core.Logger) {
	var err error

	defer func() {
		err := s.Close()
		if err != nil {
			inslog.Error("QueueExporter error on Close: ", err)
		}
	}()
	time.Sleep(30 * 1000 * 1000 * 1000) // 30 sec

	exp := s.runner.StorageExporter

	currentPulse, err := exp.GetCurrentPulse(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to get current pulse data")
		return
	}
	previousPulse := currentPulse.PulseNumber

	for true {
		currentPulse, err := exp.GetCurrentPulse(ctx)
		if err != nil {
			err = errors.Wrap(err, "failed to get current pulse data")
			return
		}

		dif := int(currentPulse.PulseNumber - previousPulse)
		if dif >= 60 {
			result, err := exp.Export(ctx, previousPulse, 10)
			if err != nil {
				if strings.Contains(err.Error(), "failed to fetch pulse data") {
					err = errors.Wrap(err, "Pulse not found.")
					return
				}
				err = errors.Wrap(err, "[ Export ]")
				return
			} else {
				err := s.send(result, inslog)
				if err != nil {
					err = errors.Wrap(err, "failed to send current pulse data")
					return
				}
			}
			previousPulse = previousPulse + 10
		} else {
			time.Sleep(5 * 1000 * 1000 * 1000) // 5 sec
		}
	}

	return
}
