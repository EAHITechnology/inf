package log

import "go.uber.org/zap"

var AccessLogger *zap.Logger
var DebugLogger *zap.Logger
var InfoLogger *zap.Logger
var WarnningLogger *zap.Logger
var ErrorLogger *zap.Logger

const (
	LOGID = "log_id"
)
