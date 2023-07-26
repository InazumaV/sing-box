//go:build !with_quic

package outbound

import (
	"context"

	"github.com/inazumav/sing-box/adapter"
	C "github.com/inazumav/sing-box/constant"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
)

func NewHysteria(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaOutboundOptions) (adapter.Outbound, error) {
	return nil, C.ErrQUICNotIncluded
}
