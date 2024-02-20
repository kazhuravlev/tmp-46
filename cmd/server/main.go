package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/kazhuravlev/sample-server/internal/api"
	"github.com/kazhuravlev/sample-server/internal/facade"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if err := cmdRun(); err != nil {
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("The sky is falling: %v\n", err)
		fmt.Println(strings.Repeat("=", 80))
	}
}

func cmdRun() error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	var flagPort int
	flag.IntVar(&flagPort, "port", 8888, "listen port")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse command flags: %w", err)
	}

	apiInst, err := api.New(logger)
	if err != nil {
		return fmt.Errorf("init api instance: %w", err)
	}

	if err := apiInst.Run(ctx); err != nil {
		return fmt.Errorf("run api instance: %w", err)
	}

	facadeInst, err := facade.New(logger, apiInst, flagPort)
	if err != nil {
		return fmt.Errorf("init facade instance: %w", err)
	}

	if err := facadeInst.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			return fmt.Errorf("run facade instance: %w", err)
		}
	}

	facadeInst.Wait()
	apiInst.Wait()

	logger.Warn("server is going to shutdown")
	logger.Info("wait all connections to stop")
	logger.Info("all connections was stopped")

	return nil
}
