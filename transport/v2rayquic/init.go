package v2rayquic

import "github.com/inazumav/sing-box/transport/v2ray"

func init() {
	v2ray.RegisterQUICConstructor(NewServer, NewClient)
}
