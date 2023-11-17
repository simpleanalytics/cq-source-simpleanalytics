package client

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/rs/zerolog"
	"github.com/simpleanalytics/cq-source-simple-analytics/internal/simpleanalytics"
)

func TestHelper(t *testing.T, table *schema.Table, ts *httptest.Server) {
	table.IgnoreInTests = false
	t.Helper()

	l := zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	sched := scheduler.NewScheduler(scheduler.WithLogger(l))

	spec := &Spec{
		UserID: "test",
		APIKey: "test",
		Websites: []WebsiteSpec{
			{
				Hostname:       "test.com",
				MetadataFields: []string{"metadata_text", "metadata_int"},
			},
		},
	}
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		t.Fatalf("failed to validate spec: %v", err)
	}

	saClient := simpleanalytics.NewClient(spec.UserID, spec.APIKey, simpleanalytics.WithBaseURL(ts.URL), simpleanalytics.WithHTTPClient(ts.Client()))
	c := New(l, *spec, saClient, nil)

	tables := schema.Tables{table}
	if err := transformers.TransformTables(tables); err != nil {
		t.Fatal(err)
	}
	messages, err := sched.SyncAll(context.Background(), c, tables)
	if err != nil {
		t.Fatalf("failed to sync: %v", err)
	}
	plugin.ValidateNoEmptyColumns(t, tables, messages)
}
