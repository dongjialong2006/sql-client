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
	var cfg *types.Options = types.NewOptions()
	fs := flag.NewFlagSetWithEnvPrefix("sql-client", "SQL_CLIENT_", flag.ContinueOnError)
	fs.StringVar(&cfg.Addr, "addr", "", "db addr or file path")
	fs.StringVar(&cfg.Type, "type", "ql", "db type")

	fs.StringVar(&cfg.Pwd, "pwd", "", "redis db pwd")
	fs.IntVar(&cfg.RedisOpt.DB, "db", 0, "redis db")
	fs.String(flag.DefaultConfigFlagname, "", "config location")

	err := fs.Parse(os.Args[1:])
	if nil != err {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	srv, err := server.New(ctx, cfg)
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
