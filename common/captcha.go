package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"net/http"
)

var store = base64Captcha.DefaultMemStore
// 获取验证码
func GetCaptcha(c *gin.Context) {
	//配置验证码
	driverString := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,     //噪点数
		ShowLineOptions: 2 | 4, //干扰线
		Length:          4,
		Source:          "123456789abcdefghijklmnopqrstuvwxwz",
		BgColor: &color.RGBA{
			R: 100,
			G: 100,
			B: 100,
			A: 100,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	var driver base64Captcha.Driver = driverString.ConvertFonts()
	//生成验证码
	cap := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cap.Generate()
	body := map[string]interface{}{"code": 1, "data": b64s, "captchaId": id, "msg": "success"}
	if err != nil {
		body = map[string]interface{}{"code": 0, "msg": err.Error()}
	}
	c.JSON(http.StatusOK, body)
}

// 自定义验证码的验证器
func VerifyCaptcha(f validator.FieldLevel) bool {
	captcha := f.Parent().FieldByName("Captcha").Interface().(string)
	captcha_id := f.Parent().FieldByName("CaptchaId").Interface().(string)

	if store.Verify(captcha_id, captcha, true) {
		return true
	}
	return false
}