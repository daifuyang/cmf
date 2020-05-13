package cmf

import (
	"fmt"
	"github.com/gincmf/cmf/controller"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

//定义路由结构体
type routerMapStruct struct {
	relativePath string
	handlers     []gin.HandlerFunc
	method       string
}

//定义Template结构体
type TemplateMapStruct struct {
	Theme     string `json:"theme"`
	ThemePath string `json:"themePath"`
	Glob      string `json:"glob"`
	Static    string `json:"static"`
}

var routerMap []routerMapStruct
var TemplateMap TemplateMapStruct

var Engine *gin.Engine
var theme,path,themePath string

func Start() {
	//注册路由
	Engine = gin.Default()
	store := cookie.NewStore([]byte("secret"))
	Engine.Use(sessions.Sessions("session", store))
	LoadTemplate() //加载模板

	registerOauthRouter() //注册OAuth2.0验证

	for _, router := range routerMap {
		switch router.method {
		case "GET":
			Engine.GET(router.relativePath, router.handlers...)
		case "POST":
			Engine.POST(router.relativePath, router.handlers...)
		case "PUT":
			Engine.PUT(router.relativePath, router.handlers...)
		case "DELETE":
			Engine.DELETE(router.relativePath, router.handlers...)
		default:
		}
	}

	//扫描主题路径
	files := scanThemeDir(path)
	for _,t := range files{
		//扫描项目模板下的全部模块
		Engine.StaticFS(t.name +"/" + "assets", http.Dir(t.path + "/public/assets" ))
	}
	//加载uploads静态资源
	Engine.StaticFS("uploads", http.Dir("public/uploads" ))
	//配置路由端口
	Engine.Run(":" + Config.App.Port) // 监听并在 0.0.0.0 上启动服务
}

func registerOauthRouter() {

}

//抛出对外注册路由方法
func Get(relativePath string, handlers ...gin.HandlerFunc) {
	routerMap = append(routerMap, routerMapStruct{relativePath, handlers, "GET"})
}

func Post(relativePath string, handlers ...gin.HandlerFunc) {
	routerMap = append(routerMap, routerMapStruct{relativePath, handlers, "POST"})
}

//处理资源控制器
func Rest(relativePath string, restController controller.RestControllerInterface) {
	if relativePath == "/" {
		routerMap = append(routerMap, routerMapStruct{"/api", []gin.HandlerFunc{restController.Get}, "GET"})
	} else{
		routerMap = append(routerMap, routerMapStruct{"/api"+ relativePath, []gin.HandlerFunc{restController.Get}, "GET"})                               //查询全部
		routerMap = append(routerMap, routerMapStruct{"/api" + relativePath + ":id", []gin.HandlerFunc{restController.Show}, "GET"})       //查询一条
		routerMap = append(routerMap, routerMapStruct{"/api" + relativePath + ":id/edit", []gin.HandlerFunc{restController.Edit}, "POST"}) //编辑一条
		routerMap = append(routerMap, routerMapStruct{"/api" + relativePath, []gin.HandlerFunc{restController.Store}, "POST"})             //新增一条
		routerMap = append(routerMap, routerMapStruct{"/api" + relativePath + ":id", []gin.HandlerFunc{restController.Delete}, "DELETE"})  //删除一条
	}
}

func LoadTemplate() {
	//加载全部主题路径
	theme = strings.TrimRight(TemplateMap.Theme, "/") + "/"
	path = strings.TrimRight(TemplateMap.ThemePath, "/") + "/"
	themePath = path + theme
	files := scanFiles(themePath)
	Engine.LoadHTMLFiles(files...)
}

type themeDirStruct struct{
	name string
	path string
}

func scanThemeDir(path string) ([]themeDirStruct){
	dirs,_ := ioutil.ReadDir(path)
	var dirList []themeDirStruct

	for _,dir := range dirs{
		if dir.IsDir() {
			dirList = append(dirList, themeDirStruct{dir.Name(),strings.TrimRight(path + dir.Name(),"/") + "/"} )
		}
	}
	return dirList
}

func scanDir(path string) ([]string,[]string){
	dirs,_ := ioutil.ReadDir(path)
	var dirList,fileList []string
	for _,dir := range dirs{
		if dir.IsDir() {
			dirList = append(dirList,strings.TrimRight(path + dir.Name(),"/") + "/")
		}else {
			fileList = append(fileList,path + dir.Name())
		}
	}
	return dirList,fileList
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
			subList := scanFiles(strings.TrimRight(dirName + file.Name(),"/") + "/")
			fileList = append(fileList,subList...)

		} else {
			suffix := strings.Split(file.Name(), ".")
			if len(suffix) > 1 && suffix[1] == "html" {
				fileList = append(fileList,dirName + file.Name())
			}
		}
	}
	return fileList
}
