package app

import (
	"context"
	"fmt"
)

// App ...
type App struct{}

// New ...
func New() *App {
	return &App{}
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	di := newContainer()

	configureLogger(di)

	if err := di.Server().Run(ctx, di.Services()...); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}
