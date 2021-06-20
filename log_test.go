package mojo_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestLogLevel(t *testing.T) {
	dateFmt := `[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]+`
	cases := []struct {
		level  string
		expect *regexp.Regexp
	}{
		{"fatal", regexp.MustCompile(fmt.Sprintf(`^\[%s\] \[fatal\] Fatal\n$`, dateFmt))},
		{"error", regexp.MustCompile(fmt.Sprintf(`^\[%s\] \[fatal\] Fatal\n\[%s\] \[error\] Error\n$`, dateFmt, dateFmt))},
		{"warn", regexp.MustCompile(fmt.Sprintf(`^\[%s\] \[fatal\] Fatal\n\[%s\] \[error\] Error\n\[%s\] \[warn\] Warn\n$`, dateFmt, dateFmt, dateFmt))},
		{"info", regexp.MustCompile(fmt.Sprintf(`^\[%s\] \[fatal\] Fatal\n\[%s\] \[error\] Error\n\[%s\] \[warn\] Warn\n\[%s\] \[info\] Info\n$`, dateFmt, dateFmt, dateFmt, dateFmt))},
		{"debug", regexp.MustCompile(fmt.Sprintf(`^\[%s\] \[fatal\] Fatal\n\[%s\] \[error\] Error\n\[%s\] \[warn\] Warn\n\[%s\] \[info\] Info\n\[%s\] \[debug\] Debug\n$`, dateFmt, dateFmt, dateFmt, dateFmt, dateFmt))},
	}
	for _, c := range cases {
		handle := strings.Builder{}
		log := mojo.Log{Handle: &handle}

		log.Level(c.level)
		log.Fatal("Fatal")
		log.Error("Error")
		log.Warn("Warn")
		log.Info("Info")
		log.Debug("Debug")

		if !c.expect.MatchString(handle.String()) {
			t.Errorf("Level %s failed: %s", c.level, handle.String())
		}
	}
}

func TestLogShort(t *testing.T) {
	cases := []struct {
		level  string
		expect string
	}{
		{"fatal", "<2>[f] Fatal\n"},
		{"error", "<2>[f] Fatal\n<3>[e] Error\n"},
		{"warn", "<2>[f] Fatal\n<3>[e] Error\n<4>[w] Warn\n"},
		{"info", "<2>[f] Fatal\n<3>[e] Error\n<4>[w] Warn\n<6>[i] Info\n"},
		{"debug", "<2>[f] Fatal\n<3>[e] Error\n<4>[w] Warn\n<6>[i] Info\n<7>[d] Debug\n"},
	}
	for _, c := range cases {
		handle := strings.Builder{}
		log := mojo.Log{Handle: &handle, Short: true}

		log.Level(c.level)
		log.Fatal("Fatal")
		log.Error("Error")
		log.Warn("Warn")
		log.Info("Info")
		log.Debug("Debug")

		if handle.String() != c.expect {
			t.Errorf("Level %s failed: %s", c.level, handle.String())
		}
	}
}

func TestLogContext(t *testing.T) {
	handle := strings.Builder{}
	parentLog := mojo.Log{Handle: &handle, Short: true}
	parentLog.Level("debug")

	log := parentLog.Context("[REQ]")
	log.Debug("Debug")

	if handle.String() != "<7>[d] [REQ] Debug\n" {
		t.Errorf("Log context failed: %s", handle.String())
	}
}
