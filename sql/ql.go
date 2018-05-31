package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"sql-client/pkg/file"
	"sql-client/pkg/show"
	"sql-client/types"
	"strings"
	"sync"

	"github.com/cznic/ql"
	"github.com/mssola/colors"
)

type QL struct {
	sync.RWMutex
	db  *ql.DB
	cfg *types.Config
	ctx context.Context
}

func NewQL(ctx context.Context, cfg *types.Config) *QL {
	srv := &QL{
		cfg: cfg,
		ctx: ctx,
	}
	go srv.stop(true)

	return srv
}

func (s *QL) Exec(value string, color *colors.Color) error {
	if "" == value || "" == s.cfg.Path {
		return nil
	}

	var err error = nil
	s.Lock()
	s.db, err = ql.OpenFile(s.cfg.Path, &ql.Options{})
	s.Unlock()
	if nil != err {
		return err
	}
	defer s.stop(false)

	rs, _, err := s.db.Run(ql.NewRWCtx(), fmt.Sprintf("BEGIN TRANSACTION; %s; COMMIT;", s.filter(value)))
	if nil != err {
		return err
	}
	if 0 == len(rs) {
		return nil
	}

	for _, tmp := range rs {
		if nil == tmp {
			continue
		}
		if err = s.parse(tmp, color); nil != err {
			return err
		}
	}

	return nil
}

func (s *QL) Stop() {
	s.stop(false)
	return
}

func (s *QL) stop(check bool) {
	if check {
		<-s.ctx.Done()
	}
	s.Lock()
	if nil != s.db {
		s.db.Close()
		s.db = nil
	}
	s.Unlock()
}

func (s *QL) filter(value string) string {
	value = strings.Replace(value, "\n", " ", -1)
	return strings.Replace(value, "\"", "'", -1)
}

func (s *QL) parse(tmp ql.Recordset, color *colors.Color) error {
	records, err := tmp.Rows(-1, 0)
	if nil != err {
		return err
	}

	if "" != s.cfg.File {
		if err = file.CreatePath(s.cfg.File); nil != err {
			return err
		}
	}

	color.Change(colors.Red, false)
	fmt.Println(color.Get(fmt.Sprintf("current query result num: %d.", len(records))))

	fields, err := tmp.Fields()
	if nil != err {
		return err
	}

	if err = s.show(fields, records, color); nil != err {
		return err
	}

	return nil
}

func (s *QL) show(fields []string, rows [][]interface{}, color *colors.Color) error {
	show.Header(fields, color)
	for i, row := range rows {
		if "" == s.cfg.File {
			show.Body(i, row, color)
			continue
		}
		data, err := json.Marshal(row)
		if nil != err {
			return err
		}
		if err = file.WriteFile(s.cfg.File, data); nil != err {
			return err
		}
	}

	return nil
}
