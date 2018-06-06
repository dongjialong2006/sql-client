# sql-client
the sql-client provide a way to access data in a ql database

QL数据库客户端

## 说明

- 本客户端支持sql语句操作。
- 该项目在mac和centos系统上均可编译，支持不同条目颜色显示。
- 当前仅支持(S)QL数据库的文件操作、redis数据库、etcd数据库和SQLLite3数据库。

## 获取

- git clone https://github.com/dongjialong2006/sql-client.git

## 编译

- make update
- make

## 配置说明

- 各项配置说明如下：
- `addr`：db file路径或访问地址
- `type`：db类型
- `pwd`：redis或etcd数据库访问密码
- `db`：指定redis数据库

## 启动

- ./sql-client -addr ./cache/agent.db

## 操作

- 如：`select * from xxx`, 按回车键即可