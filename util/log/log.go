package log

import "log"

var Logmode = DebugLog

type LogMode int

const DebugLog LogMode = 0
const InfoLog LogMode = 1
const ErrLog LogMode = 2

func Debug(v ...interface{}) {
	if Logmode > 0 {
		return
	}
	log.Println(v...)
}

func Info(v ...interface{}) {
	if Logmode > 1 {
		return
	}
	log.Println(v...)
}

func Error(v ...interface{}) {
	if Logmode > 2 {
		return
	}
	log.Println(v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}
