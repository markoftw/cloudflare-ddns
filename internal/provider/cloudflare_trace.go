package provider

import (
	"regexp"

	"github.com/favonia/cloudflare-ddns/internal/ipnet"
	"github.com/favonia/cloudflare-ddns/internal/provider/protocol"
)

var fieldIP = regexp.MustCompile(`(?m:^ip=(.*)$)`)

var defaultTraceEndpoints = []struct {
	name string
	url  string
}{
	{"cloudflare.trace(api)", "https://api.cloudflare.com/cdn-cgi/trace"},
	{"cloudflare.trace(1.1.1.1)", "https://1.1.1.1/cdn-cgi/trace"},
	{"cloudflare.trace(1.0.0.1)", "https://1.0.0.1/cdn-cgi/trace"},
	{"cloudflare.trace(cloudflare.com)", "https://cloudflare.com/cdn-cgi/trace"},
}

// NewCloudflareTrace creates a CloudflareTrace provider that tries multiple trace endpoints.
func NewCloudflareTrace() Provider {
	providers := make([]Provider, 0, len(defaultTraceEndpoints))
	for _, ep := range defaultTraceEndpoints {
		providers = append(providers, newCloudflareTraceRegexp(ep.name, ep.url))
	}
	return newFailoverProvider("cloudflare.trace", providers...)
}

// NewCloudflareTraceCustom creates a specialized CloudflareTrace provider
// with a specific URL.
func NewCloudflareTraceCustom(url string) Provider {
	return newCloudflareTraceRegexp("cloudflare.trace", url)
}

func newCloudflareTraceRegexp(name, url string) Provider {
	return protocol.Regexp{
		ProviderName: name,
		Param: map[ipnet.Type]protocol.RegexpParam{
			ipnet.IP4: {url, fieldIP},
			ipnet.IP6: {url, fieldIP},
		},
	}
}
