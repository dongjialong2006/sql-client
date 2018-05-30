package main

import (
	"fmt"
	"os"
	"os/signal"
	"sql-client/client"
	"syscall"

	"github.com/namsral/flag"
)

func main() {
	var file string
	var path string

	fs := flag.NewFlagSetWithEnvPrefix("sql-client", "SQL_CLIENT_", flag.ContinueOnError)
	fs.StringVar(&file, "file", "", "result store file")
	fs.StringVar(&path, "path", "", "db path")
	fs.String(flag.DefaultConfigFlagname, "", "config location")

	err := fs.Parse(os.Args[1:])
	if nil != err {
		fmt.Println(err)
		return
	}

	var stop = make(chan struct{})
	if err = signalNotify(stop); nil != err {
		fmt.Println(err)
		return
	}

	srv := server.New(path, file)
	defer srv.Stop()
	if err = srv.Start(stop); nil != err {
		fmt.Println(err)
	}

	return
}

func signalNotify(stop chan struct{}) error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		signal.Stop(sigChan)
		fmt.Println(fmt.Sprintf("receive stop signal:%v, the programm will be quit.", sig))
		close(stop)
	}()
	return nil
}
