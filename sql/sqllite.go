package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	// "sql-client/pkg/file"
	"sql-client/pkg/show"
	"sql-client/types"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mssola/colors"
)

type SqlLite struct {
	sync.RWMutex
	db  *sql.DB
	cfg *types.Config
	ctx context.Context
}

func NewSqlLite(ctx context.Context, cfg *types.Config) *SqlLite {
	srv := &SqlLite{
		cfg: cfg,
		ctx: ctx,
	}
	go srv.stop(true)

	return srv
}

func (s *SqlLite) Exec(value string, color *colors.Color) error {
	if "" == value {
		return nil
	}
	if "" == s.cfg.Path {
		return fmt.Errorf("sqllite db path is empty.")
	}

	var err error = nil
	s.Lock()
	s.db, err = sql.Open("sqlite3", s.cfg.Path)
	s.Unlock()
	if nil != err {
		return err
	}
	defer s.stop(false)

	value = strings.Trim(value, " ")
	stmt, err := s.db.Prepare(value)
	if nil != err {
		return err
	}

	if strings.HasPrefix(value, "select") || strings.HasPrefix(value, "SELECT") {
		rs, err := stmt.Query()
		if nil != err {
			return err
		}
		return s.parse(rs, color)
	}
	rs, err := stmt.Exec()
	if nil != err {
		return err
	}
	num, err := rs.RowsAffected()
	if nil != err {
		return err
	}
	show.Title(num)
	return nil
}

func (s *SqlLite) Stop() {
	s.stop((false))
}

func (s *SqlLite) parse(rs *sql.Rows, color *colors.Color) error {
	if nil == rs {
		return nil
	}

	cols, err := rs.Columns()
	if nil != err {
		return err
	}
	show.Header(cols)

	colTypes, err := rs.ColumnTypes()
	if nil != err {
		return err
	}

	var types = make([]string, len(cols))
	for i, stype := range colTypes {
		types[i] = stype.DatabaseTypeName()
	}

	var num = 0
	for rs.Next() {
		var values []interface{} = nil
		for i := 0; i < len(cols); i++ {
			var tmp interface{}
			values = append(values, &tmp)
		}
		if err = rs.Scan(values...); nil != err {
			return err
		}
		show.Body(num, values, types)
		num++
	}

	return nil
}

func (s *SqlLite) stop(check bool) {
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
