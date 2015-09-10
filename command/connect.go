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

/* Connect Command */
type Connect struct {
	ShellCommand
}

func (this *Connect) Name() string {
	return "CONNECT"
}

func (this *Connect) CommandCompletion() bool {
	return false
}

func (this *Connect) MinArgs() int {
	return 1
}

func (this *Connect) MaxArgs() int {
	return 1
}

func (this *Connect) ParseCommand(queryurl []string) error {
	/* Command to connect to the input query service or cluster
	   endpoint. Use the TiServer flag and set it to the value
	   of queryurl. If the command contains no input argument
	   or more than 1 argument then throw an error.
	*/
	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else {
		QUERYURL = queryurl[0]
		fmt.Println("\nCouchbase query shell connected to " + QUERYURL + " . Type Ctrl-D / \\exit / \\quit to exit.")
	}
	return nil
}

func (this *Connect) PrintHelp() {
	fmt.Println("\\CONNECT <url>")
	fmt.Println("Connect to the query service or cluster endpoint url.")
	fmt.Println("\n\t\t Default : http://localhost:8091 \n\t\t \\CONNECT https://my.secure.node.com:8093 ; \n\t\t Connects to query node at my.secure.node.com:8093 using secure https protocol.")
	fmt.Println(" For Example : \n\t \\CONNECT http://172.6.23.2:8091 ;")
	fmt.Println()
}
