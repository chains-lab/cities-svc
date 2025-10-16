package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/chains-lab/cities-svc/internal/data"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/domain/services/country"
	"github.com/chains-lab/cities-svc/internal/domain/services/invite"
	"github.com/chains-lab/cities-svc/internal/infra/jwtmanager"
	"github.com/chains-lab/cities-svc/internal/infra/usrguesser"
	"github.com/chains-lab/cities-svc/internal/rest"
	"github.com/chains-lab/cities-svc/internal/rest/controller"
	"github.com/chains-lab/cities-svc/internal/rest/middlewares"

	"github.com/chains-lab/logium"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.NewDatabase(pg)

	jwtInviteManager := jwtmanager.NewManager(cfg)
	userGuesser := usrguesser.NewService(cfg.Profile.Url, nil)

	citySvc := city.NewService(database)
	countrySvc := country.NewService(database)
	cityModerSvc := admin.NewService(database, userGuesser)
	inviteSvc := invite.NewService(database, jwtInviteManager)

	ctrl := controller.New(log, countrySvc, citySvc, cityModerSvc, inviteSvc)
	mdlv := middlewares.New(log, cityModerSvc)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
