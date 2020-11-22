package AdminControllers

import (
	"time"

	"dochub/helper"

	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type FriendController struct {
	BaseController
}

//添加以及查看友链列表
func (controller *FriendController) Get() {
	if controller.Ctx.Request.Method == "POST" {
		var fr models.Friend
		controller.ParseForm(&fr)
		fr.Status = true
		fr.TimeCreate = int(time.Now().Unix())
		if i, err := orm.NewOrm().Insert(&fr); i > 0 && err == nil {
			controller.ResponseJson(true, "友链添加成功")
		} else {
			if err != nil {
				helper.Logger.Error(err.Error())
			}
			controller.ResponseJson(false, "友链添加失败，可能您要添加的友链已存在")
		}
	} else {
		controller.Data["IsFriend"] = true
		controller.Data["Friends"], _, _ = models.NewFriend().GetListByStatus(-1)
		controller.TplName = "index.html"
	}
}
