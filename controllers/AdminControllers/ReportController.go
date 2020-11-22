package AdminControllers

import "dochub/models"

type ReportController struct {
	BaseController
}

//查看举报处理情况
func (controller *ReportController) Get() {
	controller.Data["Data"], _, _ = models.NewReport().Lists(1, 1000)
	controller.Data["IsReport"] = true
	controller.TplName = "index.html"
}
