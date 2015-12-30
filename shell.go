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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path/filepath"

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
var ServerFlag string

func init() {
	const (
		defaultserver = "http://localhost:8091/"
		usage         = "URL to the query service/cluster. \n\t\t Default : http://localhost:8091\n\n Usage: cbq \n\t\t Connects to local couchbase instance. Same as: cbq -engine=http://localhost:8091\n\t cbq -engine=http://172.23.107.18:8093 \n\t\t Connects to query node at 172.23.107.18 Port 8093 \n\t cbq -engine=https://my.secure.node.com:8093 \n\t\t Connects to query node at my.secure.node.com:8093 using secure https protocol."
	)
	flag.StringVar(&ServerFlag, "engine", defaultserver, usage)
	flag.StringVar(&ServerFlag, "e", defaultserver, "Shorthand for -engine")
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

var (
	SERVICE_URL string
	DISCONNECT  bool
	EXIT        bool
)

func main() {

	flag.Parse()
	command.W = os.Stdout

	if scriptFlag != "" {
		err := execute_input(scriptFlag, os.Stdout)
		if err != nil {
			s_err := handleError(err, ServerFlag)
			s := fmt.Sprintln(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			io.WriteString(command.W, s)
			os.Exit(1)
		}
		os.Exit(0)
	}

	/* Handle options and what they should do */

	// TODO : Readd ...
	//Taken out so as to connect to both cluster and query service
	//using go_n1ql.
	/*
		if strings.HasPrefix(ServerFlag, "http://") == false {
			ServerFlag = "http://" + ServerFlag
		}

		urlRegex := "^(https?://)[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]"
		match, _ := regexp.MatchString(urlRegex, ServerFlag)
		if match == false {
			//TODO Isha : Add error code. Throw invalid url error
			fmt.Println("Invalid url please check" + ServerFlag)
		}


		//-engine
		if strings.HasSuffix(ServerFlag, "/") == false {
			ServerFlag = ServerFlag + "/"
		}
	*/

	/* -quiet : Display Message only if flag not specified
	 */
	if !*quietFlag {
		s := fmt.Sprintln("Connect to " + ServerFlag + ". Type Ctrl-D to exit.\n")
		io.WriteString(command.W, s)
	}

	/* -version : Display the version of the shell and then exit.
	 */
	if versionFlag == true {
		dummy := []string{}
		cmd := command.Version{}
		_ = cmd.ExecCommand(dummy)
		os.Exit(0)
	}

	/* -user : Accept Admin credentials. Prompt for password and set
	   the n1ql_creds. Append to creds so that user can also define
	   bucket credentials using -credentials if they need to.
	*/
	var creds command.Credentials

	if userFlag != "" {
		s := fmt.Sprintln("Enter Password: ")
		io.WriteString(command.W, s)

		password, err := terminal.ReadPassword(0)
		if err == nil {
			if string(password) == "" {
				s_err := handleError(errors.New("Empty password string."), ServerFlag)
				s := fmt.Sprintln(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
				io.WriteString(command.W, s)
				os.Exit(1)
			} else {
				creds = append(creds, command.Credential{"user": userFlag, "pass": string(password)})
			}
		} else {
			s_err := handleError(err, ServerFlag)
			s := fmt.Sprintln(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			io.WriteString(command.W, s)
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
		io.WriteString(command.W, "No Input Credentials. In order to connect to a server with authentication, please provide credentials.")

	} else if credsFlag != "" {

		creds_ret, err := command.ToCreds(credsFlag)
		if err != nil {
			s_err := handleError(err, ServerFlag)
			s := fmt.Sprintln(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			io.WriteString(command.W, s)
		}
		for _, v := range creds_ret {
			creds = append(creds, v)
		}

	}
	//Append empty credentials. This is used for cases where one of the buckets
	//is a SASL bucket, and we need to access the other unprotected buckets.
	//CBauth works this way.

	//if credsFlag == "" && userFlag != "" {
	creds = append(creds, command.Credential{"user": "", "pass": ""})
	//}

	/* Add the credentials set by -user and -credentials to the
	   go_n1ql creds parameter.
	*/
	if creds != nil {
		ac, err := json.Marshal(creds)
		if err != nil {
			//Error while Marshalling
			s_err := handleError(err, ServerFlag)
			s := fmt.Sprintln(fgRed, "ERROR", s_err.Code(), ":", s_err, reset)
			io.WriteString(command.W, s)
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

	go_n1ql.SetPassthroughMode(true)
	//fmt.Println("Input arguments, ", os.Args)
	HandleInteractiveMode(filepath.Base(os.Args[0]))
}
