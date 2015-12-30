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

	"github.com/couchbase/query/value"
)

/* Helper function to create a stack. */
func Stack_Helper() *Stack {
	r := make(Stack, 0)
	return &r
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
