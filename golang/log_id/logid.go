package log_id

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func genLogId() string {
	var t int64 = time.Now().UnixNano() / 1000000
	var r int = rand.Intn(10000)
	return fmt.Sprintf("%d%d", t, r)
}

/*
GenCtx 生成一个新的携带logid的ctx，优先使用ctx中的logid，降级使用grpc中logid，兜底生成logid
*/
func GenCtx(ctx context.Context) context.Context {
	backgroundCtx := context.Background()
	md, ok := metadata.FromIncomingContext(ctx)
	if ctx.Value(logid) != nil {
		tempCtx := context.WithValue(backgroundCtx, logid, ctx.Value(logid).(string))
		return metadata.NewOutgoingContext(tempCtx, metadata.Pairs(logid, ctx.Value(logid).(string)))
	}
	if ok {
		if mds, ok := md[logid]; ok {
			if len(mds) != 0 {
				tempCtx := metadata.NewOutgoingContext(backgroundCtx, metadata.Pairs(logid, mds[0]))
				return context.WithValue(tempCtx, logid, mds[0])
			}
		}
	}
	logId := genLogId()
	tempCtx := metadata.NewOutgoingContext(backgroundCtx, metadata.Pairs(logid, logId))
	return context.WithValue(tempCtx, logid, logId)
}

func GenBackgroundCtx() context.Context {
	logId := genLogId()
	tempCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(logid, logId))
	return context.WithValue(tempCtx, logid, logId)
}

func GenValueCtx(key, value string) context.Context {
	tempCtx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(key, value))
	return context.WithValue(tempCtx, key, value)
}
