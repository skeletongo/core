package log

import (
	"github.com/cihub/seelog"
)

var (
	Logger seelog.LoggerInterface
)

func init() {
	var err error
	Logger, err = seelog.LoggerFromConfigAsFile("logger.xml")
	if err != nil {
		return
	}
	if err = Logger.SetAdditionalStackDepth(1); err != nil {
		panic(err)
	}
	if err = seelog.ReplaceLogger(Logger); err != nil {
		panic(err)
	}
}

func Reload(fileName string) error {
	newLogger, err := seelog.LoggerFromConfigAsFile(fileName)
	if err != nil {
		return err
	}
	Logger = newLogger
	if err = Logger.SetAdditionalStackDepth(1); err != nil {
		return err
	}
	if err = seelog.ReplaceLogger(Logger); err != nil {
		return err
	}
	return err
}

func Tracef(format string, params ...interface{}) {
	if Logger != nil {
		Logger.Tracef(format, params...)
	}
}

func Debugf(format string, params ...interface{}) {
	if Logger != nil {
		Logger.Debugf(format, params...)
	}
}

func Infof(format string, params ...interface{}) {
	if Logger != nil {
		Logger.Infof(format, params...)
	}
}

func Warnf(format string, params ...interface{}) error {
	if Logger != nil {
		return Logger.Warnf(format, params...)
	}
	return nil
}

func Errorf(format string, params ...interface{}) error {
	if Logger != nil {
		return Logger.Errorf(format, params...)
	}
	return nil
}

func Criticalf(format string, params ...interface{}) error {
	if Logger != nil {
		return Logger.Criticalf(format, params...)
	}
	return nil
}

func Trace(v ...interface{}) {
	if Logger != nil {
		Logger.Trace(v...)
	}
}

func Debug(v ...interface{}) {
	if Logger != nil {
		Logger.Debug(v...)
	}
}

func Info(v ...interface{}) {
	if Logger != nil {
		Logger.Info(v...)
	}
}

func Warn(v ...interface{}) error {
	if Logger != nil {
		return Logger.Warn(v...)
	}
	return nil
}

func Error(v ...interface{}) error {
	if Logger != nil {
		return Logger.Error(v...)
	}
	return nil
}

func Critical(v ...interface{}) error {
	if Logger != nil {
		return Logger.Critical(v...)
	}
	return nil
}

func Flush() {
	if Logger != nil {
		Logger.Flush()
	}
}