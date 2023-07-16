package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mrnbort/summer_break/api"
	"github.com/mrnbort/summer_break/processor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	Port             string        `short:"p" long:"port" description:"port" default:":8080"`
	HTTPReadTimeout  time.Duration `long:"http-read-timeout" description:"timeout for read HTTP requests" default:"5s"`
	HTTPWriteTimeout time.Duration `long:"http-write-timeout" description:"timeout for write HTTP requests" default:"30s"`
}

func main() {

	var opts options
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v", err)
		}
		os.Exit(1)
	}

	if err := run(opts); err != nil {
		log.Panicf("[ERROR] %v", err)
	}
}

func run(opts options) error {
	transactions := processor.NewProc()

	apiService := api.Service{
		Processor:    transactions,
		Port:         opts.Port,
		ReadTimeOut:  opts.HTTPReadTimeout,
		WriteTimeOut: opts.HTTPWriteTimeout,
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // cancel on SIGINT or SIGTERM
	go func() {
		sig := <-sigs
		log.Printf("received signal: %v", sig)
		cancel()
	}()

	if err := apiService.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("summer break service canceled")
			return nil
		}
		return fmt.Errorf("summer break service failed: %w", err)
	}
	return nil
}
