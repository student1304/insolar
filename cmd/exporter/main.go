package main

import (
	"fmt"
	"github.com/insolar/insolar/api"
)

func main() {

	err := api.QueueExporter()
	if err != nil {
		fmt.Println(err)
	}
}
