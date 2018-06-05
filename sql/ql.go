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

	ctx := ql.NewRWCtx()
	rs, _, err := s.db.Run(ctx, fmt.Sprintf("BEGIN TRANSACTION; %s; COMMIT;", s.filter(value)))
	if nil != err {
		return err
	}
	if 0 == len(rs) {
		show.Title(int64(ctx.RowsAffected))
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
	value = strings.Trim(value, " ")
	value = strings.Trim(value, ";")
	value = strings.Replace(value, "\n", " ", -1)
	return strings.Replace(value, "\"", "'", -1)
}

func (s *QL) parse(records ql.Recordset, color *colors.Color) error {
	if nil == records {
		return nil
	}

	rs, err := records.Rows(-1, 0)
	if nil != err {
		return err
	}

	if 0 == len(rs) {
		return nil
	}

	if "" != s.cfg.File {
		if err = file.CreatePath(s.cfg.File); nil != err {
			return err
		}
	}

	show.Title(int64(len(rs)))

	fields, err := records.Fields()
	if nil != err {
		return err
	}
	if 0 == len(fields) {
		return nil
	}

	if err = s.show(fields, rs, color); nil != err {
		return err
	}

	return nil
}

func (s *QL) show(fields []string, rows [][]interface{}, color *colors.Color) error {
	show.Header(fields)
	for i, row := range rows {
		if "" == s.cfg.File {
			show.Body(i, row, nil)
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
