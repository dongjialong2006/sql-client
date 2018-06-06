package sql

import (
	"context"
	"fmt"
	"reflect"
	"sql-client/pkg/show"
	"sql-client/types"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mssola/colors"
)

type RedisClient struct {
	sync.RWMutex
	db   redis.Conn
	ctx  context.Context
	addr string
	opts []redis.DialOption
}

func NewRedisClient(ctx context.Context, opt *types.Options) *RedisClient {
	if nil == opt.RedisOpt {
		return nil
	}

	srv := &RedisClient{
		addr: opt.Addr,
		ctx:  ctx,
	}

	srv.opts = append(srv.opts, redis.DialConnectTimeout(30*time.Second))
	srv.opts = append(srv.opts, redis.DialDatabase(opt.RedisOpt.DB))
	if "" != opt.Pwd {
		srv.opts = append(srv.opts, redis.DialPassword(opt.Pwd))
	}
	srv.opts = append(srv.opts, redis.DialReadTimeout(5*time.Second))
	srv.opts = append(srv.opts, redis.DialWriteTimeout(30*time.Second))
	if nil != opt.RedisOpt.Cfg {
		srv.opts = append(srv.opts, redis.DialTLSConfig(opt.RedisOpt.Cfg))
		srv.opts = append(srv.opts, redis.DialTLSSkipVerify(opt.RedisOpt.Vfy))
	}

	go srv.stop(true)

	return srv
}

func (r *RedisClient) Exec(value string, color *colors.Color) error {
	if "" == value {
		return nil
	}

	var err error = nil
	r.Lock()
	if nil == r.db {
		r.db, err = redis.Dial("tcp", r.addr, r.opts...)
		if nil != err {
			r.Unlock()
			return err
		}
	}
	r.Unlock()

	cmds := strings.Split(value, " ")
	if 0 == len(cmds) {
		return fmt.Errorf("command:%s format error.", value)
	}

	var args []interface{} = nil
	for _, cmd := range cmds[1:] {
		args = append(args, cmd)
	}
	resp, err := r.db.Do(cmds[0], args...)
	if nil != err {
		return err
	}

	return r.parse(resp)
}

func (r *RedisClient) parse(resp interface{}) error {
	if nil == resp {
		return nil
	}
	switch resp.(type) {
	case int, int32, int64, uint, uint32, uint64:
		show.Println(show.New(fmt.Sprintf("%d", resp), colors.Green))
	case string, []byte:
		show.Println(show.New(fmt.Sprintf("%s", resp), colors.Green))
	case []interface{}:
		var keys []string = nil
		for _, key := range resp.([]interface{}) {
			tmp := r.transfer(key)
			if "" != tmp {
				keys = append(keys, tmp)
			}
		}
		if len(keys) > 0 {
			show.Println(show.New(strings.Join(keys, " "), colors.Green))
		}
	default:
		show.Println(show.New(fmt.Sprintf("unknown resp type:%s.", reflect.TypeOf(resp).String()), colors.Red))
	}

	return nil
}

func (r *RedisClient) transfer(resp interface{}) string {
	if nil == resp {
		return ""
	}
	switch resp.(type) {
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", resp)
	case string, []byte:
		return fmt.Sprintf("%s", resp)
	default:
		return fmt.Sprintf("unknown resp type:%s.", reflect.TypeOf(resp).String())
	}

	return ""
}

func (r *RedisClient) Stop() {
	r.stop((false))
}

func (r *RedisClient) stop(check bool) {
	if check {
		<-r.ctx.Done()
	}
	r.Lock()
	if nil != r.db {
		r.db.Close()
		r.db = nil
	}
	r.Unlock()
}
