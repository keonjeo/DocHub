package HomeControllers

import (
	"net/http"
	"path/filepath"
	"strings"

	"dochub/helper"
	"github.com/astaxie/beego"
)

type StaticController struct {
	beego.Controller
}

// 将除了static之外的静态资源导向到虚拟根目录
func (controller *StaticController) Static() {
	splat := strings.TrimPrefix(controller.GetString(":splat"), "../")
	if strings.HasPrefix(splat, ".well-known") {
		http.ServeFile(controller.Ctx.ResponseWriter, controller.Ctx.Request, splat)
		return
	}
	path := filepath.Join(helper.RootPath, splat)
	http.ServeFile(controller.Ctx.ResponseWriter, controller.Ctx.Request, path)
}
