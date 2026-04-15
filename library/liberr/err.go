package liberr

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func ErrIsNil(ctx context.Context, err error, msg ...string) {
	if !g.IsNil(err) {
		if len(msg) > 0 {
			g.Log().Error(ctx, err.Error())
			// 包含原始错误信息
			panic(fmt.Sprintf("%s: %v", msg[0], err))
		} else {
			panic(err.Error())
		}
	}
}

func ValueIsNil(value interface{}, msg string) {
	if g.IsNil(value) {
		panic(msg)
	}
}