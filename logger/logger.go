package logger

import (
	"encoding/json"
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

type LoggerEntityID interface{}
type ExtraData map[string]interface{}

type SetupOptions struct {
	MaxWordSize int
	MuteEnvTest bool
	Plain       bool
}

type Log struct {
	Message  interface{}
	Level    LogLevel
	EntityId LoggerEntityID
	Extra    ExtraData
}

type Logger struct {
	MaxWordSize int
	MuteEnvTest bool
	ServiceName string
	LogLevelMax int
	Plain       bool
}

type LoggerEvent struct {
	logger   *Logger
	id       LoggerEntityID
	logLevel LogLevel
	extra    ExtraData
}

type LoggerContext struct {
	ID    LoggerEntityID
	Extra *ExtraData
}

type LoggerResult struct {
	EntityId *string   `json:"entity_id"`
	Message  string    `json:"message"`
	Level    string    `json:"level"`
	Extra    ExtraData `json:"extra"`
	Service  string    `json:"service"`
	Date     time.Time `json:"date"`
}

func New(service string, options ...SetupOptions) *Logger {
	MaxWordSize := 20
	MuteEnvTest := false
	plainTextLogs := false

	if len(options) > 0 {
		if options[0].MaxWordSize > 0 {
			MaxWordSize = options[0].MaxWordSize
		}

		if options[0].MuteEnvTest {
			MuteEnvTest = options[0].MuteEnvTest
		}

		if options[0].Plain {
			plainTextLogs = options[0].Plain
		}
	}

	return &Logger{
		MaxWordSize: MaxWordSize,
		MuteEnvTest: MuteEnvTest,
		ServiceName: service,
		LogLevelMax: 4,
		Plain:       plainTextLogs,
	}
}

// func New(service string, options *SetupOptions) *Logger {
// 	MaxWordSize := 20
// 	MuteEnvTest := false
// 	plainTextLogs := false

// 	if options != nil && options.MaxWordSize > 0 {
// 		MaxWordSize = options.MaxWordSize
// 	}

// 	if options != nil && options.MuteEnvTest {
// 		MuteEnvTest = options.MuteEnvTest
// 	}

// 	if options != nil && options.Plain {
// 		plainTextLogs = options.Plain
// 	}

// 	return &Logger{
// 		MaxWordSize: MaxWordSize,
// 		MuteEnvTest: MuteEnvTest,
// 		ServiceName: service,
// 		LogLevelMax: 4,
// 		Plain:       plainTextLogs,
// 	}
// }

func (l *Logger) SetPadSize(size int) {
	l.MaxWordSize = size
}

func (l *Logger) MuteTest() {
	l.MuteEnvTest = true
}

func (l *Logger) SetServiceName(name string) {
	l.ServiceName = name
}

func (l *Logger) Context(id LoggerEntityID, extra ...ExtraData) *LoggerEvent {
	var extraData ExtraData

	if len(extra) > 0 {
		extraData = extra[0]
	}

	return &LoggerEvent{
		logLevel: "info",
		logger:   l,
		id:       id,
		extra:    extraData,
	}
}

func (l *Logger) Info() *LoggerEvent {
	return l.getLoggerEvent("info")
}

func (l *Logger) Debug() *LoggerEvent {
	return l.getLoggerEvent("debug")
}

func (l *Logger) Warn() *LoggerEvent {
	return l.getLoggerEvent("warn")
}

func (l *Logger) Error() *LoggerEvent {
	return l.getLoggerEvent("error")
}

func (event *LoggerEvent) Info() *LoggerEvent {
	return event.logger.getLoggerEvent("info", event.id, event.extra)
}

func (event *LoggerEvent) Debug() *LoggerEvent {
	return event.logger.getLoggerEvent("debug", event.id, event.extra)
}

func (event *LoggerEvent) Warn() *LoggerEvent {
	return event.logger.getLoggerEvent("warn", event.id, event.extra)
}

func (event *LoggerEvent) Error() *LoggerEvent {
	return event.logger.getLoggerEvent("error", event.id, event.extra)
}

func (event *LoggerEvent) Extra(data ...interface{}) *LoggerEvent {
	if event.extra == nil {
		event.extra = make(ExtraData, len(data)/2)
	}

	for i := 0; i < len(data); i += 2 {
		event.extra[data[i].(string)] = data[i+1]
	}

	return event
}

func (event *LoggerEvent) Msgf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	event.logger.log(Log{
		Message:  message,
		EntityId: event.id,
		Level:    event.logLevel,
		Extra:    event.extra,
	})
}

func (event *LoggerEvent) Msg(message string) {
	event.logger.log(Log{
		Message:  message,
		EntityId: event.id,
		Level:    event.logLevel,
		Extra:    event.extra,
	})
}
func (l *Logger) log(args Log) {
	if l.Plain {
		l.logPlain(args)
		return
	}

	l.logJSON(args)
}

func (l *Logger) logJSON(args Log) {

	logLevelProps := l.getLogLevelProps(args.Level)
	prettyMessage := l.getPrettyMessage(args.Message)

	log := LoggerResult{
		EntityId: nil,
		Message:  prettyMessage,
		Level:    strings.ToLower(logLevelProps.Text),
		Extra:    args.Extra,
		Service:  l.ServiceName,
		Date:     time.Now(),
	}

	if args.EntityId != "" {
		entityId := fmt.Sprintf("%v", args.EntityId)
		log.EntityId = &entityId
	}

	jsonData, err := json.Marshal(log)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}

func (l *Logger) logPlain(args Log) {
	dateStr := time.Now().Format("02.01.2006 15:04:05")
	logLevelProps := l.getLogLevelProps(args.Level)

	prettyMessage := l.getPrettyMessage(args.Message)

	if args.Extra != nil {
		extraString := ""

		for key, value := range args.Extra {
			extraString += fmt.Sprintf("%v=%v,", key, value)
		}

		extraString = strings.TrimRight(extraString, ",")
		prettyMessage = fmt.Sprintf("[%v] %v", extraString, prettyMessage)
	}

	if args.EntityId != nil && args.EntityId != "" {
		prettyMessage = fmt.Sprintf("[%v] %v", args.EntityId, prettyMessage)
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

func (l *Logger) getLoggerEvent(level LogLevel, params ...interface{}) *LoggerEvent {
	var id string
	var extra ExtraData

	for _, param := range params {
		switch v := param.(type) {
		case string:
			id = v
		case ExtraData:
			extra = v
		}
	}

	return &LoggerEvent{
		logger:   l,
		logLevel: level,
		id:       id,
		extra:    extra,
	}
}
