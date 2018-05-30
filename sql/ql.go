package sql

import (
	"encoding/json"
	"fmt"
	"sql-client/pkg/file"
	"strings"
	"sync"

	"github.com/cznic/ql"
	"github.com/mssola/colors"
)

var rw sync.RWMutex
var ctx = ql.NewRWCtx()
var db *ql.DB = nil

func Create(path string) error {
	rw.Lock()
	defer rw.Unlock()
	var err error = nil
	db, err = ql.OpenFile(path, &ql.Options{})
	return err
}

func filter(value string) string {
	value = strings.Replace(value, "\n", " ", -1)
	return strings.Replace(value, "\"", "'", -1)
}

func Exec(value string, path string, color *colors.Color) error {
	if "" == value {
		return nil
	}

	rs, _, err := db.Run(ctx, fmt.Sprintf("BEGIN TRANSACTION; %s; COMMIT;", filter(value)))
	if nil != err {
		return err
	}
	if 0 == len(rs) {
		return nil
	}

	for _, tmp := range rs {
		if err = unMarshal(tmp, path, color); nil != err {
			return err
		}
	}

	return nil
}

func printInfo(info string, color *colors.Color) {
	if "" != info {
		fmt.Println(color.Get(info))
	}
}

func unMarshal(tmp ql.Recordset, path string, color *colors.Color) error {
	records, err := tmp.Rows(-1, 0)
	if nil != err {
		return err
	}

	if "" != path {
		if err = file.CreatePath(path); nil != err {
			return err
		}
	}

	color.Change(colors.Red, false)
	fmt.Println(color.Get(fmt.Sprintf("current query result num: %d.", len(records))))

	fields, err := tmp.Fields()
	if nil != err {
		return err
	}

	if err = show(path, fields, records, color); nil != err {
		return err
	}

	return nil
}

func show(path string, fields []string, rows [][]interface{}, color *colors.Color) error {
	header(fields, color)
	for i, row := range rows {
		if "" == path {
			draw(i, row, color)
			continue
		}
		data, err := json.Marshal(row)
		if nil != err {
			return err
		}
		if err = file.WriteFile(path, data); nil != err {
			return err
		}
	}

	return nil
}

func header(fields []string, color *colors.Color) {
	fmt.Print(">>> ")
	color.Change(colors.Red, false)
	var tmp string = color.Get("table fields: ")
	for i, field := range fields {
		color.Change(colors.Red, false)
		tmp += color.Get(fmt.Sprintf("  %s  ", field))
		if i != len(fields)-1 {
			color.Change(colors.Blue, false)
			tmp += color.Get("|")
		}
	}
	printInfo(tmp, color)
	return
}

func draw(num int, values []interface{}, color *colors.Color) {
	color.Change(colors.Red, false)
	var tmp string = color.Get(fmt.Sprintf("    row num:%d ", num))
	fmt.Print(tmp)

	tmp = ""
	for i, value := range values {
		color.Change(colors.Green, false)
		tmp += color.Get(fmt.Sprintf(" %v ", value))
		if i != len(values)-1 {
			color.Change(colors.Blue, false)
			tmp += color.Get("|")
		}
	}
	printInfo(tmp, color)
	return
}

func Close() {
	rw.Lock()
	defer rw.Unlock()
	if nil != db {
		db.Close()
		db = nil
	}
}
