package chromedp

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

const (
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	defaultWidth     = 1920
	defaultHeight    = 1080
)

// Screenshotter ...
type Screenshotter struct {
	userAgent     string
	width, height int64

	actions []chromedp.Action
}

// New ...
func New(options ...OptionFunc) *Screenshotter {
	s := &Screenshotter{
		userAgent: defaultUserAgent,
		width:     defaultWidth,
		height:    defaultHeight,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Screenshot ...
func (s *Screenshotter) Screenshot(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	actions := s.allActions()

	var buf []byte

	actions = append(actions,
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	)

	err := chromedp.Run(ctx, actions...)

	if err != nil {
		return nil, fmt.Errorf("chromedp.Run: %w", err)
	}

	return buf, nil
}

func (s *Screenshotter) allActions() []chromedp.Action {
	actions := make([]chromedp.Action, 0, len(s.actions)+2)

	actions = append(actions,
		emulation.SetUserAgentOverride(s.userAgent),
		emulation.SetDeviceMetricsOverride(s.width, s.height, 1.0, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}),
	)

	actions = append(actions, s.actions...)

	return actions
}
