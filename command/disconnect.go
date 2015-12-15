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

/* Disconnect Command */
type Disconnect struct {
	ShellCommand
}

func (this *Disconnect) Name() string {
	return "DISCONNECT"
}

func (this *Disconnect) CommandCompletion() bool {
	return false
}

func (this *Disconnect) MinArgs() int {
	return 0
}

func (this *Disconnect) MaxArgs() int {
	return 0
}

func (this *Disconnect) ParseCommand(queryurl []string) error {
	/* Command to disconnect service. Use the NoQueryService
	   flag value and the disconnect flag value to determine
	   disconnection. If the command contains an input argument
	   then throw an error.
	*/
	if len(queryurl) != 0 {
		return errors.New("Too many arguments")
	} else {
		DISCONNECT = true
		fmt.Println("\nCouchbase query shell not connected to any endpoint. Use \\CONNECT command to connect.  ")
	}
	return nil
}

func (this *Disconnect) PrintHelp() {
	fmt.Println("\\DISCONNECT;")
	fmt.Println("Disconnect from the query service or cluster endpoint url.")
	fmt.Println("\tExample : \n\t        \\DISCONNECT;")
	fmt.Println()
}
