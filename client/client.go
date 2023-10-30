package client

import (
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/rs/zerolog"
	"github.com/simpleanalytics/cq-source-simple-analytics/internal/simpleanalytics"
)

type Client struct {
	Logger   zerolog.Logger
	SAClient *simpleanalytics.Client
	Backend  state.Client
	Spec     Spec
	Website  WebsiteSpec
}

func (c *Client) ID() string {
	return strings.Join([]string{"simple-analytics", c.Website.Hostname}, ":")
}

func (c *Client) withWebsite(website WebsiteSpec) *Client {
	return &Client{
		Logger:   c.Logger.With().Str("hostname", website.Hostname).Logger(),
		SAClient: c.SAClient,
		Backend:  c.Backend,
		Spec:     c.Spec,
		Website:  website,
	}
}

func New(logger zerolog.Logger, spec Spec, services *simpleanalytics.Client, bk state.Client) *Client {
	return &Client{
		Logger:   logger,
		SAClient: services,
		Backend:  bk,
		Spec:     spec,
	}
}
