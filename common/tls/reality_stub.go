//go:build !with_reality_server

package tls

import (
	"context"

	"github.com/inazumav/sing-box/adapter"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewRealityServer(ctx context.Context, router adapter.Router, logger log.Logger, options option.InboundTLSOptions) (ServerConfig, error) {
	return nil, E.New(`reality server is not included in this build, rebuild with -tags with_reality_server`)
}
