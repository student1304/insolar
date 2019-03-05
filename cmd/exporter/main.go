package main

import (
	"fmt"
	"github.com/insolar/insolar/api"
)

type QueueExporter struct {
	Runner *api.Runner `inject:""`
}

func main() {
	queueExporter := QueueExporter{}

	storageExporterService := api.NewStorageExporterService(queueExporter.Runner)

	err := storageExporterService.QueueExporter()
	if err != nil {
		fmt.Println(err)
	}
}
