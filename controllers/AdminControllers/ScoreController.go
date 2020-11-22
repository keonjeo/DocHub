package AdminControllers

import "dochub/models"

type ScoreController struct {
	BaseController
}

//积分管理
func (controller *ScoreController) Get() {
	var log models.CoinLog
	log.Uid, _ = controller.GetInt("uid")
	log.Coin, _ = controller.GetInt("score")
	log.Log = controller.GetString("log")
	err := models.Regulate(models.GetTableUserInfo(), "Coin", log.Coin, "Id=?", log.Uid)
	if err == nil {
		err = models.NewCoinLog().LogRecord(log)
	}
	if err != nil {
		controller.ResponseJson(false, err.Error())
	}
	controller.ResponseJson(true, "积分变更成功")
}
