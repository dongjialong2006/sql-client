package server

import (
	"bufio"
	"fmt"
	"os"
	"sql-client/sql"
	"strings"

	"github.com/mssola/colors"
)

type Server struct {
	path string
	file string
}

func New(path string, file string) *Server {
	return &Server{
		path: path,
		file: file,
	}
}

func (s *Server) Start(stop chan struct{}) error {
	if "" == s.path {
		return fmt.Errorf("database file path is empty.")
	}

	if err := sql.Create(s.path); nil != err {
		return err
	}

	color := colors.Default()
	var cmd = make(chan string)
	go s.watch(cmd, color)

	for {
		select {
		case <-stop:
			return nil
		case tmp, ok := <-cmd:
			if !ok {
				return nil
			}

			if "stop" == tmp || "quit" == tmp || "close" == tmp {
				return nil
			}

			if err := sql.Exec(s.filter(tmp), s.file, color); nil != err {
				s.errInfo(err, color, false)
				continue
			}

			fmt.Print(">>> ")
		}
	}

	return nil
}

func (s *Server) filter(value string) string {
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\t", "", -1)
	return strings.Replace(value, "  ", " ", -1)
}

func (s *Server) watch(cmd chan string, color *colors.Color) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		tmp, err := reader.ReadString('\n')
		if err != nil {
			s.errInfo(err, color, true)
			continue
		}

		tmp = s.filter(tmp)
		if "" == tmp {
			continue
		}

		if len(tmp) < 5 {
			s.errInfo(fmt.Errorf("invalide input."), color, true)
			continue
		}
		cmd <- tmp
	}

	return
}

func (s *Server) errInfo(err error, color *colors.Color, before bool) {
	if before {
		fmt.Print(">>> ")
	}
	color.Change(colors.Red, false)
	fmt.Println(color.Get(fmt.Sprintf("%v", err)))
	if !before {
		fmt.Print(">>> ")
	}
}

func (s *Server) Stop() {
	sql.Close()
}
