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
	"github.com/couchbase/query/value"
	"strings"
)

type PtrStrings *[]string

func Stack_Helper() *Stack {
	r := make(Stack, 0)
	return &r
}

var (
	QueryParam map[string]*Stack = map[string]*Stack{}
	UserDefSV  map[string]*Stack = map[string]*Stack{}
	PreDefSV   map[string]*Stack = map[string]*Stack{
		"querycreds": Stack_Helper(),
		"limit":      Stack_Helper(),
		"histfile":   Stack_Helper(),
		"histsize":   Stack_Helper(),
		"autoconfig": Stack_Helper(),
	}
	NamedParam map[string]*Stack = map[string]*Stack{}
)

func init() {

	/* Populate the Predefined user variable map with default
	   values.
	*/

	v, _ := Resolve("\".cbq_history\"")
	PreDefSV["histfile"].Push(v)

	v, _ = Resolve("false")
	PreDefSV["autoconfig"].Push(v)

	/*v, _ := StrToVal("\".cbq_history\"")
	PreDefSV["querycreds"].Push(v)

	vo, _ := PreDefSV["querycreds"].Top()
	fmt.Println("LALALA", vo.Actual()) */
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

	if strings.HasPrefix(param, "-") {
		/* Then it is a query parameter. Retrieve its value and
		return.
		*/
		key := param[1:]
		v, ok := QueryParam[key]
		if ok {
			err = errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first")
		} else {
			val, err = v.Top()
		}

	} else if strings.HasPrefix(param, "$") {
		key := param[1:]
		v, ok := UserDefSV[key]
		if ok {
			err = errors.New("The" + param + " parameter doesnt have a value set. Please use the \\SET or \\PUSH command to set its value first")
		} else {
			val, err = v.Top()
		}

	} else if strings.HasPrefix(param, "-$") {
		key := param[1:]
		v, ok := NamedParam[key]
		if ok {
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
			val, err = StrToVal(param)
		}
	}
	return
}

/* The StrToVal method converts the input string into a
   value.Value type.
*/
func StrToVal(param string) (val value.Value, err error) {
	fmt.Println("Isha :: " + param)
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
		//if param != true && param != false && param !=null {
		//}
		val = value.NewValue(bytes)
		err = nil
	}

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
		err = errors.New("Stack is Empty. Please use \\PUSH")
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
