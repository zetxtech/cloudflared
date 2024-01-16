package connection

import (
	"context"
	"net"
	"net/url"

	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/cloudflare/cloudflared/edgediscovery"
)

func createSocks5UDPConnForConnIndex(proxy *url.URL, addr net.Addr) (net.PacketConn, error) {
	portMapMutex.Lock()
	defer portMapMutex.Unlock()

	dialer, err := edgediscovery.Socks5Dialer(proxy)
	if err != nil {
		return nil, err
	}

	ta, err := netapi.ParseSysAddr(addr)
	if err != nil {
		return nil, err
	}

	return dialer.PacketConn(context.TODO(), ta)
}
