package inbound

import (
	"context"
	"net"
	"os"

	"github.com/inazumav/sing-box/adapter"
	C "github.com/inazumav/sing-box/constant"
	"github.com/inazumav/sing-box/log"
	"github.com/inazumav/sing-box/option"
	"github.com/sagernet/sing-shadowsocks"
	"github.com/sagernet/sing-shadowsocks/shadowaead"
	"github.com/sagernet/sing-shadowsocks/shadowaead_2022"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/buf"
	E "github.com/sagernet/sing/common/exceptions"
	F "github.com/sagernet/sing/common/format"
	N "github.com/sagernet/sing/common/network"
)

var (
	_ adapter.Inbound           = (*ShadowsocksMulti)(nil)
	_ adapter.InjectableInbound = (*ShadowsocksMulti)(nil)
)

type ShadowsocksMulti struct {
	myInboundAdapter
	service shadowsocks.MultiService[int]
	users   []option.ShadowsocksUser
}

func newShadowsocksMulti(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowsocksInboundOptions) (*ShadowsocksMulti, error) {
	inbound := &ShadowsocksMulti{
		myInboundAdapter: myInboundAdapter{
			protocol:      C.TypeShadowsocks,
			network:       options.Network.Build(),
			ctx:           ctx,
			router:        router,
			logger:        logger,
			tag:           tag,
			listenOptions: options.ListenOptions,
		},
	}
	inbound.connHandler = inbound
	inbound.packetHandler = inbound
	var udpTimeout int64
	if options.UDPTimeout != 0 {
		udpTimeout = options.UDPTimeout
	} else {
		udpTimeout = int64(C.UDPTimeout.Seconds())
	}
	var (
		service shadowsocks.MultiService[int]
		err     error
	)
	if common.Contains(shadowaead_2022.List, options.Method) {
		service, err = shadowaead_2022.NewMultiServiceWithPassword[int](
			options.Method,
			options.Password,
			udpTimeout,
			adapter.NewUpstreamContextHandler(inbound.newConnection, inbound.newPacketConnection, inbound),
			router.TimeFunc(),
		)
	} else if common.Contains(shadowaead.List, options.Method) {
		service, err = shadowaead.NewMultiService[int](
			options.Method,
			udpTimeout,
			adapter.NewUpstreamContextHandler(inbound.newConnection, inbound.newPacketConnection, inbound))
	} else {
		return nil, E.New("unsupported method: " + options.Method)
	}
	if err != nil {
		return nil, err
	}
	err = service.UpdateUsersWithPasswords(common.MapIndexed(options.Users, func(index int, user option.ShadowsocksUser) int {
		return index
	}), common.Map(options.Users, func(user option.ShadowsocksUser) string {
		return user.Password
	}))
	if err != nil {
		return nil, err
	}
	inbound.service = service
	inbound.packetUpstream = service
	inbound.users = options.Users
	return inbound, err
}

func (h *ShadowsocksMulti) AddUsers(users []option.ShadowsocksUser) error {
	temp := make([]option.ShadowsocksUser, 0, len(h.users)+len(users))
	temp = append(temp, h.users...)
	temp = append(temp, users...)
	err := h.service.UpdateUsersWithPasswords(common.MapIndexed(temp, func(index int, user option.ShadowsocksUser) int {
		return index
	}), common.Map(temp, func(user option.ShadowsocksUser) string {
		return user.Password
	}))
	if err != nil {
		return err
	}
	h.users = temp
	return nil
}

func (h *ShadowsocksMulti) DelUsers(name []string) error {
	is := make([]int, 0, len(name))
	ulen := len(name)
	for i := range h.users {
		for _, n := range name {
			if h.users[i].Name == n {
				is = append(is, i)
				ulen--
			}
			if ulen == 0 {
				break
			}
		}
	}
	ulen = len(h.users)
	for _, i := range is {
		h.users[i] = h.users[ulen-1]
		h.users[ulen-1] = option.ShadowsocksUser{}
		h.users = h.users[:ulen-1]
		ulen--
	}
	err := h.service.UpdateUsersWithPasswords(common.MapIndexed(h.users, func(index int, user option.ShadowsocksUser) int {
		return index
	}), common.Map(h.users, func(user option.ShadowsocksUser) string {
		return user.Password
	}))
	return err
}

func (h *ShadowsocksMulti) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	return h.service.NewConnection(adapter.WithContext(log.ContextWithNewID(ctx), &metadata), conn, adapter.UpstreamMetadata(metadata))
}

func (h *ShadowsocksMulti) NewPacket(ctx context.Context, conn N.PacketConn, buffer *buf.Buffer, metadata adapter.InboundContext) error {
	return h.service.NewPacket(adapter.WithContext(ctx, &metadata), conn, buffer, adapter.UpstreamMetadata(metadata))
}

func (h *ShadowsocksMulti) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	return os.ErrInvalid
}

func (h *ShadowsocksMulti) newConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	userIndex, loaded := auth.UserFromContext[int](ctx)
	if !loaded {
		return os.ErrInvalid
	}
	user := h.users[userIndex].Name
	if user == "" {
		user = F.ToString(userIndex)
	} else {
		metadata.User = user
	}
	h.logger.InfoContext(ctx, "[", user, "] inbound connection to ", metadata.Destination)
	return h.router.RouteConnection(ctx, conn, metadata)
}

func (h *ShadowsocksMulti) newPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	userIndex, loaded := auth.UserFromContext[int](ctx)
	if !loaded {
		return os.ErrInvalid
	}
	user := h.users[userIndex].Name
	if user == "" {
		user = F.ToString(userIndex)
	} else {
		metadata.User = user
	}
	ctx = log.ContextWithNewID(ctx)
	h.logger.InfoContext(ctx, "[", user, "] inbound packet connection from ", metadata.Source)
	h.logger.InfoContext(ctx, "[", user, "] inbound packet connection to ", metadata.Destination)
	return h.router.RoutePacketConnection(ctx, conn, metadata)
}
