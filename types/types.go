package types

import (
	"crypto/tls"

	"github.com/mssola/colors"
)

const (
	QL       = "ql"
	ETCD     = "etcd"
	REDIS    = "redis"
	SQLLITE3 = "sqllite3"
)

type Options struct {
	Type     string
	Addr     string
	Pwd      string
	RedisOpt *RedisOption
}

func NewOptions() *Options {
	return &Options{
		RedisOpt: &RedisOption{},
	}
}

type RedisOption struct {
	Cfg *tls.Config
	Vfy bool
	DB  int
}

type Handle interface {
	Exec(value string, color *colors.Color) error
	Stop()
}
