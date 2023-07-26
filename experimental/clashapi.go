package experimental

import (
	"context"
	"os"

	"github.com/inazumav/sing-box/adapter"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
)

type ClashServerConstructor = func(ctx context.Context, router adapter.Router, logFactory log.ObservableFactory, options option.ClashAPIOptions) (adapter.ClashServer, error)

var clashServerConstructor ClashServerConstructor

func RegisterClashServerConstructor(constructor ClashServerConstructor) {
	clashServerConstructor = constructor
}

func NewClashServer(ctx context.Context, router adapter.Router, logFactory log.ObservableFactory, options option.ClashAPIOptions) (adapter.ClashServer, error) {
	if clashServerConstructor == nil {
		return nil, os.ErrInvalid
	}
	return clashServerConstructor(ctx, router, logFactory, options)
}
