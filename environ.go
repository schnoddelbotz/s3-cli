package main

import (
	"os"
	"strings"
)

func GetEnv(name string) *string {
	name = name + "="
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, name) {
			value := env[len(name):]
			return &value
		}
	}
	return nil
}
