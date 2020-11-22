package HomeControllers

import "github.com/astaxie/beego"

type ErrorsController struct {
	beego.Controller
}

//404
func (controller *ErrorsController) Error404() {
	referer := controller.Ctx.Request.Referer()
	controller.Layout = ""
	controller.Data["content"] = "Page Not Foud"
	controller.Data["code"] = "404"
	controller.Data["content_zh"] = "页面不存在"
	controller.Data["Referer"] = referer
	if len(referer) > 0 {
		controller.Data["IsReferer"] = true
	}
	controller.TplName = "error.html"
}

//501
func (controller *ErrorsController) Error501() {
	controller.Layout = ""
	controller.Data["code"] = "501"
	controller.Data["content"] = "Server Error"
	controller.Data["content_zh"] = "服务内部错误"
	controller.TplName = "error.html"
}

//数据库错误
func (controller *ErrorsController) ErrorDb() {
	controller.Layout = ""
	controller.Data["content"] = "Database is now down"
	controller.Data["content_zh"] = "数据库访问失败"
	controller.TplName = "error.html"
}
