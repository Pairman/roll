package log

import (
	"fmt"
	"os"
)

func Info(a ...any) {
	fmt.Print(a...)
}

func Infof(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func Infoln(a ...any) {
	fmt.Println(a...)
}

func Err(a ...any) {
	fmt.Fprint(os.Stderr, a...)
}

func Errf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Print("\n")
}

func Errln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}
