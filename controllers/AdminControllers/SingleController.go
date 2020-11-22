package AdminControllers

import (
	"net/http"
	"time"

	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type SingleController struct {
	BaseController
}

//单页列表
func (controller *SingleController) Get() {
	controller.Data["IsSingle"] = true
	controller.Data["Lists"], _, _ = models.NewPages().List(1000)
	controller.TplName = "index.html"
}

//单页编辑，只编辑文本内容
func (controller *SingleController) Edit() {
	var (
		page models.Pages
		cs   *models.CloudStore
		err  error
	)

	cs, err = models.NewCloudStore(false)
	if err != nil {
		controller.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	controller.Data["IsSingle"] = true
	alias := controller.GetString(":alias")

	if controller.Ctx.Request.Method == "POST" {
		controller.ParseForm(&page)
		page.TimeCreate = int(time.Now().Unix())
		page.Content = cs.ImageWithoutDomain(page.Content)
		_, err := orm.NewOrm().Update(&page)
		if err != nil {
			controller.ResponseJson(false, err.Error())
		}
		controller.ResponseJson(true, "更新成功")
	} else {
		page, _ = models.NewPages().One(alias)
		page.Content = cs.ImageWithDomain(page.Content)
		controller.Data["Data"] = page
		controller.TplName = "edit.html"
	}
}

//删除单页
func (controller *SingleController) Del() {
	id, _ := controller.GetInt("id")
	var page = models.Pages{Id: id}
	err := orm.NewOrm().Read(&page)
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}
	if _, err = orm.NewOrm().QueryTable(models.GetTablePages()).Filter("Id", page.Id).Delete(); err != nil {
		controller.ResponseJson(false, err.Error())
	}

	cs, _ := models.NewCloudStore(false)
	go cs.DeleteImageFromHtml(page.Content)

	controller.ResponseJson(true, "删除成功")
}
