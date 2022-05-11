package logit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	Err      *log.Logger
	Warn     *log.Logger
	Info     *log.Logger
	Debug    *log.Logger
	Critical *log.Logger
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
	IGreen  = Color("\033[1;92m%s\033[0m")
	IYellow = Color("\033[1;93m%s\033[0m")
	ICyan   = Color("\033[1;96m%s\033[0m")
)

var (
	info     = ICyan
	warn     = IYellow
	err      = Red
	debug    = Magenta
	critical = Red
)

var (
	appName *string
)

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame
}

func Init(name ...string) {
	if len(name) > 0 {
		appName = &name[0]
	} else {
		appName = nil
	}
	Info = log.New(os.Stdout, getLogType("info"), log.Lmsgprefix)
	Warn = log.New(os.Stdout, getLogType("warn"), log.Lmsgprefix)
	Err = log.New(os.Stderr, getLogType("err"), log.Lmsgprefix)
	Debug = log.New(os.Stdout, getLogType("debug"), log.Ldate|log.Ltime|log.Lshortfile)
	Critical = log.New(os.Stdout, getLogType("critical"), log.Ldate|log.Ltime|log.Llongfile)
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func PrintDebugMap(funcName string, m interface{}) {
	data, err := json.MarshalIndent(m, "", "\t|")
	if err != nil {
		log.Println("logger:: json parse err: ", err)
		return
	}
	Debug.Printf("%v:\n%v", funcName, string(data))
}

func GetCurrentFunctionName() string {
	funcName := strings.Split(getFrame(1).Function, "/")
	return fmt.Sprintf("%-25v| ", funcName[len(funcName)-1])
}

func GetCallerFunctionName() string {
	return getFrame(2).Function
}

func getLogType(logType string) string {
	switch logType {
	case "warn":
		if appName != nil {
			return warn("🟠 " + *appName + " | " + "WARN | ")
		} else {
			return warn("🟠 WARN | ")
		}
	case "err":
		if appName != nil {
			return err("❌ " + *appName + " | " + "ERR  | ")
		} else {
			return err("❌ ERR  | ")
		}
	case "debug":
		if appName != nil {
			return debug("🟣 " + *appName + " | " + "DBUG | ")
		} else {
			return debug("🟣 DBUG | ")
		}
	case "critical":
		if appName != nil {
			return critical("🔴 " + *appName + " | " + "CRIT | ")
		} else {
			return critical("🔴 CRIT | ")
		}
	default:
		if appName != nil {
			return info("🔵 " + *appName + " | " + "INFO | ")
		} else {
			return info("🔵 INFO | ")
		}
	}
}

/* Add support to add log to file
file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Info.Print(err)
}
*/
