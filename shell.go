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
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"

	"os"
	"path/filepath"
	//"reflect"
	"strings"
	"unicode"
	//"regexp"

	"github.com/couchbaselabs/go_cbq/command"
	go_n1ql "github.com/couchbaselabs/go_n1ql"
)

/*
   Command line options provided.
*/

/*
   Option        : -engine or -e
   Args          :  <url to the query service or to the cluster>
   Default value : http://localhost:8091/
   Point to the cluser/query endpoint to connect to.
*/
var TiServer string

func init() {
	const (
		defaultserver = "http://localhost:8091/"
		usage         = "URL to the query service/cluster. \n\t\t Default : http://localhost:8091\n\n Usage: cbq \n\t\t Connects to local couchbase instance. Same as: cbq -engine=http://localhost:8091\n\t cbq -engine=http://172.23.107.18:8093 \n\t\t Connects to query node at 172.23.107.18 Port 8093 \n\t cbq -engine=https://my.secure.node.com:8093 \n\t\t Connects to query node at my.secure.node.com:8093 using secure https protocol."
	)
	flag.StringVar(&TiServer, "engine", defaultserver, usage)
	flag.StringVar(&TiServer, "e", defaultserver, "Shorthand for -engine")
}

/*
   Option        : -no-engine or -ne
   Default value : false
   Enable/Disable startup connection to a query service/cluster endpoint.
*/
var NoQueryService bool

func init() {
	const (
		defaultval = false
		usage      = "Start shell without connecting to a query service/cluster endpoint. \n\t\t Default : false \n\t\t Possible Values : true/false"
	)
	flag.BoolVar(&NoQueryService, "no-engine", defaultval, usage)
	flag.BoolVar(&NoQueryService, "ne", defaultval, " Shorthand for -no-engine")
}

/*
   Option        : -quiet
   Default value : false
   Enable/Disable startup connection message for the shell
*/
var quietFlag = flag.Bool("quiet", false, "Enable/Disable startup connection message for the shell \n\t\t Default : false \n\t\t Possible Values : true/false")

/*
   Option        : -timeout or -t
   Args          : <timeout value>
   Default value : "0ms"
   Query timeout parameter.
*/

var timeoutFlag string

func init() {
	const (
		defaultval = ""
		usage      = "Query timeout parameter. Units are mandatory. For Example : \"10ms\". \n\t\t Valid Units : ns (nanoseconds), us (microseconds), ms (milliseconds), s (seconds), m (minutes), h (hours) "
	)
	flag.StringVar(&timeoutFlag, "timeout", defaultval, usage)
	flag.StringVar(&timeoutFlag, "t", defaultval, " Shorthand for -timeout")
}

/*
   Option        : -user or -u
   Args          : Login username
   Login credentials for users. The shell will prompt for the password.
*/

var userFlag string

func init() {
	const (
		defaultval = ""
		usage      = "Username \n\t For Example : -u=Administrator"
	)
	flag.StringVar(&userFlag, "user", defaultval, usage)
	flag.StringVar(&userFlag, "u", defaultval, " Shorthand for -credentials")

}

/*
   Option        : -credentials or -c
   Args          : A list of credentials, in the form of user/password objects.
   Login credentials for users as well as SASL Buckets.
*/

var credsFlag string

func init() {
	const (
		defaultval = ""
		usage      = "A list of credentials, in the form user:password. \n\t For Example : Administrator:password, beer-sample:asdasd"
	)
	flag.StringVar(&credsFlag, "credentials", defaultval, usage)
	flag.StringVar(&credsFlag, "c", defaultval, " Shorthand for -credentials")

}

/*
   Option        : -version or -v
   Shell Version
*/

var versionFlag bool

func init() {
	const (
		usage = "Shell Version \n\t Usage: -version"
	)
	flag.BoolVar(&versionFlag, "version", false, usage)
	flag.BoolVar(&versionFlag, "v", false, "Shorthand for -version")

}

/*
   Option        : -script or -s
   Args          : <query>
   Single command mode
*/

var scriptFlag string

func init() {
	const (
		defaultval = ""
		usage      = "Single command mode. Execute input command and exit shell. \n\t For Example : -script=\"select * from system:keyspaces\""
	)
	flag.StringVar(&scriptFlag, "script", defaultval, usage)
	flag.StringVar(&scriptFlag, "s", defaultval, " Shorthand for -script")

}

/*
   Option        : -pretty
   Default value : false
   Pretty print output
*/

var prettyFlag = flag.Bool("pretty", true, "Pretty print the output.")

/*
   Option        : -exit-on-error
   Default value : false
   Exit shell after first error encountered.
*/

var errorExitFlag = flag.Bool("exit-on-error", false, "Exit shell after first error encountered.")

/*
   Option        : -file or -f
   Args          : <filename>
   Input file to run queries from. Exit after the queries are run.
*/

var inputFlag string

func init() {
	const (
		defaultval = ""
		usage      = "File to load commands from. \n\t For Example : -file=temp.txt"
	)
	flag.StringVar(&inputFlag, "file", defaultval, usage)
	flag.StringVar(&inputFlag, "f", defaultval, " Shorthand for -file")

}

/*
   Option        : -ouput or -o
   Args          : <filename>
   Output file to send results of queries to.
*/

var outputFlag string

func init() {
	const (
		defaultval = ""
		usage      = "File to output commands and their results. \n\t For Example : -output=temp.txt"
	)
	flag.StringVar(&outputFlag, "output", defaultval, usage)
	flag.StringVar(&outputFlag, "o", defaultval, " Shorthand for -output")

}

/*
   Option        : -log-file or -l
   Args          : <filename>
   Log commands for session.
*/

var logFlag string

func init() {
	const (
		defaultval = ""
		usage      = "File to log commands into. \n\t For Example : -log-file=temp.txt"
	)
	flag.StringVar(&logFlag, "log-file", defaultval, usage)
	flag.StringVar(&logFlag, "l", defaultval, " Shorthand for -log-file")

}

/* Define credentials as user/pass and convert into
   JSON object credentials
*/
type Credentials map[string]string
type MyCred []Credentials

var (
	QUERYURL   string
	DISCONNECT bool
	EXIT       bool
)

func main() {

	flag.Parse()

	if scriptFlag != "" {
		err := execute_query(scriptFlag, os.Stdout)
		if err != nil {
			s_err := handleError(err, TiServer)
			fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
		}
		os.Exit(1)
	}

	/* Handle options and what they should do */

	// TODO : Readd ...
	//Taken out so as to connect to both cluster and query service
	//using go_n1ql.
	/*
		if strings.HasPrefix(TiServer, "http://") == false {
			TiServer = "http://" + TiServer
		}

		urlRegex := "^(https?://)[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]"
		match, _ := regexp.MatchString(urlRegex, TiServer)
		if match == false {
			//TODO Isha : Add error code. Throw invalid url error
			fmt.Println("Invalid url please check" + TiServer)
		}


		//-engine
		if strings.HasSuffix(TiServer, "/") == false {
			TiServer = TiServer + "/"
		}
	*/

	/* -quiet : Display Message only if flag not specified
	 */
	if !*quietFlag {
		fmt.Printf("Couchbase query shell connected to %v . Type Ctrl-D to exit.\n", TiServer)
	}

	/* -version : Display the version of the shell and then exit.
	 */
	if versionFlag == true {
		dummy := []string{}
		cmd := command.Version{}
		_ = cmd.ParseCommand(dummy)
		os.Exit(0)
	}

	/* -user : Accept Admin credentials. Prompt for password and set
	   the n1ql_creds. Append to creds so that user can also define
	   bucket credentials using -credentials if they need to.
	*/
	var creds MyCred

	if userFlag != "" {
		fmt.Println("Enter Password: ")
		password, err := terminal.ReadPassword(0)
		if err == nil {
			if string(password) == "" {
				s_err := handleError(errors.New("Endered empty password string."), TiServer)
				fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
				os.Exit(1)
			} else {
				creds = append(creds, Credentials{"user": userFlag, "pass": string(password)})
			}
		} else {
			s_err := handleError(err, TiServer)
			fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			os.Exit(1)
		}
	}

	/* -credentials : Accept credentials to pass to the n1ql endpoint.
	   Ensure that the user inputs credentials in the form a:b.
	   It is important to apend these credentials to those given by
	   -user.
	*/
	if userFlag == "" && credsFlag == "" {
		/* No credentials exist. This can still be used to connect to
		   un-authenticated servers.
		*/
		fmt.Println("No Input Credentials. In order to connect to a server with authentication, please provide credentials.")

	} else if credsFlag != "" {
		//Handle the input string of credentials.
		//The string needs to be parsed into a byte array so as to pass to go_n1ql.
		cred := strings.Split(credsFlag, ",")

		/* Append input credentials in [{"user": <username>, "pass" : <password>}]
		   format as expected by go_n1ql creds.
		*/
		for _, i := range cred {
			up := strings.Split(i, ":")
			if len(up) < 2 {
				// One of the input credentials is incorrect
				err := errors.New("Username or Password missing in -credentials/-c option. Please check")

				s_err := handleError(err, TiServer)
				fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
				os.Exit(1)

			} else {
				creds = append(creds, Credentials{"user": up[0], "pass": up[1]})
			}
		}
	}

	/* Add the credentials set by -user and -credentials to the
	   go_n1ql creds parameter.
	*/
	if creds != nil {
		ac, err := json.Marshal(creds)
		if err != nil {
			//Error while Marshalling
			s_err := handleError(err, TiServer)
			fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			os.Exit(1)
		}
		go_n1ql.SetQueryParams("creds", string(ac))
	}

	if timeoutFlag != "0ms" {
		go_n1ql.SetQueryParams("timeout", timeoutFlag)
	}

	if inputFlag != "" {
		//Read each line from the file and call execute query

	}

	//fmt.Println("Input arguments, ", os.Args)
	HandleInteractiveMode(filepath.Base(os.Args[0]))
}

func execute_query(line string, w io.Writer) error {

	command.W = w

	if DISCONNECT == true || NoQueryService == true {
		if strings.HasPrefix(strings.ToLower(line), "\\connect") {
			NoQueryService = false
			command.DISCONNECT = false
			DISCONNECT = false
		}
	}

	// Set query parameters
	//fmt.Println("Timeout value :" + timeoutFlag.String())
	//os.Setenv("n1ql_timeout", timeoutFlag.String())

	if strings.HasPrefix(line, "\\\\") {
		commandkey := line[2:]
		commandkey = strings.TrimSpace(commandkey)
		//commandkey = commandkey[0]

		//fmt.Println("Alias: ", commandkey)
		//fmt.Println("Alias: ", reflect.TypeOf(commandkey))

		val, ok := command.AliasCommand[commandkey]

		if !ok {
			return errors.New("\nAlias " + commandkey + " doesnt exist." + val + "\n")
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

func N1QLCommandParser(line string, n1ql *sql.DB, w io.Writer) error {
	if strings.HasPrefix(strings.ToLower(line), "create") {
		_, err := n1ql.Query(line)
		if err != nil {
			return err
		}
	} else {

		rows, err := n1ql.Query(line)

		if err != nil {
			return err

		} else {
			iter := 0

			var werr error
			_, werr = io.WriteString(w, "\n \"results\" :  [ ")

			for rows.Next() {
				var results *json.RawMessage

				if iter == 0 {
					iter++
				} else {
					_, werr = io.WriteString(w, ", \n")
				}
				if err := rows.Scan(&results); err != nil {
					return err
				}
				b, err := results.MarshalJSON()
				if err != nil {
					return err
				}

				if *prettyFlag == true {
					var dat map[string]interface{}
					if err := json.Unmarshal(b, &dat); err != nil {
						return err
					}

					b, err = json.MarshalIndent(dat, "", "  ")
					if err != nil {
						return err
					}
				}

				_, werr = io.WriteString(w, string(b))
			}
			err = rows.Close()
			if err != nil {
				return err
			}

			_, werr = io.WriteString(w, " ] \n")

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
