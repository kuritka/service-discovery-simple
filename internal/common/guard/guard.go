// Package guard panics when error occurs
package guard

import (
	"fmt"
	"os"

	"github.com/enescakir/emoji"
	"github.com/logrusorgru/aurora"
)

// FailOnError panics when error occurs.
func FailOnError(err error, message string, a ...interface{}) {
	if err != nil {
		m := fmt.Sprintf(message, a...)
		if message != "" {
			print(fmt.Sprintf("%s %s\n%s\n", emoji.RedCircle, aurora.BrightRed(m), aurora.Yellow(err.Error())))
			os.Exit(1)
		}
		print(fmt.Sprintf("%s %s\n", emoji.RedCircle, aurora.Yellow(err.Error())))
		os.Exit(1)
	}
}

// Message write message to output on the screen. Nothing else...
func Message(message string, a ...interface{}) {
	m := fmt.Sprintf(message, a...)
	fmt.Printf("%s %s\n", emoji.YellowSquare, aurora.Yellow(m))
}
