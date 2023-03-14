package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"goadmin/admins/controllers"
	"goadmin/admins/middleware"
	"goadmin/common"
	"goadmin/modes"
	"html/template"
)
//该文件初始化，会自动调用
func initAdmin(router *gin.Engine)  {

	router.SetFuncMap(template.FuncMap{
		"has_powa":modes.HasPowerAction, //判断操作权限
		"has_pows":modes.HasPowerStr, //判断权限使用
		"has_pow":common.HasPower, //选择权限的时候使用
		"GetAdminName":middleware.GetAdminName,
		"GetAdminId":middleware.GetAdminId,
		"GetAdminTruename":middleware.GetAdminTruename,
		"AdminsUname":modes.AdminsUname,
		"GetMenus":middleware.GetMenus,
		"GetCurrentMenu":middleware.GetCurrentMenu,
	})
	//加载静态文件目录
	router.Static("static","static/admins")
	router.Static("uploads","static/uploads")
	//加载模板目录
	router.LoadHTMLGlob("views/admins/**/*")

	//加载自动验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("capt",common.VerifyCaptcha)
	}

	//加载session
	var stores = cookie.NewStore([]byte("fYEd7Io048lal6T9"))
	//使用session中间件
	router.Use(sessions.Sessions("sessionid", stores))

}

func AdminRouter(router *gin.Engine)  {

	initAdmin(router)

	//登录
	router.GET("admin/login",controllers.AdminLogin)
	router.POST("admin/dologin",controllers.AdminLoginAction)
	router.GET("admin/getCaptcha",common.GetCaptcha)


	//路由组，使用中间件
	admin:=router.Group("admin/", middleware.CheckLogin)
	{
		admin.GET("index",controllers.AdminIndex) //首页
		admin.GET("loginout",controllers.LoginOut) //退出
		//系统设置
		admin.GET("system_setting",controllers.Setting) //系统设置
		admin.POST("system_setting_save",controllers.DoSetting) //系统设置
		//栏目管理
		admin.GET("system_menus",controllers.Menus) //栏目管理
		admin.POST("system_menus_save",controllers.MenusSave) //栏目管理
		admin.POST("system_menus_del",controllers.MenusDel) //栏目管理
		admin.GET("system_menus_edit",controllers.MenusEdit) //栏目管理
		//管理员管理
		admin.GET("admin_list",controllers.AdminList) //管理员列表
		admin.GET("admin_add",controllers.AdminAdd) //管理员列表
		admin.GET("admin_edit",controllers.AdminEdit) //管理员列表
		admin.POST("admin_save",controllers.AdminSave) //管理员列表
		admin.GET("admin_del",controllers.AdminDel) //管理员列表
		admin.GET("admin_status",controllers.AdminStatus) //管理员列表
		admin.GET("admin_info",controllers.AdminInfo) //管理员列表
		admin.GET("admin_change_pass",controllers.AdminChangePass) //修改密码
		admin.POST("admin_save_pass",controllers.AdminSavePass) //修改密码

		//权限设置
		admin.GET("admin_power",controllers.AdminPower) //管理员权限
		admin.POST("admin_power_save",controllers.AdminPowerSave) //管理员权限保存
		//系统日志管理
		admin.GET("admin_log",controllers.AdminLogList) //日志列表
		admin.POST("admin_log_del",controllers.AdminLogDel) //删除日志
		admin.POST("admin_log_delbatch",controllers.AdminLogDelBatch) //批量删除日志
		//文章管理
		admin.GET("news_cate_list",controllers.NewsCateList) //分类列表
		admin.GET("news_cate_add",controllers.NewsCateAdd) //添加分类
		admin.POST("news_cate_save",controllers.NewsCateSave) //保存分类
		admin.GET("news_cate_del",controllers.NewsCateDel) //删除分类

		admin.GET("news_list",controllers.NewsList) //文章列表
		admin.GET("news_add",controllers.NewsAdd) //添加文章
		admin.POST("news_save",controllers.NewsSave) //保存文章
		admin.GET("news_del",controllers.NewsDel) //删除文章
		admin.GET("news_trash_list",controllers.NewsTrashList) //回收站列表
		admin.GET("news_restore",controllers.NewsRestore) //恢复
		admin.GET("news_del_real",controllers.NewsDelReal) //永久删除

		admin.POST("fileupload",common.FileUpload) //文件上传
		admin.Any("imgupload",common.ImgUpload) //图片上传

		admin.GET("users_cate_list",controllers.UsersCateList) //用户角色管理
		admin.GET("users_cate_add",controllers.UsersCateAdd) //添加，编辑角色
		admin.POST("users_cate_save",controllers.UsersCateSave) //保存角色
		admin.GET("users_cate_del",controllers.UsersCateDel) //删除角色

		admin.GET("users_list",controllers.UsersList) //会员列表
		admin.GET("users_add",controllers.UsersAdd) //添加会员
		admin.GET("users_edit",controllers.UsersEdit) //修改会员
		admin.GET("users_show",controllers.UsersInfo) //会员信息
		admin.POST("users_save",controllers.UsersSave) //保存会员
		admin.GET("users_del",controllers.UsersDel) //删除会员
		//积分管理
		admin.GET("users_socre",controllers.SocreList) //积分列表
		admin.GET("users_socre_change",controllers.SocreChange) //变更积分
		admin.POST("users_socre_save",controllers.SocreSave) //保存积分
		admin.GET("get_user_score",controllers.GetUserScore) //获取积分

		admin.GET("users_log_list",controllers.UsersLogList) //日志列表
		admin.GET("users_log_del",controllers.UsersLogDel) //删除日志

	}



}
