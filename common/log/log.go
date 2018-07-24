package log

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	//"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/modlog"
)

var logPrefixFormatter = "[%s] "

type CliLoggingProvider struct {
}

func (p *CliLoggingProvider) GetLogger(module string) api.Logger {
	return &CliLogger{module: module}
}

type CliLogger struct {
	module string
}

func (l *CliLogger) Fatal(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Fatal(args...)
}
func (l *CliLogger) Fatalf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Fatalf(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Fatalln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Fatalln(args...)
}
func (l *CliLogger) Panic(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Exit(args...)
}
func (l *CliLogger) Panicf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Exitf(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Panicln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Exitln(args...)
}
func (l *CliLogger) Print(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Info(args...)
}
func (l *CliLogger) Printf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Infof(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Println(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Infoln(args...)
}
func (l *CliLogger) Debug(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Warning(args...)
}
func (l *CliLogger) Debugf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Warningf(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Debugln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Warningln(args...)
}
func (l *CliLogger) Info(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Info(args...)
}
func (l *CliLogger) Infof(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Infof(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Infoln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Infoln(args...)
}
func (l *CliLogger) Warn(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Warning(args...)
}
func (l *CliLogger) Warnf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Warningf(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Warnln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Warningln(args...)
}
func (l *CliLogger) Error(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Error(args...)
}
func (l *CliLogger) Errorf(format string, v ...interface{}) {
	args := append([]interface{}{}, l.module)
	args = append(args, v...)
	glog.Errorf(logPrefixFormatter+format, args...)
}
func (l *CliLogger) Errorln(v ...interface{}) {
	args := append([]interface{}{}, fmt.Sprintf(logPrefixFormatter, l.module))
	args = append(args, v...)
	glog.Errorln(args...)
}
