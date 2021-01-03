/**
** @创建时间: 2020/12/1 10:38 下午
** @作者　　: return
** @描述　　:
 */
package log

import (
	"fmt"
	"github.com/gincmf/cmf/data"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)


func Error(message string)  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	prefix := fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(91),"[Error]")
	logger(prefix,message)
}

func Info(message string)  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	prefix := "[Info]"
	logger(prefix,message)
}

func logger(prefix string, message string)  {
	_, file, line, _ := runtime.Caller(2)
	fmt.Println(prefix,file +" "+ strconv.Itoa(line),message)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	currentPath := strings.Replace(dir, "\\", "/", -1)
	_, err = os.Stat(currentPath+"/log/pay.log")
	if os.IsNotExist(err) {
		err = os.MkdirAll(currentPath+"/log/",os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
		fp, err := os.Create(currentPath+"/log/pay.log")  // 如果文件已存在，会将文件清空。
		if err != nil {
			fmt.Println("文件创建失败。")
		}
		defer fp.Close()  //关闭文件，释放资源。
	}

	f, err := os.OpenFile(currentPath+"/log/pay.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, prefix +" "+ time.Now().Format(data.TimeLayout) +" "+ message)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
