/* Copyright 2017 dzthink
license that can be found in the LICENSE file.
*/

//provide a cross platform for some syscall api

package boot

import (
    "os"
    "path"
    "io/ioutil"
    "strconv"
    "logger"
    "syscall"
    "os/signal"
)

func(b *Booter) logProcessInfo() {
    //记录进程id
    os.MkdirAll(path.Dir(b.Conf.String("process::pid")), 0755)
    os.Create(b.Conf.String("process::pid"))
    pid := os.Getpid()
    ioutil.WriteFile(b.Conf.String("process::pid"), []byte(strconv.Itoa(pid)), 0755)

    logger.Info("application process init success,pid:%d,user:%s,group:%s",
        pid, b.Conf.String("process::user"), b.Conf.String("process::group"))
}

func(b *Booter) processSignal() {
    signal.Notify(b.sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
    for{
        msg := <-b.sigs
        switch msg {
        default:

        //case syscall.SIG:
            //reload
            //b.App.Reload(b.Conf)
        case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
            logger.Info("application stoping, signal[%v]", msg)
            b.App.Stop()
            return
        }
    }
}




