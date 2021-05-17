package api

import (
	"github.com/astaxie/beego/validation"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"net/http"
	"pocassist/ui"
	"strings"
)

const (
	SuccessCode = 1
	ErrCode = 0
)

func init() {
	// release 如果需要debug 此处改为 gin.DebugMode
	gin.SetMode(gin.DebugMode)
}

// API Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// Err 通用错误处理
func ErrResp(errStr string) (int,Response) {
	res := Response{
		Code: ErrCode,
		Data: nil,
		Error:  errStr,
	}
	return http.StatusOK, res
}

// SuccessResp 通用处理
func SuccessResp(data interface{}) (int,Response) {
	res := Response{
		Code:  SuccessCode,
		Data:  data,
		Error: "",
	}
	return http.StatusOK, res
}

func DealValidError(valid validation.Validation) string {
	errStr := "参数校验不通过:"
	for _, err := range valid.Errors {
		errStr += err.Message + ";"
	}
	return errStr
}


type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) > len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

// BinaryFileSystem ...
func BinaryFileSystem(root string) *binaryFileSystem {
	return &binaryFileSystem{
		fs: &assetfs.AssetFS{
			Asset:     ui.Asset,
			AssetDir:  ui.AssetDir,
			AssetInfo: ui.AssetInfo,
			Prefix:    root,
		},
	}
}


func Route(port string) {
	router := gin.Default()
	// web
	//router.Use(static.Serve("/static", BinaryFileSystem("ui/build")))
	//router.StaticFS("/css", BinaryFileSystem("ui/build/static"))
	//router.StaticFS("/js", BinaryFileSystem("ui/build/static"))
	//router.StaticFS("/media", BinaryFileSystem("ui/build/static"))
	//
	//router.GET("/", func(c *gin.Context) {
	//	c.Writer.WriteHeader(200)
	//	indexHtml, _ := ui.Asset("ui/build/index.html")
	//	_, _ = c.Writer.Write(indexHtml)
	//	c.Writer.Header().Add("Accept", "text/html")
	//	c.Writer.Flush()
	//})
	router.StaticFS("/ui", BinaryFileSystem("ui/build"))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/ui")
	})

	// api
	router.POST("/api/v1/user/login", GetAuth)
	pluginRoutes := router.Group("/api/v1/poc")
	pluginRoutes.Use(JWT())
	{
		// all
		pluginRoutes.GET("/", GetPlugins)
		// 增
		pluginRoutes.POST("/", CreatePlugin)
		// 改
		pluginRoutes.PUT("/:id/", UpdatePlugin)
		// 详情
		pluginRoutes.GET("/:id/", GetPlugin)
		// 删
		pluginRoutes.DELETE("/:id/", DeletePlugin)
		// 运行
		pluginRoutes.POST("/run/", RunPlugin)
	}

	vulRoutes := router.Group("/api/v1/vul")
	vulRoutes.Use(JWT())
	{
		// basic
		vulRoutes.GET("/basic/", GetBasic)
		// all
		vulRoutes.GET("/", GetVuls)
		// 增
		vulRoutes.POST("/", CreateVul)
		// 改
		vulRoutes.PUT("/:id/", UpdateVul)
		// 详情
		vulRoutes.GET("/:id/", GetVul)
		// 删
		vulRoutes.DELETE("/:id/", DeleteVul)
	}

	appRoutes := router.Group("/api/v1/product")
	appRoutes.Use(JWT())
	{
		// all
		appRoutes.GET("/", GetWebApps)
		// 增
		appRoutes.POST("/", CreateWebApp)
	}



	userRoutes := router.Group("/api/v1/user")
	userRoutes.Use(JWT())
	{
		userRoutes.POST("/self/resetpwd/", SelfResetPassword)
		userRoutes.GET("/info", SelfGetInfo)
		userRoutes.GET("/logout", SelfLogout)
	}

	router.Run(":" + port)
}
