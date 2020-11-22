package AdminControllers

import (
	"fmt"

	"dochub/helper"

	"strings"

	"dochub/helper/conv"
	"dochub/models"
)

//IT文库注册会员管理

type UserController struct {
	BaseController
}

func (controller *UserController) Prepare() {
	controller.BaseController.Prepare()
	controller.Data["IsUser"] = true
}

//用户列表
func (controller *UserController) List() {
	var (
		condition []string
		listRows  = 10
		id        = 0
		p         = 1
		username  string
	)
	//path中的参数
	params := conv.Path2Map(controller.GetString(":splat"))

	//页码处理
	if _, ok := params["p"]; ok {
		p = helper.Interface2Int(params["p"])
	} else {
		p, _ = controller.GetInt("p")
	}
	p = helper.NumberRange(p, 1, 1000000)

	//搜索的用户id处理
	if _, ok := params["id"]; ok {
		id = helper.Interface2Int(params["id"])
	} else {
		id, _ = controller.GetInt("id")
	}
	if id > 0 {
		condition = append(condition, fmt.Sprintf("i.Id=%v", id))
		controller.Data["Id"] = id
	}

	//搜索的用户名处理
	if _, ok := params["username"]; ok {
		username = params["username"]
	} else {
		username = controller.GetString("username")
	}
	if len(username) > 0 {
		condition = append(condition, fmt.Sprintf(`u.Username like "%v"`, "%"+username+"%"))
		controller.Data["Username"] = username
	}

	data, totalRows, err := models.NewUser().UserList(p, listRows, "", "*", strings.Join(condition, " and "))
	if err != nil {
		controller.Ctx.WriteString(err.Error())
		return
	}

	controller.Data["Page"] = helper.Paginations(6, totalRows, listRows, p, "/admin/user/", "id", id, "username", username)
	controller.Data["Users"] = data
	controller.Data["ListRows"] = listRows
	controller.Data["TotalRows"] = totalRows
	controller.TplName = "list.html"
}
