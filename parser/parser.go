package parser

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"database/sql"
	"github.com/couchbase/query/shell/cbq/command"
	_ "github.com/couchbaselabs/go_n1ql"
)

const (
	QUERYURL   = command.QUERYURL
	DISCONNECT = command.DISCONNECT
)

//var ShellCmdList = []command.ShellCommand{command.Connect, command, Version{}}

func N1QLCommandParser(line string, n1ql *sql.DB, w io.Writer) error {
	rows, err := n1ql.Query(line)

	if err != nil {
		log.Fatal(err)
	}

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
	Cmd := command.Connect{}
	err := Cmd.ParseCommand(cmd_args[1:])
	return err
}
