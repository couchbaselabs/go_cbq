//  Copyright (c) 2015-2016 Couchbase, Inc.
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
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"

	"github.com/couchbase/query/errors"
	"github.com/couchbaselabs/go_cbq/command"
)

/*
This method executes the input command or statement. It
returns an error code and optionally a non empty error message.
*/
func execute_input(line string, w io.Writer) (int, string) {

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
			return errors.NO_SUCH_ALIAS, " : " + commandkey + "\n"
		}

		err_code, err_str := execute_input(val, w)
		/* Error handling for Shell errors and errors recieved from
		   go_n1ql.
		*/
		if err_code != 0 {
			return err_code, err_str
		}

	} else if strings.HasPrefix(line, "\\") {
		//This block handles the shell commands
		err_code, err_str := ExecShellCmd(line)
		if err_code != 0 {
			return err_code, err_str
		}

	} else {
		//This block handles N1QL statements
		// If connected to a query service then NoQueryService == false.
		if NoQueryService == true {
			//Not connected to a query service
			return errors.NO_CONNECTION, ""
		} else {
			/* Try opening a connection to the endpoint. If successful, ping.
			   If successful execute the n1ql command. Else try to connect
			   again.
			*/
			n1ql, err := sql.Open("n1ql", ServerFlag)
			if err != nil {
				return errors.GO_N1QL_OPEN, ""
			} else {
				//Successfully logged into the server
				err_code, err_str := ExecN1QLStmt(line, n1ql, w)
				if err_code != 0 {
					return err_code, err_str
				}
			}

		}
	}

	return 0, ""
}

func WriteHelper(rows *sql.Rows, columns []string, values, valuePtrs []interface{}, rownum int) ([]byte, int, string) {
	//Scan the values into the respective columns
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, errors.ROWS_SCAN, err.Error()
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
			err = json.Unmarshal(b, &parsed)
			if err != nil {
				return nil, errors.JSON_UNMARSHAL, err.Error()
			}

			//Fill up final result object
			dat[col] = parsed

		} else {
			continue
		}

		//Remove one level of nesting for the results when we have only 1 column to project.
		if len(columns) == 1 {
			c, err = dat[col].MarshalJSON()
			if err != nil {
				return nil, errors.JSON_MARSHAL, err.Error()
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
				return nil, errors.JSON_MARSHAL, err.Error()
			}

		}

	} else {
		if len(columns) != 1 {
			b, err = json.Marshal(dat)
			if err != nil {
				return nil, errors.JSON_MARSHAL, err.Error()
			}
		} else {
			b = c
		}

	}

	if *prettyFlag == true {
		var data map[string]interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			return nil, errors.JSON_UNMARSHAL, err.Error()
		}

		b, err = json.MarshalIndent(data, "        ", "    ")
		if err != nil {
			return nil, errors.JSON_MARSHAL, err.Error()
		}
	}

	return b, 0, ""
}

func ExecN1QLStmt(line string, n1ql *sql.DB, w io.Writer) (int, string) {
	//if strings.HasPrefix(strings.ToLower(line), "prepare") {

	rows, err := n1ql.Query(line)

	if err != nil {
		return errors.GON1QL_QUERY, err.Error()

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

				extras, err_code, err_string := WriteHelper(rows, columns, values, valuePtrs, rownum)

				if extras == nil && err_code != 0 {
					return err_code, err_string
				}

				var dat map[string]interface{}

				if err := json.Unmarshal(extras, &dat); err != nil {
					return errors.JSON_UNMARSHAL, err.Error()
				}

				_, werr = io.WriteString(w, "    \"requestID\": \""+dat["requestID"].(string)+"\",\n")

				jsonString, err := json.MarshalIndent(dat["signature"], "        ", "    ")

				if err != nil {
					return errors.JSON_MARSHAL, err.Error()
				}
				_, werr = io.WriteString(w, "    \"signature\": "+string(jsonString)+",\n")
				_, werr = io.WriteString(w, "    \"results\" : [\n\t")
				status = dat["status"].(string)
				rownum++
				continue
			}

			if rownum == 1 {

				// Get the second row to post process as the metrics

				var err_code int
				var err_string string
				metrics, err_code, err_string = WriteHelper(rows, columns, values, valuePtrs, rownum)

				if metrics == nil && err_code != 0 {
					return err_code, err_string
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

			result, err_code, err_string := WriteHelper(rows, columns, values, valuePtrs, rownum)
			if result == nil && err_code != 0 {
				return err_code, err_string
			}

			_, werr = io.WriteString(w, string(result))

		}

		err = rows.Close()
		if err != nil {
			return errors.ROWS_CLOSE, err.Error()
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
			return errors.WRITER_OUTPUT, werr.Error()
		}
	}

	return 0, ""
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

func ExecShellCmd(line string) (int, string) {

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
			return errors.UNBALANCED_PAREN, ""
		}

	}

	cmd_args := strings.Split(line, " ")

	//Lookup Command from function registry

	Cmd, ok := command.COMMAND_LIST[cmd_args[0]]
	if ok == true {
		err_code, err_str := Cmd.ExecCommand(cmd_args[1:])
		if err_code != 0 {
			return err_code, err_str
		}
	} else {
		return errors.NO_SUCH_COMMAND, ""
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
	return 0, ""
}
