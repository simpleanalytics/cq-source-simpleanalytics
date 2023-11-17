package main

import (
	"context"
	"log"

	"github.com/cloudquery/plugin-sdk/v4/serve"
	"github.com/simpleanalytics/cq-source-simple-analytics/resources/plugin"
)

func main() {
	p := serve.Plugin(plugin.Plugin())
	if err := p.Serve(context.Background()); err != nil {
		log.Fatalf("failed to serve plugin: %v", err)
	}
}
