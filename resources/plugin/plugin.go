package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	internalPlugin "github.com/simpleanalytics/cq-source-simple-analytics/plugin"
)

func Plugin() *plugin.Plugin {
	return plugin.NewPlugin(
		internalPlugin.Name,
		internalPlugin.Version,
		Configure,
		plugin.WithKind(internalPlugin.Kind),
		plugin.WithTeam(internalPlugin.Team),
	)
}
