package slogsyslog

import (
	"bufio"
	"bytes"
	"log/slog"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)

}

func TestHandler(t *testing.T) {
	defer goleak.VerifyNone(t)
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	opt := Option{
		Level:  slog.LevelDebug,
		Writer: w,
	}

	handler := opt.NewSyslogHandler()
	slog.SetDefault(slog.New(handler))

	logMsg := "test"
	slog.Info(logMsg)
	time.Sleep(time.Second)

	if err := w.Flush(); err != nil {
		t.Error(err)
		return
	}

	expectedByteLen := 99

	rb := make([]byte, buf.Len())
	r := bufio.NewReader(&buf)
	if _, err := r.Read(rb); err != nil {
		t.Errorf("Failed read logs from buffer: %s", err)
		return
	}

	logMsgAr := strings.Split(string(rb), " ")
	actualByteLen, err := strconv.Atoi(logMsgAr[0])
	if err != nil {
		t.Errorf("First word in log not word: %s", err)
		return
	}

	if actualByteLen != expectedByteLen {
		t.Errorf("Expected log len %d, actual %d", actualByteLen, expectedByteLen)
		return
	}

	actualLogMsg := logMsgAr[len(logMsgAr)-1]
	if actualLogMsg != logMsg {
		t.Errorf("Expected log message `%s`, actual `%s`", logMsg, actualLogMsg)
		return
	}
}
