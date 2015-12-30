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
	"strings"
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

func (this *Pop) ExecCommand(args []string) error {

	var err error

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments")

	} else if len(args) < this.MinArgs() {
		return errors.New("Too few arguments")

	} else if len(args) == 0 {
		/* For \Pop with no input arguments, Pop the top value
		on the stack for every variable.
		*/

		//Named Parameters
		err = Popparam_Helper(NamedParam)
		if err != nil {
			return err
		}

		//Query Parameters
		err = Popparam_Helper(QueryParam)
		if err != nil {
			return err
		}

		//User Defined Session Variables
		err = Popparam_Helper(UserDefSV)
		if err != nil {
			return err
		}

		//Predefined Session Variables
		err = Popparam_Helper(PreDefSV)
		if err != nil {
			return err
		}

	} else {
		//Check what kind of parameter needs to be popped

		if strings.HasPrefix(args[0], "-$") {
			// For Named Parameters
			vble := args[0]
			vble = vble[2:]

			err = PopValue_Helper(false, NamedParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(args[0], "-") {
			// For query parameters
			vble := args[0]
			vble = vble[1:]

			err = PopValue_Helper(false, QueryParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(args[0], "$") {
			// For User defined session variables
			vble := args[0]
			vble = vble[1:]

			err = PopValue_Helper(false, UserDefSV, vble)
			if err != nil {
				return err
			}

		} else {
			// For Predefined session variables
			vble := args[0]

			err = PopValue_Helper(false, PreDefSV, vble)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Pop) PrintHelp() {
	io.WriteString(W, "\\POP [<parameter>]")
	printDesc(this.Name())
	io.WriteString(W, "\n")
}

/* Push value from the Top of the stack onto the parameter stack.
   This is used by the \POP command with no arguments.
*/
func Popparam_Helper(param map[string]*Stack) (err error) {
	for _, v := range param {
		t, err := v.Top()
		if err != nil {
			return err
		}
		v.Push(t)
	}
	return
}
