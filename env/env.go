// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package env has the utilities to load environment variables
package env

import (
	"bufio"
	"os"
	"strings"

	"github.com/cuttle-ai/brain/log"
)

//LoadEnv will load the environment variables from .env file in the directory
func LoadEnv(l log.Log) {
	/*
	 * We will load the environment variables file
	 * We will scan the file
	 * Read the env in the file
	 */
	f, err := os.Open(".env")
	if err != nil {
		l.Error("error while opening the environment variables file .env", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		envs := strings.Split(scanner.Text(), "=")
		if len(envs) != 2 {
			continue
		}
		os.Setenv(envs[0], envs[1])
	}

}
