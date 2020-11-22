package main

import (
	"fmt"
	"log"

	"dochub/controllers/HomeControllers"
	"dochub/helper"
	"dochub/models"
	_ "dochub/routers"

	"github.com/astaxie/beego"
)

//初始化函数
func init() {

	fmt.Println("")
	fmt.Println("Powered By dochub")
	fmt.Println("Version:", helper.VERSION)
	fmt.Println("")

	//sitemap静态目录
	beego.SetStaticPath("/sitemap", "sitemap")

	//初始化日志
	helper.InitLogs()

	//初始化分词器
	go func() {
		helper.Segmenter.LoadDictionary("./dictionary/dictionary.txt")
		beego.Info("==程序启动完毕==")
	}()

	var err error
	err = beego.AddFuncMap("TimestampFormat", helper.TimestampFormat)
	if err != nil {
		log.Fatalf("fail to helper.TimestampFormat")
	}
	err = beego.AddFuncMap("Interface2Int", helper.Interface2Int)
	if err != nil {
		log.Fatalf("fail to helper.Interface2Int")
	}
	err = beego.AddFuncMap("Interface2String", helper.Interface2String)
	if err != nil {
		log.Fatalf("fail to helper.Interface2String")
	}
	err = beego.AddFuncMap("Default", helper.Default)
	if err != nil {
		log.Fatalf("fail to helper.Default")
	}
	err = beego.AddFuncMap("FormatByte", helper.FormatByte)
	if err != nil {
		log.Fatalf("fail to helper.FormatByte")
	}
	err = beego.AddFuncMap("CalcInt", helper.CalcInt)
	if err != nil {
		log.Fatalf("fail to helper.CalcInt")
	}
	err = beego.AddFuncMap("StarVal", helper.StarVal)
	if err != nil {
		log.Fatalf("fail to helper.StarVal")
	}
	err = beego.AddFuncMap("Equal", helper.Equal)
	if err != nil {
		log.Fatalf("fail to helper.Equal")
	}
	err = beego.AddFuncMap("SimpleList", models.NewDocument().TplSimpleList) //简易的文档列表
	if err != nil {
		log.Fatalf("fail to models.NewDocument().TplSimpleList")
	}
	err = beego.AddFuncMap("HandlePageNum", helper.HandlePageNum) //处理文档页码为0的显示问题
	if err != nil {
		log.Fatalf("fail to helper.HandlePageNum")
	}
	err = beego.AddFuncMap("DoesCollect", models.DoesCollect) //判断用户是否已收藏了该文档
	if err != nil {
		log.Fatalf("fail to models.DoesCollect")
	}
	err = beego.AddFuncMap("DoesSign", models.NewSign().DoesSign) //用户今日是否已签到
	if err != nil {
		log.Fatalf("fail to models.NewSign().DoesSign")
	}
	err = beego.AddFuncMap("Friends", models.NewFriend().Friends) //友情链接
	if err != nil {
		log.Fatalf("fail to models.NewFriend().Friends")
	}
	err = beego.AddFuncMap("CategoryName", models.NewCategory().GetTitleById) //根据分类id获取分类名称
	if err != nil {
		log.Fatalf("fail to models.NewCategory().GetTitleById")
	}
	err = beego.AddFuncMap("IsIllegal", models.NewDocument().IsIllegal) //根据md5判断文档是否是非法文档
	if err != nil {
		log.Fatalf("fail to models.NewDocument().IsIllegal")
	}
	err = beego.AddFuncMap("IsRemark", models.NewDocumentRemark().IsRemark) //根据文档是否存在备注
	if err != nil {
		log.Fatalf("fail to  models.NewDocumentRemark().IsRemark")
	}
	err = beego.AddFuncMap("BuildURL", helper.BuildURL) //创建URL
	if err != nil {
		log.Fatalf("fail to helper.BuildURL")
	}
	err = beego.AddFuncMap("HeightLight", helper.HeightLight) //高亮
	if err != nil {
		log.Fatalf("fail to helper.HeightLight")
	}
	err = beego.AddFuncMap("ReportReason", models.NewSys().GetReportReason) //举报原因
	if err != nil {
		log.Fatalf("fail to models.NewSys().GetReportReason")
	}
	err = beego.AddFuncMap("GetDescByMd5", models.NewDocText().GetDescByMd5)
	if err != nil {
		log.Fatalf("fail to models.NewDocText().GetDescByMd5")
	}
	err = beego.AddFuncMap("GetDescByDsId", models.NewDocText().GetDescByDsId)
	if err != nil {
		log.Fatalf("fail to models.NewDocText().GetDescByDsId")
	}
	err = beego.AddFuncMap("GetDescByDid", models.NewDocText().GetDescByDid)
	if err != nil {
		log.Fatalf("fail tomodels.NewDocText().GetDescByDid")
	}
	err = beego.AddFuncMap("DefaultImage", models.GetImageFromCloudStore) //获取默认图片
	if err != nil {
		log.Fatalf("fail to models.GetImageFromCloudStore")
	}
}

func main() {
	//定义错误和异常处理控制器
	beego.ErrorController(&HomeControllers.ErrorsController{})
	beego.Run()
}
