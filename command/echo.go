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
	//"strings"

	"github.com/couchbase/query/value"
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
	//var quotestr = false

	if len(queryurl) > this.MaxArgs() {

		return errors.New("Too many arguments")
	} else if len(queryurl) < this.MinArgs() {

		return errors.New("Too few arguments")
	} else {

		for _, val := range queryurl {

			fmt.Println("DEBUG : value going in to Resolve : ", val)
			v, err := Resolve(val)
			if err != nil {
				return err
			}

			fmt.Println("DEBUG : value coming out of Resolve : ", v)

			if v.Type() == value.STRING {

				io.WriteString(W, fmt.Sprintf("%s", v))
				io.WriteString(W, " ")
				fmt.Println("In here  ")
			} else {
				//if !strings.HasPrefix(val, "$") &&
				//	!strings.HasPrefix(val, "-$") &&
				//	!strings.HasPrefix(val, "-") &&
				//	v.Type() == value.BINARY {
				//	return errors.New("Unbalanced input string " + val + " ")
				//}

				tmp, err := ValToStr(v)
				fmt.Println("DEBUG : value coming out of ValtoStr : ", tmp)

				if err != nil {
					return err
				}
				io.WriteString(W, string(tmp))
				io.WriteString(W, " ")
				fmt.Println("NOOO In here  ")
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
