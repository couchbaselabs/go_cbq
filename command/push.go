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
	"strings"

	go_n1ql "github.com/couchbaselabs/go_n1ql"
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
	var err error

	if len(queryurl) > this.MaxArgs() {
		return errors.New("Too many arguments")
	} else if len(queryurl) == 1 {
		return errors.New("Too few arguments")
	} else if len(queryurl) == 0 {
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
		//Check what kind of parameter needs to be set.

		if strings.HasPrefix(queryurl[0], "-$") {
			// For Named Parameters
			vble := queryurl[0]
			vble = vble[2:]

			err = PushValue_Helper(false, NamedParam, vble, queryurl[1])
			if err != nil {
				return err
			}
			//fmt.Println("DEBUG ISHA ", NamedParam[vble])

		} else if strings.HasPrefix(queryurl[0], "-") {
			// For query parameters
			vble := queryurl[0]
			vble = vble[1:]

			err = PushValue_Helper(false, QueryParam, vble, queryurl[1])
			if err != nil {
				return err
			}

			if vble == "creds" {
				/* Define credentials as user/pass and convert into
				   JSON object credentials
				*/

				var creds MyCred

				creds_ret, err := ToCreds(queryurl[1])
				if err != nil {
					return err
				}

				for _, v := range creds_ret {
					creds = append(creds, v)
				}

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

		} else if strings.HasPrefix(queryurl[0], "$") {
			// For User defined session variables
			vble := queryurl[0]
			vble = vble[1:]

			err = PushValue_Helper(false, UserDefSV, vble, queryurl[1])
			if err != nil {
				return err
			}

		} else {
			// For Predefined session variables
			vble := queryurl[0]

			err = PushValue_Helper(false, PreDefSV, vble, queryurl[1])
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (this *Push) PrintHelp() {
	fmt.Println("\\PUSH [<parameter> <value>]")
	fmt.Println("Push the value of the given parameter to the input parameter stack")
	fmt.Println("<parameter> = <prefix><name>")
	fmt.Println(" For Example : \n\t \\PUSH -$r 9.5 \n\t \\PUSH $Val -$r ; \n\t \\PUSH ;")
	fmt.Println()
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
