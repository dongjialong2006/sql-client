package types

import (
	"github.com/mssola/colors"
)

const (
	QL       = "ql"
	SQLLITE3 = "sqllite3"
)

type Config struct {
	Addr string
	Path string
	File string
	Type string
}

type Handle interface {
	Exec(value string, color *colors.Color) error
	Stop()
}
