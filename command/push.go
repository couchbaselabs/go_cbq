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

/* Push Command */
type Push struct {
	ShellCommand
}

func (this *Push) Name() string {
	return "PUSH"
}

func (this *Push) CommandCompletion() bool {
	return false
}

func (this *Push) MinArgs() int {
	return 0
}

func (this *Push) MaxArgs() int {
	return 2
}

func (this *Push) ExecCommand(args []string) error {
	/* Command to set the value of the given parameter to
	   the input value. The top value of the parameter stack
	   is modified. If the command contains no input argument
	   or more than 1 argument then throw an error.
	*/
	var err error = nil

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments")

	} else if len(args) == 1 {
		return errors.New("Too few arguments")

	} else if len(args) == 0 {
		/* For \PUSH with no input arguments, push the top value
		on the stack for every variable.
		*/

		//Named Parameters
		err = Pushparam_Helper(NamedParam)
		if err != nil {
			return err
		}

		//Query Parameters
		err = Pushparam_Helper(QueryParam)
		if err != nil {
			return err
		}

		//User Defined Session Variables
		err = Pushparam_Helper(UserDefSV)
		if err != nil {
			return err
		}

		//Predefined Session Variables
		err = Pushparam_Helper(PreDefSV)
		if err != nil {
			return err
		}

	} else {
		//Check what kind of parameter needs to be pushed.
		err = PushOrSet(args, false)
		if err != nil {
			return err
		}
	}
	return err
}

func (this *Push) PrintHelp(desc bool) {
	io.WriteString(W, "\\PUSH \n\\PUSH <parameter> <value>\n")
	if desc {
		printDesc(this.Name())
	}
	io.WriteString(W, "\n")
}

/* Push value from the Top of the stack onto the parameter stack.
   This is used by the \PUSH command with no arguments.
*/
func Pushparam_Helper(param map[string]*Stack) (err error) {
	for _, v := range param {
		t, err := v.Top()
		if err != nil {
			return err
		}
		v.Push(t)
	}
	return
}
