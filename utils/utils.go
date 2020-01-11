package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" + "0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func PrintRed(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(red(msg))
}

func PrintTeal(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(teal(msg))
}

func PrintGreen(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(green(msg))
}

func PrintYellow(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(yellow(msg))
}

func PrintMagenta(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(magenta(msg))
}

func PrintPurple(m string, args ...interface{}) {
	msg := fmt.Sprintf(m, args...)
	fmt.Printf(purple(msg))
}

// private

var (
	teal    = colorize("\033[1;36m%s\033[0m")
	red     = colorize("\033[1;31m%s\033[0m")
	green   = colorize("\033[1;32m%s\033[0m")
	yellow  = colorize("\033[1;33m%s\033[0m")
	magenta = colorize("\033[1;35m%s\033[0m")
	purple  = colorize("\033[1;34m%s\033[0m")
)

func colorize(color string) func(...interface{}) string {
	colorized := func(args ...interface{}) string {
		return fmt.Sprintf(color,
			fmt.Sprint(args...))
	}

	return colorized
}
