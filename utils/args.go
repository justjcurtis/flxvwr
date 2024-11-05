package utils

import "os"

func GetArgs() []string {
	args := os.Args[1:]
	stdin, err := ReadStdin()
	if err != nil {
		return args
	}
	return append(args, stdin...)
}
