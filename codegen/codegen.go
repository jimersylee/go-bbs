package main

import (
	"github.com/jimersylee/go-bbs/utils/simple"

	"github.com/jimersylee/go-bbs/model"
)

func main() {
	simple.Generate("./", "github.com/mlogclub/mlog", simple.GetGenerateStruct(&model.SysConfig{}))
}
