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
	"math"
)

/* Unalias Command */
type Unalias struct {
	ShellCommand
}

func (this *Unalias) Name() string {
	return "Unalias"
}

func (this *Unalias) CommandCompletion() bool {
	return false
}

func (this *Unalias) MinArgs() int {
	return 1
}

func (this *Unalias) MaxArgs() int {
	return math.MaxInt64
}

func (this *Unalias) ParseCommand(queryurl []string) error {
	ferr := ""
	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments.")

	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")

	} else {

		for _, k := range queryurl {
			_, ok := AliasCommand[k]
			if ok {
				delete(AliasCommand, k)
			} else {
				ferr = ferr + fmt.Errorf("Alias ", k, " doest exist.\n").Error()
			}
		}

	}
	if ferr != "" {
		return fmt.Errorf("%s", ferr)
	}

	return nil

}

func (this *Unalias) PrintHelp() {
	fmt.Println("\\UNALIAS <alias name>")
	fmt.Println("Delete the alias given by <alias name>.")
	fmt.Println("\tExample : \n\t        \\Unalias serverversion;\n\t        \\Unalias subcommand1 subcommand2 serverversion;")
	fmt.Println()
}
