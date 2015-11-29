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
	"encoding/json"
	"errors"
	"fmt"
	go_n1ql "github.com/couchbaselabs/go_n1ql"
	"strings"
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
	//fmt.Println("Isha Queryurl ", queryurl)
	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {
		return errors.New("Too few arguments")
	} else {
		//Check what kind of parameter needs to be set.
		// For query parameters
		if strings.HasPrefix(queryurl[0], "-$") {
			// For Named Parameters
			vble := queryurl[0]
			vble = vble[2:]

			err = PushValue_Helper(true, NamedParam, vble, queryurl[1])
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(queryurl[0], "-") {
			// For query parameters
			vble := queryurl[0]
			vble = vble[1:]

			err = PushValue_Helper(true, QueryParam, vble, queryurl[1])
			if err != nil {
				return err
			}

			if vble == "creds" {
				/* Define credentials as user/pass and convert into
				   JSON object credentials
				*/
				type Credentials map[string]string
				type MyCred []Credentials
				var creds MyCred

				cred := strings.Split(queryurl[1], ",")

				/* Append input credentials in [{"user": <username>, "pass" : <password>}]
				   format as expected by go_n1ql creds.
				*/
				for _, i := range cred {
					up := strings.Split(i, ":")
					if len(up) < 2 {
						// One of the input credentials is incorrect
						err := errors.New("Username or Password missing in -credentials/-c option. Please check")
						return err
					} else {
						creds = append(creds, Credentials{"user": up[0], "pass": up[1]})
					}
				}
				creds = append(creds, Credentials{"user": "", "pass": ""})
				ac, err := json.Marshal(creds)
				if err != nil {
					return err
				}
				go_n1ql.SetQueryParams("creds", string(ac))

			} else {
				v, e := QueryParam[vble].Top()
				if e != nil {
					return err
				}

				val, err := ValToStr(v)
				if err != nil {
					return err
				}
				fmt.Println("DEBUG : QUERYPARAM : ", vble, " VALUE : ", val)
				val = strings.Replace(val, "\"", "", 2)
				go_n1ql.SetQueryParams(vble, val)
			}

			//QueryParam["credentials"]

		} else if strings.HasPrefix(queryurl[0], "$") {
			// For User defined session variables
			vble := queryurl[0]
			vble = vble[1:]

			err = PushValue_Helper(true, UserDefSV, vble, queryurl[1])
			if err != nil {
				return err
			}

		} else {
			// For Predefined session variables
			vble := queryurl[0]

			err = PushValue_Helper(true, PreDefSV, vble, queryurl[1])
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (this *Set) PrintHelp() {
	fmt.Println("\\SET <parameter> <value>")
	fmt.Println("Set the value of the given parameter to the input value")
	fmt.Println("<parameter> = <prefix><name>")
	fmt.Println(" For Example : \n\t \\SET -$r 9.5 \n\t \\SET $Val -$r ;")
	fmt.Println()
}
