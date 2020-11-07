package bootstrap

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gincmf/cmf/controller"
	"github.com/gincmf/cmf/data"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"strings"
)

type GroupMapStruct data.GroupMapStruct

var (
	g errgroup.Group
)

var routerMap []data.RouterMapStruct
var TemplateMap data.TemplateMapStruct

var Engine *gin.Engine
var theme, path, themePath string


func Start(){
	//server := &http.Server{
	//	Addr:         ":"+ config.App.Port,
	//	Handler:      register(),
	//}
	//
	//g.Go(func() error {
	//	return server.ListenAndServe()
	//})
	//
	//// 捕获err
	//if err := g.Wait(); err != nil {
	//	fmt.Println("Get errors: ", err)
	//}else {
	//	fmt.Println("Get all num successfully!")
	//}


	register()

}
// func register () http.Handler{
func register () {
	//注册路由
	Engine = gin.Default()

	var store cookie.Store
	var err error
	sessionStore := cookie.NewStore([]byte(config.Database.AuthCode))
	if config.Redis.Host != "" && config.Redis.Enabled {
		store,err = redis.NewStore(10, "tcp", config.Redis.Host+":"+config.Redis.Port, config.Redis.Port, []byte(config.Database.AuthCode))
		if err != nil {
			fmt.Println("[ERROR]", fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(91), err.Error()))
			store = sessionStore
		}
	}else{
		store = sessionStore
	}
	Engine.Use(sessions.Sessions("mySession", store))
	LoadTemplate() //加载模板
	Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	rangeRouter(routerMap)

	// oauth.RegisterOauthRouter(Engine, Db, Conf()) //注册OAuth2.0验证
	// oauth.RegisterTenantRouter(Engine, Db, Conf())

	//扫描主题路径
	files := scanThemeDir(path)
	for _, t := range files {
		//扫描项目模板下的全部模块
		Engine.StaticFS(path+"/"+t.name+"/"+"assets", http.Dir(t.path+"/public/assets"))
	}
	//加载uploads静态资源
	Engine.StaticFS("/uploads", http.Dir("public/uploads"))

	_ = Engine.Run(":" + config.App.Port)
	//配置路由端口
	// return Engine
}

func rangeRouter(routerMap []data.RouterMapStruct) {
	for _, router := range routerMap {
		switch router.Method {
		case "GET":
			Engine.GET(router.RelativePath, router.Handlers...)
		case "POST":
			Engine.POST(router.RelativePath, router.Handlers...)
		case "PUT":
			Engine.PUT(router.RelativePath, router.Handlers...)
		case "DELETE":
			Engine.DELETE(router.RelativePath, router.Handlers...)
		default:
		}
	}
}

//抛出对外注册路由方法
func Get(relativePath string, handlers ...gin.HandlerFunc) {
	routerMap = append(routerMap, data.RouterMapStruct{RelativePath: relativePath, Handlers: handlers, Method: "GET"})
}

func Post(relativePath string, handlers ...gin.HandlerFunc) {
	routerMap = append(routerMap, data.RouterMapStruct{RelativePath: relativePath, Handlers: handlers, Method: "POST"})
}

//处理资源控制器
func Rest(relativePath string, restController controller.RestControllerInterface, handlers ...gin.HandlerFunc) {
	if relativePath == "/" {
		routerMap = append(routerMap,  data.RouterMapStruct{Handlers: []gin.HandlerFunc{restController.Get}, Method: "GET"})
	} else {
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: relativePath, Handlers: append(handlers, restController.Get), Method: "GET"}) //查询全部
		rPath := strings.TrimRight(relativePath, "/") + "/"
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: rPath + ":id", Handlers: append(handlers, restController.Show), Method: "GET"})      //查询一条
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: rPath + ":id", Handlers: append(handlers, restController.Edit), Method: "POST"})     //编辑一条
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: relativePath, Handlers: append(handlers, restController.Store), Method: "POST"})     //新增一条
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: relativePath, Handlers: append(handlers, restController.Delete), Method: "DELETE"})  //删除一条
		routerMap = append(routerMap,  data.RouterMapStruct{RelativePath: rPath + ":id", Handlers: append(handlers, restController.Delete), Method: "DELETE"}) //删除一条
	}
}

// 路由组
func Group(relativePath string, handlers... gin.HandlerFunc) GroupMapStruct {
	return  GroupMapStruct{
		RelativePath:     relativePath,
		Handlers: handlers,
	}
}

func (group *GroupMapStruct) Rest(relativePath string, restController controller.RestControllerInterface, handlers ...gin.HandlerFunc) {
	// 临时赋值
	rPath := ""
	if group.RelativePath != "" {
		rPath = strings.TrimRight(group.RelativePath, "")
	}
	handlers = append(group.Handlers,handlers...)
	Rest(rPath + relativePath,restController,handlers...)
}

func (group *GroupMapStruct) Get(relativePath string, handlers ...gin.HandlerFunc) {
	rPath := ""
	if group.RelativePath != "" {
		rPath = strings.TrimRight("/"+group.RelativePath, "/") + "/"
	}
	handlers = append(group.Handlers,handlers...)
	Get(rPath + relativePath,handlers...)
}

func (group *GroupMapStruct) Post(relativePath string, handlers ...gin.HandlerFunc) {
	rPath := ""
	if group.RelativePath != "" {
		rPath = strings.TrimRight("/"+group.RelativePath, "/") + "/"
	}
	handlers = append(group.Handlers,handlers...)
	Post(rPath + relativePath,handlers...)
}

func LoadTemplate() {
	//加载全部主题路径
	theme = strings.TrimRight(TemplateMap.Theme, "/") + "/"
	path = strings.TrimRight(TemplateMap.ThemePath, "/") + "/"
	themePath = path + theme
	files := scanFiles(themePath)
	Engine.LoadHTMLFiles(files...)
}

type themeDirStruct struct {
	name string
	path string
}

func scanThemeDir(path string) []themeDirStruct {
	dirs, _ := ioutil.ReadDir(path)
	var dirList []themeDirStruct

	for _, dir := range dirs {
		if dir.IsDir() {
			dirList = append(dirList, themeDirStruct{dir.Name(), strings.TrimRight(path+dir.Name(), "/") + "/"})
		}
	}
	return dirList
}

func scanDir(path string) ([]string, []string) {
	dirs, _ := ioutil.ReadDir(path)
	var dirList, fileList []string
	for _, dir := range dirs {
		if dir.IsDir() {
			dirList = append(dirList, strings.TrimRight(path+dir.Name(), "/")+"/")
		} else {
			fileList = append(fileList, path+dir.Name())
		}
	}
	return dirList, fileList
}

// 递归扫描目录
func scanFiles(dirName string) []string {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			subList := scanFiles(strings.TrimRight(dirName+file.Name(), "/") + "/")
			fileList = append(fileList, subList...)

		} else {
			suffix := strings.Split(file.Name(), ".")
			if len(suffix) > 1 && suffix[1] == "html" {
				fileList = append(fileList, dirName+file.Name())
			}
		}
	}
	return fileList
}
