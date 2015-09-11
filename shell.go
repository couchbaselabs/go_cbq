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
	"strings"
	"time"
	//"regexp"

	"github.com/couchbaselabs/go_cbq/command"
	_ "github.com/couchbaselabs/go_n1ql"
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
   Default value : "2ms"
   Query timeout parameter.
*/

var timeoutFlag time.Duration

func init() {
	const (
		defaultval = 0 * time.Minute
		usage      = "Query timeout parameter. Units are mandatory. \n\t\t Default : \"2ms\" \n\t\t Valid Units : ns (nanoseconds), us (microseconds), ms (milliseconds), s (seconds), m (minutes), h (hours) "
	)
	flag.DurationVar(&timeoutFlag, "timeout", defaultval, usage)
	flag.DurationVar(&timeoutFlag, "t", defaultval, " Shorthand for -timeout")
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

var prettyFlag = flag.Bool("pretty", false, "Pretty print the output.")

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
	flag.StringVar(&scriptFlag, "file", defaultval, usage)
	flag.StringVar(&scriptFlag, "f", defaultval, " Shorthand for -file")

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
	flag.StringVar(&scriptFlag, "output", defaultval, usage)
	flag.StringVar(&scriptFlag, "o", defaultval, " Shorthand for -output")

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
	flag.StringVar(&scriptFlag, "log-file", defaultval, usage)
	flag.StringVar(&scriptFlag, "l", defaultval, " Shorthand for -log-file")

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

		os.Setenv("n1ql_creds", string(ac))
	}

	HandleInteractiveMode(filepath.Base(os.Args[0]))
}

func execute_query(line string, w io.Writer) error {

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

	if strings.HasPrefix(line, "\\") == true {
		err := ShellCommandParser(line)
		if err != nil {
			s_err := handleError(err, TiServer)
			fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
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
				s_err := handleError(err, TiServer)
				fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
				fmt.Println(fgRed, "Error in sql Open", reset)
			} else {
				//Successfully logged into the server
				err = n1ql.Ping()
				if err != nil {
					s_err := handleError(err, TiServer)
					fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
					fmt.Println(fgRed, "Error in sql Ping", reset)
				} else {
					//fSuccessfully Pinged

					err := N1QLCommandParser(line, n1ql, w)
					if err != nil {
						s_err := handleError(err, TiServer)
						fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
					}
				}
			}

		} else {
			//Not connected to a query service
			err := errors.New("Not connected to any instance. Use \\CONNECT shell command to connect to an instance.")
			s_err := handleError(err, TiServer)
			fmt.Println(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
		}
	}

	return nil
}

func N1QLCommandParser(line string, n1ql *sql.DB, w io.Writer) error {
	if strings.HasPrefix(strings.ToLower(line), "create") {
		_, err := n1ql.Exec(line)
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
				var dat map[string]interface{}
				if err := json.Unmarshal(b, &dat); err != nil {
					return err
				}
				if *prettyFlag {
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

func ShellCommandParser(line string) error {

	line = strings.ToLower(line)
	line = strings.TrimSpace(line)

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
