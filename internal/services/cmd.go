package services

import (
	"context"
	"sync"

	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/config"
	"github.com/chains-lab/cities-svc/internal/services/rest"
	"github.com/chains-lab/cities-svc/internal/services/rest/handlers"
	"github.com/chains-lab/cities-svc/internal/services/rest/middlewares"
	"github.com/chains-lab/logium"
)

func StartServices(ctx context.Context, cfg config.Config, log logium.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	restApi := rest.NewRest(cfg, log)

	run(func() {
		mid := middlewares.NewAdapter(cfg, log, app)
		handl := handlers.NewAdapter(cfg, log, app)

		restApi.Api(ctx, mid, handl)
	})
}
