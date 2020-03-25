package httpLog

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"ws/conf"
)
var DefaultWriter io.Writer = os.Stdout

func Log(fn http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		_body,_:=ioutil.ReadAll(r.Body)
		debugPrintRoute(r.Method,string(_body))
		fn(w,r)
	}
}
func debugPrintRoute(httpMethod,httpBody string) {
	if IsDebugging() {
		debugPrint("%-6s %-25s --> %s\n", httpMethod,"",httpBody)
	}
}
func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[gows-debug] "+format, values...)
	}
}
func IsDebugging()(isDebug bool)  {
	if conf.Config().Common.Env=="dev" {
		isDebug=true
	}
	return
}