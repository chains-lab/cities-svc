package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	cmd2 "github.com/chains-lab/cities-svc/cmd"
	"github.com/chains-lab/cities-svc/cmd/migrations"
	"github.com/chains-lab/cities-svc/internal"

	"github.com/chains-lab/logium"
)

func Run(args []string) bool {
	cfg := internal.LoadConfig()

	log := logium.NewLogger(cfg.Log.Level, cfg.Log.Format)
	log.Info("Starting server...")

	var (
		service = kingpin.New("chains-auth", "")
		runCmd  = service.Command("run", "run command")

		serviceCmd     = runCmd.Command("service", "run service")
		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	cmd, err := service.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		cmd2.StartServices(ctx, cfg, log, &wg)
	case migrateUpCmd.FullCommand():
		err = migrations.MigrateUp(cfg.Database.SQL.URL)
	case migrateDownCmd.FullCommand():
		err = migrations.MigrateDown(cfg.Database.SQL.URL)
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All api stopped")
	}

	return true
}
