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
	"io"

	"github.com/couchbase/query/errors"
)

/* Set Command */
type Set struct {
	ShellCommand
}

func (this *Set) Name() string {
	return "SET"
}

func (this *Set) CommandCompletion() bool {
	return false
}

func (this *Set) MinArgs() int {
	return 2
}

func (this *Set) MaxArgs() int {
	return MAX_ARGS
}

func (this *Set) ExecCommand(args []string) (int, string) {
	/* Command to set the value of the given parameter to
	   the input value. The top value of the parameter stack
	   is modified. If the command contains no input argument
	   or more than 1 argument then throw an error.
	*/

	if len(args) > this.MaxArgs() {
		return errors.TOO_MANY_ARGS, ""
	} else if len(args) < this.MinArgs() {
		return errors.TOO_FEW_ARGS, ""
	} else {
		//Check what kind of parameter needs to be set.
		err_code, err_str := PushOrSet(args, true)
		if err_code != 0 {
			return err_code, err_str
		}
	}
	return 0, ""
}

func (this *Set) PrintHelp(desc bool) (int, string) {
	_, werr := io.WriteString(W, "\\SET <parameter> <value>\n")
	if desc {
		err_code, err_str := printDesc(this.Name())
		if err_code != 0 {
			return err_code, err_str
		}
	}
	_, werr = io.WriteString(W, "\n")
	if werr != nil {
		return errors.WRITER_OUTPUT, werr.Error()
	}
	return 0, ""
}
