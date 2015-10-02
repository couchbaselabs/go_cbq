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
	"strings"
)

/* Pop Command */
type Pop struct {
	ShellCommand
}

func (this *Pop) Name() string {
	return "Pop"
}

func (this *Pop) CommandCompletion() bool {
	return false
}

func (this *Pop) MinArgs() int {
	return 0
}

func (this *Pop) MaxArgs() int {
	return 2
}

func (this *Pop) ParseCommand(queryurl []string) error {

	var err error
	//fmt.Println("Isha Queryurl ", queryurl)
	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else if len(queryurl) == 0 {
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

		if strings.HasPrefix(queryurl[0], "-$") {
			// For Named Parameters
			vble := queryurl[0]
			vble = vble[2:]

			err = PopValue_Helper(false, NamedParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(queryurl[0], "-") {
			// For query parameters
			vble := queryurl[0]
			vble = vble[1:]

			err = PopValue_Helper(false, QueryParam, vble)
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(queryurl[0], "$") {
			// For User defined session variables
			vble := queryurl[0]
			vble = vble[1:]

			err = PopValue_Helper(false, UserDefSV, vble)
			if err != nil {
				return err
			}

		} else {
			// For Predefined session variables
			vble := queryurl[0]

			err = PopValue_Helper(false, PreDefSV, vble)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Pop) PrintHelp() {
	fmt.Println("\\Pop [<parameter>]")
	fmt.Println("Pop the value of the given parameter from the input parameter stack")
	fmt.Println("<parameter> = <prefix><name>")
	fmt.Println(" For Example : \n\t \\Pop -$r \n\t \\Pop $Val ; \n\t \\Pop ;")
	fmt.Println()
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
