package main

import (
	"os"
	"testing"

	"github.com/jetstack/cert-manager/test/acme/dns"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRun(t *testing.T) {
	s := NewAlibabaDNSSolver()
	fixture := dns.NewFixture(s,
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetResolvedFQDN("x."),
		dns.SetUseAuthoritative(true),
		dns.SetManifestPath("testdata/alidns"),
		dns.SetBinariesPath("./_out/kubebuilder/bin"),
	)
	fixture.RunConformance(t)
}
