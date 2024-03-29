package services

import (
	"testing"

	"github.com/simpleanalytics/cq-source-simple-analytics/client"
)

func TestEvents(t *testing.T) {
	ts := testServer(t, "events.ndjson")
	defer ts.Close()
	client.TestHelper(t, Events(), ts)
}
