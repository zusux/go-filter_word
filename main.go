package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/importcjj/sensitive"
	"net/http"
	"os"
	"sync"
)

var(
	m *sync.RWMutex
	Filter *sensitive.Filter
)

func init()  {
	m = new (sync.RWMutex)
	Filter = sensitive.New()
	err := Filter.LoadWordDict("dict/dict.txt")
	if err != nil{
		panic(err)
	}
	err = Filter.LoadWordDict("dict/add.txt")
	if err != nil{
		panic(err)
	}
}

func main()  {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/filter/:content", filter)
	r.GET("/find/:content", find)
	r.GET("/add/:content", add)
	r.Run(":8880")
}


func filter(c *gin.Context)  {
	content := c.Param("content")
	word := Filter.Replace(content,'*')
	c.JSON(http.StatusOK,gin.H{
		"status":1,
		"message":"过滤成功",
		"data":word,
	})
}
func find(c *gin.Context)  {
	content := c.Param("content")
	res := Filter.FindAll(content)
	status := 0
	message := "未查找到敏感词"
	if len(res) > 0{
		status = 1
		message = "查找到敏感词"
	}
	c.JSON(http.StatusOK,gin.H{
		"status":status,
		"message":message,
		"data":res,
	})
}
func add(c *gin.Context)  {
	content := c.Param("content")
	Filter.AddWord(content)
	c.JSON(http.StatusOK,gin.H{
		"status":1,
		"message":"添加成功",
		"data":content,
	})
	m.Lock()
	writeToFile(content)
	m.Unlock()
}

func writeToFile(content string)  {
	//追加到 add.txt
	f,err := os.OpenFile("dict/add.txt",os.O_WRONLY,0644)
	defer f.Close()
	if err != nil{
		fmt.Errorf("打开文件 add.txt 失败 %v",err)
		return
	}
	//查找文件末尾偏移量
	n,_ := f.Seek(0,2)
	//从文件末尾写入内容
	if n>0 {
		_,err = f.WriteAt([]byte("\n"+content),n)
	}else{
		_,err = f.WriteAt([]byte(content),n)
	}
}
