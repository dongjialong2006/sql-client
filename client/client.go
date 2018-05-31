package server

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sql-client/sql"
	"sql-client/types"
	"strings"

	"github.com/mssola/colors"
)

type Server struct {
	ctx    context.Context
	handle types.Handle
}

func New(ctx context.Context, cfg *types.Config) (*Server, error) {
	srv := &Server{
		ctx: ctx,
	}
	switch cfg.Type {
	case types.QL:
		srv.handle = sql.NewQL(ctx, cfg)
	default:
		return nil, fmt.Errorf("unknown db driver type:%s.", cfg.Type)
	}

	return srv, nil
}

func (s *Server) Start() error {
	color := colors.Default()
	var cmd = make(chan string)
	go s.watch(cmd, color)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case tmp, ok := <-cmd:
			if !ok {
				return nil
			}

			if "stop" == tmp || "quit" == tmp || "close" == tmp {
				return nil
			}

			if err := s.handle.Exec(s.filter(tmp), color); nil != err {
				s.errInfo(err, color, false)
				continue
			}

			fmt.Print(">>> ")
		}
	}

	return nil
}

func (s *Server) Stop() {
	s.handle.Stop()
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
