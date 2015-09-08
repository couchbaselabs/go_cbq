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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/couchbase/query/shell/cbq/command"
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
		usage = "Shell Version \n\t -version"
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
		usage      = "Single command mode. Execute input command and exit shell. \n\t -script=\"select * from system:keyspaces\""
	)
	flag.StringVar(&scriptFlag, "script", defaultval, usage)
	flag.StringVar(&scriptFlag, "s", defaultval, " Shorthand for -script")

}

/*
   Option        : -prompt
   Args          : <character>
   Default value : ">"
   Query prompt character.
*/

var promptFlag = flag.String("prompt", ">", "Set the character for the command prompt")

/*
   Option        : -pretty
   Default value : false
   Pretty print output
*/

var prettyFlag = flag.Bool("pretty", false, "Set the character for the command prompt")

/*
   Option        : -file or -f
   Args          : <filename>
   Input file to run queries from. Exit after the queries are run.
*/

/*
   Option        : -ouput or -o
   Args          : <filename>
   Output file to send results of queries to.
*/

var outputFlag string

func init() {
	const (
		defaultval = ""
		usage      = "File to output commands and their results. \n\t -output=temp.txt"
	)
	flag.StringVar(&scriptFlag, "output", defaultval, usage)
	flag.StringVar(&scriptFlag, "o", defaultval, " Shorthand for -output")

}

/*
   Option        : -log-file or -l
   Args          : <filename>
   Log commands for session.
*/

type Credentials map[string]string
type MyCred []Credentials

var in int

var (
	QUERYURL   string
	DISCONNECT bool
)

func main() {

	//QUERYURL = command.QUERYURL

	flag.Parse()

	// Check if the Connect command has been issued and
	// what the new url to connect to is.

	/*if QUERYURL != "" {
		TiServer = QUERYURL
	}*/

	fmt.Println("Came in here no : ", in)
	in = in + 1

	/* Handle options and what they should do */

	if strings.Contains(TiServer, ":8091") == false &&
		strings.Contains(TiServer, ":9000") == false &&
		strings.Contains(TiServer, ":8093") == false &&
		strings.Contains(TiServer, ":9499") == false {

		if strings.HasSuffix(TiServer, ":") == false {
			TiServer = TiServer + ":8091"
		} else {
			TiServer = TiServer + "8091"
		}

	}

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

	//-quiet
	if !*quietFlag {
		fmt.Printf("Couchbase query shell connected to %v . Type Ctrl-D to exit.\n", TiServer)
	}

	//-credentials
	if credsFlag == "" {
		// Isha TODO : Check if credentials exist and then appropriately throw error.
		fmt.Println("Empty " + credsFlag)
	} else {
		//Handle the input string of credentials.
		//The string needs to be parsed into a byte array so as to pass to go_n1ql.
		fmt.Println("Input credentials string " + credsFlag)

		cred := strings.Split(credsFlag, ",")
		//fmt.Println(cred)
		var creds MyCred

		for _, i := range cred {
			up := strings.Split(i, ":")
			//fmt.Println(up)

			creds = append(creds, Credentials{"user": up[0], "pass": up[1]})
		}
		ac, err := json.Marshal(creds)
		//fmt.Println(string(ac))
		//fmt.Println(creds)
		if err != nil {
			fmt.Println(err)
			//	return
		}

		os.Setenv("n1ql_creds", string(ac))
	}

	fmt.Println("Main tiserver : " + TiServer)

	HandleInteractiveMode(filepath.Base(os.Args[0]))
}

func execute_query(line string, w io.Writer) error {

	fmt.Println("This is the execute query server tiserver : " + TiServer)
	fmt.Println("This is the parser QUERYURL : " + QUERYURL)

	fmt.Println("IMP DISCONNECT : ", DISCONNECT, " IMP Noqueryservice :", NoQueryService)

	if DISCONNECT == true {
		fmt.Println("LINE : "+line+"   NoQ: ", NoQueryService)
		if strings.HasPrefix(strings.ToLower(line), "\\connect") {
			NoQueryService = false
			command.DISCONNECT = false
			DISCONNECT = false
		}
	}

	if !NoQueryService {

		fmt.Println(NoQueryService)
		n1ql, err := sql.Open("n1ql", TiServer)
		if err != nil {
			//log.Fatal(err)
			fmt.Println("Error in sql Open")
		} else {
			fmt.Println("successfully logged")
		}

		err = n1ql.Ping()
		if err != nil {
			//log.Fatal(err)
			fmt.Println("Error in sql Ping")
		} else {
			fmt.Println("successfully pinged")
		}

		// Set query parameters
		//fmt.Println("Timeout value :" + timeoutFlag.String())
		//os.Setenv("n1ql_timeout", timeoutFlag.String())
		if strings.HasPrefix(line, "\\") == true {
			err = ShellCommandParser(line)

		} else {
			err = N1QLCommandParser(line, n1ql, w)
		}
		if err != nil {
			//log.Fatal(err)
			fmt.Println("Error in N1QLCOmmandParser or ShellCommandParser")
		}

	} else {
		//Isha TODO : Add this to the error handling. This is temporary until the \CONNECT command is implemented.
		io.WriteString(w, "\nNot connected to any instance. Use \\CONNECT shell command to connect to an instance.\n")
	}
	return nil
}

func N1QLCommandParser(line string, n1ql *sql.DB, w io.Writer) error {
	rows, err := n1ql.Query(line)

	if err != nil {
		//log.Fatal(err)
		fmt.Println("Error in n1ql.Query")
		log.Print("Weeeee  ", err)

	} else {

		defer rows.Close()
		iter := 0
		for rows.Next() {
			var results *json.RawMessage
			//var results string
			if iter == 0 {
				iter++
			} else {
				io.WriteString(w, ",")
				io.WriteString(w, "\n")
			}
			if err := rows.Scan(&results); err != nil {
				return err
			}
			b, _ := results.MarshalJSON()
			var dat map[string]interface{}
			if err := json.Unmarshal(b, &dat); err != nil {
				return err
			}
			c, err := json.MarshalIndent(dat, "", "  ")
			if err != nil {
				return err
			}
			//reader := bytes.NewReader(c)
			//io.Copy(w, reader)
			io.WriteString(w, string(c))
			//io.WriteString(w, ",")
		}
		io.WriteString(w, "\n")
	}

	return nil
}

func ShellCommandParser(line string) error {

	line = strings.ToLower(line)

	/*if strings.HasPrefix(line, "\\connect") {
			cmd_args := strings.Split(line, " ")
			&Connect{}.ParseCommand(cmd_args)
	}
	if strings.HasPrefix(line, "\\disconnect") {

	}*/

	cmd_args := strings.Split(line, " ")

	//Lookup Command from function registry
	//Cmd = cmd_args[0]

	var err error

	if strings.HasPrefix(line, "\\connect") {
		Cmd := command.Connect{}
		err = Cmd.ParseCommand(cmd_args[1:])

	} else if strings.HasPrefix(line, "\\disconnect") {
		Cmd := command.Disconnect{}
		err = Cmd.ParseCommand(nil)

	} else if strings.HasPrefix(line, "\\help") {
		Cmd := command.Help{}
		err = Cmd.ParseCommand(cmd_args[1:])

	} else if strings.HasPrefix(line, "\\version") {
		Cmd := command.Version{}
		err = Cmd.ParseCommand()

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

	return err
}
