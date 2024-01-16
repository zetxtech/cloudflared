package edgediscovery

import (
	"context"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/socks5/client"
)

func Dial(ctx context.Context, proxyUrl *url.URL, addr net.Addr) (net.Conn, error) {
	dialer, err := Socks5Dialer(proxyUrl)
	if err != nil {
		return nil, err
	}

	ta, err := netapi.ParseSysAddr(addr)
	if err != nil {
		return nil, err
	}

	return dialer.Conn(ctx, ta)
}

func Socks5Dialer(proxy *url.URL) (netapi.Proxy, error) {
	password, _ := proxy.User.Password()
	return client.Dial(proxy.Hostname(), proxy.Port(),
		proxy.User.Username(), password), nil
}

// FromEnvironmentUsing returns the dialer specify by the proxy-related
// variables in the environment and makes underlying connections
// using the provided forwarding Dialer (for instance, a *net.Dialer
// with desired configuration).
func FromEnvironmentUsing() (*url.URL, bool) {
	allProxy := allProxyEnv.Get()
	if len(allProxy) == 0 {
		return nil, false
	}

	proxyURL, err := url.Parse(allProxy)
	if err != nil {
		return nil, false
	}

	if strings.ToLower(proxyURL.Scheme) != "socks5" || proxyURL.Port() == "" {
		return nil, false
	}

	return proxyURL, true
}

var (
	allProxyEnv = &envOnce{
		names: []string{"ALL_PROXY", "all_proxy"},
	}
	noProxyEnv = &envOnce{
		names: []string{"NO_PROXY", "no_proxy"},
	}
)

// envOnce looks up an environment variable (optionally by multiple
// names) once. It mitigates expensive lookups on some platforms
// (e.g. Windows).
// (Borrowed from net/http/transport.go)
type envOnce struct {
	names []string
	once  sync.Once
	val   string
}

func (e *envOnce) Get() string {
	e.once.Do(e.init)
	return e.val
}

func (e *envOnce) init() {
	for _, n := range e.names {
		e.val = os.Getenv(n)
		if e.val != "" {
			return
		}
	}
}

// reset is used by tests
func (e *envOnce) reset() {
	e.once = sync.Once{}
	e.val = ""
}
