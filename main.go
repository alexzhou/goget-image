package main

import (
	"flag"
	"os"
)
import "fmt"

// 定义命令行参数对应的变量，这三个变量都是指针类型
var cliDir = flag.String("out", "./", "Input dir for saving images,default is the command dir")
var cliUrl = flag.String("url", "", "Input the web page url")

// 定义一个值类型的命令行参数变量，在 Init() 函数中对其初始化
// 因此，命令行参数对应变量的定义和初始化是可以分开的
func Init() {
	//flag.IntVar(&cliFlag, "flagname", 1234, "Just for demo")
}

func download(url string, dir string) {

}

func main() {
	// 初始化变量 cliFlag
	Init()
	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	} else if *cliUrl == "" {
		fmt.Println("should set a web page url")
		flag.PrintDefaults()
		os.Exit(1)
	}

	Download(*cliUrl, *cliDir)

}
