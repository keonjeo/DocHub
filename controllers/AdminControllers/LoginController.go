package AdminControllers

import (
	"html/template"

	"time"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	BaseController
}

//重置prepare方法，移除模板继承
func (controller *LoginController) Prepare() {
	controller.EnableXSRF = false
	//设置默认模板
	TplTheme := "default"
	controller.TplPrefix = "Admin/" + TplTheme + "/Login/"
	controller.Layout = ""
	//当前模板静态文件
	controller.Data["TplStatic"] = "/static/Admin/" + TplTheme
	controller.AdminId = helper.Interface2Int(controller.GetSession("AdminId"))
}

//登录后台
func (controller *LoginController) Login() {
	controller.EnableXSRF = true
	controller.Data["Sys"], _ = models.NewSys().Get()
	if controller.Ctx.Request.Method == "GET" {
		controller.Xsrf()
		controller.TplName = "index.html"
	} else {
		var (
			msg   string = "登录失败，用户名或密码不正确"
			admin models.Admin
		)
		controller.ParseForm(&admin)
		if admin, err := models.NewAdmin().Login(admin.Username, admin.Password, admin.Code); err == nil && admin.Id > 0 {
			controller.SetSession("AdminId", admin.Id)
			controller.ResponseJson(true, "登录成功")
		} else {
			controller.ResponseJson(false, msg)
		}
	}
}

//更新登录密码
func (controller *LoginController) UpdatePwd() {
	if controller.AdminId > 0 {
		PwdOld := controller.GetString("password_old")
		PwdNew := controller.GetString("password_new")
		PwdEnsure := controller.GetString("password_ensure")
		if PwdOld == PwdNew || PwdNew != PwdEnsure {
			controller.ResponseJson(false, "新密码不能与原密码相同，且确认密码必须与新密码一致")
		} else {
			var admin = models.Admin{Password: helper.MD5Crypt(PwdOld)}
			if orm.NewOrm().Read(&admin, "Password"); admin.Id > 0 {
				admin.Password = helper.MD5Crypt(PwdNew)
				if rows, err := orm.NewOrm().Update(&admin); rows > 0 {
					controller.ResponseJson(true, "密码更新成功")
				} else {
					controller.ResponseJson(false, "密码更新失败："+err.Error())
				}

			} else {
				controller.ResponseJson(false, "原密码不正确")
			}
		}
	} else {
		controller.Error404()
	}
}

//更新管理员信息
func (controller *LoginController) UpdateAdmin() {
	if controller.AdminId > 0 {
		code := controller.GetString("code")
		if code == "" {
			controller.ResponseJson(false, "登录验证码不能为空")
		}

		username := controller.GetString("username")
		if username == "" {
			controller.ResponseJson(false, "登录用户名不能为空")
		}

		email := controller.GetString("email")

		admin := models.Admin{
			Id:       controller.AdminId,
			Username: username,
			Code:     code,
			Email:    email,
		}

		if _, err := orm.NewOrm().Update(&admin, "Code", "Email", "Username"); err == nil {
			controller.ResponseJson(true, "资料更新成功")
		} else {
			controller.ResponseJson(false, "资料更新失败："+err.Error())
		}

	} else {
		controller.Error404()
	}
}

//退出登录
func (controller *LoginController) Logout() {
	controller.DelSession("AdminId")
	controller.Redirect("/admin/login?t="+time.Now().String(), 302)
}

//防止跨站攻击，在有表单的控制器中调用
func (controller *LoginController) Xsrf() {
	//使用的时候，直接在模板表单添加{{.xsrfdata}}
	controller.Data["xsrfdata"] = template.HTML(controller.XSRFFormHTML())
}
