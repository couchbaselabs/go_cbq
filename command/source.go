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
	//"github.com/sbinet/liner"
)

/* Source Command */
type Source struct {
	ShellCommand
}

func (this *Source) Name() string {
	return "SOURCE"
}

func (this *Source) CommandCompletion() bool {
	return false
}

func (this *Source) MinArgs() int {
	return 1
}

func (this *Source) MaxArgs() int {
	return 1
}

func (this *Source) ExecCommand(args []string) error {
	/* Command to load a file into the shell.
	 */
	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments")

	} else if len(args) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else {
		/* This case needs to be handled in the ShellCommand
		   in the main package, since we need to run each
		   query as it is being read. Otherwise, if we load it
		   into a buffer, we restrict the number of queries that
		   can be loaded from the file.
		*/
		FILE_INPUT = true
	}
	return nil
}

func (this *Source) PrintHelp(desc bool) {
	io.WriteString(W, "\\SOURCE <filename>\n")
	if desc {
		printDesc(this.Name())
	}
	io.WriteString(W, "\n")
}
