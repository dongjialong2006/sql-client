# sql-client
the sql-client provide a way to access data in a ql database

QL数据库客户端

## 说明

- 本客户端用于连接`github.com/cznic/ql`数据库文件的客户端，支持sql语句操作。
- 该项目在mac和centos系统上均可编译，支持不同条目颜色显示。
- 当前仅支持(S)QL数据库的文件操作，内存库操作后续更新。

## 获取

- git clone https://github.com/dongjialong2006/sql-client.git

## 编译

- make update
- make

## 配置说明

- 各项配置说明如下：
- `file`：查询结果存储到文件的路径
- `path`：database file路径
- `driver`：database驱动

## 启动

- ./sql-client -path ./cache/agent.db

## 操作

- 如：`select * from xxx`, 按回车键即可