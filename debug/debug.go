package debug

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func DebugPrint(format string, values ...interface{}) {
	if gin.IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(gin.DefaultWriter, "[GATEWAY-debug] "+format, values...)
	}
}
