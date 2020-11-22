package HomeControllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"time"

	"dochub/helper"
	"dochub/helper/conv"
	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type UploadController struct {
	BaseController
}

//分词
func (controller *UploadController) SegWord() {
	var wds string
	if controller.IsLogin > 0 {
		wds = helper.SegWord(controller.GetString("word"))
	}
	controller.ResponseJson(true, "分词成功", wds)
}

//文档上传页面
func (controller *UploadController) Get() {
	cond := orm.NewCondition().And("status", 1)
	data, _, _ := models.GetList(models.GetTableCategory(), 1, 2000, cond, "Sort", "Title")
	controller.Xsrf()
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Upload", "文档上传-文档分享", "文档上传,文档分享", "文档上传-文档分享", controller.Sys.Site)
	controller.Data["Cates"], _ = conv.InterfaceToJson(data)
	controller.Data["json"] = data
	controller.Data["IsUpload"] = true
	controller.Data["PageId"] = "wenku-upload"
	controller.Data["MaxSize"] = models.NewSys().GetByField("MaxFile").MaxFile
	controller.TplName = "index.html"
}

//文档执行操作
//处理流程：
//1、检测用户是否已登录，未登录不允许上传
//2、检测是否存在了该文档的md5，如果已存在，则根据md5查询存储在文档存档表中的数据；如果文档已经在文档存储表中存在，则该文档不需要再获取封面、大小、页码等数据
//3、检测文档格式是否符合要求。
//4、计算文档md5，然后根据md5再比对一次文档是否在存档表中存在
//5、文档未存在，则将文档数据录入文档存储表(document_store)
//6、执行文档转pdf，并获取文档页数、封面、摘要等
//7、获取文档大小
func (controller *UploadController) Post() {
	var (
		ext  string //文档扩展名
		dir  = fmt.Sprintf("./uploads/%v/%v", time.Now().Format("2006/01/02"), controller.IsLogin)
		form models.FormUpload
		err  error
	)

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "您当前未登录，请先登录")
	}

	controller.ParseForm(&form)

	//文件在文档库中未存在，则接收文件并做处理
	f, fh, err := controller.GetFile("File")
	if err == nil {
		defer f.Close()
		os.MkdirAll(dir, os.ModePerm)
		ext = strings.ToLower(filepath.Ext(fh.Filename))
		if _, ok := helper.AllowedUploadDocsExt[ext]; !ok {
			controller.ResponseJson(false, "您上传的文档格式不正确，请上传正确格式的文档")
		}
		file := fmt.Sprintf("%v-%v-%v", controller.IsLogin, time.Now().UnixNano(), ext)
		form.TmpFile = filepath.Join(dir, file)
		form.Ext = ext
		err = controller.SaveToFile("File", form.TmpFile)
		if err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, "文件保存失败")
		}
	}

	// 文档处理
	err = models.DocumentProcess(controller.IsLogin, form)
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}
	controller.ResponseJson(true, "恭喜您，文档上传成功")
}
