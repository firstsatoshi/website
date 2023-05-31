package startup

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"sync"
	"syscall"

	"github.com/firstsatoshi/website/common/task"
	"github.com/zeromicro/go-zero/core/logx"
)

func TaskStartup(tasks []task.Task) {

	// Only single proccess ,  we use exclusive file lock
	if true {
		exeName, _ := os.Executable()
		_, exeFileName := path.Split(exeName)
		lockFile := fmt.Sprintf(".%v-lock.pid", exeFileName)
		user, err := user.Current()
		if err == nil {
			lockFile = path.Join(user.HomeDir, lockFile)
		} else {
			lockFile = path.Join("./", lockFile)
		}

		lock, err := os.Create(lockFile)
		if err != nil {
			fmt.Printf("create file lock failed, error:%v\n", err.Error())
			os.Exit(1)
		}
		defer os.Remove(lockFile)
		defer lock.Close()

		// LOCK_EX  exclusive lock
		err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {

			fmt.Printf("Dulpicate process! Please run `ps aux | grep %v` to check!  error:%v\n", exeFileName, err.Error())
			os.Exit(1) // we use Exit(1), don't run defer function
		}

		// unlock
		defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	}

	// ========
	wg := sync.WaitGroup{}

	// sub-goroutine for withdraw task
	for _, t := range tasks {
		wg.Add(1)
		go func(tsk task.Task) {
			defer wg.Done()
			tsk.Start()
		}(t)
	}

	// Mainproces block to deal with KILL signal from user
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-chSignal:
		logx.Info("Receive exit signals, now waiting other sub-goroutines to exit...")
		// close other sub-goroutines
		for _, t := range tasks {
			t.Stop()
		}
	}

	// wait all of sub-goroutines
	wg.Wait()
	logx.Info("All of Sub-goroutines exited, man processs Exit now")
	return
}
