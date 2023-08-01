//go:build !with_quic

package inbound

import (
	"context"

	"github.com/inazumav/sing-box/adapter"
	C "github.com/inazumav/sing-box/constant"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
)

func NewTUIC(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.TUICInboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}
