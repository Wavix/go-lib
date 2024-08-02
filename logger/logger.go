package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Color string

type LoggerResponse interface{}

type LogLevel string

type LogLevelResponse struct {
	Text       string
	LevelColor color.Attribute
	TextColor  color.Attribute
}

type LoggerEventID interface{}

type SetupOptions struct {
	MaxWordSize  int
	MuteEnvTest  bool
	EventIDLimit int
}

type Log struct {
	Message  interface{}
	Level    LogLevel
	EventIDs []LoggerEventID
}

type Logger struct {
	MaxWordSize int
	MuteEnvTest bool
	ServiceName string
	LogLevelMax int
}

func New(service string, options *SetupOptions) *Logger {
	MaxWordSize := 20
	MuteEnvTest := false

	if options != nil && options.MaxWordSize > 0 {
		MaxWordSize = options.MaxWordSize
	}

	if options != nil && options.MuteEnvTest {
		MuteEnvTest = options.MuteEnvTest
	}

	return &Logger{
		MaxWordSize: MaxWordSize,
		MuteEnvTest: MuteEnvTest,
		ServiceName: service,
		LogLevelMax: 4,
	}
}

func (l *Logger) SetPadSize(size int) {
	l.MaxWordSize = size
}

func (l *Logger) MuteTest() {
	l.MuteEnvTest = true
}

func (l *Logger) SetServiceName(name string) {
	l.ServiceName = name
}

func (l *Logger) Info(message interface{}, args ...LoggerEventID) {
	l.log(Log{
		Message:  message,
		EventIDs: args,
		Level:    "info",
	})
}

func (l *Logger) Warn(message interface{}, args ...LoggerEventID) {
	l.log(Log{
		Message:  message,
		EventIDs: args,
		Level:    "warn",
	})
}

func (l *Logger) Error(message interface{}, args ...LoggerEventID) {
	l.log(Log{
		Message:  message,
		EventIDs: args,
		Level:    "error",
	})
}

func (l *Logger) log(args Log) {
	var eventID LoggerEventID
	if len(args.EventIDs) > 0 {
		eventID = args.EventIDs[0]
	}

	prettyMessage := l.getPrettyMessage(args.Message)

	dateStr := time.Now().Format("02.01.2006 15:04:05")
	logLevelProps := l.getLogLevelProps(args.Level)

	spaces := l.getSpaces(l.ServiceName)
	serviceColor := l.getServiceColor(l.ServiceName)

	spacesAfterLevel := strings.Repeat(" ", l.LogLevelMax-len(logLevelProps.Text)+1)

	if eventID != nil {
		eventIDStr := fmt.Sprintf("%v", eventID)
		prettyMessage = fmt.Sprintf("[%v] %v", eventIDStr, prettyMessage)
	}

	color.New(color.FgCyan).Print(dateStr)
	color.New(logLevelProps.LevelColor).Printf(" %s %s", logLevelProps.Text, spacesAfterLevel)

	color.New(serviceColor).Printf("[%s]:%s", l.ServiceName, spaces)
	color.New(logLevelProps.TextColor).Println(prettyMessage)
}

func (l *Logger) getPrettyMessage(message interface{}) string {
	var msgStr string
	switch val := message.(type) {
	case string:
		msgStr = val
	default:
		return ""
	}
	if msgStr == "" {
		return ""
	}
	return strings.ToUpper(msgStr[:1]) + msgStr[1:]
}

func (l *Logger) getLogLevelProps(level LogLevel) LogLevelResponse {
	var levelColor color.Attribute
	var textColor color.Attribute
	var text string

	switch level {
	case "info":
		text, levelColor, textColor = "INFO", color.FgGreen, color.FgWhite
	case "warn":
		text, levelColor, textColor = "WARN", color.FgYellow, color.FgYellow
	case "error":
		text, levelColor, textColor = "ERROR", color.FgRed, color.FgRed
	default:
		text, levelColor, textColor = "DEBUG", color.FgWhite, color.FgWhite
	}

	return LogLevelResponse{
		Text:       text,
		LevelColor: levelColor,
		TextColor:  textColor,
	}
}

func (l *Logger) getServiceColor(service string) color.Attribute {
	colors := []color.Attribute{
		color.FgRed,
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
		color.FgWhite,
	}

	sum := 0
	for _, r := range service {
		sum += int(r)
	}

	return colors[sum%len(colors)]
}

func (l *Logger) getSpaces(service string) string {
	spaces := ""
	for i := len(service); i < l.MaxWordSize; i++ {
		spaces += " "
	}
	return spaces
}
