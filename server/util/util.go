package util

import "os"

// LoadArg from command or use defaultValue if not present
func LoadArg(command, shortCommand, defaultValue string) (value string) {
	value = defaultValue
	args := os.Args
	for i, arg := range args {
		switch {
		case arg == command, arg == shortCommand:
			value = args[i+1]
		}
	}
	return
}
