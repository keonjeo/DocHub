package HomeControllers

import (
	"fmt"

	"dochub/helper"
	"dochub/models"
	"github.com/astaxie/beego/orm"
)

type CollectController struct {
	BaseController
}

//收藏文档
func (controller *CollectController) Get() {
	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "您当前未登录，请先登录")
	}

	cid, _ := controller.GetInt("Cid")
	did, _ := controller.GetInt("Did")

	if cid == 0 || did == 0 {
		controller.ResponseJson(false, "收藏失败：参数不正确")
	}

	collect := models.Collect{Did: did, Cid: cid}
	rows, err := orm.NewOrm().Insert(&collect)
	if err != nil {
		helper.Logger.Error("SQL执行失败：%v", err.Error())
	}

	if err != nil || rows == 0 {
		controller.ResponseJson(false, "收藏失败：您已收藏过该文档")
	}

	//文档被收藏的数量+1
	models.Regulate(models.GetTableDocumentInfo(), "Ccnt", 1, fmt.Sprintf("`Id`=%v", did))

	//收藏夹的文档+1
	models.Regulate(models.GetTableCollectFolder(), "Cnt", 1, fmt.Sprintf("`Id`=%v", cid))

	controller.ResponseJson(true, "恭喜您，文档收藏成功。")
}

//收藏夹列表
func (controller *CollectController) FolderList() {
	uid, _ := controller.GetInt("uid")
	if uid < 1 {
		uid = controller.IsLogin
	}

	if uid == 0 {
		controller.ResponseJson(false, "获取收藏夹失败：请先登录")
	}

	lists, rows, err := models.GetList(models.GetTableCollectFolder(), 1, 100, orm.NewCondition().And("Uid", uid), "-Id")
	if err != nil {
		helper.Logger.Error(err.Error())
	}
	if rows > 0 && err == nil {
		controller.ResponseJson(true, "收藏夹获取获取成功", lists)
	}
	controller.ResponseJson(false, "暂时没有收藏夹，请先在会员中心创建收藏夹")
}

//取消收藏文档
func (controller *CollectController) CollectCancel() {
	if controller.IsLogin == 0 {
		controller.ResponseJson(false, "您当前未登录，请先登录")
	}

	cid, _ := controller.GetInt("Cid")
	did, _ := controller.GetInt("Did")
	if cid == 0 || did == 0 {
		controller.ResponseJson(false, "收藏失败：参数不正确")
	}

	if err := models.NewCollect().Cancel(did, cid, controller.IsLogin); err != nil {
		helper.Logger.Error(err.Error())
		controller.ResponseJson(false, "移除收藏文档失败，可能您为收藏该文档")
	}

	//文档被收藏的数量-1
	models.Regulate(models.GetTableDocumentInfo(), "Ccnt", -1, "`Id`=?", did)

	//收藏夹的文档-1
	models.Regulate(models.GetTableCollectFolder(), "Cnt", -1, "`Id`=?", cid)

	controller.ResponseJson(true, "恭喜您，删除收藏文档成功")
}
