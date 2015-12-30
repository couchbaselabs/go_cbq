//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"

	"github.com/couchbaselabs/go_cbq/command"
)

func execute_input(line string, w io.Writer) error {

	command.W = w

	if DISCONNECT == true || NoQueryService == true {
		if strings.HasPrefix(strings.ToLower(line), "\\connect") {
			NoQueryService = false
			command.DISCONNECT = false
			DISCONNECT = false
		}
	}

	if strings.HasPrefix(line, "\\\\") {
		// This block handles aliases
		commandkey := line[2:]
		commandkey = strings.TrimSpace(commandkey)

		val, ok := command.AliasCommand[commandkey]

		if !ok {
			return errors.New("Alias " + commandkey + " doesnt exist." + val + "\n")
		}

		err := execute_input(val, w)
		/* Error handling for Shell errors and errors recieved from
		   go_n1ql.
		*/
		if err != nil {
			return err
		}

	} else if strings.HasPrefix(line, "\\") {
		//This block handles the shell commands
		err := ExecShellCmd(line)
		if err != nil {
			return err
		}

	} else {
		//This block handles N1QL statements
		// If connected to a query service then NoQueryService == false.
		if NoQueryService == true {
			//Not connected to a query service
			err := errors.New("Not connected to any instance. Use \\CONNECT shell command to connect to an instance.")
			return err
		} else {
			/* Try opening a connection to the endpoint. If successful, ping.
			   If successful execute the n1ql command. Else try to connect
			   again.
			*/
			n1ql, err := sql.Open("n1ql", ServerFlag)
			if err != nil {
				tmpstr := fmt.Sprintln(fgRed, "Error in sql Open", reset)
				io.WriteString(w, tmpstr)
				return err
			} else {
				//Successfully logged into the server
				err := ExecN1QLStmt(line, n1ql, w)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func WriteHelper(rows *sql.Rows, columns []string, values, valuePtrs []interface{}, rownum int) ([]byte, error) {
	//Scan the values into the respective columns
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	dat := map[string]*json.RawMessage{}
	var c []byte = nil
	var b []byte = nil
	var err error = nil

	for i, col := range columns {
		var parsed *json.RawMessage

		val := values[i]
		b, _ := val.([]byte)

		if string(b) != "" {
			//Parse the sub values of the main map first.
			json.Unmarshal(b, &parsed)

			//Fill up final result object
			dat[col] = parsed

		} else {
			continue
		}

		//Remove one level of nesting for the results when we have only 1 column to project.
		if len(columns) == 1 {
			c, err = dat[col].MarshalJSON()
			if err != nil {
				return nil, err
			}
		}

	}

	b = nil
	err = nil

	// The first and second row represent the metadata. Because of the
	// way the rows are returned we need to create a map with the
	// correct data.
	if rownum == 0 || rownum == 1 {
		keys := make([]string, 0, len(dat))
		for key, _ := range dat {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		if keys != nil {
			map_value := dat[keys[0]]
			b, err = map_value.MarshalJSON()
			if err != nil {
				return nil, err
			}

		}

	} else {
		if len(columns) != 1 {
			b, err = json.Marshal(dat)
			if err != nil {
				return nil, err
			}
		} else {
			b = c
		}

	}

	if *prettyFlag == true {
		var data map[string]interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, err
		}

		b, err = json.MarshalIndent(data, "        ", "    ")
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func ExecN1QLStmt(line string, n1ql *sql.DB, w io.Writer) error {
	//if strings.HasPrefix(strings.ToLower(line), "prepare") {

	rows, err := n1ql.Query(line)

	if err != nil {
		return err

	} else {
		iter := 0
		rownum := 0

		var werr error
		status := ""
		var metrics []byte
		metrics = nil

		// Multi column projection
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		//Check if spacing is enough
		_, werr = io.WriteString(w, "\n{\n")

		for rows.Next() {

			for i, _ := range columns {
				valuePtrs[i] = &values[i]
			}

			if rownum == 0 {

				// Get the first row to post process.

				extras, err := WriteHelper(rows, columns, values, valuePtrs, rownum)

				if extras == nil && err != nil {
					return err
				}

				var dat map[string]interface{}

				if err := json.Unmarshal(extras, &dat); err != nil {
					return err
				}

				_, werr = io.WriteString(w, "    \"requestID\": \""+dat["requestID"].(string)+"\",\n")

				jsonString, err := json.MarshalIndent(dat["signature"], "        ", "    ")

				if err != nil {
					return err
				}
				_, werr = io.WriteString(w, "    \"signature\": "+string(jsonString)+",\n")
				_, werr = io.WriteString(w, "    \"results\" : [\n\t")
				status = dat["status"].(string)
				rownum++
				continue
			}

			if rownum == 1 {

				// Get the second row to post process as the metrics
				var err error
				metrics, err = WriteHelper(rows, columns, values, valuePtrs, rownum)

				if metrics == nil && err != nil {
					return err
				}

				//Wait until all the rows have been written to write the metrics.
				rownum++
				continue
			}

			if iter == 0 {
				iter++
			} else {
				_, werr = io.WriteString(w, ", \n\t")
			}

			result, err := WriteHelper(rows, columns, values, valuePtrs, rownum)
			if result == nil && err != nil {
				return err
			}

			_, werr = io.WriteString(w, string(result))

		}

		err = rows.Close()
		if err != nil {
			return err
		}

		//Suffix to result array
		_, werr = io.WriteString(w, "\n    ],")

		//Write the status and the metrics
		if status != "" {
			_, werr = io.WriteString(w, "\n    \"status\": \""+status+"\"")
		}
		if metrics != nil {
			_, werr = io.WriteString(w, ",\n    \"metrics\": ")
			_, werr = io.WriteString(w, string(metrics))
		}

		_, werr = io.WriteString(w, "\n}\n")

		// For any captured write error
		if werr != nil {
			return err
		}
	}

	return nil
}

/* From
http://intogooglego.blogspot.com/2015/05/day-6-string-minifier-remove-whitespaces.html
*/
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

func ExecShellCmd(line string) error {

	arg1 := strings.Split(line, " ")
	arg1str := strings.ToLower(arg1[0])

	line = arg1str + " " + strings.Join(arg1[1:], " ")
	line = strings.TrimSpace(line)

	line = stringMinifier(line)

	// Handle input strings to \echo command.
	if strings.HasPrefix(line, "\\echo") {

		count_param := strings.Count(line, "\"")
		count_param_bs := strings.Count(line, "\\\"")

		if count_param%2 == 0 && count_param_bs%2 == 0 {
			r := strings.NewReplacer("\\\"", "\\\"", "\"", "")
			line = r.Replace(line)

		} else {
			return errors.New("Unbalanced Paranthesis in input.")
		}

	}

	cmd_args := strings.Split(line, " ")

	//Lookup Command from function registry
	var err error

	Cmd, ok := command.COMMAND_LIST[cmd_args[0]]
	if ok == true {
		err = Cmd.ExecCommand(cmd_args[1:])
		if err != nil {
			return err
		}
	} else {
		return errors.New("Command doesnt exist. Use help for command help.")
	}

	if (strings.HasPrefix(line, "\\source") || strings.HasPrefix(line, "\\load")) &&
		command.FILE_INPUT == true {

		fmt.Println("ISHA DEBUG : FILENAME ", command.FILE_INPUT)
	}

	SERVICE_URL = command.SERVICE_URL

	if SERVICE_URL != "" {
		ServerFlag = SERVICE_URL
		command.SERVICE_URL = ""
		SERVICE_URL = ""
	}

	DISCONNECT = command.DISCONNECT
	if DISCONNECT == true {
		NoQueryService = true

	}

	EXIT = command.EXIT
	return nil
}
