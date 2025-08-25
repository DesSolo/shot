package app

import (
	"log"
	"log/slog"
	"os"

	"google.golang.org/grpc"

	"shot/internal/app/screenshot"
	"shot/internal/pkg/config"
	"shot/internal/pkg/provider"
	"shot/internal/pkg/screenshotters/chromedp"
	"shot/pkg/interceptor"
	"shot/pkg/server"
)

type container struct {
	config   *config.Config
	server   *server.Server
	services []server.Service
	provider *provider.Provider
}

func newContainer() *container {
	return &container{}
}

func (c *container) Config() *config.Config {
	if c.config == nil {
		configFilePath := os.Getenv("CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "config.yml"
		}

		cfg, err := config.NewConfigFromFile(configFilePath)
		if err != nil {
			log.Fatalf("failed to load new config: %v", err)
		}

		c.config = cfg
	}

	return c.config
}

func (c *container) Provider() *provider.Provider {
	if c.provider == nil {
		options := c.Config().Screenshotter

		var defaultScreenshotter provider.Screenshotter

		defaultScreenshotter = chromedp.New()

		if options.Default != nil {
			defaultScreenshotter = loadScreenshotter(*options.Default)
		}

		p := provider.New(defaultScreenshotter)

		for host, site := range options.Sites {
			if !site.Enabled {
				continue
			}

			if site.Host == "" {
				site.Host = host
			}

			p.RegisterScreenshotter(host, loadScreenshotter(site))
		}

		c.provider = p
	}

	return c.provider
}

// Server ...
func (c *container) Server() *server.Server {
	if c.server == nil {
		c.server = server.New(server.NewDefaultOptions(
			server.WithGRPCServerOptions(
				grpc.ChainUnaryInterceptor(
					interceptor.Logging(slog.Default()),
					interceptor.Validation,
				),
			),
		))
	}

	return c.server
}

// Services ...
func (c *container) Services() []server.Service {
	if len(c.services) == 0 {
		c.services = append(c.services, screenshot.New(
			c.Provider(),
		))
	}

	return c.services
}
