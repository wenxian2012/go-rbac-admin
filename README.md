# Go Web RBAC Admin 

![语言](https://img.shields.io/badge/language-goland1.2-blue.svg)
![base](https://img.shields.io/badge/base-gin-blue.svg)
![base](https://img.shields.io/badge/base-casbin-blue.svg)

> 一个Go Web Api 服务, 包含 用户、权限、菜单、动作、资源、JWT等，可以用于快速构建私活项目的RBAC后台

## 目录结构
* conf: 用于存储配置文件
* docs: 文档
    * sql执行命令
    * API注释
* dto: 数据传输对象
* logs: 日志
* middleware:应用中间件
    * inject 初始化对象
    * jwt
    * permission  权限验证
* models: 应用数据库模型
* pkg: 第三方包
* routers: 路由逻辑处理
* service: 逻辑处理
* test: 单元测试
    
## API文档

> http://127.0.0.1:8000/swagger/index.html

## 部署

### 支持

- 部署 Mysql

### 库

Create a **go database** and import [SQL](https://github.com/wenxian2012/go-rbac-admin/blob/master/docs/sql/go.sql)

创建一个库 go,然后导入sql,创建表！

### 配置文件

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 127.0.0.1:3306
Name = go
TablePrefix = go_
```

### 安装部署
```

yum install go -y 


export GOPROXY=https://goproxy.io
go get github.com/wenxian2012/go-rbac-admin
cd $GOPATH/src/github.com/wenxian2012/go-rbac-admin
go build main.go
go run  main.go 
```


### 热编译(开发时使用)
```bash

go get github.com/silenceper/gowatch

gowatch   
```

## Features
```
- RESTful API
- Gorm
- logging
- Jwt-go
- Swagger
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
```

## 特别感谢

```
本项目主要参考了:
https://github.com/EDDYCJY/go-gin-example  包含更多的例子，上传文件图片等。本项目进行了增改。
https://github.com/LyricTian/gin-admin     主要为 RBAC 表、逻辑设计。
https://github.com/wenxian2012/go-rbac-admin     主要为 gin+ casbin例子。
```

## 其他
```shell
## 更新API文档
swag init 
```