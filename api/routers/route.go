package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/middleware/jwt"
	"github.com/jweny/pocassist/api/routers/v1/auth"
	"github.com/jweny/pocassist/api/routers/v1/plugin"
	"github.com/jweny/pocassist/api/routers/v1/vulnerability"
	"github.com/jweny/pocassist/api/routers/v1/webapp"
	"github.com/jweny/pocassist/pkg/conf"
	"net/http"
)

func Setup() {
	gin.SetMode(conf.GlobalConfig.ServerConfig.RunMode)
}


func InitRouter(port string) {
	router := gin.Default()

	router.StaticFS("/ui", BinaryFileSystem("web/build"))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/ui")
	})

	// api
	router.POST("/api/v1/user/login", auth.Login)
	pluginRoutes := router.Group("/api/v1/poc")
	pluginRoutes.Use(jwt.JWT())
	{
		// all
		pluginRoutes.GET("/", plugin.Get)
		// 增
		pluginRoutes.POST("/", plugin.Add)
		// 改
		pluginRoutes.PUT("/:id/", plugin.Update)
		// 详情
		pluginRoutes.GET("/:id/", plugin.Detail)
		// 删
		pluginRoutes.DELETE("/:id/", plugin.Delete)
		// 测试单个poc
		pluginRoutes.POST("/run/", plugin.Test)
		//// 批量测试poc
		//pluginRoutes.POST("/runs", plugin.RunPlugins)
	}

	vulRoutes := router.Group("/api/v1/vul")
	vulRoutes.Use(jwt.JWT())
	{
		// basic
		vulRoutes.GET("/basic/", vulnerability.Basic)
		// all
		vulRoutes.GET("/", vulnerability.Get)
		// 增
		vulRoutes.POST("/", vulnerability.Create)
		// 改
		vulRoutes.PUT("/:id/", vulnerability.Update)
		// 详情
		vulRoutes.GET("/:id/", vulnerability.Detail)
		// 删
		vulRoutes.DELETE("/:id/", vulnerability.Delete)
	}

	appRoutes := router.Group("/api/v1/product")
	appRoutes.Use(jwt.JWT())
	{
		// all
		appRoutes.GET("/", webapp.Get)
		// 增
		appRoutes.POST("/", webapp.Create)
	}

	userRoutes := router.Group("/api/v1/user")
	userRoutes.Use(jwt.JWT())
	{
		userRoutes.POST("/self/resetpwd/", auth.Reset)
		userRoutes.GET("/info", auth.Self)
		userRoutes.GET("/logout", auth.Logout)
	}

	// todo scan add jwt
	scanRoutes := router.Group("/api/vi/scan")
	{
		scanRoutes.POST("")
	}

	router.Run(":" + port)
}
