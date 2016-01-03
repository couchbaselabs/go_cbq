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
	"io"
	"strconv"
	"strings"

	"github.com/couchbase/query/value"
	go_n1ql "github.com/couchbaselabs/go_n1ql"
	"github.com/sbinet/liner"
)

//type PtrStrings *[]string

var (
	QueryParam map[string]*Stack = map[string]*Stack{}
	NamedParam map[string]*Stack = map[string]*Stack{}
	UserDefSV  map[string]*Stack = map[string]*Stack{}
	PreDefSV   map[string]*Stack = map[string]*Stack{
		"querycreds": Stack_Helper(),
		"limit":      Stack_Helper(),
		"histfile":   Stack_Helper(),
		"histsize":   Stack_Helper(),
		"autoconfig": Stack_Helper(),
		"state":      Stack_Helper(),
	}
)

type Credential map[string]string
type Credentials []Credential

var creds Credentials

func init() {

	/* Populate the Predefined user variable map with default
	   values.
	*/

	var err error

	err = PushValue_Helper(false, PreDefSV, "histfile", "\".cbq_history\"")
	if err != nil {
		io.WriteString(W, err.Error()+"\n")
	}
	err = PushValue_Helper(false, PreDefSV, "autoconfig", "false")
	if err != nil {
		io.WriteString(W, err.Error()+"\n")
	}

	histlim := int(liner.HistoryLimit)
	err = PushValue_Helper(false, PreDefSV, "histsize", strconv.Itoa(histlim))
	if err != nil {
		io.WriteString(W, err.Error()+"\n")
	}

	err = PushValue_Helper(false, PreDefSV, "limit", "0")
	if err != nil {
		io.WriteString(W, err.Error()+"\n")
	}
}

/* The Resolve method is used to evaluate the input parameter
   to the \SET / \PUSH / \POP / \UNSET and \ECHO commands. It
   takes in a string, and resolves it to the appropriate value.
   The input string can be broadly classified into 2 categories,
   1. Parameters (here we will need to read the top value from
   the parameter stack)
   2. Actual values that can be converted to value.Value using
   the StrToVal method.
*/
func Resolve(param string) (val value.Value, err error) {

	/* Parse the input string to check whether it is a parameter
	   or a value. If it is a parameter, then we parse it
	   appropriately to check which stacks top value needs to be
	   returned.
	*/

	param = strings.TrimSpace(param)

	if strings.HasPrefix(param, "\\\\") {
		/* It is a Command alias */
		key := param[2:]
		st_val, ok := AliasCommand[key]
		if !ok {
			err = errors.New("Command for " + key + " does not exist. Please use \\ALIAS to create a command alias.\n")
		} else {

			//Quote input properly so that resolve returns string and not binary.
			if !strings.HasPrefix(st_val, "\"") {
				st_val = "\"" + st_val + "\""
			}
			val, err = StrToVal(st_val)
		}

	} else if strings.HasPrefix(param, "-$") {
		key := param[2:]
		v, ok := NamedParam[key]
		if !ok {
			err = errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first")
		} else {
			val, err = v.Top()
		}

	} else if strings.HasPrefix(param, "-") {
		/* Then it is a query parameter. Retrieve its value and
		return.
		*/
		key := param[1:]
		v, ok := QueryParam[key]

		if !ok {
			err = errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first")
		} else {
			val, err = v.Top()
		}

	} else if strings.HasPrefix(param, "$") {
		key := param[1:]

		v, ok := UserDefSV[key]
		if !ok {
			err = errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first")
		} else {
			val, err = v.Top()
		}

	} else {

		/* There can be two possibilities. 1. Its a Predefined
		   Session Parameter. In this case we lookup its value
		   and return that. 2. It is a value, in which case we
		   directly convert it to a value.Value type and return
		   it.
		*/

		v, ok := PreDefSV[param]
		if ok {
			val, err = v.Top()
		} else {
			if !strings.HasPrefix(param, "\"") {
				param = "\"" + param + "\""
			}
			val, err = StrToVal(param)

		}
	}
	return
}

/* The StrToVal method converts the input string into a
   value.Value type.
*/
func StrToVal(param string) (val value.Value, err error) {

	param = strings.TrimSpace(param)

	if strings.HasPrefix(param, "\"") {
		if strings.HasSuffix(param, "\"") {
			param = param[1 : len(param)-1]
		}
	}

	bytes := []byte(param)

	val = value.NewValue(bytes)

	if val.Type() == value.BINARY {
		param = "\"" + param + "\""
		bytes := []byte(param)
		val = value.NewValue(bytes)
	}

	err = nil
	return

}

/* The ValToStr method converts the input value into a
   string type.
*/
func ValToStr(item value.Value) (param string, err error) {

	//IshaFix : Call String() method in value.Value once it is added
	bytes, err := item.MarshalJSON()
	if err != nil {
		param = ""
	}
	param = string(bytes)
	err = nil

	return
}

/* Helper function to push or set a value in a stack. */
func PushValue_Helper(set bool, param map[string]*Stack, vble, value string) (err error) {

	st_Val, ok := param[vble]

	v, err := Resolve(value)
	if err != nil {
		return err
	} else {
		//Stack already exists
		if ok {
			if set == true {
				err = st_Val.SetTop(v)
				if err != nil {
					return err
				}
			} else if set == false {
				st_Val.Push(v)
			}

		} else {
			/* If the stack for the input variable is empty then
			   push the current value onto the variable stack.
			*/
			param[vble] = Stack_Helper()
			param[vble].Push(v)
		}
	}
	return

}

/* Helper function to pop or unset a value in a stack. */
func PopValue_Helper(unset bool, param map[string]*Stack, vble string) (err error) {

	st_Val, ok := param[vble]

	if unset == true {
		// Unset the enire stack for given parameter
		if ok {
			for st_Val.Len() > 0 {
				_, err := st_Val.Pop()
				if err != nil {
					return err
				}
			}
			//While unsetting also delete the stack for the
			//given variable.
			delete(param, vble)
		} else {
			err = errors.New("Parameter does not exist")
		}
	} else {
		//To pop a value from the input stack
		if ok {
			_, err = st_Val.Pop()
		} else {
			err = errors.New("Parameter does not exist")
		}
	}
	return

}

func ToCreds(credsFlag string) (Credentials, error) {

	//Handle the input string of credentials.
	//The string needs to be parsed into a byte array so as to pass to go_n1ql.
	cred := strings.Split(credsFlag, ",")
	var creds Credentials
	creds = append(creds, Credential{"user": "", "pass": ""})

	/* Append input credentials in [{"user": <username>, "pass" : <password>}]
	format as expected by go_n1ql creds.
	*/
	for _, i := range cred {
		up := strings.Split(i, ":")
		if len(up) < 2 {
			// One of the input credentials is incorrect
			err := errors.New("Username or Password missing in -credentials/-c option. Please check")
			return nil, err
		} else {
			creds = append(creds, Credential{"user": up[0], "pass": up[1]})
		}
	}
	return creds, nil

}

func PushOrSet(args []string, pushvalue bool) error {
	// Check what kind of parameter needs to be set or pushed
	// depending on the pushvalue boolean value.
	var err error = nil
	if strings.HasPrefix(args[0], "-$") {
		// For Named Parameters
		vble := args[0]
		vble = vble[2:]

		err = PushValue_Helper(pushvalue, NamedParam, vble, args[1])
		if err != nil {
			return err
		}
		//Pass the named parameters to the rest api using the SetQueryParams method
		v, e := NamedParam[vble].Top()
		if e != nil {
			return err
		}

		val, err := ValToStr(v)
		if err != nil {
			return err
		}

		val = strings.Replace(val, "\"", "", 2)
		vble = "$" + vble
		go_n1ql.SetQueryParams(vble, val)

	} else if strings.HasPrefix(args[0], "-") {
		// For query parameters
		vble := args[0]
		vble = vble[1:]

		err = PushValue_Helper(pushvalue, QueryParam, vble, args[1])
		if err != nil {
			return err
		}

		if vble == "creds" {
			// Define credentials as user/pass and convert into
			//   JSON object credentials

			var creds Credentials

			creds_ret, err := ToCreds(args[1])
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

			val = strings.Replace(val, "\"", "", 2)
			go_n1ql.SetQueryParams(vble, val)
		}

	} else if strings.HasPrefix(args[0], "$") {
		// For User defined session variables
		vble := args[0]
		vble = vble[1:]

		err = PushValue_Helper(pushvalue, UserDefSV, vble, args[1])
		if err != nil {
			return err
		}

	} else {
		// For Predefined session variables
		vble := args[0]

		err = PushValue_Helper(pushvalue, PreDefSV, vble, args[1])
		if err != nil {
			return err
		}
	}
	return err
}

func printDesc(cmdname string) {

	switch cmdname {
	case "ALIAS":
		io.WriteString(W, "Create an alias for input. <command> = <shell command> or <query statement>\n")
		io.WriteString(W, "\tExample : \n\t        \\ALIAS serverversion \"select version(), min_version()\" ;\n\t        \\ALIAS \"\\SET -max-parallelism 8\";\n")

	case "CONNECT":
		io.WriteString(W, "Connect to the query service or cluster endpoint url.\n")
		io.WriteString(W, "Default : http://localhost:8091\n")
		io.WriteString(W, "\tExample : \n\t        \\CONNECT http://172.6.23.2:8091 ; \n\t         \\CONNECT https://my.secure.node.com:8093 ;\n")

	case "COPYRIGHT":
		io.WriteString(W, "Print Couchbase Copyright information\n")
		io.WriteString(W, "\tExample : \n\t        \\COPYRIGHT;\n")

	case "DISCONNECT":
		io.WriteString(W, "Disconnect from the query service or cluster endpoint url.\n")
		io.WriteString(W, "\tExample : \n\t        \\DISCONNECT;")

	case "ECHO":
		io.WriteString(W, "Echo the value of the input. <arg> = <prefix><name> (a parameter) or \n <arg> = <alias> (command alias) or \n <arg> = <input> (any input statement) \n")
		io.WriteString(W, "\tExample : \n\t        \\ECHO -$r ;\n\t        \\ECHO \\Com; \n")

	case "EXIT":
		io.WriteString(W, "Exit the shell\n")
		io.WriteString(W, "\tExample : \n\t        \\EXIT; \n\t        \\QUIT;\n")

	case "HELP":
		io.WriteString(W, "The input arguments are shell commands. If a * is input then the command displays HELP information for all input shell commands.\n")
		io.WriteString(W, "\tExample : \n\t        \\HELP VERSION; \n\t        \\HELP EXIT DISCONNECT VERSION; \n\t        \\HELP;\n")

	case "POP":
		io.WriteString(W, "Pop the value of the given parameter from the input parameter stack. <parameter> = <prefix><name>\n")
		io.WriteString(W, "\tExample : \n\t        \\Pop -$r ;\n\t        \\Pop $Val ; \n\t        \\Pop ;\n")

	case "PUSH":
		io.WriteString(W, "Push the value of the given parameter to the input parameter stack. <parameter> = <prefix><name>\n")
		io.WriteString(W, "\tExample : \n\t        \\PUSH -$r 9.5 ;\n\t        \\PUSH $Val -$r; \n\t        \\PUSH ;\n")

	case "SET":
		io.WriteString(W, "Set the value of the given parameter to the input value. <parameter> = <prefix><name>\n")
		io.WriteString(W, "\tExample : \n\t        \\SET -$r 9.5 ;\n\t        \\SET $Val -$r ;\n")

	case "SOURCE":
		io.WriteString(W, "Load input file into shell\n")
		io.WriteString(W, " For Example : \n\t \\SOURCE temp1.txt ;\n")

	case "UNALIAS":
		io.WriteString(W, "Delete the alias given by <alias name>.\n")
		io.WriteString(W, "\tExample : \n\t        \\UNALIAS serverversion;\n\t        \\UNALIAS subcommand1 subcommand2 serverversion;\n")

	case "UNSET":
		io.WriteString(W, "Unset the value of the given parameter. <parameter> = <prefix><name> \n")
		io.WriteString(W, "\tExample : \n\t        \\Unset -$r ;\n\t        \\Unset $Val ;\n")

	case "VERSION":
		io.WriteString(W, "Print the Shell Version\n")
		io.WriteString(W, "\tExample : \n\t        \\VERSION;\n")

	default:
		io.WriteString(W, "IshaFix : Does not exist\n")

	}

}
