//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package command

const (
	SHELL_VERSION = "1.0"
)

var (
	QUERYURL   = ""
	DISCONNECT = false
	EXIT       = false
)

/* Command registry : List of Shell Commands supported by cbq */
var COMMAND_LIST = map[string]ShellCommand{

	/* Connection Management */
	"\\connect":    &Connect{},
	"\\disconnect": &Disconnect{},
	"\\exit":       &ExitorQuit{},
	"\\quit":       &ExitorQuit{},

	/* Shell and Server Information */
	"\\help":      &Help{},
	"\\version":   &Version{},
	"\\copyright": &Copyright{},
}

/*
	Interface to be implemented by shell commands.
*/
type ShellCommand interface {
	/* Name of the comand */
	Name() string
	/* Return true if included in shell command completion */
	CommandCompletion() bool
	/* Returns the Minimum number of input arguments required by the function */
	MinArgs() int
	/* Returns the Maximum number of input arguments allowed by the function */
	MaxArgs() int
	/* Method that implements the functionality*/
	ParseCommand(v []string) error
	/* Print Help information for command and its usage with an example */
	PrintHelp()
}
