# zgo/http


## zgo/http/httplog

一行代码记录所有http日志

只需要在import中添加`_ "github.com/JoveYu/zgo/http/httplog/patch"` 即可自动打印http库的请求日志

```go
package main

import (
	"net/http"

	_ "github.com/JoveYu/zgo/http/httplog/patch"
	"github.com/JoveYu/zgo/log"
)

func main() {
	log.Install("stdout")
	http.Get("http://baidu.com")
}
```

## zgo/http/httpclient

扩展标准库的client，添加一些常用的函数，默认开启httplog

