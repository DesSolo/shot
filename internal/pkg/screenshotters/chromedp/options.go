package chromedp

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// OptionFunc ...
type OptionFunc func(*Screenshotter)

// WithUserAgent ...
func WithUserAgent(userAgent string) OptionFunc {
	return func(p *Screenshotter) {
		p.userAgent = userAgent
	}
}

// WithResolution ...
func WithResolution(width, height int64) OptionFunc {
	return func(p *Screenshotter) {
		p.width = width
		p.height = height
	}
}

// WithCookies ...
func WithCookies(cookies ...Cookie) OptionFunc {
	return func(p *Screenshotter) {
		p.actions = append(p.actions, chromedp.ActionFunc(func(ctx context.Context) error {
			for _, cookie := range cookies {
				if err := network.SetCookie(cookie.Key, cookie.Value).WithDomain(cookie.Domain).Do(ctx); err != nil {
					return fmt.Errorf("network.SetCookie %q: %w", cookie.Key, err)
				}
			}

			return nil
		}))
	}
}

// WithHeaders ...
func WithHeaders(headers ...Header) OptionFunc {
	netHeaders := make(network.Headers, len(headers))
	for _, header := range headers {
		netHeaders[header.Key] = header.Value
	}

	return func(p *Screenshotter) {
		p.actions = append(p.actions, network.SetExtraHTTPHeaders(netHeaders))
	}
}
