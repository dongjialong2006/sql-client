package sql

import (
	"context"
	"fmt"
	"sql-client/pkg/show"
	"sql-client/types"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/mssola/colors"
)

type EtcdClient struct {
	sync.RWMutex
	db  client.KeysAPI
	ctx context.Context
	cfg client.Config
}

func NewEtcdClient(ctx context.Context, opt *types.Options) *EtcdClient {
	if nil == opt.RedisOpt {
		return nil
	}

	srv := &EtcdClient{
		cfg: client.Config{
			Endpoints:               []string{opt.Addr},
			Transport:               client.DefaultTransport,
			Password:                opt.Pwd,
			HeaderTimeoutPerRequest: 3 * time.Second,
		},
		ctx: ctx,
	}

	return srv
}

func (e *EtcdClient) Exec(value string, color *colors.Color) error {
	if "" == value {
		return nil
	}

	e.Lock()
	if nil == e.db {
		conn, err := client.New(e.cfg)
		if err != nil {
			e.Unlock()
			return err
		}
		e.db = client.NewKeysAPI(conn)
	}
	e.Unlock()

	resp, err := e.run(strings.Split(value, " "))
	if nil != err {
		return err
	}

	show.TitlePrintln(show.New(fmt.Sprintf("%s", resp), colors.Green))
	return nil
}

func (e *EtcdClient) run(cmds []string) (*client.Response, error) {
	var resp *client.Response = nil
	var err error = nil
	switch cmds[0] {
	case "set", "SET":
		if 3 != len(cmds) {
			return nil, fmt.Errorf("set command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Set(e.ctx, cmds[1], cmds[2], &client.SetOptions{})
	case "get", "GET":
		if 2 != len(cmds) {
			return nil, fmt.Errorf("get command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Get(e.ctx, cmds[1], &client.GetOptions{
			Sort: true,
		})
	case "update", "UPDATE":
		if 3 != len(cmds) {
			return nil, fmt.Errorf("update command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Update(e.ctx, cmds[1], cmds[2])
	case "delete", "DELETE":
		if 2 != len(cmds) {
			return nil, fmt.Errorf("delete command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Delete(e.ctx, cmds[1], &client.DeleteOptions{})
	case "create", "CREATE":
		if 3 != len(cmds) {
			return nil, fmt.Errorf("create command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Create(e.ctx, cmds[1], cmds[2])
	default:
		return nil, fmt.Errorf("unknown db operation type:%s.", cmds[0])
	}

	return resp, err
}

func (e *EtcdClient) Stop() {
	<-e.ctx.Done()
	if nil != e.db {
		e.db = nil
	}
}
