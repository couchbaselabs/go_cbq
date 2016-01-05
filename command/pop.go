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
	"strings"

	"github.com/couchbase/query/errors"
)

/* Pop Command */
type Pop struct {
	ShellCommand
}

func (this *Pop) Name() string {
	return "POP"
}

func (this *Pop) CommandCompletion() bool {
	return false
}

func (this *Pop) MinArgs() int {
	return 0
}

func (this *Pop) MaxArgs() int {
	return 1
}

func (this *Pop) ExecCommand(args []string) (int, string) {

	if len(args) > this.MaxArgs() {
		return errors.TOO_MANY_ARGS, ""

	} else if len(args) < this.MinArgs() {
		return errors.TOO_FEW_ARGS, ""

	} else if len(args) == 0 {
		/* For \Pop with no input arguments, Pop the top value
		on the stack for every variable. Dont return errors in
		this case as any one of these stacks can be empty.
		*/

		//Named Parameters
		Popparam_Helper(NamedParam)

		//Query Parameters
		Popparam_Helper(QueryParam)

		//User Defined Session Variables
		Popparam_Helper(UserDefSV)

		//Predefined Session Variables
		Popparam_Helper(PreDefSV)

	} else {
		//Check what kind of parameter needs to be popped

		if strings.HasPrefix(args[0], "-$") {
			// For Named Parameters
			vble := args[0]
			vble = vble[2:]

			err_code, err_string := PopValue_Helper(false, NamedParam, vble)
			if err_code != 0 {
				return err_code, err_string
			}

		} else if strings.HasPrefix(args[0], "-") {
			// For query parameters
			vble := args[0]
			vble = vble[1:]

			err_code, err_string := PopValue_Helper(false, QueryParam, vble)
			if err_code != 0 {
				return err_code, err_string
			}

		} else if strings.HasPrefix(args[0], "$") {
			// For User defined session variables
			vble := args[0]
			vble = vble[1:]

			err_code, err_string := PopValue_Helper(false, UserDefSV, vble)
			if err_code != 0 {
				return err_code, err_string
			}

		} else {
			// For Predefined session variables
			vble := args[0]

			err_code, err_string := PopValue_Helper(false, PreDefSV, vble)
			if err_code != 0 {
				return err_code, err_string
			}
		}
	}
	return 0, ""
}

func (this *Pop) PrintHelp(desc bool) (int, string) {
	_, werr := io.WriteString(W, "\\POP \n\\POP <parameter>\n")
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

/* Pop the top value of the parameter stack.
   This is used by the \POP command with no arguments.
*/
func Popparam_Helper(param map[string]*Stack) (int, string) {
	for _, v := range param {
		_, err_code, err_str := v.Pop()
		if err_code != 0 {
			return err_code, err_str
		}
	}
	return 0, ""
}
