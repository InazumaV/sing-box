package inbound

import (
	"context"

	"github.com/inazumav/sing-box/adapter"
	C "github.com/inazumav/sing-box/constant"
	"github.com/inazumav/sing-box/experimental/libbox/platform"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func New(ctx context.Context, router adapter.Router, logger log.ContextLogger, options option.Inbound, platformInterface platform.Interface) (adapter.Inbound, error) {
	if options.Type == "" {
		return nil, E.New("missing inbound type")
	}
	switch options.Type {
	case C.TypeShadowsocks:
		return NewShadowsocks(ctx, router, logger, options.Tag, options.ShadowsocksOptions)
	case C.TypeVMess:
		return NewVMess(ctx, router, logger, options.Tag, options.VMessOptions)
	case C.TypeTrojan:
		return NewTrojan(ctx, router, logger, options.Tag, options.TrojanOptions)
	case C.TypeHysteria:
		return NewHysteria(ctx, router, logger, options.Tag, options.HysteriaOptions)
	case C.TypeVLESS:
		return NewVLESS(ctx, router, logger, options.Tag, options.VLESSOptions)
	case C.TypeTUIC:
		return NewTUIC(ctx, router, logger, options.Tag, options.TUICOptions)
	default:
		return nil, E.New("unknown inbound type: ", options.Type)
	}
}
