package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sql-client/client"
	"sql-client/types"
	"syscall"

	"github.com/namsral/flag"
)

func main() {
	var cfg types.Config
	fs := flag.NewFlagSetWithEnvPrefix("sql-client", "SQL_CLIENT_", flag.ContinueOnError)
	fs.StringVar(&cfg.Addr, "addr", "", "db addr")
	fs.StringVar(&cfg.File, "file", "", "result store file")
	fs.StringVar(&cfg.Path, "path", "", "db path")
	fs.StringVar(&cfg.Type, "driver", "ql", "db path")
	fs.String(flag.DefaultConfigFlagname, "", "config location")

	err := fs.Parse(os.Args[1:])
	if nil != err {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	srv, err := server.New(ctx, &cfg)
	if nil != err {
		fmt.Println(err)
		return
	}

	if err = signalNotify(cancel); nil != err {
		fmt.Println(err)
		return
	}

	if err = srv.Start(); nil != err {
		fmt.Println(err)
	}

	return
}

func signalNotify(cancel context.CancelFunc) error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		signal.Stop(sigChan)
		fmt.Println(fmt.Sprintf("receive stop signal:%v, the programm will be quit.", sig))
		cancel()
	}()
	return nil
}
