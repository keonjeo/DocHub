package HomeControllers

import (
	"fmt"

	"strings"

	"time"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ViewController struct {
	BaseController
}

func (controller *ViewController) Get() {
	id, _ := controller.GetInt(":id")
	if id < 1 {
		controller.Redirect("/", 302)
		return
	}

	doc, err := models.NewDocument().GetById(id)

	// 文档不存在、查询错误、被删除，报 404
	if err != nil || doc.Id <= 0 || doc.Status < models.DocStatusConverting {
		controller.Abort("404")
	}

	var cates []models.Category
	cates, _ = models.NewCategory().GetCategoriesById(doc.Cid, doc.ChanelId, doc.Pid)
	breadcrumb := make(map[string]models.Category)
	for _, cate := range cates {
		switch cate.Id {
		case doc.ChanelId:
			breadcrumb["Chanel"] = cate
		case doc.Pid:
			breadcrumb["Parent"] = cate
		case doc.Cid:
			TimeStart := int(time.Now().Unix()) - controller.Sys.TimeExpireHotspot
			controller.Data["Hots"], _, _ = models.NewDocument().SimpleList(fmt.Sprintf("di.Cid=%v and di.TimeCreate>%v", doc.Cid, TimeStart), 10, "Dcnt")
			controller.Data["Latest"], _, _ = models.NewDocument().SimpleList(fmt.Sprintf("di.Cid=%v", doc.Cid), 10, "Id")
			breadcrumb["Child"] = cate
		}
	}
	controller.Data["Breadcrumb"] = breadcrumb

	models.Regulate(models.GetTableDocumentInfo(), "Vcnt", 1, "`Id`=?", id)

	pageShow := 5
	if doc.Page > pageShow {
		controller.Data["PreviewPages"] = make([]string, pageShow)
	} else {
		controller.Data["PreviewPages"] = make([]string, doc.Page)
	}
	controller.Data["PageShow"] = pageShow

	controller.Xsrf()
	if controller.Data["Comments"], _, err = models.NewDocumentComment().GetCommentList(id, 1, 10); err != nil {
		helper.Logger.Error(err.Error())
	}

	content := models.NewDocText().GetDescByMd5(doc.Md5, 5000)
	seoTitle := fmt.Sprintf("%v - %v · %v · %v ", doc.Title, breadcrumb["Chanel"].Title, breadcrumb["Parent"].Title, breadcrumb["Child"].Title)
	seoKeywords := fmt.Sprintf("%v,%v,%v,", breadcrumb["Chanel"].Title, breadcrumb["Parent"].Title, breadcrumb["Child"].Title) + doc.Keywords
	seoDesc := beego.Substr(doc.Description+content, 0, 255)
	controller.Data["Seo"] = models.NewSeo().GetByPage("PC-View", seoTitle, seoKeywords, seoDesc, controller.Sys.Site)
	controller.Data["Content"] = content
	controller.Data["Reasons"] = models.NewSys().GetReportReasons()
	controller.Data["IsViewer"] = true
	controller.Data["PageId"] = "wenku-content"
	controller.Data["Doc"] = doc

	doc.Ext = strings.TrimLeft(doc.Ext, ".")
	if doc.Page == 0 { //不能预览的文档
		controller.Data["OnlyCover"] = true
		controller.TplName = "disabled.html"
	} else {
		controller.Data["ViewAll"] = doc.PreviewPage == 0 || doc.PreviewPage >= doc.Page
		controller.TplName = "svg.html"
	}

}

//文档下载
func (controller *ViewController) Download() {
	id, _ := controller.GetInt(":id")
	if id <= 0 {
		controller.ResponseJson(false, "文档id不正确")
	}

	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "请先登录")
	}

	link, err := models.NewUser().CanDownloadFile(controller.IsLogin, id)
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}
	controller.ResponseJson(true, "下载链接获取成功", map[string]interface{}{"url": link})
}

//是否可以免费下载
func (controller *ViewController) DownFree() {
	if controller.IsLogin > 0 {
		did, _ := controller.GetInt("id")
		if free := models.NewFreeDown().IsFreeDown(controller.IsLogin, did); free {
			controller.ResponseJson(true, fmt.Sprintf("您上次下载过当前文档，且仍在免费下载有效期(%v天)内，本次下载免费", controller.Sys.FreeDay))
		}
	}
	controller.ResponseJson(false, "不能免费下载，不在免费下载期限内")
}

//文档评论
func (controller *ViewController) Comment() {
	id, _ := controller.GetInt(":id")
	score, _ := controller.GetInt("Score")
	answer := controller.GetString("Answer")
	if answer != controller.Sys.Answer {
		controller.ResponseJson(false, "请输入正确的答案")
	}
	if id > 0 {
		if controller.IsLogin > 0 {
			if score < 1 || score > 5 {
				controller.ResponseJson(false, "请给文档评分")
			} else {
				comment := models.DocumentComment{
					Uid:        controller.IsLogin,
					Did:        id,
					Content:    controller.GetString("Comment"),
					TimeCreate: int(time.Now().Unix()),
					Status:     true,
					Score:      score * 10000,
				}
				cnt := strings.Count(comment.Content, "") - 1
				if cnt > 255 || cnt < 8 {
					controller.ResponseJson(false, "评论内容限8-255个字符")
				} else {
					_, err := orm.NewOrm().Insert(&comment)
					if err != nil {
						controller.ResponseJson(false, "发表评论失败：每人仅限给每个文档点评一次")
					} else {
						//文档评论人数增加
						sql := fmt.Sprintf("UPDATE `%v` SET `Score`=(`Score`*`ScorePeople`+%v)/(`ScorePeople`+1),`ScorePeople`=`ScorePeople`+1 WHERE Id=%v", models.GetTableDocumentInfo(), comment.Score, comment.Did)
						_, err := orm.NewOrm().Raw(sql).Exec()
						if err != nil {
							helper.Logger.Error(err.Error())
						}
						controller.ResponseJson(true, "恭喜您，评论发表成功")
					}
				}
			}
		} else {
			controller.ResponseJson(false, "评论失败，您当前处于未登录状态，请先登录")
		}
	} else {
		controller.ResponseJson(false, "评论失败，参数不正确")
	}
}

//获取评论列表
func (controller *ViewController) GetComment() {
	p, _ := controller.GetInt("p", 1)
	did, _ := controller.GetInt("did")
	if p > 0 && did > 0 {
		if rows, _, err := models.NewDocumentComment().GetCommentList(did, p, 10); err != nil {
			helper.Logger.Error(err.Error())
			controller.ResponseJson(false, "评论列表获取失败")
		} else {
			controller.ResponseJson(true, "评论列表获取成功", rows)
		}
	}
}
