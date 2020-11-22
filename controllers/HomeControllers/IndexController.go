package HomeControllers

import (
	"fmt"

	"strings"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type IndexController struct {
	BaseController
}

func (controller *IndexController) Get() {

	//获取横幅
	controller.Data["Banners"], _, _ = models.GetList(models.GetTableBanner(), 1, 100, orm.NewCondition().And("status", 1), "Sort")

	//判断用户是否已登录，如果已登录，则返回用户信息
	if controller.IsLogin > 0 {
		users, rows, err := models.NewUser().UserList(1, 1, "", "*", "i.`Id`=?", controller.IsLogin)
		if err != nil {
			helper.Logger.Error(err.Error())
		}
		if rows > 0 {
			controller.Data["User"] = users[0]
		} else {
			//如果用户不存在，则重置cookie
			controller.IsLogin = 0
			controller.ResetCookie()
		}
		controller.Data["LoginUid"] = controller.IsLogin
	} else {
		controller.Xsrf()
	}

	modelCate := models.NewCategory()
	//首页分类显示
	_, controller.Data["Cates"] = modelCate.GetAll(true)
	controller.Data["Latest"], _, _ = models.NewDocument().SimpleList(fmt.Sprintf("d.`Id` in(%v)", strings.Trim(controller.Sys.Trends, ",")), 5)
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Index", "文库首页", "文库首页", "文库首页", controller.Sys.Site)
	controller.Data["IsHome"] = true
	controller.Data["PageId"] = "wenku-index"
	controller.TplName = "index.html"
}
