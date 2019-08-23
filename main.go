package main // import "github.com/tttlkkkl/alidns"

import (
	"os"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
)

var GroupName string

func main() {
	GroupName = os.Getenv("GROUP_NAME")

	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}
	s := NewAlibabaDNSSolver()
	cmd.RunWebhookServer(GroupName, s)
}
