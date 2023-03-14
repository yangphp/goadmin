package modes

type SysSetting struct {
	Id int `form:"id"`
	SysTitle string `form:"sys_title"`
	SysCompany string `form:"sys_company"`
	SysQq string `form:"sys_qq"`
	SysPhone string `form:"sys_phone"`
	SysDev string `form:"sys_dev"`
	SysBeian string `form:"sys_beian"`
}

//获取设置
func GetSetting() SysSetting {
	setting := SysSetting{}
	DB.Model(SysSetting{}).First(&setting)

	return setting
}

//保存设置
func SaveSetting(setting *SysSetting) (bool,string) {

	res := DB.Model(SysSetting{}).Where("id = ?",setting.Id).Save(setting)
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"保存成功"
}


