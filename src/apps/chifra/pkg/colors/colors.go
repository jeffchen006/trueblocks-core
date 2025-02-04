// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package colors

var Off = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var White = "\033[37m"
var Black = "\033[30m"
var Bright = "\033[1m"
var BrightRed = (Bright + Red)
var BrightGreen = (Bright + Green)
var BrightYellow = (Bright + Yellow)
var BrightBlue = (Bright + Blue)
var BrightMagenta = (Bright + Magenta)
var BrightCyan = (Bright + Cyan)
var BrightWhite = (Bright + White)
var BrightBlack = (Bright + Black)

func ColorsOff() {
	Off = ""
	Red = ""
	Green = ""
	Yellow = ""
	Blue = ""
	Magenta = ""
	Cyan = ""
	White = ""
	Black = ""
	Bright = ""
	BrightRed = (Bright + Red)
	BrightGreen = (Bright + Green)
	BrightYellow = (Bright + Yellow)
	BrightBlue = (Bright + Blue)
	BrightMagenta = (Bright + Magenta)
	BrightCyan = (Bright + Cyan)
	BrightWhite = (Bright + White)
	BrightBlack = (Bright + Black)
}
