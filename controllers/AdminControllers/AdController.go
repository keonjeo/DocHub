package AdminControllers

type AdController struct {
	BaseController
}

func (controller *AdController) Get() {
	controller.Data["IsAd"] = true
	controller.TplName = "index.html"
}
