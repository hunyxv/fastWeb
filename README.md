# fastWeb
基于 [fasthttp](https://github.com/valyala/fasthttp) 和 [httpreuter](https://github.com/julienschmidt/httprouter) 的简单 web 框架。

```go
package main

import (
	"fastweb"
	"time"
)

func HelloWorld(ctx fastweb.Context) {
	ctx.SetBodyStrf(200, "[%s] Hello World!\n", time.Now())
}

func Query(ctx fastweb.Context) {
	var req struct {
		Field1 string    `valid:"field1,required,maxlength=10,minlength=3"`
		Field2 int       `valid:"field2,required"`
		Field3 time.Time `valid:"field3,format=2006-01-02 15:04:05"`
		Field4 string    `valid:"field4,re=a\\d{3}b,required"`
		Field5 string    `valid:"field5,strip"`
		Field6 float32   `valid:"a"`
		Field7 float64
		Field8 int32
	}

	err := ctx.QueryParams(&req)
	if err != nil {
		panic(err)
	}
	ctx.SetBodyStrf(200, "[%s] query: %+v\n", time.Now(), req)
}

func URLParams(ctx fastweb.Context) {
	ctx.SetBodyStrf(200, "[%s] params: %v\n", time.Now(), ctx.URLParams())
}

func main() {
	engine := fastweb.New()
	engine.GET("/", HelloWorld)
	group := engine.Group("/query")
	group.GET("/", Query)
	group.GET("/:id", URLParams)

	engine.Run("0.0.0.0:8080", fastweb.WithName("app"))
}
```
