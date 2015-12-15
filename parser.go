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
	"strings"
	"unicode"

	"github.com/couchbaselabs/go_cbq/command"
)

func execute_query(line string, w io.Writer) error {

	command.W = w

	if DISCONNECT == true || NoQueryService == true {
		if strings.HasPrefix(strings.ToLower(line), "\\connect") {
			NoQueryService = false
			command.DISCONNECT = false
			DISCONNECT = false
		}
	}

	if strings.HasPrefix(line, "\\\\") {
		commandkey := line[2:]
		commandkey = strings.TrimSpace(commandkey)

		val, ok := command.AliasCommand[commandkey]

		if !ok {
			return errors.New("Alias " + commandkey + " doesnt exist." + val + "\n")
		}

		err := execute_query(val, w)
		/* Error handling for Shell errors and errors recieved from
		   go_n1ql.
		*/
		if err != nil {
			return err
		}

	} else if strings.HasPrefix(line, "\\") {
		err := ShellCommandParser(line)
		if err != nil {
			return err
		}

	} else {
		// If connected to a query service then NoQueryService == false.
		if !NoQueryService {
			/* Try opening a connection to the endpoint. If successful, ping.
			   If successful execute the n1ql command. Else try to connect
			   again.
			*/
			n1ql, err := sql.Open("n1ql", TiServer)
			if err != nil {
				fmt.Println(fgRed, "Error in sql Open", reset)
				return err
			} else {
				//Successfully logged into the server
				err = n1ql.Ping()
				if err != nil {
					fmt.Println(fgRed, "Error in sql Ping", reset)
					fmt.Println(fgRed, err, reset)
					return err

				} else {
					//fSuccessfully Pinged

					err := N1QLCommandParser(line, n1ql, w)
					if err != nil {
						return err
					}
				}
			}

		} else {
			//Not connected to a query service
			err := errors.New("Not connected to any instance. Use \\CONNECT shell command to connect to an instance.")
			return err
		}
	}

	return nil
}

func WriteHelper(rows *sql.Rows) ([]byte, error) {
	var results *json.RawMessage
	if err := rows.Scan(&results); err != nil {
		return nil, err
	}
	b, err := results.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if *prettyFlag == true {
		var dat map[string]interface{}
		if err := json.Unmarshal(b, &dat); err != nil {
			return nil, err
		}

		b, err = json.MarshalIndent(dat, "", "  ")
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func N1QLCommandParser(line string, n1ql *sql.DB, w io.Writer) error {
	if strings.HasPrefix(strings.ToLower(line), "prepare") {
		_, err := n1ql.Query(line)
		fmt.Println("Im in here")
		if err != nil {
			return err
		}
	} else {

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

			//Check if spacing is enough
			_, werr = io.WriteString(w, "\n{\n")

			for rows.Next() {

				if rownum == 0 {
					rownum++

					// Get the first row to post process.

					extras, err := WriteHelper(rows)

					if extras == nil && err != nil {
						return err
					}

					var dat map[string]interface{}

					if err := json.Unmarshal(extras, &dat); err != nil {
						panic(err)
					}

					_, werr = io.WriteString(w, "\"requestID\": \""+dat["requestID"].(string)+"\",\n")

					jsonString, err := json.MarshalIndent(dat["signature"], "", "  ")
					if err != nil {
						return err
					}
					_, werr = io.WriteString(w, "\"signature\": "+string(jsonString)+",\n")
					_, werr = io.WriteString(w, "\"results\" : [\n\t")
					status = dat["status"].(string)
					continue
				}

				if rownum == 1 {
					rownum++

					// Get the second row to post process as the metrics
					var err error
					metrics, err = WriteHelper(rows)

					if metrics == nil && err != nil {
						return err
					}

					//Wait until all the rows have been written to write the metrics.
					continue
				}

				if iter == 0 {
					iter++
				} else {
					_, werr = io.WriteString(w, ", \n\t")
				}

				result, err := WriteHelper(rows)
				if result == nil && err != nil {
					return err
				}

				_, werr = io.WriteString(w, "\t"+string(result))

			}
			err = rows.Close()
			if err != nil {
				return err
			}

			_, werr = io.WriteString(w, "\n\t],\n\t")
			//Write the status and the metrics
			if status != "" {
				_, werr = io.WriteString(w, "\n\"status\": "+status)
			}
			if metrics != nil {
				_, werr = io.WriteString(w, ",\n\"metrics\": ")
				_, werr = io.WriteString(w, "\t"+string(metrics))
			}

			_, werr = io.WriteString(w, "\n}\n\n")

			// For any captured write error
			if werr != nil {
				return err
			}
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

func ShellCommandParser(line string) error {

	line = strings.ToLower(line)
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
		err = Cmd.ParseCommand(cmd_args[1:])
		if err != nil {
			return err
		}
	} else {
		return errors.New("Command doesnt exist. Use help for command help.")
	}

	if (strings.HasPrefix(line, "\\source") || strings.HasPrefix(line, "\\load")) &&
		command.FILEINPUT == true {

		fmt.Println("ISHA DEBUG : FILENAME ", command.FILEINPUT)
	}

	QUERYURL = command.QUERYURL

	if QUERYURL != "" {
		TiServer = QUERYURL
		command.QUERYURL = ""
		QUERYURL = ""
	}

	DISCONNECT = command.DISCONNECT
	if DISCONNECT == true {
		NoQueryService = true

	}

	EXIT = command.EXIT
	return nil
}
