package AdminControllers

import (
	"fmt"
	"os"
	"time"

	"dochub/helper"
	"dochub/models"
)

type KindEditorController struct {
	BaseController
}

//上传。这里是后台使用的，不限文件类型
func (controller *KindEditorController) Upload() {
	f, fh, err := controller.GetFile("imgFile")
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}
	defer f.Close()
	now := time.Now()
	dir := fmt.Sprintf("uploads/kindeditor/%v", now.Format("2006/01/02"))
	os.MkdirAll(dir, 0777)
	ext := helper.GetSuffix(fh.Filename, ".")
	filename := "article." + helper.MD5Crypt(fmt.Sprintf("%v-%v-%v", now, fh.Filename, controller.AdminId)) + "." + ext
	//存储文件
	tmpFile := dir + "/" + filename
	err = controller.SaveToFile("imgFile", tmpFile)
	if err != nil {
		controller.Response(map[string]interface{}{"message": err.Error(), "error": 1})
	}
	defer os.RemoveAll(tmpFile)

	var cs *models.CloudStore
	if cs, err = models.NewCloudStore(false); err != nil {
		controller.Response(map[string]interface{}{"message": err.Error(), "error": 1})
	}

	//将文件上传到OSS
	err = cs.Upload(tmpFile, filename)
	if err == nil {
		controller.Response(map[string]interface{}{"url": cs.GetSignURL(filename), "error": 0})
	}
	controller.Response(map[string]interface{}{"message": err.Error(), "error": 1})
}
