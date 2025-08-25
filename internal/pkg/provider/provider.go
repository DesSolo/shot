package provider

import (
	"context"
	"fmt"
	"net/url"
	"sync"
)

// Screenshotter ...
type Screenshotter interface {
	Screenshot(ctx context.Context, url string) ([]byte, error)
}

// Provider ...
type Provider struct {
	defaultScreenshotter Screenshotter
	screenshotters       map[string]Screenshotter
	mux                  sync.RWMutex
}

// New ...
func New(screenshotter Screenshotter) *Provider {
	return &Provider{
		defaultScreenshotter: screenshotter,
		screenshotters:       make(map[string]Screenshotter),
	}
}

// RegisterScreenshotter ...
func (p *Provider) RegisterScreenshotter(host string, s Screenshotter) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.screenshotters[host] = s
}

// Image ...
func (p *Provider) Image(ctx context.Context, url string) ([]byte, error) {
	screenshotter, err := p.resolveScreenshotter(url)
	if err != nil {
		return nil, fmt.Errorf("resolveScreenshotter: %w", err)
	}

	screenshot, err := screenshotter.Screenshot(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}

	return screenshot, nil
}

func (p *Provider) resolveScreenshotter(uri string) (Screenshotter, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	p.mux.RLock()
	defer p.mux.RUnlock()

	screenshotter, ok := p.screenshotters[u.Host]
	if !ok {
		return p.defaultScreenshotter, nil
	}

	return screenshotter, nil
}
