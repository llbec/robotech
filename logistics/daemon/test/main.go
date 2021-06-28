package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	c, e := background("test.log", false)
	if e != nil {
		fmt.Println(e)
	}
	if c != nil {
		fmt.Printf("close forward func!\n")
		return
	}
	//以下是业务代码
	log.Println(os.Getpid(), "业务代码开始...")
	time.Sleep(time.Second * 20) //休眠20秒
	for i := 0; i < 20; i++ {
		log.Printf("%d: count %d\n", os.Getpid(), i)
	}
	log.Println(os.Getpid(), "业务代码结束")
}

var runIdx int = 0               //background调用计数
const ENV_NAME = "XW_DAEMON_IDX" //环境变量名

//@link https://www.zhihu.com/people/zh-five
func background(logFile string, isExit bool) (*exec.Cmd, error) {
	//判断子进程还是父进程
	runIdx++
	envIdx, err := strconv.Atoi(os.Getenv(ENV_NAME))
	fmt.Printf("runidx(%d), envidx(%d)\n", runIdx, envIdx)
	if err != nil {
		envIdx = 0
	}
	if runIdx <= envIdx { //子进程, 退出
		fmt.Printf("quit runidx(%d), envidx(%d)\n", runIdx, envIdx)
		return nil, nil
	}

	/*以下是父进程执行的代码*/

	//因为要设置更多的属性, 这里不使用`exec.Command`方法, 直接初始化`exec.Cmd`结构体
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,      //注意,此处是包含程序名的
		Env:  os.Environ(), //父进程中的所有环境变量
	}

	//为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ENV_NAME, runIdx))

	//若有日志文件, 则把子进程的输出导入到日志文件
	if logFile != "" {
		stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(os.Getpid(), ": 打开日志文件错误:", err)
			return nil, err
		}
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	//异步启动子进程
	fmt.Printf("start cmd: %v\n", cmd)
	err = cmd.Start()
	if err != nil {
		log.Println(os.Getpid(), "启动子进程失败:", err)
		return nil, err
	} else {
		//执行成功
		log.Println(os.Getpid(), ":", "启动子进程成功:", "->", cmd.Process.Pid, "\n ")
	}

	//若启动子进程成功, 父进程是否直接退出
	if isExit {
		os.Exit(0)
	}

	return cmd, nil
}
