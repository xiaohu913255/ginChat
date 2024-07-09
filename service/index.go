package service

import (
	"github.com/gin-gonic/gin"
	"text/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} HelloWorld
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html")	//模板解析，可传递多个文件
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")	//模板渲染，将data参数值作为数据传递给客户端呈现
	//c.JSON(200, gin.H{
	//	"message": "Hello World!!",
	//})
}
