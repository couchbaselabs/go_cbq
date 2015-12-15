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
	return 2
}

func (this *Set) ParseCommand(queryurl []string) error {
	/* Command to set the value of the given parameter to
	   the input value. The top value of the parameter stack
	   is modified. If the command contains no input argument
	   or more than 1 argument then throw an error.
	*/
	var err error

	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else {
		//Check what kind of parameter needs to be set.
		err = PushOrSet(queryurl, true)
		if err != nil {
			return err
		}
	}
	return err
}

func (this *Set) PrintHelp() {
	fmt.Println("\\SET <parameter> <value>")
	fmt.Println("Set the value of the given parameter to the input value. <parameter> = <prefix><name>")
	fmt.Println("\tExample : \n\t        \\SET -$r 9.5 ;\n\t        \\SET $Val -$r ;")
	fmt.Println()
}
