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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsonrpc "github.com/gorilla/rpc/v2/json2"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

// StorageExporterService is a service that provides API for exporting storage data.
type StorageQueueExporterService struct {
	runner  *Runner
	brokers []string
	topic   string
	writer  *kafka.Writer
}

// NewStorageExporterService creates new StorageExporter service instance.
func NewStorageQueueExporterService(runner *Runner, brokers []string, topic string) *StorageQueueExporterService {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &StorageQueueExporterService{runner: runner, brokers: brokers, topic: topic, writer: w}
}

func (s *StorageQueueExporterService) send(result *core.StorageExportResult) error {
	msgs := []kafka.Message{}
	for k, v := range result.Data {
		msgs = append(msgs,
			kafka.Message{
				Key:   []byte(k),
				Value: v.([]byte),
			})
	}

	err := s.writer.WriteMessages(context.Background(), msgs...)

	return err
}

func (s *StorageQueueExporterService) Exporter(r *http.Request, args *StorageExporterArgs, reply *StorageExporterReply) error {
	s.QueueExporter()
	return nil
}

func (s *StorageQueueExporterService) QueueExporter() error {
	defer s.writer.Close()

	exp := s.runner.StorageExporter
	ctx := context.TODO()
	inslog := inslogger.FromContext(ctx)

	currentPulse, err := exp.GetCurrentPulse(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current pulse data")
	}
	previousPulse := currentPulse

	for true {
		currentPulse, err := exp.GetCurrentPulse(ctx)
		inslog.Info("currentPulse  ..." + strconv.Itoa(int(currentPulse.PulseNumber)))
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
				err := s.send(result)
				if err != nil {
					return errors.Wrap(err, "failed to send current pulse data")
				}
			}
			previousPulse = currentPulse
		} else {
			inslog.Info("QueueExporter sleep   ")
			time.Sleep(10000)
		}
	}

	return nil
}
