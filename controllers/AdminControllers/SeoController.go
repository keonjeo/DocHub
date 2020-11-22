package AdminControllers

import "dochub/models"

type SeoController struct {
	BaseController
}

func (controller *SeoController) Get() {
	controller.Data["Data"], _, _ = models.GetList(models.GetTableSeo(), 1, 50, nil, "-IsMobile")
	controller.Data["IsSeo"] = true
	controller.TplName = "index.html"
}

func (controller *SeoController) UpdateSitemap() {
	go models.NewSeo().BuildSitemap()
	controller.ResponseJson(true, "Sitemap更新已提交后台执行，请耐心等待")
}
