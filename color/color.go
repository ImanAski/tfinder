package color

import (
	"runtime"
	"tfinder/config"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

var NameToCode = map[string]string{
	"red":    Red,
	"green":  Green,
	"yellow": Yellow,
	"blue":   Blue,
	"purple": Purple,
	"cyan":   Cyan,
}

func Colorize(text, pattern string, cfg config.Config) string {
	if runtime.GOOS == "windows" {
		// TODO: should see the windows compability
		// FIXME: hello
		return text
	}

	colorName, exists := cfg.Colors[pattern]
	if !exists {
		colorName = "blue"
	}

	colorCode, exists := NameToCode[colorName]
	if !exists {
		colorCode = Blue
	}

	return colorCode + text + Reset
}
