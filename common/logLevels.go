package common

type LogLevel struct {
	NumericLevel int
	Type         string
}

var ELogLevel = LogLevel{}
var MLogLevel = map[string]LogLevel{
	"info":    ELogLevel.Information(),
	"warning": ELogLevel.Warning(),
	"error":   ELogLevel.Error(),
	"fatal":   ELogLevel.Fatal(),
	"none":    ELogLevel.None(),
}

func CanLog(MinLevel, ToLog LogLevel) bool {
	return ToLog.NumericLevel >= MinLevel.NumericLevel
}

func (l LogLevel) Information() LogLevel {
	return LogLevel{
		NumericLevel: 0,
		Type:         "INFO",
	}
}

func (l LogLevel) Warning() LogLevel {
	return LogLevel{
		NumericLevel: 1,
		Type:         "WARNING",
	}
}

func (l LogLevel) Error() LogLevel {
	return LogLevel{
		NumericLevel: 2,
		Type:         "ERROR",
	}
}

func (l LogLevel) Fatal() LogLevel {
	return LogLevel{
		NumericLevel: 3,
		Type:         "FATAL",
	}
}

func (l LogLevel) None() LogLevel {
	return LogLevel{
		NumericLevel: 4,
		Type:         "",
	}
}

func (l LogLevel) String() string {
	return l.Type
}
