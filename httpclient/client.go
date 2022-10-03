package httpclient

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var client *http.Client
var once sync.Once

func GetHTTPClient(timeout time.Duration) *http.Client {
	if client == nil {
		once.Do(func() {
			netTransport := &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   time.Second,
					KeepAlive: 0,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
				IdleConnTimeout:     0,
				MaxIdleConnsPerHost: 50000,
				MaxIdleConns:        50000,
			}

			client = &http.Client{
				Timeout:   timeout,
				Transport: netTransport,
			}
		})
	}

	return client

}
