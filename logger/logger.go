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
	Message interface{}
	Level   LogLevel
	EventId LoggerEventID
}

type Logger struct {
	MaxWordSize int
	MuteEnvTest bool
	ServiceName string
	LogLevelMax int
}

type LoggerEvent struct {
	Logger   *Logger
	ID       LoggerEventID
	LogLevel LogLevel
}

type LoggerEventParams struct {
	ID LoggerEventID
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

func (l *Logger) Info(params *LoggerEventParams) *LoggerEvent {
	return l.getLoggerEvent("info", params)
}

func (l *Logger) Debug(params *LoggerEventParams) *LoggerEvent {
	return l.getLoggerEvent("debug", params)
}

func (l *Logger) Warn(params *LoggerEventParams) *LoggerEvent {
	return l.getLoggerEvent("warn", params)
}

func (l *Logger) Error(params *LoggerEventParams) *LoggerEvent {
	return l.getLoggerEvent("error", params)
}

func (event *LoggerEvent) Msgf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	event.Logger.log(Log{
		Message: message,
		EventId: event.ID,
		Level:   event.LogLevel,
	})
}

func (event *LoggerEvent) Msg(message string) {
	event.Logger.log(Log{
		Message: message,
		EventId: event.ID,
		Level:   event.LogLevel,
	})
}

func (l *Logger) log(args Log) {
	dateStr := time.Now().Format("02.01.2006 15:04:05")
	logLevelProps := l.getLogLevelProps(args.Level)

	prettyMessage := l.getPrettyMessage(args.Message)
	if args.EventId != nil {
		prettyMessage = fmt.Sprintf("[%v] %v", args.EventId, prettyMessage)
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	levelColor := color.New(logLevelProps.LevelColor).SprintFunc()
	serviceColor := color.New(l.getServiceColor(l.ServiceName)).SprintFunc()
	textColor := color.New(logLevelProps.TextColor).SprintFunc()

	spaces := l.getSpaces(l.ServiceName)
	spacesAfterLevel := strings.Repeat(" ", l.LogLevelMax-len(logLevelProps.Text)+1)

	fmt.Printf(
		"%s %s %s[%s]:%s%s\n",
		cyan(dateStr),
		levelColor(logLevelProps.Text),
		spacesAfterLevel,
		serviceColor(l.ServiceName),
		spaces,
		textColor(prettyMessage),
	)
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

func (l *Logger) getLoggerEvent(level LogLevel, params *LoggerEventParams) *LoggerEvent {
	var eventID LoggerEventID
	if params != nil && params.ID != nil {
		eventID = params.ID
	}

	return &LoggerEvent{
		Logger:   l,
		LogLevel: level,
		ID:       eventID,
	}
}
