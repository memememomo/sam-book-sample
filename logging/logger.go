package logging

import (
	"flag"
	"fmt"
	"sam-book-sample/settings"

	"github.com/golang/glog"
)

func Init() bool {
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()
	return true
}

func DumpForDebug(v interface{}) {
	if settings.Env().IsDebug() {
		glog.InfoDepth(1, fmt.Sprintf("%#v", v))
	}
}
