package main

import "os"

// LoadArg from command or use defaultValue if not present
func LoadArg(command, defaultValue string) (value string) {
	value = defaultValue
	args := os.Args
	for i, arg := range args {
		switch {
		case arg == command:
			value = args[i+1]
		}
	}
	return
}
