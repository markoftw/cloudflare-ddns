//nolint:testpackage // accessing unexported helper ensures thorough coverage
package provider

import (
	"context"
	"io"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/favonia/cloudflare-ddns/internal/ipnet"
	"github.com/favonia/cloudflare-ddns/internal/pp"
)

type stubProvider struct {
	name string
	ip   netip.Addr
	ok   bool
}

func (s stubProvider) Name() string {
	return s.name
}

func (s stubProvider) GetIP(context.Context, pp.PP, ipnet.Type) (netip.Addr, bool) {
	if s.ok {
		return s.ip, true
	}
	return netip.Addr{}, false
}

func TestFailoverProviderPrefersFirstSuccess(t *testing.T) {
	t.Parallel()

	ppfmt := pp.New(io.Discard, false, pp.Quiet)
	ctx := context.Background()

	expected := netip.MustParseAddr("1.2.3.4")
	failover := newFailoverProvider("test",
		stubProvider{name: "primary", ip: expected, ok: true},
		stubProvider{name: "secondary", ip: netip.MustParseAddr("5.6.7.8"), ok: true},
	)

	ip, ok := failover.GetIP(ctx, ppfmt, ipnet.IP4)
	require.True(t, ok)
	require.Equal(t, expected, ip)
}

func TestFailoverProviderFallsBack(t *testing.T) {
	t.Parallel()

	ppfmt := pp.New(io.Discard, false, pp.Quiet)
	ctx := context.Background()

	expected := netip.MustParseAddr("2.2.2.2")
	failover := newFailoverProvider("test",
		stubProvider{name: "primary", ip: netip.Addr{}, ok: false},
		stubProvider{name: "secondary", ip: expected, ok: true},
	)

	ip, ok := failover.GetIP(ctx, ppfmt, ipnet.IP6)
	require.True(t, ok)
	require.Equal(t, expected, ip)
}

func TestFailoverProviderAllFail(t *testing.T) {
	t.Parallel()

	ppfmt := pp.New(io.Discard, false, pp.Quiet)
	ctx := context.Background()

	failover := newFailoverProvider("test",
		stubProvider{name: "primary", ip: netip.Addr{}, ok: false},
		stubProvider{name: "secondary", ip: netip.Addr{}, ok: false},
	)

	_, ok := failover.GetIP(ctx, ppfmt, ipnet.IP4)
	require.False(t, ok)
}
