package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/chains-lab/cities-svc/internal"
	"github.com/chains-lab/cities-svc/internal/domain/services/admin"
	"github.com/chains-lab/cities-svc/internal/domain/services/city"
	"github.com/chains-lab/cities-svc/internal/events/publisher"
	"github.com/chains-lab/cities-svc/internal/repo"

	"github.com/chains-lab/cities-svc/internal/domain/services/invite"
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

	database := repo.NewDatabase(pg)

	eventPublish := publisher.New(cfg.Kafka.Broker)

	citySvc := city.NewService(database, eventPublish)
	cityAdminSvc := admin.NewService(database, eventPublish)
	inviteSvc := invite.NewService(database, eventPublish)

	ctrl := controller.New(log, citySvc, cityAdminSvc, inviteSvc)
	mdlv := middlewares.New(log)

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
