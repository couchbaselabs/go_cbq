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
	"io"
)

/* Unalias Command */
type Unalias struct {
	ShellCommand
}

func (this *Unalias) Name() string {
	return "UNALIAS"
}

func (this *Unalias) CommandCompletion() bool {
	return false
}

func (this *Unalias) MinArgs() int {
	return 1
}

func (this *Unalias) MaxArgs() int {
	return MAX_ARGS
}

func (this *Unalias) ExecCommand(args []string) error {

	//Cascade errors for non-existing alias into final error message
	ferr := ""

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments.")

	} else if len(args) < this.MinArgs() {
		return errors.New("Too few arguments")

	} else {

		// Range over input aliases amd delete if they exist.
		for _, k := range args {
			_, ok := AliasCommand[k]
			if ok {
				delete(AliasCommand, k)
			} else {
				ferr = ferr + errors.New("Alias "+k+" doest exist.\n").Error()
			}
		}

	}
	if ferr != "" {
		return errors.New(ferr)
	}

	return nil

}

func (this *Unalias) PrintHelp() {
	io.WriteString(W, "\\UNALIAS <alias name>...")
	printDesc(this.Name())
	io.WriteString(W, "\n")
}
