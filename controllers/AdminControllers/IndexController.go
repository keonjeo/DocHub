package AdminControllers

import "github.com/astaxie/beego"

type IndexController struct {
	BaseController
}

func (controller *IndexController) Get() {
	controller.Data["BeegoVersion"] = beego.VERSION
	controller.Data["IsIndex"] = true
	controller.TplName = "index.html"
}
