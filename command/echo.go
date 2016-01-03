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
	"io"

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
	return MAX_ARGS
}

func (this *Echo) ExecCommand(args []string) error {

	if len(args) > this.MaxArgs() {
		return errors.New("Too many arguments")

	} else if len(args) < this.MinArgs() {
		return errors.New("Too few arguments")

	} else {

		// Range over the input arguments to echo.
		for _, val := range args {

			// Resolve each value to return a value.Value.
			v, err := Resolve(val)
			if err != nil {
				return err
			}

			// If the value type is string then output it directly.
			if v.Type() == value.STRING {

				io.WriteString(W, v.Actual().(string))
				io.WriteString(W, " ")

			} else {
				// Convert non string values to string and then output.

				tmp, err := ValToStr(v)

				if err != nil {
					return err
				}
				io.WriteString(W, tmp)
				io.WriteString(W, " ")
			}
		}
	}

	io.WriteString(W, "\n")
	return nil

}

func (this *Echo) PrintHelp(desc bool) {
	io.WriteString(W, "\\ECHO <arg>...\n")
	if desc {
		printDesc(this.Name())
	}
	io.WriteString(W, "\n")
}
