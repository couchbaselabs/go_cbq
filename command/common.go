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
	//"reflect"
	//"strconv"
	"strings"

	"github.com/couchbase/query/value"
	//"github.com/sbinet/liner"
)

type PtrStrings *[]string

/* Helper function to create a stack. */
func Stack_Helper() *Stack {
	r := make(Stack, 0)
	return &r
}

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
	}
)

type Credentials map[string]string
type MyCred []Credentials

var creds MyCred

func init() {

	/* Populate the Predefined user variable map with default
	   values.
	*/

	//err := PushValue_Helper(false, PreDefSV, "histfile", "\".cbq_history\"")
	//if err != nil {
	//	fmt.Println(err)
	//}

	/*v, _ := Resolve("\".cbq_history\"")
	PreDefSV["histfile"].Push(v)

	v, _ = Resolve("false")
	PreDefSV["autoconfig"].Push(v)

	histlim := int(liner.HistoryLimit)
	v, _ = Resolve(strconv.Itoa(histlim))

	PreDefSV["histsize"].Push(v)

	v, _ = Resolve("nil")
	PreDefSV["querycreds"].Push(v)

	v, _ = Resolve("0")
	PreDefSV["limit"].Push(v)
	*/
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
	//fmt.Println("Res inp ", param)
	param = strings.TrimSpace(param)

	if strings.HasPrefix(param, "\\\\") {
		/* It is a Command alias */
		key := param[2:]
		st_val, ok := AliasCommand[key]
		if !ok {
			err = errors.New("Command for " + key + " does not exist. Please use \\ALIAS to create a command alias.\n")
		} else {

			//st_val = "\"" + st_val + "\""
			//fmt.Println("DEBUG Resolve : ", st_val, reflect.TypeOf(st_val))

			val, err = StrToVal(st_val)

			//fmt.Println("DEBUG Resolve : ", st_val, val.Type())
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
		fmt.Println(key)
		v, ok := QueryParam[key]

		if !ok {
			fmt.Println(errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first"))
		} else {
			val, err = v.Top()
		}
		//fmt.Println("Res inp ", val)

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
				//fmt.Println("Came in here")
			}
			val, err = StrToVal(param)
			//fmt.Println("Test", param, val.Type())
		}
	}
	return
}

/* The StrToVal method converts the input string into a
   value.Value type.
*/
func StrToVal(param string) (val value.Value, err error) {
	//fmt.Println("Isha :: " + param)
	param = strings.TrimSpace(param)
	bytes := []byte(param)

	switch bytes[0] {

	case '{':
		var p map[string]interface{}
		err = json.Unmarshal(bytes, &p)
		if err != nil {
			return value.ZERO_VALUE, err
		}
		val = value.NewValue(p)

	case '[':
		//type sliceValue []interface{}
		var p []interface{}
		err = json.Unmarshal(bytes, &p)
		if err != nil {
			return value.ZERO_VALUE, err
		}
		val = value.NewValue(p)

		//For strings, number, boolean, null and binary
	default:

		val = value.NewValue(bytes)
		err = nil
	}

	return

}

/*func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return string(b)
}*/

/* The ValToStr method converts the input value into a
   string type.
*/
func ValToStr(item value.Value) (param string, err error) {
	//fmt.Println(item.Type())

	//bytes, err := json.MarshalIndent(item, "    ", "    ")

	//if item.Type() == value.BINARY {

	//	param = B2S(item.Actual().([]uint8))
	//	err = nil
	//} else {
	bytes, err := item.MarshalJSON()
	if err != nil {
		param = ""
	}
	param = string(bytes)
	err = nil
	//}

	return
}

/* Stack methods to be used for session parameters */
type Stack []value.Value

/* Push input value val onto the stack */
func (stack *Stack) Push(val value.Value) {
	*stack = append(*stack, val)
}

/* Return the top element in the stack. If the stack
   is empty then return ZERO_VALUE.
*/
func (stack *Stack) Top() (val value.Value, err error) {
	if stack.Len() == 0 {
		val = nil
		err = errors.New("Stack is Empty")
	} else {
		x := stack.Len() - 1
		val = (*stack)[x]
		err = nil
	}

	return
}

func (stack *Stack) SetTop(v value.Value) (err error) {
	if stack.Len() == 0 {
		fmt.Println(errors.New("Stack is Empty. Please use \\PUSH"))
	} else {
		x := stack.Len() - 1
		(*stack)[x] = v
		err = nil
	}
	return
}

/* Delete the top element in the stack. If the stack
   is empty then print err stack empty
*/
func (stack *Stack) Pop() (val value.Value, err error) {
	if stack.Len() == 0 {
		val = nil
		err = errors.New("Stack is Empty. Cannot Pop()")
	} else {
		x := stack.Len() - 1
		val = (*stack)[x]
		*stack = (*stack)[:x]
		err = nil
	}

	return
}

func (stack *Stack) Len() int {
	return len(*stack)
}

/* Helper function to push or set a value in a stack. */
func PushValue_Helper(set bool, param map[string]*Stack, vble, value string) (err error) {

	st_Val, ok := param[vble]

	v, err := Resolve(value)
	if err != nil {
		return err
	} else {
		if ok {
			//fmt.Println("Returned val from Resolve   ", v)
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

	if unset == false {
		//To pop a value from the input stack
		if ok {
			_, err = st_Val.Pop()
		} else {
			err = errors.New("Parameter does not exist")
		}
	} else {
		// Unset the enire stack for given parameter
		if ok {
			for st_Val.Len() > 0 {
				_, err := st_Val.Pop()
				if err != nil {
					return err
				}
			}
		} else {
			err = errors.New("Parameter does not exist")
		}
	}
	return

}

func ToCreds(credsFlag string) (MyCred, error) {

	//Handle the input string of credentials.
	//The string needs to be parsed into a byte array so as to pass to go_n1ql.
	cred := strings.Split(credsFlag, ",")
	var creds MyCred
	creds = append(creds, Credentials{"user": "", "pass": ""})

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
			creds = append(creds, Credentials{"user": up[0], "pass": up[1]})
		}
	}
	return creds, nil

}
