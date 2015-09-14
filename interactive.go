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
	"strings"

	"github.com/couchbase/query/errors"
	"github.com/couchbaselabs/go_cbq/command"
	"github.com/sbinet/liner"
)

/* The following values define the query prompt for cbq.
   The expected end of line character is a ;.
*/
const (
	QRY_EOL     = ";"
	QRY_PROMPT1 = "> "
	QRY_PROMPT2 = "   > "
)

/* The following variables are used to display the error
   messages in red text and then reset the terminal prompt
   color.
*/
var reset = "\x1b[0m"
var fgRed = "\x1b[31m"
var first = false

/* The handleError method creates the error using the methods
   defined in the n1ql errors package. This is where all the
   shell errors are handled.
*/
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

/* This method is used to handle user interaction with the
   cli. After combining the multi line input, it is sent to
   the executequery method which parses and executes the
   input command. In the event an error is returned from the
   query execution, it is printed in red. The input prompt is
   the name of the executable.
*/
func HandleInteractiveMode(prompt string) {

	/* Find the HOME environment variable. If it isnt set then
	   try USERPROFILE for windows. If neither is found then
	   the cli cant find the history file to read from.
	*/
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
		if homeDir == "" {
			fmt.Printf("Unable to determine home directory, history file disabled\n")
		}
	}

	/* Create a new liner */
	var liner = liner.NewLiner()
	defer liner.Close()

	/* Load history from Home directory
	   TODO : Once Histfile and Histsize are introduced then change this code
	*/
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

		/* Check for shell comments : -- and #. Add them to the history
		   but do not send them to be parsed.
		*/
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			UpdateHistory(liner, homeDir, line)
			continue
		}

		// Building query string mode: set prompt, gather current line
		fullPrompt = QRY_PROMPT2
		queryLines = append(queryLines, line)

		/* If the current line ends with a QRY_EOL, join all query lines,
		   trim off trailing QRY_EOL characters, and submit the query string.
		*/
		if strings.HasSuffix(line, QRY_EOL) {
			queryString := strings.Join(queryLines, " ")
			for strings.HasSuffix(queryString, QRY_EOL) {
				queryString = strings.TrimSuffix(queryString, QRY_EOL)
			}
			if queryString != "" {
				UpdateHistory(liner, homeDir, queryString+QRY_EOL)
				err = execute_query(queryString, os.Stdout)
				/* Error handling for Shell errors and errors recieved from
				   go_n1ql.
				*/
				if err != nil {
					s_err := handleError(err, TiServer)
					fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
					if *errorExitFlag == true {
						if first == false {
							first = true
							fmt.Println("Exiting on first error encountered")
							liner.Close()
							os.Clearenv()
							os.Exit(1)
						}
					}
				}

				/* For the \EXIT and \QUIT shell commands we need to
				   make sure that we close the liner and then exit. In
				   the event an error is returned from executequery after
				   the \EXIT command, then handle the error and exit with
				   exit code 1 (which is for general errors).
				*/
				if EXIT == true && err == nil {
					command.EXIT = false
					liner.Close()
					os.Exit(0)
				} else if EXIT == true && err != nil {
					command.EXIT = false
					liner.Close()
					os.Exit(1)
				}

			}

			// reset state for multi-line query
			queryLines = []string{}
			fullPrompt = prompt + QRY_PROMPT1
		}
	}

}

/* If ^C is pressed then Abort the shell. This is
   provided by the liner package.
*/
func signalCatcher(liner *liner.State) {
	liner.SetCtrlCAborts(false)

}
