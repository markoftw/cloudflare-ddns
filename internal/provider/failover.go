package provider

import (
	"context"
	"net/netip"

	"github.com/favonia/cloudflare-ddns/internal/ipnet"
	"github.com/favonia/cloudflare-ddns/internal/pp"
)

type failoverProvider struct {
	providerName string
	providers    []Provider
}

func newFailoverProvider(providerName string, providers ...Provider) Provider {
	filtered := make([]Provider, 0, len(providers))
	for _, p := range providers {
		if p != nil {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return nil
	}
	if len(filtered) == 1 {
		return filtered[0]
	}

	return &failoverProvider{
		providerName: providerName,
		providers:    filtered,
	}
}

func (f *failoverProvider) Name() string {
	return f.providerName
}

func (f *failoverProvider) GetIP(ctx context.Context, ppfmt pp.PP, ipNet ipnet.Type) (netip.Addr, bool) {
	if len(f.providers) == 0 {
		ppfmt.Noticef(pp.EmojiImpossible, "Failover provider %s is misconfigured with zero backends", f.providerName)
		return netip.Addr{}, false
	}

	for idx, p := range f.providers {
		ip, ok := p.GetIP(ctx, ppfmt, ipNet)
		if ok {
			if idx > 0 {
				ppfmt.Noticef(pp.EmojiSwitch, "Switched to fallback provider %s for %s detection", Name(p), ipNet.Describe())
			}
			return ip, true
		}

		if idx < len(f.providers)-1 {
			next := f.providers[idx+1]
			ppfmt.Noticef(pp.EmojiWarning, "Provider %s failed to detect %s; trying %s next", Name(p), ipNet.Describe(), Name(next))
		}
	}

	return netip.Addr{}, false
}
