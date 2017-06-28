/* Copyright 2017 dzthink
license that can be found in the LICENSE file.
*/

//Package boot implement core feature of the gboot below:
//1. init log,signal
//2. read and parse the config file
package boot

import(
	"os"
	"sync"
	"flag"
	"syscall"
	"github.com/dzthink/gboot/logger"
)

type Booter struct {
	App Application        //启动应用实例
	Conf ConfigInterface
	Args []string           //命令行参数
	wait sync.WaitGroup     //wait for all gorouting to finish
	sigs chan os.Signal

}

func NewBooter(app Application) *Booter {
	return &Booter{
		App:app,
		sigs:make(chan os.Signal),
	}
}

func(b *Booter) Run() {
	//process panic, trigger signal SIGUSER2
	defer func() {
		if err := recover(); err != nil {
			logger.Fatal("application crashed ,error[%v]", err)
			b.sigs <- syscall.SIGINT
		}
	}()


	//parse the command line parameters
	b.init()
	b.wait.Add(1)
	logger.Info("inital signal hanlers for application")
	//start handle the signals
	go (func(){
		b.processSignal()
		b.wait.Done()
	})()


	//start application
	b.wait.Add(1)
	logger.Info("starting application")
	go (func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Fatal("application crashed ,error[%v]", err)
			}
			b.wait.Done()
			b.sigs <- syscall.SIGINT
		}()
		err := b.App.Start(b.Conf)
		if err != nil {
			logger.Fatal("Start fail ..., error[%v]", err)
		}
	})()

	b.wait.Wait();
}

func(b *Booter) init() {
	b.parseArgs(os.Args)
	logger.InitLog(b.Conf.String("log::log"),b.Conf.String("log::level"))
	logger.Info("application log initial in file %s with level %s", b.Conf.String("log::log"), b.Conf.String("log::level"))
	b.logProcessInfo()
}
func(b *Booter) parseArgs(args []string){
	confFile := flag.String("C", "config/"+b.App.Name()+".ini", "config file path")
	flag.Parse();

	config, err := NewConfig(confFile)
	if err != nil {
		//conf parse error,exiting
	}
	if !config.Has("process::pid") {
		config.Set("process::pid", "var" + getPathSepartor() + b.App.Name() + ".pid")
	}

	if !config.Has("process::user") {
		config.Set("process::user", "nobody")
	}

	if !config.Has("process::group") {
		config.Set("process::group", "nobody")
	}

	if !config.Has("log::log") {
		config.Set("log::log", "log" + getPathSepartor() + b.App.Name() + ".log")
	}

	if !config.Has("log::level") {
		config.Set("log::level", "info")
	}
	b.Conf = config
}


func getPathSepartor() string {
	if os.IsPathSeparator('\\') {
		return "\\"
	} else {
		return "/"
	}
}


