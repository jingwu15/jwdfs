### JW-DFS 分布式文件系统

```
JW-DFS 分布式文件系统 是为了解决业务环境中因负责均衡，有多台客户机上传文件到中心文件服务器的业务。
JW-DFS 使用 json 作为配置文件，默认 /etc/jwdfs.json，易于使用
```

#### 1. 安装
```
# 1. 安装信赖包
sh ./install.sh

# 2. 构建项目
go build -ldflags "-s -w" jwdfs jwdfs
```

#### 2. 配置 /etc/jwdfs.json
```
{
    "server" : {
        "host":"127.0.0.1",
        "port":"8058",
        "updir":"/data/jwdfs"
    },
    "client": {
        "host":"127.0.0.1",
        "port":"8058",
        "downdir":"/tmp/download"
    }
}
```

#### 3. 使用
```
[root@master jwdfs]$ ./jwdfs
JW-DFS is a file upload/download server. You can use it to many computer.

Usage:
  ./jwdfs [command]

Available Commands:
  client      upload file to JW-DFS
  help        Help about any command
  server      the JW-DFS Server
  version     show version for JW-DFS

Flags:
  -d, --debug string   debug, default false (default "false")
  -h, --help           help for ./jwdfs

Use "./jwdfs [command] --help" for more information about a command.
```
