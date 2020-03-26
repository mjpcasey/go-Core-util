# RPC

## 调用方式

```go
type Request struct {
	Num  int
}

type Response struct {
	Num int
}
```

## server
```go
type TestReceiver struct {}

func (r *TestReceiver) TestCall(request *Request, reply *Response) (err error) {
	return nil
}

var conf rpc.ServerConfig

err = app.GetConfig().Scan("server", &conf)
if err == nil {
    err = app.GetApp().GetRpcManager().NewServer(conf, &TestReceiver{}) // 注册竞价核心rpc server
}
```

## client
````go
var conf rpc.ClientConfig
var cPath = "client"

err = app.GetConfig().Scan(cPath, &conf)
if err == nil {
    client, err = app.GetApp().GetRpcManager().NewClient(conf)
    if err == nil {
        err = client.Start()
    }
    
    if err == nil {
    	var request = Request{Num: 0}
        var response Response

        err := client.Call("TestReceiver.TestCall", &request, &response)
    }
} else {
    err = fmt.Errorf("读取配置[%s]错误: %s", cPath, err)
}
````

