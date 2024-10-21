package web

import (
	"github.com/jimersylee/go-bbs/utils/oss"
	"github.com/jimersylee/go-bbs/utils/session"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"

	"github.com/jimersylee/go-bbs/utils/simple"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

const avatarMaxBytes int64 = 1024 * 1024 // 1M

type UploadController struct {
	Ctx *context.Context
}

func (this *UploadController) Post() *simple.JsonResult {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		return simple.Error(simple.ErrorNotLogin)
	}

	file, header, err := this.Ctx.FormFile("image")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	defer file.Close()

	if header.Size > avatarMaxBytes {
		return simple.ErrorMsg("图片不能超过1M")
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	logrus.Info("上传文件：", header.Filename, " size:", header.Size)

	url, err := oss.UploadImage(fileBytes)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 编辑器中上传
func (this *UploadController) PostEditor() {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		_ = this.Ctx.JSON(iris.Map{
			"msg":  "请先登录",
			"code": 1,
		})
		return
	}

	maxSize := this.Ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := this.Ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		this.Ctx.StatusCode(iris.StatusInternalServerError)
		this.Ctx.WriteString(err.Error())
		return
	}

	var errFiles []string
	var succMap map[string]string

	errFiles = make([]string, 0)
	succMap = make(map[string]string)

	form := this.Ctx.Request().MultipartForm
	files := form.File["file[]"]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}
		url, err := oss.UploadImage(fileBytes)
		if err != nil {
			logrus.Error(err)
			errFiles = append(errFiles, file.Filename)
			continue
		}

		succMap[file.Filename] = url
	}

	_ = this.Ctx.JSON(iris.Map{
		"msg":  "",
		"code": 0,
		"data": iris.Map{
			"errFiles": errFiles,
			"succMap":  succMap,
		},
	})
	return

}

func (this *UploadController) PostFetch() {
	user := session.GetCurrentUser(this.Ctx)
	if user == nil {
		_ = this.Ctx.JSON(iris.Map{
			"msg":  "请先登录",
			"code": 1,
		})
		return
	}

	var data map[string]string
	data = make(map[string]string)

	err := this.Ctx.ReadJSON(&data)
	if err != nil {
		_ = this.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
		})
		return
	}

	url := data["url"]
	output, err := oss.CopyImage(url)
	if err != nil {
		_ = this.Ctx.JSON(iris.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}
	_ = this.Ctx.JSON(iris.Map{
		"msg":  "",
		"code": 0,
		"data": iris.Map{
			"url": output,
		},
	})
}
