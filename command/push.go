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

func (this *Push) ParseCommand(queryurl []string) error {
	/* Command to set the value of the given parameter to
	   the input value. The top value of the parameter stack
	   is modified. If the command contains no input argument
	   or more than 1 argument then throw an error.
	*/
	fmt.Println("Isha Queryurl ", queryurl)
	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else {
		//Check what kind of parameter needs to be set.
		// For query parameters
		if strings.HasPrefix(queryurl[0], "-") {

			vble := queryurl[0]
			vble = vble[1:]

			v, _ := Resolve(queryurl[1])

			st_val, ok := QueryParam[vble]
			if ok {
				st_val.Push(v)

			} else {
				QueryParam[vble] = Stack_Helper()
				QueryParam[vble].Push(v)
			}

			//tmp, _ := QueryParam[vble].Top()
			fmt.Println(*QueryParam[vble])
		} else if strings.HasPrefix(queryurl[0], "$") {
			// For User defined session variables
		} else if strings.HasPrefix(queryurl[0], "-$") {
			// For Named Parameters
		} else {
			// For Predefined session variables

		}
	}
	return nil
}

func (this *Push) PrintHelp() {
	fmt.Println("\\PUSH [<parameter> <value>]")
	fmt.Println("Set the value of the given parameter to the input value")
	fmt.Println("<parameter> = <prefix><name>")
	fmt.Println(" For Example : \n\t \\SET -$r 9.5 \n\t \\SET $Val -$r ;")
	fmt.Println()
}
