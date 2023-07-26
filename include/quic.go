//go:build with_quic

package include

import (
	_ "github.com/inazumav/sing-box/transport/v2rayquic"
	_ "github.com/sagernet/sing-dns/quic"
)

const WithQUIC = true
