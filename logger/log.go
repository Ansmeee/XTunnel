package logger

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/glog"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var Logger *glog.Logger

func Init() {
	logger := glog.New()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot find home dir")
	}

	if err = logger.SetConfigWithMap(g.Map{
		"path":               filepath.Join(homeDir, "XTunnel", "logs"),
		"level":              "all",
		"timeFormat":         "2006-01-02 15:04:05",
		"stStatus":           0,
		"writerColorEnable":  true,
		"rotateSize":         "500M", // 日志文件大小切割
		"rotateExpire":       "1d",   // 日志文件切割最大有效期
		"rotateBackupLimit":  50,     // 日志文件最大备份数量
		"rotateBackupExpire": "3d",   // 日志文件备份最大有效期
		"StdoutPrint":        true,
	}); err != nil {
		log.Fatalf("InitLog err: %s", err.Error())
	}

	Logger = logger
}

type LogData struct {
	Timestamp string
	TraceID   string
	SpanID    string
	Tag       string
	Data      interface{}
}

func buildData(ctx context.Context, tag string, args interface{}) LogData {
	traceID := gtrace.GetTraceID(ctx)
	spanID := gtrace.GetSpanID(ctx)
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Microsecond), 10) //纳秒/1000=微秒

	data := LogData{
		Timestamp: timestamp,
		TraceID:   traceID,
		SpanID:    spanID,
		Data:      args,
		Tag:       tag,
	}

	return data
}

func Info(ctx context.Context, tag string, args interface{}) {
	Logger.Skip(1).Line(true).Stack(false).Info(ctx, buildData(ctx, tag, args))
}

func Error(ctx context.Context, tag string, args interface{}) {
	Logger.Skip(1).Line(true).Stack(false).Error(ctx, buildData(ctx, tag, args))
}
