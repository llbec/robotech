package daemon

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
)

func NewDaemon(p int, f ProcFunc) *Daemon {
	return &Daemon{port: p, proc: f}
}

func (dmon *Daemon) Run(logfile string) {
	cmd := background(logfile)
	if cmd != nil {
		return
	}

	e := Lock(os.Getpid())
	if e != nil {
		log.Printf("Process %d close by %v\n", os.Getpid(), e)
		return
	}
	defer Unlock()

	dmon.ch = make(chan int, 1)
	chQuit := make(chan int, 1)
	go dmon.proc(dmon.ch, chQuit)

	rpc.Register(dmon)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", dmon.port))
	if e != nil {
		panic(e)
	}
	defer l.Close()
	go http.Serve(l, nil)

	<-chQuit
	log.Printf("Process %d exit\n", os.Getpid())
}

func (dmon *Daemon) Stop() {
	dmon.Input(-1)
}

func (dmon *Daemon) Input(sig int) {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("127.0.0.1:%d", dmon.port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	err = client.Call("Daemon.Cmd", sig, &reply)
	if err != nil {
		log.Fatal("Daemon.Cmd error:", err)
	}
}

// http client call
//
func (dmon *Daemon) Cmd(arg int, rep *int) error {
	dmon.ch <- arg
	return nil
}

//@link https://github.com/zh-five/xdaemon.git
//@link https://zhuanlan.zhihu.com/p/146192035
var runIdx int = 0               //background调用计数
const ENV_NAME = "MF_DAEMON_IDX" //环境变量名

func background(logFile string) *exec.Cmd {
	//判断子进程还是父进程
	runIdx++
	envIdx, err := strconv.Atoi(os.Getenv(ENV_NAME))
	if err != nil {
		envIdx = 0
	}
	if runIdx <= envIdx { //子进程, 退出
		return nil
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
			log.Printf("%d fatal: logfile %v\n", os.Getpid(), err)
			return nil
		}
		//log.SetOutput(stdout) // 将文件设置为log输出的文件
		//log.SetPrefix("[pricefallback]")
		//log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	//异步启动子进程
	err = cmd.Start()
	if err != nil {
		log.Println(os.Getpid(), "start process failed:", err)
		return nil
	} else {
		//执行成功
		log.Printf("%d : start %s -> %d\n", os.Getpid(), cmd.Path, cmd.Process.Pid)
	}

	return cmd
}
