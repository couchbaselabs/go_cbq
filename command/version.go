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

/* Version Command */
type Version struct {
	ShellCommand
}

func (this *Version) Name() string {
	return "VERSION"
}

func (this *Version) CommandCompletion() bool {
	return false
}

func (this *Version) MinArgs() int {
	return 0
}

func (this *Version) MaxArgs() int {
	return 0
}

func (this *Version) ParseCommand(queryurl []string) error {
	/* Print the shell version. If the command contains an input
	   argument then throw an error.
	*/
	if len(queryurl) != 0 {
		return errors.New("Too many arguments")
	} else {
		fmt.Println("SHELL VERSION : " + SHELL_VERSION)
		fmt.Println("Use N1QL commands select version() or select min_version() to display server version.")
	}
	return nil
}

func (this *Version) PrintHelp() {
	fmt.Println("\\VERSION")
	fmt.Println("Print the Shell Version")
	fmt.Println("\tExample : \n\t        \\VERSION;")
	fmt.Println()
}
