package sql

import (
	"context"
	"fmt"
	"sql-client/pkg/show"
	"sql-client/types"
	"sync"

	"github.com/cznic/ql"
	"github.com/mssola/colors"
)

type QLClient struct {
	sync.RWMutex
	db  *ql.DB
	cfg *types.Options
	ctx context.Context
}

func NewQLClient(ctx context.Context, cfg *types.Options) *QLClient {
	srv := &QLClient{
		cfg: cfg,
		ctx: ctx,
	}
	go srv.stop(true)

	return srv
}

func (s *QLClient) Exec(value string, color *colors.Color) error {
	if "" == value || "" == s.cfg.Addr {
		return nil
	}

	var err error = nil
	s.Lock()
	s.db, err = ql.OpenFile(s.cfg.Addr, &ql.Options{})
	s.Unlock()
	if nil != err {
		return err
	}
	defer s.stop(false)

	ctx := ql.NewRWCtx()
	rs, _, err := s.db.Run(ctx, fmt.Sprintf("BEGIN TRANSACTION; %s; COMMIT;", value))
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

func (s *QLClient) Stop() {
	s.stop(false)
	return
}

func (s *QLClient) stop(check bool) {
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

func (s *QLClient) parse(records ql.Recordset, color *colors.Color) error {
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

func (s *QLClient) show(fields []string, rows [][]interface{}, color *colors.Color) error {
	show.Header(fields)
	for i, row := range rows {
		show.Body(i, row, nil)
	}

	return nil
}
