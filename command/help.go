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
	"math"
	"strings"
)

/* Help Command */
type Help struct {
	ShellCommand
}

func (this *Help) Name() string {
	return "HELP"
}

func (this *Help) CommandCompletion() bool {
	return false
}

func (this *Help) MinArgs() int {
	return 0
}

func (this *Help) MaxArgs() int {
	return math.MaxInt16
}

func (this *Help) ParseCommand(v []string) error {
	/* Input Command : \HELP;
	   Print Help information for all commands. */
	if len(v) == 0 {
		fmt.Println("Help Information for all Shell Commands")
		for _, val := range COMMAND_LIST {
			val.PrintHelp()
		}
	} else {
		/* Input Command : \HELP SET \VERSION;
		   Print help information for input shell commands. The commands
		   need not contain the \ prefix. Return an error if the Command
		   doesnt exist. */
		for _, val := range v {
			if strings.HasPrefix(val, "\\") == false {
				val = "\\" + val
			}
			Cmd, ok := COMMAND_LIST[val]
			if ok == true {
				Cmd.PrintHelp()
			} else {
				fmt.Println("Command doesnt exist. Use \\HELP; to list help for all shell commands.")
				return errors.New("Command doesnt exist")
			}
		}

	}
	return nil
}

func (this *Help) PrintHelp() {
	fmt.Println("\\HELP [<args>]")
	fmt.Println("The input arguments are shell commands. If a * is input then the command displays HELP information for all input shell commands")
	fmt.Println(" For Example : \n\t \\HELP VERSION; \n\t \\HELP EXIT DISCONNECT VERSION; \n\t \\HELP;")
	fmt.Println()
}
