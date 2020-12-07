package decoder

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

  "github.com/EAHITechnology/inf/golang/log_id"
	"github.com/gorilla/schema"
	"golang.org/x/net/context"
)

type Decoder struct {
	decoder *schema.Decoder
}

func NewDecoder() *Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return &Decoder{decoder: decoder}
}

/*
解析body参数，同时获取一个带有logid的ctx
*/
func (this *Decoder) ParseArg(r *http.Request, dst interface{}) (context.Context, error) {
	var ctx context.Context
	if logId := r.Header.Get(LOGID); logId != "" {
		ctx = log_id.GenValueCtx(LOGID, logId)
	} else {
		ctx = log_id.GenBackgroundCtx()
	}
	if dst != nil {
		b, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(b, dst); err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

/*
快速解析header参数到目标结构体,结构体以json为tag解析，header中不应该有map否则解析将失败
*/
func (this *Decoder) ParseHeader(r *http.Request, dst interface{}) error {
	headerMap := make(map[string]interface{})
	for key := range r.Header {
		val := r.Header.Get(key)
		headerMap[key] = val
	}
	if err := FillStruct(headerMap, dst); err != nil {
		return err
	}
	return nil
}

/*
使用Decoder获取url参数需要传入参数目标结构体指针dst和原始数据r.URL.Query()；
在结构体tag中尽量指明标签schema否则将按照驼峰进行解析，以免丢失字段
*/
func (this *Decoder) DecodeUrl(dst interface{}, src map[string][]string) error {
	if err := this.decoder.Decode(dst, src); err != nil {
		return err
	}
	return nil
}
