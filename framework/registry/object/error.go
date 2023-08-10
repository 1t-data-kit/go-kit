package object

import (
	"github.com/1t-data-kit/go-kit/base"
)

var (
	notSupportError = base.NewError("framework.registry.object.not-support", "不支持的注册对象类型(仅支持struct)")
)
