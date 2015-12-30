//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package command

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

/* Alias Command */
type Alias struct {
	ShellCommand
}

func (this *Alias) Name() string {
	return "ALIAS"
}

func (this *Alias) CommandCompletion() bool {
	return false
}

func (this *Alias) MinArgs() int {
	return 2
}

func (this *Alias) MaxArgs() int {
	return MAX_ALIASES
}

func (this *Alias) ExecCommand(args []string) error {

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments.")

	} else if len(args) < this.MinArgs() {

		if len(args) == 0 {
			// \ALIAS without input args lists the aliases present.
			if len(AliasCommand) == 0 {
				io.WriteString(W, "There are no defined command aliases. Use \\ALIAS <name> <value> to define an alias.\n")
			}

			for k, v := range AliasCommand {

				tmp := fmt.Sprintf("%-14s %-14s\n", k, v)
				io.WriteString(W, tmp)
			}

		} else {
			// Error out if it has 1 argument.
			return errors.New("Too few arguments")
		}

	} else {
		// Concatenate the elements of args with separator " "
		// to give the input value
		value := strings.Join(args[1:], " ")

		//Add this to the map for Aliases
		key := args[0]

		//Aliases can be replaced.
		if key != "" {
			AliasCommand[key] = value
		}

	}
	return nil

}

func (this *Alias) PrintHelp(desc bool) {
	io.WriteString(W, "\\ALIAS \n\\ALIAS <command name> <command>\n")
	if desc {
		printDesc(this.Name())
	}
	io.WriteString(W, "\n")
}
