package tls

import (
	"context"
	"net"

	"github.com/inazumav/sing-box/adapter"
	C "github.com/inazumav/sing-box/constant"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
	aTLS "github.com/sagernet/sing/common/tls"
)

func NewServer(ctx context.Context, router adapter.Router, logger log.Logger, options option.InboundTLSOptions) (ServerConfig, error) {
	if !options.Enabled {
		return nil, nil
	}
	if options.Reality != nil && options.Reality.Enabled {
		return NewRealityServer(ctx, router, logger, options)
	} else {
		return NewSTDServer(ctx, router, logger, options)
	}
}

func ServerHandshake(ctx context.Context, conn net.Conn, config ServerConfig) (Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, C.TCPTimeout)
	defer cancel()
	return aTLS.ServerHandshake(ctx, conn, config)
}
