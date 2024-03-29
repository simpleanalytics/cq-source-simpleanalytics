package client

import "github.com/cloudquery/plugin-sdk/v4/schema"

func WebsiteMultiplex(meta schema.ClientMeta) []schema.ClientMeta {
	var l = make([]schema.ClientMeta, 0)
	client := meta.(*Client)
	for _, website := range client.Spec.Websites {
		l = append(l, client.withWebsite(website))
	}
	return l
}
