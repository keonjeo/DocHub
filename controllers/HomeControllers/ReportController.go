package HomeControllers

import (
	"time"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type ReportController struct {
	BaseController
}

//举报
func (controller *ReportController) Get() {
	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "您当前未登录，请先登录")
	}

	reason, _ := controller.GetInt("Reason")
	did, _ := controller.GetInt("Did")

	if reason == 0 || did == 0 {
		controller.ResponseJson(false, "举报失败，请选择举报原因")
	}

	t := int(time.Now().Unix())
	report := models.Report{Status: false, Did: did, TimeCreate: t, TimeUpdate: t, Uid: controller.IsLogin, Reason: reason}
	rows, err := orm.NewOrm().Insert(&report)
	if err != nil {
		helper.Logger.Error("SQL执行失败：%v", err.Error())
	}
	if err != nil || rows == 0 {
		controller.ResponseJson(false, "举报失败：您已举报过该文档")
	}
	controller.ResponseJson(true, "恭喜您，举报成功，我们将在24小时内对您举报的内容进行处理。")
}
