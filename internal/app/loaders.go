package app

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"shot/internal/pkg/config"
	"shot/internal/pkg/provider"
	"shot/internal/pkg/screenshotters/chromedp"
)

func configureLogger(di *container) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.Level(di.Config().Logger.Level),
			}),
		),
	)
}

func loadScreenshotter(screenshotter config.Screenshotter) provider.Screenshotter {
	var options []chromedp.OptionFunc

	if screenshotter.UserAgent != "" {
		options = append(options, chromedp.WithUserAgent(screenshotter.UserAgent))
	}

	if screenshotter.Resolution != nil {
		options = append(options, chromedp.WithResolution(screenshotter.Resolution.Width, screenshotter.Resolution.Height))
	}

	if len(screenshotter.Cookies) > 0 {
		cookies := make([]chromedp.Cookie, 0, len(screenshotter.Cookies))
		for _, cookie := range screenshotter.Cookies {
			parts := strings.Split(cookie, "=")
			if len(parts) != 2 {
				continue
			}

			if screenshotter.Host == "" {
				log.Fatalf("host must be specified for use cookies %s", cookie)
			}

			cookies = append(cookies, chromedp.Cookie{
				Domain: screenshotter.Host,
				Key:    parts[0],
				Value:  parts[1],
			})
		}
		options = append(options, chromedp.WithCookies(cookies...))
	}

	if len(screenshotter.Headers) > 0 {
		headers := make([]chromedp.Header, 0, len(screenshotter.Headers))
		for key, value := range screenshotter.Headers {
			headers = append(headers, chromedp.Header{
				Key:   key,
				Value: value,
			})
		}
		options = append(options, chromedp.WithHeaders(headers...))
	}

	return chromedp.New(options...)
}
