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
	"io"
	"math"
	"strings"
)

/* Echo Command */
type Echo struct {
	ShellCommand
}

func (this *Echo) Name() string {
	return "ECHO"
}

func (this *Echo) CommandCompletion() bool {
	return false
}

func (this *Echo) MinArgs() int {
	return 1
}

func (this *Echo) MaxArgs() int {
	return math.MaxInt64
}

func (this *Echo) ParseCommand(queryurl []string) error {

	if len(queryurl) > this.MaxArgs() {

		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {

		return errors.New("Too few arguments")
	} else {
		//fmt.Println("AAAAA: ", queryurl)

		for _, val := range queryurl {

			if strings.HasPrefix(val, "\\") {

				//Command aliases
				key := val[2:]
				st_val, ok := AliasCommand[key]
				if !ok {
					err := errors.New("\nCommand for " + key + " does not exist. Please use \\ALIAS to create a command alias.\n")
					io.WriteString(W, err.Error())
					continue
				} else {
					io.WriteString(W, st_val)
					io.WriteString(W, "\t")
				}

			} else if strings.HasPrefix(val, "-$") {

				//Named Parameters
				key := val[2:]
				st_val, ok := NamedParam[key]
				if !ok {
					err := errors.New("\nNamed Parameter -$" + key + " does not exist. Please use \\SET or \\PUSH to create it.\n")
					io.WriteString(W, err.Error())
					continue
				}

				v, err := st_val.Top()
				if err != nil {
					return err
				}

				tmp, err := ValToStr(v)
				if err != nil {
					return err
				}

				io.WriteString(W, tmp)
				io.WriteString(W, "\t")

			} else if strings.HasPrefix(val, "-") {

				//Query Parameters
				key := val[1:]

				st_val, ok := QueryParam[key]
				if !ok {
					err := errors.New("\nQuery Parameter -" + key + " does not exist. Please use \\SET or \\PUSH to create it.\n")
					io.WriteString(W, err.Error())
					continue
				}

				v, err := st_val.Top()
				if err != nil {
					return err
				}

				tmp, err := ValToStr(v)
				if err != nil {
					return err
				}

				io.WriteString(W, tmp)
				io.WriteString(W, "\t")

			} else if strings.HasPrefix(val, "$") {
				//User Defined Session Variables

				key := val[1:]

				st_val, ok := UserDefSV[key]
				if !ok {
					err := errors.New("\nUser defined Shell Parameter $" + key + " does not exist. Please use \\SET or \\PUSH to create it.\n")
					io.WriteString(W, err.Error())
					continue
				}

				v, err := st_val.Top()
				if err != nil {
					return err
				}

				tmp, err := ValToStr(v)
				if err != nil {
					return err
				}

				io.WriteString(W, tmp)
				io.WriteString(W, "\t")

			} else if strings.HasPrefix(val, "\"") {
				/* When we want to echo input statements, parse it
				   till we see another ".
				*/

			} else {
				//Predefined Session Variables and generic input

				/* In this case it can either be a predefined session
				   variable since they dont have prefixes, or can be
				   a random string that the user wants to echo.
				*/
				st_val, ok := PreDefSV[val]
				if !ok {
					/* It isnt a Predefined shell parameter, which means
					that it is a generic input.
					*/
					io.WriteString(W, val)
					io.WriteString(W, "\t")

				} else {
					//  It is a Predefined Shell Parameter
					v, err := st_val.Top()
					if err != nil {
						return err
					}

					tmp, err := ValToStr(v)
					if err != nil {
						return err
					}

					io.WriteString(W, tmp)
					io.WriteString(W, "\t")

				}

			}

		}
	}

	io.WriteString(W, "\n")
	return nil

}

func (this *Echo) PrintHelp() {
	fmt.Println("\\ECHO <arg>")
	fmt.Println("Echo the value of the input")
	fmt.Println(" <arg> = <prefix><name> (a parameter)or \n <arg> = <alias> or (command alias)\n <arg> = <input> (any input statement) ")
	fmt.Println(" For Example : \n  \\ECHO -$r ;\n  \\ECHO \\Com; \n  ")
	fmt.Println()
}
