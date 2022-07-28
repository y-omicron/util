package Log

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"os"
	"sync"
	"time"
)

type Level int

// Available logging levels.
const (
	LevelDebug Level = iota
	LevelTrace
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

var timeFormatLayout = "2006-01-02 15:04:05"

type Context struct {
	wg       sync.WaitGroup
	ctx      context.Context
	can      context.CancelFunc
	msg      chan string
	fMsg     chan string
	IsFile   bool
	LogLevel Level
}

var logContext *Context

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	logContext = &Context{
		ctx:    ctx,
		can:    cancel,
		msg:    make(chan string, 10),
		fMsg:   make(chan string, 10),
		IsFile: false,
	}
}
func New(logLevel Level, OnLogFile bool, name ...string) {
	selfCtx, cancel := context.WithCancel(logContext.ctx)
	logContext.LogLevel = logLevel
	// 文件日志输出
	if OnLogFile {
		logContext.IsFile = true
		wf, err := os.OpenFile(name[0], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			panic(err)
		}
		go func() {
		loop:
			for true {
				select {
				case msg := <-logContext.fMsg:
					_, err := wf.Write([]byte(msg + "\n"))
					if err != nil {
						panic(err)
					}
					logContext.wg.Done()
				case <-selfCtx.Done():
					break loop
				}
			}
			// Drain the msgChan at best effort
			for {
				if len(logContext.msg) == 0 {
					break
				}
				_, err := wf.Write([]byte(<-logContext.fMsg))
				logContext.wg.Done()
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	// 终端日志输出
	go func() {
	loop:
		for true {
			select {
			case msg := <-logContext.msg:
				color.Println(msg)
				logContext.wg.Done()
			case <-selfCtx.Done():
				cancel()
				break loop
			}
		}
		// Drain the msgChan at best effort
		for {
			if len(logContext.msg) == 0 {
				break
			}
			color.Println(<-logContext.msg)
			logContext.wg.Done()
		}
	}()
}
func Close() {
	logContext.wg.Wait()
}
func Warning(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ info ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)

	}
	if logContext.LogLevel < LevelWarning {
		return
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=FF9900>[ warning ] %s </>", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
func Info(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ info ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	if LevelInfo < logContext.LogLevel {
		return
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=99CC66>[ info ]</> %s", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
func Debug(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ debug ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	if LevelDebug < logContext.LogLevel {
		return
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=FF6666>[ debug ]</> %s", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
func Trace(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ trace ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	if LevelTrace < logContext.LogLevel {
		return
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=0066CC>[ trace ]</> %s", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
func Error(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ error ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=FFCCCC>[ error ] %s</>", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
func Fatal(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ fatal ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <fg=FF0033>[ fatal ] %s</>", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
	os.Exit(0)
}
func Success(cStr string, args ...any) {
	logContext.wg.Add(1)
	if logContext.IsFile {
		fMsg := fmt.Sprintf("%s [ fatal ] %s", time.Now().Format(timeFormatLayout), cStr)
		logContext.fMsg <- fmt.Sprintf(fMsg, args...)
	}
	if LevelInfo < logContext.LogLevel {
		return
	}
	msg := fmt.Sprintf("<fg=696969>%s</> <bg=CCFF99>[ Success ] %s</>", time.Now().Format(timeFormatLayout), cStr)
	logContext.msg <- fmt.Sprintf(msg, args...)
	logContext.wg.Add(1)
}
