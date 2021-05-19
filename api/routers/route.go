package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jweny/pocassist/api/middleware/jwt"
	"github.com/jweny/pocassist/api/routers/v1"
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
	router.POST("/api/v1/user/login", v1.GetAuth)
	pluginRoutes := router.Group("/api/v1/poc")
	pluginRoutes.Use(jwt.JWT())
	{
		// all
		pluginRoutes.GET("/", v1.GetPlugins)
		// 增
		pluginRoutes.POST("/", v1.CreatePlugin)
		// 改
		pluginRoutes.PUT("/:id/", v1.UpdatePlugin)
		// 详情
		pluginRoutes.GET("/:id/", v1.GetPlugin)
		// 删
		pluginRoutes.DELETE("/:id/", v1.DeletePlugin)
		// 运行
		pluginRoutes.POST("/run/", v1.RunPlugin)
	}

	vulRoutes := router.Group("/api/v1/vul")
	vulRoutes.Use(jwt.JWT())
	{
		// basic
		vulRoutes.GET("/basic/", v1.GetBasic)
		// all
		vulRoutes.GET("/", v1.GetVuls)
		// 增
		vulRoutes.POST("/", v1.CreateVul)
		// 改
		vulRoutes.PUT("/:id/", v1.UpdateVul)
		// 详情
		vulRoutes.GET("/:id/", v1.GetVul)
		// 删
		vulRoutes.DELETE("/:id/", v1.DeleteVul)
	}

	appRoutes := router.Group("/api/v1/product")
	appRoutes.Use(jwt.JWT())
	{
		// all
		appRoutes.GET("/", v1.GetWebApps)
		// 增
		appRoutes.POST("/", v1.CreateWebApp)
	}



	userRoutes := router.Group("/api/v1/user")
	userRoutes.Use(jwt.JWT())
	{
		userRoutes.POST("/self/resetpwd/", v1.SelfResetPassword)
		userRoutes.GET("/info", v1.SelfGetInfo)
		userRoutes.GET("/logout", v1.SelfLogout)
	}

	router.Run(":" + port)
}
