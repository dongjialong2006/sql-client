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
	db     client.KeysAPI
	ctx    context.Context
	cfg    client.Config
	auth   client.AuthAPI
	member client.MembersAPI
	role   client.AuthRoleAPI
	user   client.AuthUserAPI
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
		e.auth = client.NewAuthAPI(conn)
		e.member = client.NewMembersAPI(conn)
		e.role = client.NewAuthRoleAPI(conn)
		e.user = client.NewAuthUserAPI(conn)
	}
	e.Unlock()

	return e.run(strings.Split(value, " "))
}

func (e *EtcdClient) parse(resp *client.Response) string {
	if nil == resp {
		return ""
	}

	if nil == resp.Node {
		return fmt.Sprintf("cluster id:%s, index:%d", resp.ClusterID, resp.Index)
	}

	if resp.Node.TTL > 0 {
		return fmt.Sprintf("cluster id:%s, index:%d, key:%s, ttl:%d", resp.ClusterID, resp.Index, resp.Node.Key, resp.Node.TTL)
	}

	return fmt.Sprintf("cluster id:%s, index:%d, key:%s", resp.ClusterID, resp.Index, resp.Node.Key)
}

func (e *EtcdClient) run(cmds []string) error {
	var err error = nil
	switch cmds[0] {
	case "user", "USER":
		err = e.users(cmds)
	case "role", "ROLE":
		err = e.roles(cmds)
	case "memeber", "MEMBER":
		err = e.members(cmds)
	case "auth", "AUTH":
		err = e.auths(cmds)
	default:
		err = e.keys(cmds)
	}

	return err
}

func (e *EtcdClient) keys(cmds []string) error {
	var resp *client.Response = nil
	var err error = nil
	switch strings.ToLower(cmds[0]) {
	case "set":
		if 3 != len(cmds) {
			return fmt.Errorf("set command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Set(e.ctx, cmds[1], cmds[2], &client.SetOptions{})
	case "mkdir", "setdir", "updatedir":
		if 3 != len(cmds) {
			return fmt.Errorf("set command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Set(e.ctx, cmds[1], cmds[2], &client.SetOptions{
			Dir: true,
		})
	case "get", "GET":
		if 2 != len(cmds) {
			return fmt.Errorf("get command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Get(e.ctx, cmds[1], &client.GetOptions{
			Sort: true,
		})
	case "update":
		if 3 != len(cmds) {
			return fmt.Errorf("update command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Update(e.ctx, cmds[1], cmds[2])
	case "delete", "rm":
		if 2 != len(cmds) {
			return fmt.Errorf("delete command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Delete(e.ctx, cmds[1], &client.DeleteOptions{})
	case "rmdir":
		if 2 != len(cmds) {
			return fmt.Errorf("delete command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Delete(e.ctx, cmds[1], &client.DeleteOptions{
			Dir: true,
		})
	case "mk":
		if 3 != len(cmds) {
			return fmt.Errorf("create command:%s format error.", strings.Join(cmds, " "))
		}
		resp, err = e.db.Create(e.ctx, cmds[1], cmds[2])
	case "watch":
		return fmt.Errorf("the system do not persist this command.")
	}

	show.Println(show.New(e.parse(resp), colors.Green))
	return err
}

func (e *EtcdClient) users(cmds []string) error {
	if len(cmds) >= 3 {
		return fmt.Errorf("user command:%s format error.", strings.Join(cmds, " "))
	}

	var err error = nil
	switch strings.ToLower(cmds[1]) {
	case "add":
		if 3 == len(cmds) {
			err = e.user.AddUser(e.ctx, cmds[2], "")
		} else {
			err = e.user.AddUser(e.ctx, cmds[2], cmds[3])
		}
	case "remove":
		err = e.user.RemoveUser(e.ctx, cmds[2])
	case "passwd":
		_, err = e.user.ChangePassword(e.ctx, cmds[2], cmds[3])
	case "get":
		user, err := e.user.GetUser(e.ctx, cmds[2])
		if nil != err {
			return err
		}
		e.parseUser(user)
	case "list":
		users, err := e.user.ListUsers(e.ctx)
		if nil != err {
			return err
		}
		if len(users) > 0 {
			show.Println(show.New(strings.Join(users, " "), colors.Green))
		}
	case "grant":
		_, err = e.user.GrantUser(e.ctx, cmds[2], cmds[3:])
	case "revoke":
		_, err = e.user.RevokeUser(e.ctx, cmds[2], cmds[3:])
	}

	return err
}

func (e *EtcdClient) parseUser(user *client.User) {
	if nil == user {
		return
	}

	show.Println(show.New(fmt.Sprintf("user:%s, pwd:%s, rols:%v", user.User, user.Password, user.Roles), colors.Green))
}

func (e *EtcdClient) auths(cmds []string) error {
	var err error = nil
	switch strings.ToLower(cmds[1]) {
	case "enable":
		err = e.auth.Enable(e.ctx)
	case "disable":
		err = e.auth.Disable(e.ctx)
	}

	return err
}

func (e *EtcdClient) roles(cmds []string) error {
	var err error = nil
	if len(cmds) >= 2 {
		return fmt.Errorf("role command:%s format error.", strings.Join(cmds, " "))
	}
	switch strings.ToLower(cmds[1]) {
	case "add":
		err = e.role.AddRole(e.ctx, cmds[2])
	case "get":
		role, err := e.role.GetRole(e.ctx, cmds[2])
		if nil != err {
			return err
		}
		e.parseRole(role)
	case "list":
		roles, err := e.role.ListRoles(e.ctx)
		if nil != err {
			return err
		}
		if len(roles) > 0 {
			show.Println(show.New(strings.Join(roles, " "), colors.Green))
		}
	case "remove":
		err = e.role.RemoveRole(e.ctx, cmds[2])
	case "grant":
		role, err := e.role.GrantRoleKV(e.ctx, cmds[2], cmds[3:], client.ReadWritePermission)
		if nil != err {
			return err
		}
		e.parseRole(role)
	case "revoke":
		role, err := e.role.RevokeRoleKV(e.ctx, cmds[2], cmds[3:], client.ReadWritePermission)
		if nil != err {
			return err
		}
		e.parseRole(role)
	}

	return err
}

func (e *EtcdClient) parseRole(role *client.Role) {
	if nil == role {
		return
	}

	var tmp string = fmt.Sprintf("role:%s", role.Role)
	if nil != role.Grant {
		if len(role.Grant.KV.Read) > 0 {
			tmp = fmt.Sprintf("%s, grant read:%s", tmp, strings.Join(role.Grant.KV.Read, " "))
		}
		if len(role.Grant.KV.Write) > 0 {
			tmp = fmt.Sprintf("%s, grant write:%s", tmp, strings.Join(role.Grant.KV.Write, " "))
		}
	}

	if len(role.Permissions.KV.Read) > 0 {
		tmp = fmt.Sprintf("%s, permissions read:%s", tmp, strings.Join(role.Permissions.KV.Read, " "))
	}
	if len(role.Permissions.KV.Write) > 0 {
		tmp = fmt.Sprintf("%s, permissions write:%s", tmp, strings.Join(role.Permissions.KV.Write, " "))
	}

	if nil != role.Revoke {
		if len(role.Revoke.KV.Read) > 0 {
			tmp = fmt.Sprintf("%s, revoke read:%s", tmp, strings.Join(role.Revoke.KV.Read, " "))
		}
		if len(role.Revoke.KV.Write) > 0 {
			tmp = fmt.Sprintf("%s, revoke write:%s", tmp, strings.Join(role.Revoke.KV.Write, " "))
		}
	}
	show.Println(show.New(tmp, colors.Green))
}

func (e *EtcdClient) members(cmds []string) error {
	if len(cmds) >= 2 {
		return fmt.Errorf("member command:%s format error.", strings.Join(cmds, " "))
	}
	var err error = nil
	switch strings.ToLower(cmds[1]) {
	case "add":
		member, err := e.member.Add(e.ctx, cmds[2])
		if nil != err {
			return err
		}
		e.parseMember(member)
	case "list":
		members, err := e.member.List(e.ctx)
		if nil != err {
			return err
		}
		e.parseMembers(members)
	case "update":
		err = e.member.Update(e.ctx, cmds[2], cmds[3:])
	case "remove":
		err = e.member.Remove(e.ctx, cmds[2])
	case "leader":
		member, err := e.member.Leader(e.ctx)
		if nil != err {
			return err
		}
		e.parseMember(member)
	}

	return err
}

func (e *EtcdClient) parseMembers(members []client.Member) {
	if 0 == len(members) {
		return
	}

	var infos []*show.ShowInfo = nil
	for i, member := range members {
		var tmp = fmt.Sprintf("member id:%s, name:%s", member.ID, member.Name)
		if len(member.ClientURLs) > 0 {
			tmp = fmt.Sprintf("%s, client urls:%s", tmp, strings.Join(member.ClientURLs, " "))
		}
		if len(member.PeerURLs) > 0 {
			tmp = fmt.Sprintf("%s, peer urls:%s", tmp, strings.Join(member.PeerURLs, " "))
		}
		infos = append(infos, show.New(tmp, colors.Green))
		if i != len(members)-1 {
			infos = append(infos, show.New("|", colors.Blue))
		}
	}

	show.PrintListln(infos)
}

func (e *EtcdClient) parseMember(member *client.Member) {
	if nil == member {
		return
	}

	var tmp string = fmt.Sprintf("member id:%s, name:%s", member.ID, member.Name)
	if len(member.ClientURLs) > 0 {
		tmp = fmt.Sprintf("%s, client urls:%s", tmp, strings.Join(member.ClientURLs, " "))
	}
	if len(member.PeerURLs) > 0 {
		tmp = fmt.Sprintf("%s, peer urls:%s", tmp, strings.Join(member.PeerURLs, " "))
	}
	show.Println(show.New(tmp, colors.Green))
}

func (e *EtcdClient) Stop() {
	<-e.ctx.Done()
	if nil != e.db {
		e.db = nil
	}
}
