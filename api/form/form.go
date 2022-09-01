package form

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func FormValueInt(ctx *gin.Context, name string) (int, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, errors.New(fmt.Sprintf("unable to find param value '%s'", name))
	}
	return strconv.Atoi(str)
}

func FormValueIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx *gin.Context, name string) (int64, error) {
	str := ctx.PostForm(name)
	if str == "" {
		return 0, errors.New(fmt.Sprintf("unable to find param value '%s'", name))
	}
	return strconv.ParseInt(str, 10, 64)
}

func FormValueInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}
