package global

import (
	"fmt"
	"io"
	"os"
)

var Writer io.Writer = os.Stdout

func Printf(format string, values ...any) {
	fmt.Fprintf(Writer, format, values...)
}
