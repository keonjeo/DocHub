package HomeControllers

import (
	"html/template"
	"strings"

	"fmt"
	"time"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Output struct {
	status int
	msg    string
}

type BaseController struct {
	beego.Controller
	TplTheme  string //模板主题
	TplStatic string //模板静态文件
	IsLogin   int    //用户是否已登录
	Sys       models.Sys
	Out       Output
}

//初始化函数
func (controller *BaseController) Prepare() {
	ctrl, _ := controller.GetControllerAndAction()
	ctrl = strings.TrimSuffix(ctrl, "Controller")

	//设置默认模板
	controller.TplTheme = "default"
	controller.TplPrefix = "Home/" + controller.TplTheme + "/" + ctrl + "/"
	controller.Layout = "Home/" + controller.TplTheme + "/layout.html"

	//防止跨站攻击
	//检测用户是否已经在cookie存在登录
	controller.checkCookieLogin()

	//初始化
	controller.Data["LoginUid"] = controller.IsLogin
	//当前模板静态文件
	controller.Data["TplStatic"] = "/static/Home/" + controller.TplTheme

	version := helper.VERSION
	if helper.Debug { //debug模式下，每次更新js
		version = fmt.Sprintf("%v.%v", version, time.Now().Unix())
	}
	controller.Sys, _ = models.NewSys().Get()
	controller.Data["Version"] = version
	controller.Data["Sys"] = controller.Sys
	controller.Data["Chanels"] = models.NewCategory().GetByPid(0, true)
	controller.Data["Pages"], _, _ = models.NewPages().List(beego.AppConfig.DefaultInt("pageslimit", 6), 1)
	controller.Data["AdminId"] = helper.Interface2Int(controller.GetSession("AdminId"))
	controller.Data["CopyrightDate"] = time.Now().Format("2006")

	controller.Data["PreviewDomain"] = ""

	if cs, err := models.NewCloudStore(false); err == nil {
		controller.Data["PreviewDomain"] = cs.GetPublicDomain()
	} else {
		helper.Logger.Error(err.Error())
	}

}

//是否已经登录，如果已登录，则返回用户的id
func (controller *BaseController) CheckLogin() int {
	uid := controller.GetSession("uid")
	if uid != nil {
		id, ok := uid.(int)
		if ok && id > 0 {
			return id
		}
	}
	return 0
}

//防止跨站攻击，在有表单的控制器放大中调用，不要直接在base控制器中调用，因为用户每访问一个页面都重新刷新cookie了
func (controller *BaseController) Xsrf() {
	//使用的时候，直接在模板表单添加{{.xsrfdata}}
	controller.Data["xsrfdata"] = template.HTML(controller.XSRFFormHTML())
}

//检测用户登录的cookie是否存在
func (controller *BaseController) checkCookieLogin() {
	secret := beego.AppConfig.DefaultString("CookieSecret", helper.DEFAULT_COOKIE_SECRET)
	timestamp, ok := controller.GetSecureCookie(secret, "uid")
	if !ok {
		return
	}
	uid, ok := controller.Ctx.GetSecureCookie(secret+timestamp, "token")
	if !ok || len(uid) == 0 {
		controller.ResetCookie()
	}

	if controller.IsLogin = helper.Interface2Int(uid); controller.IsLogin > 0 {
		if info := models.NewUser().UserInfo(controller.IsLogin); info.Status == false {
			//被封禁的账号，重置cookie
			controller.ResetCookie()
		}
	}
}

//重置cookie
func (controller *BaseController) ResetCookie() {
	controller.Ctx.SetCookie("uid", "")
	controller.Ctx.SetCookie("token", "")
}

//设置用户登录的cookie，其实uid是时间戳的加密，而token才是真正的uid
//@param            uid         interface{}         用户UID
func (controller *BaseController) SetCookieLogin(uid interface{}) {
	secret := beego.AppConfig.DefaultString("CookieSecret", helper.DEFAULT_COOKIE_SECRET)
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	expire := 3600 * 24 * 365
	controller.Ctx.SetSecureCookie(secret, "uid", timestamp, expire)
	controller.Ctx.SetSecureCookie(secret+timestamp, "token", fmt.Sprintf("%v", uid), expire)
}

//校验文档是否已经存在
func (controller *BaseController) DocExist() {
	if models.NewDocument().IsExistByMd5(controller.GetString("md5")) > 0 {
		controller.ResponseJson(true, "文档存在")
	}
	controller.ResponseJson(false, "文档不存在")
}

//响应json
func (controller *BaseController) ResponseJson(isSuccess bool, msg string, data ...interface{}) {
	status := 0
	if isSuccess {
		status = 1
	}
	ret := map[string]interface{}{"status": status, "msg": msg}
	if len(data) > 0 {
		ret["data"] = data[0]
	}
	controller.Data["json"] = ret
	controller.ServeJSON()
	controller.StopRun()
}

//单页
func (controller *BaseController) Pages() {
	alias := controller.GetString(":page")
	page, err := models.NewPages().One(alias)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.Abort("404")
	}
	if page.Id == 0 || page.Status == false {
		controller.Abort("404")
	}
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Pages", page.Title, page.Keywords, page.Description, controller.Sys.Site)
	page.Vcnt += 1
	orm.NewOrm().Update(&page, "Vcnt")
	cs, _ := models.NewCloudStore(false)
	page.Content = cs.ImageWithDomain(page.Content)

	controller.Data["Page"] = page
	controller.Data["Lists"], _, _ = models.NewPages().List(20, 1)
	controller.Data["PageId"] = "wenku-content"
	controller.TplName = "pages.html"
}
