package logger

import(
	"log"
	"os"
	"path"
	"strings"
	"fmt"
)
type Log interface {
	Debug(v interface{}, tag ...string)
	Info(v interface{}, tag ...string)
	Fatal(v interface{}, err error, tag ...string)
	Close()
}

const (
	DEBUG = iota
	INFO
	FATAL
)
type Logger struct {
	logger *log.Logger
	level int32
	file *os.File
}

func NewLogger(file string, level string) (*Logger, error) {
	os.MkdirAll(path.Dir(file), 0755)
	_, fileErr := os.Stat(file)
	if os.IsNotExist(fileErr) {
		os.Create(file)
	}
	logFile, err := os.OpenFile(file, os.O_WRONLY | os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	logger := log.New(logFile, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Llongfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Llongfile)
	simpleLogger := &Logger{
		logger : logger,
		level : levelStringToInt32(level),
		file : logFile,
	}
	return simpleLogger, nil
}

func levelStringToInt32(level string) int32 {
	s := strings.ToLower(level)
	if strings.EqualFold("debug", s) {
		return DEBUG
	} else if strings.EqualFold("info", s) {
		return INFO
	} else if strings.EqualFold("fatal", s) {
		return FATAL
	} else {
		return INFO
	}
}
func(l *Logger) Debug(v string, params ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Output(3, fmt.Sprintf(v, params...))
	}
}

func(l *Logger) Info(v string, params ...interface{}) {
	if l.level <= INFO {
		l.logger.Output(3, fmt.Sprintf(v, params...))
	}
}


func(l *Logger) Fatal(v string,  params ...interface{}) {
	if l.level <= FATAL {
		l.logger.Output(3, fmt.Sprintf(v, params...))
	}
}

func(l *Logger) Close() {
	l.file.Close()
}

var _LOGGER *Logger

func InitLog(file string, level string) {
	var err error
	_LOGGER, err  = NewLogger(file, level)
	if err != nil {
		fmt.Println(fmt.Sprintf("log initial fail with parameter[%s, %s], message is %s", file, level, err.Error()))
	}
}

func Debug(v string, params ...interface{}) {
	if _LOGGER == nil {
		fmt.Println("warning:logger inti fail")
		fmt.Println(fmt.Sprintf(v, params...))
		return
	}
	_LOGGER.Debug(v, params...)
}

func Info(v string, params ...interface{}) {
	if _LOGGER == nil {
		fmt.Println("warning:logger inti fail")
		fmt.Println(fmt.Sprintf(v, params...))
		return
	}
	_LOGGER.Info(v, params...)
}

func Fatal(v string, params ...interface{}) {
	if _LOGGER == nil {
		fmt.Println("warning:logger inti fail")
		fmt.Println(fmt.Sprintf(v, params...))
		return
	}
	_LOGGER.Fatal(v, params...)
}
