package HomeControllers

import (
	"fmt"
	"path/filepath"

	"github.com/astaxie/beego"

	"strings"

	"time"

	"os"

	"dochub/helper"
	"dochub/helper/conv"
	"dochub/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type UserController struct {
	BaseController
}

func (controller *UserController) Prepare() {
	controller.BaseController.Prepare()
	controller.Xsrf()
}

//会员中心
func (controller *UserController) Get() {
	uid, _ := controller.GetInt(":uid")
	path := controller.GetString(":splat")
	params := conv.Path2Map(path)
	//排序
	sort := "new"
	if param, ok := params["sort"]; ok {
		sort = param
	}
	//页码
	p := 1
	if page, ok := params["p"]; ok {
		p = helper.Interface2Int(page)
		if p < 1 {
			p = 1
		}
	}

	switch sort {
	case "dcnt":
		sort = "dcnt"
	case "score":
		sort = "score"
	case "vcnt":
		sort = "vcnt"
	case "ccnt":
		sort = "ccnt"
	default:
		sort = "new"
	}
	//显示风格
	style := "list"
	if s, ok := params["style"]; ok {
		style = s
	}
	if style != "th" {
		style = "list"
	}
	//cid:collect folder id ,收藏夹id
	cid := 0
	if s, ok := params["cid"]; ok {
		cid = helper.Interface2Int(s)
	}
	if p < 1 {
		p = 1
	}
	if uid < 1 {
		uid = controller.IsLogin
	}
	controller.Data["Uid"] = uid

	if uid <= 0 {
		controller.Redirect("/user/login", 302)
		return
	}

	listRows := 16
	user, rows, err := models.NewUser().GetById(uid)
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	if rows == 0 {
		controller.Redirect("/", 302)
		return
	}

	if cid > 0 {
		sql := fmt.Sprintf("select Title,Cnt from %v where Id=? limit 1", models.GetTableCollectFolder())
		var params []orm.Params
		orm.NewOrm().Raw(sql, cid).Values(&params)
		if len(params) == 0 {
			controller.Redirect(fmt.Sprintf("/user/%v/collect", uid), 302)
			return
		}

		controller.Data["Folder"] = params[0]
		fields := "di.Id,di.`Uid`, di.`Cid`, di.`TimeCreate`, di.`Dcnt`, di.`Vcnt`, di.`Ccnt`, di.`Score`, di.`Status`, di.`ChanelId`, di.`Pid`,c.Title Category,u.Username,d.Title,ds.`Md5`, ds.`Ext`, ds.`ExtCate`, ds.`ExtNum`, ds.`Page`, ds.`Size`"
		sqlFormat := `
							select %v from %v di left join %v u on di.Uid=u.Id
							left join %v clt on clt.Did=di.Id
							left join %v d on d.Id=di.Id
							left join %v c on c.Id=di.cid
							left join %v ds on ds.Id=di.DsId
							where %v order by %v limit %v,%v
							`
		sql = fmt.Sprintf(sqlFormat,
			fields,
			models.GetTableDocumentInfo(),
			models.GetTableUser(),
			models.GetTableCollect(),
			models.GetTableDocument(),
			models.GetTableCategory(),
			models.GetTableDocumentStore(),
			fmt.Sprintf("clt.Cid=%v", cid),
			"clt.Id desc",
			(p-1)*listRows, listRows,
		)

		var data []orm.Params
		orm.NewOrm().Raw(sql).Values(&data)
		controller.Data["Lists"] = data
		controller.Data["Page"] = helper.Paginations(6, helper.Interface2Int(params[0]["Cnt"]), listRows, p, fmt.Sprintf("/user/%v/doc/cid/%v", user["Id"], cid), "sort", sort, "style", style)
	} else {
		controller.Data["Lists"], _, _ = models.GetDocList(uid, 0, 0, 0, p, listRows, sort, 1, 0)
		controller.Data["Page"] = helper.Paginations(6, helper.Interface2Int(user["Document"]), listRows, p, fmt.Sprintf("/user/%v/doc", user["Id"]), "sort", sort, "style", style)
	}

	controller.Data["Tab"] = "doc"
	controller.Data["Cid"] = cid
	controller.Data["User"] = user
	controller.Data["PageId"] = "wenku-user"
	controller.Data["IsUser"] = true
	controller.Data["Sort"] = sort
	controller.Data["Style"] = style
	controller.Data["P"] = p
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Ucenter-Doc", "文档列表-会员中心-"+user["Username"].(string), "会员中心,文档列表,"+user["Username"].(string), "文档列表-会员中心-"+user["Username"].(string), controller.Sys.Site)
	controller.Data["Ranks"], _, err = models.NewUser().UserList(1, 8, "i.Document desc", "u.Id,u.Username,u.Avatar,u.Intro,i.Document", "i.Status=1")
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	controller.TplName = "index.html"

}

//金币记录
func (controller *UserController) Coin() {
	uid, _ := controller.GetInt(":uid")
	p, _ := controller.GetInt("p", 1)
	if p < 1 {
		p = 1
	}
	if uid < 1 {
		uid = controller.IsLogin
	}

	if uid <= 0 {
		controller.Redirect("/user/login", 302)
		return
	}

	listRows := 16
	lists, _, _ := models.GetList(models.GetTableCoinLog(), p, listRows, orm.NewCondition().And("Uid", uid), "-Id")
	if p > 1 { // 当页码大于0，则以 JSON 返回数据
		controller.ResponseJson(true, "数据获取成功", lists)
	}

	user, rows, err := models.NewUser().GetById(uid)
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	if rows == 0 {
		controller.Redirect("/", 302)
		return
	}

	controller.Data["Lists"] = lists
	controller.Data["User"] = user
	controller.Data["PageId"] = "wenku-user"
	controller.Data["Tab"] = "coin"
	controller.Data["IsUser"] = true
	controller.Data["Ranks"], _, err = models.NewUser().UserList(1, 8, "i.Document desc", "u.Id,u.Username,u.Avatar,u.Intro,i.Document", "i.Status=1")
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Ucenter-Coin", "财富记录—会员中心-"+user["Username"].(string), "会员中心,财富记录,"+user["Username"].(string), "财富记录—会员中心-"+user["Username"].(string), controller.Sys.Site)
	controller.TplName = "coin.html"

}

// 收藏夹
func (controller *UserController) Collect() {
	controller.Data["Tab"] = "collect"
	action := controller.GetString("action")
	uid, _ := controller.GetInt(":uid")
	p, _ := controller.GetInt("p", 1)
	if p < 1 {
		p = 1
	}
	if uid < 1 {
		uid = controller.IsLogin
	}

	if uid <= 0 {
		controller.Redirect("/user/login", 302)
		return
	}

	listRows := 100
	lists, _, _ := models.GetList(models.GetTableCollectFolder(), p, listRows, orm.NewCondition().And("Uid", uid), "-Id")
	if p > 1 { // 页码大于1，以 JSON 返回数据
		controller.ResponseJson(true, "数据获取成功", lists)
	}

	user, rows, err := models.NewUser().GetById(uid)
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	if rows == 0 {
		controller.Redirect("/", 302)
		return
	}
	controller.Data["Lists"] = lists
	controller.Data["User"] = user
	controller.Data["PageId"] = "wenku-user"
	controller.Data["IsUser"] = true
	controller.Data["Uid"] = uid
	controller.Data["Ranks"], _, err = models.NewUser().UserList(1, 8, "i.Document desc", "u.Id,u.Username,u.Avatar,u.Intro,i.Document", "i.Status=1")
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	controller.TplName = "collect.html"
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Ucenter-Folder", "收藏夹—会员中心-"+user["Username"].(string), "会员中心,收藏夹,"+user["Username"].(string), "收藏夹—会员中心-"+user["Username"].(string), controller.Sys.Site)
	if action == "edit" {
		controller.Data["Edit"] = true
	} else {
		controller.Data["Edit"] = false
	}
}

//用户登录
func (controller *UserController) Login() {

	if controller.IsLogin > 0 {
		controller.Redirect("/user", 302)
		return
	}

	// GET 请求
	if controller.Ctx.Request.Method == "GET" {
		controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Login", "会员登录", "会员登录", "会员登录", controller.Sys.Site)
		controller.Data["IsUser"] = true
		controller.Data["PageId"] = "wenku-reg"
		controller.TplName = "login.html"
		return
	}

	type Post struct {
		Email, Password string
	}

	var post struct {
		Email, Password string
	}

	controller.ParseForm(&post)
	valid := validation.Validation{}
	res := valid.Email(post.Email, "Email")
	if !res.Ok {
		controller.ResponseJson(false, "登录失败，邮箱格式不正确")
	}

	ModelUser := models.NewUser()
	users, rows, err := ModelUser.UserList(1, 1, "", "", "u.`email`=? and u.`password`=?", post.Email, helper.MD5Crypt(post.Password))
	if rows == 0 || err != nil {
		if err != nil {
			helper.Logger.Error(err.Error())
		}
		controller.ResponseJson(false, "登录失败，邮箱或密码不正确")
	}

	user := users[0]
	controller.IsLogin = helper.Interface2Int(user["Id"])

	if controller.IsLogin > 0 {
		//查询用户有没有被封禁
		if info := ModelUser.UserInfo(controller.IsLogin); info.Status == false { //被封禁了
			controller.ResponseJson(false, "登录失败，您的账号已被管理员禁用")
		}
		controller.BaseController.SetCookieLogin(controller.IsLogin)
		controller.ResponseJson(true, "登录成功")
	}
	controller.ResponseJson(false, "登录失败，未知错误！")
}

//用户退出登录
func (controller *UserController) Logout() {
	controller.ResetCookie()
	if v, ok := controller.Ctx.Request.Header["X-Requested-With"]; ok && v[0] == "XMLHttpRequest" {
		controller.ResponseJson(true, "退出登录成功")
	}
	controller.Redirect("/", 302)
}

//会员注册[GET/POST]
func (controller *UserController) Reg() {
	if controller.IsLogin > 0 {
		controller.Redirect("/user", 302)
		return
	}

	if controller.Ctx.Request.Method == "GET" {
		controller.Data["IsUser"] = true
		controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Login", "会员注册", "会员注册", "会员注册", controller.Sys.Site)
		controller.Data["PageId"] = "wenku-reg"
		if controller.Sys.IsCloseReg {
			controller.TplName = "regclose.html"
		} else {
			controller.TplName = "reg.html"
		}
		return
	}

	if controller.Sys.IsCloseReg {
		controller.ResponseJson(false, "注册失败，站点已关闭注册功能")
	}

	//先验证邮箱验证码是否正确
	email := controller.GetString("email")
	code := controller.GetString("code")
	if controller.Sys.CheckRegEmail {
		sessEmail := fmt.Sprintf("%v", controller.GetSession("RegMail"))
		sessCode := fmt.Sprintf("%v", controller.GetSession("RegCode"))
		if sessEmail != email || sessCode != code {
			controller.ResponseJson(false, "邮箱验证码不正确，请重新输入或重新获取")
		}
	}

	// 注册
	err, uid := models.NewUser().Reg(
		email,
		controller.GetString("username"),
		controller.GetString("password"),
		controller.GetString("repassword"),
		controller.GetString("intro"),
	)
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}

	models.Regulate(models.GetTableSys(), "CntUser", 1, "Id=1") //站点用户数量增加
	controller.IsLogin = uid
	controller.SetCookieLogin(uid)
	controller.ResponseJson(true, "会员注册成功")
}

// 发送邮件
func (controller *UserController) SendMail() {
	beego.Info(controller.Ctx.GetCookie("_xsrf"))

	//发送邮件的类型：注册(reg)和找回密码(findpwd)
	t := controller.GetString("type")
	if t != "reg" && t != "findpwd" {
		controller.ResponseJson(false, "邮件发送类型不正确")
	}

	valid := validation.Validation{}
	email := controller.GetString("email")
	res := valid.Email(email, "mail")
	if res.Error != nil || !res.Ok {
		controller.ResponseJson(false, "邮箱格式不正确")
	}

	//检测邮箱是否已被注册
	ModelUser := models.NewUser()
	user := ModelUser.GetUserField(orm.NewCondition().And("email", email))

	//注册邮件
	if t == "reg" {
		if user.Id > 0 {
			controller.ResponseJson(false, "该邮箱已经被注册会员")
		}

		code := helper.RandStr(6, 0)
		err := models.NewEmail().SendMail(email, fmt.Sprintf("%v会员注册验证码", controller.Sys.Site), strings.Replace(controller.Sys.TplEmailReg, "{code}", code, -1))
		if err != nil {
			helper.Logger.Error("邮件发送失败：%v", err.Error())
			controller.ResponseJson(false, "邮件发送失败，请联系管理员检查邮箱配置是否正确")
		}

		controller.SetSession("RegMail", email)
		controller.SetSession("RegCode", code)
		controller.ResponseJson(true, "邮件发送成功，请打开邮箱查看验证码")
	}

	// 找回密码
	if user.Id == 0 {
		controller.ResponseJson(false, "邮箱不存在")
	}

	code := helper.RandStr(6, 0)
	err := models.NewEmail().SendMail(email, fmt.Sprintf("%v找回密码验证码", controller.Sys.Site), strings.Replace(controller.Sys.TplEmailFindPwd, "{code}", code, -1))
	if err != nil {
		helper.Logger.Error("邮件发送失败：%v", err.Error())
		controller.ResponseJson(false, "邮件发送失败，请联系管理员检查邮箱配置是否正确")
	}

	controller.SetSession("FindPwdMail", email)
	controller.SetSession("FindPwdCode", code)
	controller.ResponseJson(true, "邮件发送成功，请打开邮箱查看验证码")
}

//会员签到，增加金币
func (controller *UserController) Sign() {

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "签到失败，请先登录")
	}

	var data = models.Sign{
		Uid:  controller.IsLogin,
		Date: time.Now().Format("20060102"),
	}
	_, err := orm.NewOrm().Insert(&data)
	if err != nil {
		controller.ResponseJson(false, "签到失败，您今天已签到")
	}

	if err = models.Regulate(models.GetTableUserInfo(), "Coin", controller.Sys.Sign, fmt.Sprintf("Id=%v", controller.IsLogin)); err == nil {
		log := models.CoinLog{
			Uid:  controller.IsLogin,
			Coin: controller.Sys.Sign,
			Log:  fmt.Sprintf("签到成功，获得 %v 个金币", controller.Sys.Sign),
		}
		models.NewCoinLog().LogRecord(log)
	}
	controller.ResponseJson(true, fmt.Sprintf("恭喜您，今日签到成功，领取了 %v 个金币", controller.Sys.Sign))
}

// 检测用户是否已登录
func (controller *UserController) CheckLogin() {
	if controller.BaseController.IsLogin > 0 {
		controller.ResponseJson(true, "已登录")
	}
	controller.ResponseJson(false, "您当前处于未登录状态，请先登录")
}

// 创建收藏夹
func (controller *UserController) CreateCollectFolder() {

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "您当前未登录，请先登录")
	}

	cover := ""
	timestamp := int(time.Now().Unix())

	//文件在文档库中未存在，则接收文件并做处理
	f, fh, err := controller.GetFile("Cover")
	if err == nil {
		defer f.Close()
		slice := strings.Split(fh.Filename, ".")
		ext := slice[len(slice)-1]
		dir := fmt.Sprintf("./uploads/%v/%v/", time.Now().Format("2006-01-02"), controller.IsLogin)
		os.MkdirAll(dir, 0777)
		file := helper.MD5Crypt(fmt.Sprintf("%v-%v-%v", timestamp, controller.IsLogin, fh.Filename)) + "." + ext

		tmpFile := dir + file
		err = controller.SaveToFile("Cover", tmpFile)
		if err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, "封面保存失败")
		}
		defer os.RemoveAll(tmpFile)

		if err = helper.CropImage(tmpFile, helper.CoverWidth, helper.CoverHeight); err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, "封面裁剪失败")
		}

		//将图片移动到OSS
		var cs *models.CloudStore
		if cs, err = models.NewCloudStore(false); err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, "连接云存储失败")
		}
		if err = cs.Upload(tmpFile, file); err != nil {
			helper.Logger.Error(err.Error())
		} else {
			cover = file
		}
	}

	// 收藏夹
	folder := models.CollectFolder{
		Uid:         controller.IsLogin,
		Title:       controller.GetString("Title"),
		Description: controller.GetString("Description"),
		TimeCreate:  int(time.Now().Unix()),
		Cnt:         0,
		Cover:       cover,
	}

	// 收藏夹 Id 大于0，则表示编辑收藏夹
	folder.Id, _ = controller.GetInt("Id")

	if folder.Id > 0 { // 编辑收藏夹
		cols := []string{"Title", "Description"}
		if len(cover) > 0 {
			cols = append(cols, "Cover")
		}
		if _, err = orm.NewOrm().Update(&folder, cols...); err == nil {
			controller.ResponseJson(true, "收藏夹编辑成功")
		}
	} else { // 创建收藏夹
		if _, err = orm.NewOrm().Insert(&folder); err == nil { //收藏夹数量+1
			models.Regulate(models.GetTableUserInfo(), "Collect", 1, "Id=?", controller.IsLogin)
			controller.ResponseJson(true, "收藏夹创建成功")
		}
	}

	if err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "操作失败，请重试")
	}

	controller.ResponseJson(true, "操作成功")
}

// 找回密码
func (controller *UserController) FindPwd() {
	if controller.IsLogin > 0 {
		controller.Redirect("/user", 302)
		return
	}

	if controller.Ctx.Request.Method == "GET" {
		controller.Data["Seo"] = models.NewSeo().GetByPage("PC-Findpwd", "找回密码", "找回密码", "找回密码", controller.Sys.Site)
		controller.Data["IsUser"] = true
		controller.Data["PageId"] = "wenku-reg"
		controller.TplName = "findpwd.html"
		return
	}

	rules := map[string][]string{
		"username":   {"required", "mincount:2", "maxcount:16"},
		"email":      {"required", "email"},
		"code":       {"required", "len:6"},
		"password":   {"required", "mincount:6"},
		"repassword": {"required", "mincount:6"},
	}

	params, errs := helper.Valid(controller.Ctx.Request.Form, rules)
	if len(errs) > 0 {
		if _, ok := errs["username"]; ok {
			controller.ResponseJson(false, "用户名限2-16个字符")
		}
		if _, ok := errs["email"]; ok {
			controller.ResponseJson(false, "邮箱格式不正确")
		}
		if _, ok := errs["code"]; ok {
			controller.ResponseJson(false, "请输入6位验证码")
		}
		if _, ok := errs["password"]; ok {
			controller.ResponseJson(false, "密码长度，至少6个字符")
		}
		if _, ok := errs["repassword"]; ok {
			controller.ResponseJson(false, "密码长度，至少6个字符")
		}
	}

	//校验验证码和邮箱是否匹配
	if fmt.Sprintf("%v", controller.GetSession("FindPwdMail")) != params["email"].(string) || fmt.Sprintf("%v", controller.GetSession("FindPwdCode")) != params["code"].(string) {
		controller.ResponseJson(false, "验证码不正确，修改密码失败")
	}
	pwd := helper.MD5Crypt(params["password"].(string))
	repwd := helper.MD5Crypt(params["repassword"].(string))
	if pwd != repwd {
		controller.ResponseJson(false, "确认密码和密码不一致")
	}

	user := models.NewUser().GetUserField(orm.NewCondition().And("Email", params["email"]))
	if user.Id == 0 || user.Username != params["username"].(string) {
		controller.ResponseJson(false, "重置密码失败，用户名与邮箱不匹配")
	}

	_, err := models.UpdateByIds("user", "Password", pwd, user.Id)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "重置密码失败，请刷新页面重试")
	}
	controller.DelSession("FindPwdMail")
	controller.DelSession("FindPwdCode")
	controller.ResponseJson(true, "重置密码成功，请重新登录")
}

//删除文档
func (controller *UserController) DocDel() {

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "请先登录")
	}

	docid, _ := controller.GetInt(":doc")
	if docid == 0 {
		controller.ResponseJson(false, "删除失败，文档不存在")
	}

	err := models.NewDocumentRecycle().RemoveToRecycle(controller.IsLogin, true, docid)
	if err != nil {
		helper.Logger.Error("删除失败：%v", err.Error())
		controller.ResponseJson(false, "删除失败，文档不存在")
	}

	controller.ResponseJson(true, "删除成功")
}

//文档编辑
func (controller *UserController) DocEdit() {

	if controller.IsLogin == 0 {
		controller.Redirect("/user", 302)
	}

	docId, _ := controller.GetInt(":doc")
	if docId == 0 {
		controller.Redirect("/user", 302)
	}

	info := models.DocumentInfo{Id: docId}
	err := orm.NewOrm().Read(&info)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.Redirect("/user", 302)
	}

	if info.Uid != controller.IsLogin { // 文档所属用户id与登录的用户id不一致
		controller.Redirect("/user", 302)
	}

	doc := models.Document{Id: docId}

	// POST
	if controller.Ctx.Request.Method == "POST" {
		ruels := map[string][]string{
			"Title":  {"required", "unempty"},
			"Chanel": {"required", "gt:0", "int"},
			"Pid":    {"required", "gt:0", "int"},
			"Cid":    {"required", "gt:0", "int"},
			"Tags":   {"required"},
			"Intro":  {"required"},
			"Price":  {"required", "int"},
		}
		params, errs := helper.Valid(controller.Ctx.Request.Form, ruels)
		if len(errs) > 0 {
			controller.ResponseJson(false, "参数错误")
		}
		doc.Title = params["Title"].(string)
		doc.Keywords = params["Tags"].(string)
		doc.Description = params["Intro"].(string)
		info.Pid = params["Pid"].(int)
		info.Cid = params["Cid"].(int)
		info.ChanelId = params["Chanel"].(int)
		info.Price = params["Price"].(int)
		info.TimeUpdate = int(time.Now().Unix())
		orm.NewOrm().Update(&doc, "Title", "Keywords", "Description")
		orm.NewOrm().Update(&info, "Pid", "Cid", "ChanelId", "Price")
		//原分类-1
		models.Regulate(models.GetTableCategory(), "Cnt", -1, fmt.Sprintf("Id in(%v,%v,%v)", info.ChanelId, info.Cid, info.Pid))
		//新分类+1
		models.Regulate(models.GetTableCategory(), "Cnt", 1, fmt.Sprintf("Id in(%v,%v,%v)", params["Chanel"], params["Cid"], params["Pid"]))
		controller.ResponseJson(true, "文档编辑成功")
	}

	// GET
	err = orm.NewOrm().Read(&doc)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.Redirect("/user", 302)
	}

	cond := orm.NewCondition().And("status", 1)
	data, _, _ := models.GetList(models.GetTableCategory(), 1, 2000, cond, "sort")
	controller.Data["User"], _, _ = models.NewUser().GetById(controller.IsLogin)
	controller.Data["Ranks"], _, err = models.NewUser().UserList(1, 8, "i.Document desc", "u.Id,u.Username,u.Avatar,u.Intro,i.Document", "i.Status=1")
	controller.Data["IsUser"] = true
	controller.Data["Cates"], _ = conv.InterfaceToJson(data)
	controller.Data["json"] = data
	controller.Data["PageId"] = "wenku-user"
	controller.Data["Info"] = info
	controller.Data["Doc"] = doc
	controller.Data["Tab"] = "doc"
	controller.TplName = "edit.html"
}

//删除收藏(针对收藏夹)
func (controller *UserController) CollectFolderDel() {
	cid, _ := controller.GetInt(":cid")
	if cid > 0 && controller.IsLogin > 0 {
		err := models.NewCollect().DelFolder(cid, controller.IsLogin)
		if err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, err.Error())
		}
		controller.ResponseJson(true, "收藏夹删除成功")
	}
	controller.ResponseJson(false, "删除失败，参数错误")
}

//取消收藏(针对文档)
func (controller *UserController) CollectCancel() {
	cid, _ := controller.GetInt(":cid")
	did, _ := controller.GetInt(":did")
	if err := models.NewCollect().Cancel(did, cid, controller.IsLogin); err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "移除收藏失败，可能您为收藏该文档")
	}
	controller.ResponseJson(true, "移除收藏成功")
}

//更换头像
func (controller *UserController) Avatar() {

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "请先登录")
	}

	dir := fmt.Sprintf("./uploads/%v/%v", time.Now().Format("2006-01-02"), controller.IsLogin)
	os.MkdirAll(dir, 0777)
	f, fh, err := controller.GetFile("Avatar")
	if err != nil {
		helper.Logger.Error("用户(%v)更新头像失败：%v", controller.IsLogin, err.Error())
		controller.ResponseJson(false, "头像文件上传失败")
	}
	defer f.Close()

	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(fh.Filename), "."))
	if !(ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif") {
		controller.ResponseJson(false, "头像图片格式只支持jpg、jpeg、png和gif")
	}

	tmpFile := dir + "/" + helper.MD5Crypt(fmt.Sprintf("%v-%v-%v", fh.Filename, controller.IsLogin, time.Now().Unix())) + "." + ext
	saveFile := helper.MD5Crypt(tmpFile) + "." + ext
	err = controller.SaveToFile("Avatar", tmpFile)
	if err != nil {
		helper.Logger.Error("用户(%v)头像保存失败：%v", controller.IsLogin, err.Error())
		controller.ResponseJson(false, "头像文件保存失败")
	}

	//头像裁剪
	if err = helper.CropImage(tmpFile, helper.AvatarWidth, helper.AvatarHeight); err != nil {
		helper.Logger.Error("图片裁剪失败：%v", err.Error())
	}

	var cs *models.CloudStore
	if cs, err = models.NewCloudStore(false); err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "内部服务错误：云存储连接失败")
	}

	err = cs.Upload(tmpFile, saveFile)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "头像文件保存失败")
	}
	os.RemoveAll(tmpFile)

	//查询数据库用户数据
	var user = models.User{Id: controller.IsLogin}
	orm.NewOrm().Read(&user)
	oldAvatar := user.Avatar
	user.Avatar = saveFile
	rows, err := orm.NewOrm().Update(&user, "Avatar")
	if rows > 0 && err == nil {
		controller.ResponseJson(true, "头像更新成功")
		go cs.Delete(oldAvatar)
	}
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	controller.ResponseJson(false, "头像更新失败")
}

//编辑个人信息
func (controller *UserController) Edit() {
	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "请先登录")
	}
	changepwd := false
	cols := []string{"Intro"}
	rules := map[string][]string{
		"OldPassword": {"required"},
		"NewPassword": {"required"},
		"RePassword":  {"required"},
		"Intro":       {"required"},
	}
	params, errs := helper.Valid(controller.Ctx.Request.Form, rules)
	if len(errs) > 0 {
		controller.ResponseJson(false, "参数不正确")
	}
	var user = models.User{Id: controller.IsLogin}
	orm.NewOrm().Read(&user)
	if len(params["OldPassword"].(string)) > 0 || len(params["NewPassword"].(string)) > 0 || len(params["RePassword"].(string)) > 0 {
		if len(params["NewPassword"].(string)) < 6 || len(params["RePassword"].(string)) < 6 {
			controller.ResponseJson(false, "密码长度必须至少6个字符")
		}
		opwd := helper.MD5Crypt(params["OldPassword"].(string))
		npwd := helper.MD5Crypt(params["NewPassword"].(string))
		rpwd := helper.MD5Crypt(params["RePassword"].(string))
		if user.Password != opwd {
			controller.ResponseJson(false, "原密码不正确")
		}
		if npwd != rpwd {
			controller.ResponseJson(false, "确认密码和新密码必须一致")
		}
		if opwd == npwd {
			controller.ResponseJson(false, "确认密码不能与原密码相同")
		}
		user.Password = rpwd
		cols = append(cols, "Password")
		changepwd = true
	}
	user.Intro = params["Intro"].(string)
	affected, err := orm.NewOrm().Update(&user, cols...)
	if err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "设置失败，请刷新页面重试")
	}
	if affected == 0 {
		controller.ResponseJson(true, "设置失败，可能您未对内容做更改")
	}
	if changepwd {
		controller.ResetCookie()
		controller.ResponseJson(true, "设置成功，您设置了新密码，请重新登录")
	}
	controller.ResponseJson(true, "设置成功")
}
