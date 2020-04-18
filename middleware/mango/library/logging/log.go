package logging

import (
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/conf"
	"github.com/x554462/gin-example/middleware/mango/library/utils"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Level int

var (
	fileName   string
	fileHandle *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	lock       sync.Mutex
	logger     *log.Logger
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func init() {
	var err error
	fileName = getLogFileName()
	fileHandle, err = utils.MustOpenFile(fileName, getLogFilePath())
	if err != nil {
		log.Printf("logging err: %v", err)
	}

	logger = log.New(fileHandle, DefaultPrefix, log.LstdFlags)
}

func checkLogFile() {
	var err error
	f := getLogFileName()
	if f != fileName {
		fileName = f
		fileHandle, err = utils.MustOpenFile(fileName, getLogFilePath())
		if err != nil {
			log.Printf("logging err: %v", err)
		}
		logger = log.New(fileHandle, DefaultPrefix, log.LstdFlags)
	}
}

func Debug(v ...interface{}) {
	output(getPrefix(DEBUG, DefaultCallerDepth), v...)
}

func Info(v ...interface{}) {
	output(getPrefix(INFO, DefaultCallerDepth), v...)
}

func Warn(v ...interface{}) {
	output(getPrefix(WARNING, DefaultCallerDepth), v...)
}

func Error(v ...interface{}) {
	output(getPrefix(ERROR, DefaultCallerDepth), v...)
}

func Fatal(v ...interface{}) {
	output(getPrefix(FATAL, DefaultCallerDepth), v...)
}

func ErrorWithPrefix(prefix string, v ...interface{}) {
	output(fmt.Sprintf("%s[%s]", levelPrefix(ERROR), prefix), v...)
}

func output(prefix string, v ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	checkLogFile()
	logger.SetPrefix(prefix)
	logger.Println(v)
}

func getPrefix(level Level, callerDepth int) string {
	_, file, line, ok := runtime.Caller(callerDepth)
	if ok {
		return fmt.Sprintf("%s[%s:%d]", levelPrefix(level), filepath.Base(file), line)
	} else {
		return levelPrefix(level)
	}
}

func levelPrefix(level Level) string {
	return fmt.Sprintf("[%s]", levelFlags[level])
}

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", conf.ServerConf.RuntimePath, conf.ServerConf.LogPath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.log", conf.ServerConf.LogName, time.Now().Format("2006-01-02"))
}
