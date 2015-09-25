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
	"math"
	"strings"
)

/* Alias Command */
type Alias struct {
	ShellCommand
}

func (this *Alias) Name() string {
	return "Alias"
}

func (this *Alias) CommandCompletion() bool {
	return false
}

func (this *Alias) MinArgs() int {
	return 2
}

func (this *Alias) MaxArgs() int {
	return math.MaxInt64
}

func (this *Alias) ParseCommand(queryurl []string) error {

	fmt.Println("AAAAA: ", queryurl)

	if len(queryurl) > this.MaxArgs() {

		return errors.New("Too many arguments. Quote second input argument")
	} else if len(queryurl) < this.MinArgs() {
		if len(queryurl) == 0 {
			// \ALIAS without input args lists the aliases present.
			if len(AliasCommand) == 0 {
				io.WriteString(W, "There are no defined command aliases. Use \\ALIAS <name> <value> to define.\n")
			}
			io.WriteString(W, "Alias \t Value\n")
			for k, v := range AliasCommand {

				tmp := fmt.Sprintf("%-14s %-14s\n", k, v)
				io.WriteString(W, tmp)
			}

		} else {
			return errors.New("Too few arguments")
		}

	} else {
		value := strings.Join(queryurl[1:], " ")

		fmt.Println("BBBBB: ", value)

		fmt.Println("AAAAA: ", queryurl)

		//Add this to the map for Aliases
		key := queryurl[0]
		_, ok := AliasCommand[key]
		if !ok {
			AliasCommand[key] = value
		} else {
			return errors.New("Alias " + key + " already exists.\n")
		}

	}
	return nil

}

func (this *Alias) PrintHelp() {
	fmt.Println("\\ALIAS <command name> <command>")
	fmt.Println("Create a command alias for a shell command or query.")
	fmt.Println(" <command> = <shell command> or \t <query statement>")
	fmt.Println(" For Example : \n  \\ALIAS serverversion \"select version(), min_version()\" ;\n  \\ALIAS \"\\SET -max-parallelism 8\"")
	fmt.Println()
}
