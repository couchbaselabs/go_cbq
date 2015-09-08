//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/couchbase/query/errors"
	//wordwrap "github.com/kr/text"
	"github.com/sbinet/liner"
)

const (
	QRY_EOL     = ";"
	QRY_PROMPT1 = "> "
	QRY_PROMPT2 = "   > "
)

var once = 0
var reset = "\x1b[0m"
var fgRed = "\x1b[31m"

func handleError(err error, tiServer string) errors.Error {

	if strings.Contains(strings.ToLower(err.Error()), "connection refused") {
		return errors.NewShellErrorCannotConnect("Unable to connect to query service " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "unsupported protocol") {
		return errors.NewShellErrorUnsupportedProtocol("Unsupported Protocol Scheme " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "no such host") {
		return errors.NewShellErrorNoSuchHost("No such Host " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "unknown port tcp") {
		return errors.NewShellErrorUnknownPorttcp("Unknown port " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "no host in request url") {
		return errors.NewShellErrorNoHostInRequestUrl("No Host in request URL " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "no route to host") {
		return errors.NewShellErrorNoRouteToHost("No Route to host " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "operation timed out") {
		return errors.NewShellErrorOperationTimeout("Operation timed out. Check query service url " + tiServer)
	} else if strings.Contains(strings.ToLower(err.Error()), "network is unreachable") {
		return errors.NewShellErrorUnreachableNetwork("Network is unreachable " + tiServer)
	} else {
		return errors.NewError(err, "")
	}
}

func HandleInteractiveMode(prompt string) {

	fmt.Println("Interactive tiserver : " + TiServer)

	fmt.Println(prompt)
	// try to find a HOME environment variable
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		// then try USERPROFILE for Windows
		homeDir = os.Getenv("USERPROFILE")
		if homeDir == "" {
			fmt.Printf("Unable to determine home directory, history file disabled\n")
		}
	}

	var liner = liner.NewLiner()
	defer liner.Close()

	LoadHistory(liner, homeDir)

	go signalCatcher(liner)

	// state for reading a multi-line query
	queryLines := []string{}
	fullPrompt := prompt + QRY_PROMPT1
	for {
		line, err := liner.Prompt(fullPrompt)
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Building query string mode: set prompt, gather current line
		fullPrompt = QRY_PROMPT2
		queryLines = append(queryLines, line)

		// If the current line ends with a QRY_EOL, join all query lines,
		// trim off trailing QRY_EOL characters, and submit the query string:
		if strings.HasSuffix(line, QRY_EOL) {
			queryString := strings.Join(queryLines, " ")
			for strings.HasSuffix(queryString, QRY_EOL) {
				queryString = strings.TrimSuffix(queryString, QRY_EOL)
			}
			if queryString != "" {
				//queryString = wordwrap.Wrap(queryString, 10)
				UpdateHistory(liner, homeDir, queryString+QRY_EOL)
				fmt.Println("Called Isha ", once)
				once = once + 1
				err = execute_query(queryString, os.Stdout)
				if err != nil {
					s_err := handleError(err, TiServer)
					fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
				}
			}
			// reset state for multi-line query
			queryLines = []string{}
			fullPrompt = prompt + QRY_PROMPT1
		}
	}

}

/**
 *  Attempt to clean up after ctrl-C otherwise
 *  terminal is left in bad shape
 */
func signalCatcher(liner *liner.State) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	liner.Close()
	os.Exit(0)
}
