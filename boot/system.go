// Copyright 2017 dzthink
// license that can be found in the LICENSE file.

//provide a cross platform for some syscall api

// +build !windows

package boot
import (
    "os"
    "path"
    "io/ioutil"
    "strconv"
    "github.com/dzthink/gboot/logger"
    "syscall"
    "os/signal"
    "os/user"
)

func(b *Booter) logProcessInfo() {
    //记录进程id
    os.MkdirAll(path.Dir(b.Conf.String("process::pid")), 0755)
    os.Create(b.Conf.String("process::pid"))
    pid := os.Getpid()
    ioutil.WriteFile(b.Conf.String("process::pid"), []byte(strconv.Itoa(pid)), 0755)

    //设置用户和组信息
    u, err := user.Lookup(b.Conf.String("process::user"))
    if err != nil {
        logger.Fatal("error occurs while init user for application", err, "boot")
    }
    uid, _ := strconv.Atoi(u.Uid)
    syscall.Setuid(uid)
    g, gerr := user.LookupGroup(b.Conf.String("process::group"))
    if gerr != nil {
        logger.Fatal("error occurs while init group for application", err, "boot")
    }
    gid, _ := strconv.Atoi(g.Gid)
    syscall.Setgid(gid);
    logger.Info("application process init success,pid:%d,user:%s,group:%s",
        pid, b.Conf.String("process::user"), b.Conf.String("process::group"))
}

func(b *Booter) processSignal() {
    signal.Notify(b.sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt)
    for{
        msg := <-b.sigs
        switch msg {
        default:

        case syscall.SIGUSR1:
            //reload
            b.App.Reload(b.Conf)
        case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
            logger.Info("application stoping, signal[%v]", msg)
            b.App.Stop()
            return
        }
    }
}




