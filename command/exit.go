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

/* Exit and Quit Commands */
type Exit struct {
	ShellCommand
}

func (this *Exit) Name() string {
	return "EXIT"
}

func (this *Exit) CommandCompletion() bool {
	return false
}

func (this *Exit) MinArgs() int {
	return 0
}

func (this *Exit) MaxArgs() int {
	return 0
}

func (this *Exit) ExecCommand(args []string) error {
	/* Command to Exit the shell. We set the EXIT flag to true.
	Once this command is processed, and executequery returns to
	HandleInteractiveMode, handle errors (if any) and then exit
	with the correct exit status. If the command contains an
	input argument then throw an error.
	*/
	if len(args) != 0 {
		return errors.New("Too many arguments")
	} else {
		io.WriteString(W, "\n Exiting the shell.")
		EXIT = true
	}
	return nil
}

func (this *Exit) PrintHelp() {
	io.WriteString(W, "\\EXIT; OR QUIT;")
	printDesc(this.Name())
	io.WriteString(W, "\n")
}
