/*
 * filename   : logging.go
 * created at : 2014-11-08 19:16:54
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

/*
Features:
  1. Mutiple log levels with separated destination and format
  2. Template powered log format
*/

package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"
	"time"
)

type LogLevel int

type LogWriter struct {
	writer   io.Writer
	level    LogLevel
	template *template.Template
}

type LogData struct {
	Time, Level, Message string
}

type Logger struct {
	defaultLogWriter                                  io.Writer
	name                                              string
	outputLevel                                       LogLevel
	timeFormat                                        string
	DEBUG, INFO, NOTICE, WARN, ERROR, CRITICAL, FATAL *LogWriter
}

const (
	LevelFatal LogLevel = iota
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug

	DefaultLogFormat = "{{.Time}} [{{.Level}}] {{.Message}}"
	DefaultLogLevel  = LevelNotice
	EOL              = "\n"
)

var (
	LevelNames = map[LogLevel]string{
		LevelDebug:    "DEBUG",
		LevelInfo:     "INFO",
		LevelNotice:   "NOTICE",
		LevelWarning:  "WARN",
		LevelError:    "ERROR",
		LevelCritical: "CRIT",
		LevelFatal:    "FATAL",
	}
	DefaultLogWriter = os.Stdout

    TemplateDebug = template.Must(template.New("ColoredDebug").Parse("{{.Time}} \033[37m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateInfo = template.Must(template.New("ColoredInfo").Parse("{{.Time}} \033[32m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateNotice = template.Must(template.New("ColoredNotice").Parse("{{.Time}} \033[34m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateWarning = template.Must(template.New("ColoredWarning").Parse("{{.Time}} \033[33m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateError = template.Must(template.New("ColoredError").Parse("{{.Time}} \033[31m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateCritical = template.Must(template.New("ColoredCritical").Parse("{{.Time}} \033[35m[{{.Level}}] {{.Message}} \033[0m"))
    TemplateFatal = template.Must(template.New("ColoredFatal").Parse("{{.Time}} \033[30;41m[{{.Level}}] {{.Message}} \033[0m"))
)

func LevelString(level LogLevel) string {
	return LevelNames[level]
}

func NewLogger(name string) *Logger {
	if logger, err := createLogger(name, DefaultLogFormat); err == nil {
        logger.DEBUG.SetTemplate(TemplateDebug)
        logger.INFO.SetTemplate(TemplateInfo)
        logger.NOTICE.SetTemplate(TemplateNotice)
        logger.WARN.SetTemplate(TemplateWarning)
        logger.ERROR.SetTemplate(TemplateError)
        logger.CRITICAL.SetTemplate(TemplateCritical)
        logger.FATAL.SetTemplate(TemplateFatal)
        return logger
    } else {
        panic(err)
    }
}

func NewPlainLogger(name string) *Logger {
	if logger, err := createLogger(name, DefaultLogFormat); err == nil {
        return logger
    } else {
        panic(err)
    }
}

func createLogger(name string, format string) (*Logger, error) {
	logger := &Logger{
		name:             name,
		timeFormat:       time.ANSIC,
		outputLevel:      DefaultLogLevel,
		defaultLogWriter: DefaultLogWriter,
		DEBUG:            &LogWriter{level: LevelDebug},
		INFO:             &LogWriter{level: LevelInfo},
		NOTICE:           &LogWriter{level: LevelNotice},
		WARN:             &LogWriter{level: LevelWarning},
		ERROR:            &LogWriter{level: LevelError},
		CRITICAL:         &LogWriter{level: LevelCritical},
		FATAL:            &LogWriter{level: LevelFatal},
	}

	logger.SetWriter(logger.defaultLogWriter)

	if err := logger.SetFormat(DefaultLogFormat); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) foreachLogWriter(f func(*LogWriter)) {
	v := reflect.Indirect(reflect.ValueOf(l))
	t := v.Type()
	needle := reflect.TypeOf(&LogWriter{})
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == needle {
			f(v.Field(i).Interface().(*LogWriter))
		}
	}
}

func (l *Logger) OutputLevel() LogLevel {
	return l.outputLevel
}

func (l *Logger) SetOutputLevel(level LogLevel) {
	l.outputLevel = level
	l.SetWriter(l.defaultLogWriter)
}

func (l *Logger) SetFormat(format string) error {
	if tmpl, err := template.New(l.name).Parse(format); err == nil {
        l.foreachLogWriter(func(w *LogWriter) {
            w.SetTemplate(tmpl)
        })
        return nil
    } else {
        return err
    }
}

func (l *Logger) SetTimeFormat(format string) {
	l.timeFormat = format
}

func (l *Logger) SetWriter(writer io.Writer) {
	l.foreachLogWriter(func(w *LogWriter) {
		if w.level > l.outputLevel {
			w.SetWriter(ioutil.Discard)
		} else {
			w.SetWriter(writer)
		}
	})

    // set default log writer to writer. The writer will be used when user
    // increase log output level.
    l.defaultLogWriter = writer
}

func (w *LogWriter) Printf(format string, v ...interface{}) {
	data := &LogData{
		Time:    time.Now().Format(time.ANSIC),
		Level:   LevelString(w.level),
		Message: fmt.Sprintf(format, v...),
	}
	w.template.Execute(w.writer, data)
	w.writer.Write([]byte(EOL))
}

func (w *LogWriter) Print(v ...interface{}) {
	data := &LogData{
		Time:    time.Now().Format(time.ANSIC),
		Level:   LevelString(w.level),
		Message: fmt.Sprint(v...),
	}
	w.template.Execute(w.writer, data)
	w.writer.Write([]byte(EOL))
}

func (w *LogWriter) SetWriter(writer io.Writer) {
	w.writer = writer
}

func (w *LogWriter) SetTemplate(tmpl *template.Template) {
    w.template = tmpl
}
