
# zgo/web

## 介绍

对自带的http模块的适当封装，尽量贴近HTTP的本质，方便使用

1. 引入ctx 提供常用的工具函数 直接使用 参考`context.go`
2. 基于正则表达式的路由(30行实现)，灵活简单，支持url配置输入参数
3. 灵活的中间件支持，目前自带CORS JSONP 中间件
4. 完全兼容http库，支持不使用路由，独立使用context，使用`ContextHandler`
5. 支持对请求报文进行调试 使用`web.DefaultServer.Debug=true`

## 文档

[https://godoc.org/github.com/JoveYu/zgo/web](https://godoc.org/github.com/JoveYu/zgo/web)

## TODO

1. 后续可以考虑支持多backend 比如fasthttp 目前对极致性能需求不大
2. 目前正则路由已经比较灵活，后续有需要在考虑更强大正则
3. 控制内存分配，引入pool

## Example

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/JoveYu/zgo/log"
	"github.com/JoveYu/zgo/web"
)

func ping(ctx web.Context) {
	ctx.Abort(403, fmt.Sprintf("%s pong", ctx.Method()))
}

func hello(ctx web.Context) {
	name := ctx.Param("name")

	ctx.WriteHeader(200)
	ctx.WriteString("hello ")
	ctx.WriteString(name)
}

func add(ctx web.Context) {
	astr := ctx.Param("a")
	bstr := ctx.Param("b")

	a, _ := strconv.ParseInt(astr, 10, 64)
	b, _ := strconv.ParseInt(bstr, 10, 64)

	ctx.WriteHeader(200)

	ctx.WriteString(fmt.Sprintf("%d", a+b))
}

func redir(ctx web.Context) {
	ctx.Redirect(302, "http://baidu.com")
}

func query(ctx web.Context) {
	query := ctx.GetQuery("test")
	ctx.WriteHeader(200)
	ctx.WriteJSON(map[string]string{"test": query})
}

func main() {
	log.Install("stdout")

	// curl /ping  ->  pong
	web.GET("^/ping$", ping)
	web.POST("^/ping$", ping)

	// curl /params/world  ->  hello world
	web.GET("^/params/(?P<name>\\w+)$", hello)

	// curl /add/1/2  ->  3
	web.GET("^/add/(?P<a>\\d+)/(?P<b>\\d+)$", add)

	// curl /redir  ->  to baidu
	web.GET("^/redir$", redir)

	// curl /query?test=123  -> {"test":"123"}
	web.GET("^/query$", query)

	web.Run("127.0.0.1:7000")
}

```


