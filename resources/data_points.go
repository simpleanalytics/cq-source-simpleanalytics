package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/cq-source-simple-analytics/client"
	"github.com/cloudquery/cq-source-simple-analytics/internal/simpleanalytics"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"golang.org/x/sync/errgroup"
)

const tableDataPoints = "simple_analytics_data_points"

func DataPoints() *schema.Table {
	return &schema.Table{
		Name:        tableDataPoints,
		Description: "https://docs.simpleanalytics.com/api/export-data-points",
		Resolver:    fetchDataPoints,
		Multiplex:   client.WebsiteMultiplex,
		Transform: transformers.TransformWithStruct(
			&simpleanalytics.DataPoint{},
			transformers.WithPrimaryKeys("UUID", "Hostname"),
		),
		Columns: []schema.Column{
			{
				Name:     "metadata",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Metadata"),
			},
		},
		IsIncremental: true,
	}
}

func fetchDataPoints(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	c := meta.(*client.Client)

	// Set start time according to these priorities:
	// 1. backend state
	// 2. start_time from plugin spec (which defaults to 2018)
	start := c.Spec.StartTime()
	if c.Backend != nil {
		value, err := c.Backend.Get(ctx, tableDataPoints, c.ID())
		if err != nil {
			return fmt.Errorf("failed to get cursor from backend: %w", err)
		}
		if value != "" {
			c.Logger.Info().Str("cursor", value).Msg("cursor found")
			start, err = time.Parse(client.AllowedTimeLayout, value)
			if err != nil {
				return fmt.Errorf("failed to parse cursor from backend: %w", err)
			}
		}
	}
	end := c.Spec.EndTime()
	c.Logger.Info().Time("start", start).Time("end", end).Msg("fetching data points")

	// Stream data points from Simple Analytics, from start time to now.
	fields := make([]string, len(simpleanalytics.AllExportFields))
	copy(fields, simpleanalytics.AllExportFields)
	for _, field := range c.Website.MetadataFields {
		fields = append(fields, "metadata."+field)
	}
	opts := simpleanalytics.ExportOptions{
		Hostname: c.Website.Hostname,
		Start:    start,
		End:      end,
		Fields:   fields,
	}
	g, gctx := errgroup.WithContext(ctx)
	var ch = make(chan simpleanalytics.DataPoint)
	g.Go(func() error {
		defer close(ch)
		return c.SAClient.Export(gctx, opts, ch)
	})
	for range ch {
		v := <-ch
		if v.UUID == "" {
			// ignore values without UUID
			continue
		}
		res <- v
	}
	err := g.Wait()
	if err != nil {
		return fmt.Errorf("failed to fetch data points: %w", err)
	}

	// Save cursor state to the backend.
	if c.Backend != nil {
		// We subtract WindowOverlapSeconds from the end time to allow delayed data points
		// to be fetched on the next sync. This will cause some duplicates, but
		// allows us to guarantee at-least-once delivery. Duplicates can be removed
		// by using overwrite-delete-stale write mode, by de-duplicating in queries,
		// or by running a post-processing step.
		newCursor := end.Add(time.Duration(c.Spec.WindowOverlapSeconds) * time.Second).Format(client.AllowedTimeLayout)
		err = c.Backend.Set(ctx, tableDataPoints, c.ID(), newCursor)
		if err != nil {
			return fmt.Errorf("failed to save cursor to backend: %w", err)
		}
		c.Logger.Info().Str("cursor", newCursor).Msg("cursor updated")
	}
	return nil
}
