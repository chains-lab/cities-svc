package cli

import (
	"context"
	"sync"

	"github.com/chains-lab/cities-svc/internal/api/rest"
	"github.com/chains-lab/cities-svc/internal/api/rest/handlers"
	"github.com/chains-lab/cities-svc/internal/app"
	"github.com/chains-lab/cities-svc/internal/config"
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

	restSVC := rest.NewRest(cfg, log)

	run(func() {
		handl := handlers.NewAdapter(cfg, log, app)

		restSVC.Router(ctx, cfg, handl)
	})
}
