package logex

var logMap = make(map[string]*Logger)

func add(l *Logger) {

	logMap[l.name] = l
}

func SetLevelByString(name string, level string) {
	l, ok := logMap[name]
	if !ok {
		return
	}

	switch level {
	case "debug":
		l.level = LEVEL_DEBUG
	case "info":
		l.level = LEVEL_INFO
	case "warn":
		l.level = LEVEL_WARN
	case "error":
		l.level = LEVEL_ERROR
	case "fatal":
		l.level = LEVEL_FATAL
	}

}
