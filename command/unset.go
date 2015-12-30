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

/* Unset Command */
type Unset struct {
	ShellCommand
}

func (this *Unset) Name() string {
	return "UNSET"
}

func (this *Unset) CommandCompletion() bool {
	return false
}

func (this *Unset) MinArgs() int {
	return 1
}

func (this *Unset) MaxArgs() int {
	return 1
}

func (this *Unset) ExecCommand(args []string) error {
	/* Command to Unset the value of the given parameter.
	 */
	var err error

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments")

	} else if len(args) < this.MinArgs() {
		return errors.New("Too few arguments")

	} else {
		//Check what kind of parameter needs to be Unset.
		// For query parameters
		if strings.HasPrefix(args[0], "-$") {
			// For Named Parameters
			vble := args[0]
			vble = vble[2:]

			err = PopValue_Helper(true, NamedParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(args[0], "-") {
			// For query parameters
			vble := args[0]
			vble = vble[1:]

			err = PopValue_Helper(true, QueryParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(args[0], "$") {
			// For User defined session variables
			vble := args[0]
			vble = vble[1:]

			err = PopValue_Helper(true, UserDefSV, vble)
			if err != nil {
				return err
			}

		} else {
			// For Predefined session variables
			vble := args[0]

			err = PopValue_Helper(true, PreDefSV, vble)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (this *Unset) PrintHelp(desc bool) {
	io.WriteString(W, "\\UNSET <parameter>\n")
	if desc {
		printDesc(this.Name())
	}
	io.WriteString(W, "\n")
}
