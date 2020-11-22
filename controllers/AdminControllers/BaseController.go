package AdminControllers

import (
	"strings"

	"dochub/models"

	"time"

	"fmt"

	"dochub/helper"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	TplTheme  string //模板主题
	TplStatic string //模板静态文件
	AdminId   int    //管理员是否已登录，如果已登录，则管理员ID大于0
	Sys       models.Sys
}

//初始化函数
func (controller *BaseController) Prepare() {
	var ok bool
	controller.Sys, _ = models.NewSys().Get()
	//检测是否已登录，未登录则跳转到登录页
	AdminId := controller.GetSession("AdminId")
	controller.AdminId, ok = AdminId.(int)
	controller.Data["Admin"], _ = models.NewAdmin().GetById(controller.AdminId)
	if !ok || controller.AdminId == 0 {
		controller.Redirect("/admin/login", 302)
		return
	}

	version := helper.VERSION
	if helper.Debug {
		version = fmt.Sprintf("%v.%v", version, time.Now().Unix())
	}
	controller.Data["Version"] = version
	//后台关闭XSRF功能
	controller.EnableXSRF = false
	ctrl, _ := controller.GetControllerAndAction()
	ctrl = strings.TrimSuffix(ctrl, "Controller")
	//设置默认模板
	controller.TplTheme = "default"
	controller.TplPrefix = "Admin/" + controller.TplTheme + "/" + ctrl + "/"
	controller.Layout = "Admin/" + controller.TplTheme + "/layout.html"
	//当前模板静态文件
	controller.Data["TplStatic"] = "/static/Admin/" + controller.TplTheme
	//controller.Data["PreviewDomain"] = beego.AppConfig.String("oss::PreviewUrl")
	if cs, err := models.NewCloudStore(false); err == nil {
		controller.Data["PreviewDomain"] = cs.GetPublicDomain()
	} else {
		helper.Logger.Error(err.Error())
		controller.Data["PreviewDomain"] = ""
	}
	controller.Data["Sys"] = controller.Sys
	controller.Data["Title"] = "文库系统管理后台"
	controller.Data["Lang"] = "zh-CN"
}

//自定义的文档错误
func (controller *BaseController) ErrorDiy(status, redirect, msg interface{}, timewait int) {
	controller.Data["status"] = status
	controller.Data["redirect"] = redirect
	controller.Data["msg"] = msg
	controller.Data["timewait"] = timewait
	controller.TplName = "error_diy.html"
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

//404
func (controller *BaseController) Error404() {
	controller.Layout = ""
	controller.Data["content"] = "Page Not Foud"
	controller.Data["code"] = "404"
	controller.Data["content_zh"] = "页面被外星人带走了"
	controller.TplName = "error.html"
}

//501
func (controller *BaseController) Error501() {
	controller.Layout = ""
	controller.Data["code"] = "501"
	controller.Data["content"] = "Server Error"
	controller.Data["content_zh"] = "服务器被外星人戳炸了"
	controller.TplName = "error.html"
}

//数据库错误
func (controller *BaseController) ErrorDb() {
	controller.Layout = ""
	controller.Data["content"] = "Database is now down"
	controller.Data["content_zh"] = "数据库被外星人抢走了"
	controller.TplName = "error.html"
}

//更新内容
func (controller *BaseController) Update() {
	id := strings.Split(controller.GetString("id"), ",")
	i, err := models.UpdateByIds(controller.GetString("table"), controller.GetString("field"), controller.GetString("value"), id)
	ret := map[string]interface{}{"status": 0, "msg": "更新失败，可能您未对内容作更改"}
	if i > 0 && err == nil {
		ret["status"] = 1
		ret["msg"] = "更新成功"
	}
	if err != nil {
		ret["msg"] = err.Error()
	}
	controller.Data["json"] = ret
	controller.ServeJSON()
}

//删除内容
func (controller *BaseController) Del() {
	id := strings.Split(controller.GetString("id"), ",")
	i, err := models.DelByIds(controller.GetString("table"), id)
	ret := map[string]interface{}{"status": 0, "msg": "删除失败，可能您要删除的内容已经不存在"}
	if i > 0 && err == nil {
		ret["status"] = 1
		ret["msg"] = "删除成功"
	}
	if err != nil {
		ret["msg"] = err.Error()
	}
	controller.Data["json"] = ret
	controller.ServeJSON()
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

//响应json
func (controller *BaseController) Response(data map[string]interface{}) {
	controller.Data["json"] = data
	controller.ServeJSON()
	controller.StopRun()
}
