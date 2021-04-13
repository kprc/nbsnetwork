package logger

import "log"

const(
	LogDebug int = iota
	LogInfo
	LogWarn
	LogErr
)



type logger struct {
	level int
}

type Logger interface {
	Println(v ...interface{})
	DebugPrintln(v ...interface{})
	InfoPrintln(v ...interface{})
	WarnPrintln(v ...interface{})
	ErrPrintln(v ...interface{})

}

var New = func(level int) *logger{
	return &logger{
		level: level,
	}
}

func (l *logger)Println(v ...interface{})  {
	if l.level >= LogDebug{
		log.Println(v...)
	}
}

func (l *logger)DebugPrintln(v ...interface{})  {
	if l.level >= LogDebug{
		log.Println(v...)
	}
}

func (l *logger)InfoPrintln(v ...interface{})  {
	if l.level >= LogInfo{
		log.Println(v...)
	}
}

func (l *logger)WarnPrintln(v ...interface{})  {
	if l.level >= LogWarn{
		log.Println(v...)
	}
}

func (l *logger)ErrPrintln(v ...interface{})  {
	if l.level >= LogErr{
		log.Println(v...)
	}
}


