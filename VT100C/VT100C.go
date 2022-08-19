package VT100C

import (
	"fmt"
)

//\033[8m    // 消隐

const (
	CleanAttribute = "\x1B[0m"
	Highlight      = "\x1B[1m"
	Underline      = "\x1B[4m"
	Twinkle        = "\x1B[5m"
	ReDisplay      = "\x1B[7m"
	ShowCursor     = "\x1B[?25l"
	HideCursor     = "\x1B[?25h"
)

type Control string

const (
	Up    Control = "A"
	Down  Control = "B"
	Left  Control = "C"
	Right Control = "D"
)

func Set(x, y int) {
	fmt.Printf("\x1b[%d;%dH", x, y)
}
func Save() {
	fmt.Printf("\x1b[s")
}
func Recv() {
	fmt.Printf("\x1B[u")
}
func Move(ControlType Control, n int) {
	fmt.Printf("\x1B[%d%s", n, ControlType)
}
func CleanLine(n int) {
	for i := 0; i < n; i++ {
		Move(Up, 1)
		// 2K清除整行
		fmt.Printf("\x1B[2K")
	}
}
func CleanSnap() {
	fmt.Printf("\x1B[2J")
}
