package command

import (
	"fmt"
)

/* Connect Command */
type Connect struct {
	ShellCommand
}

func (this *Connect) Name() string {
	return "CONNECT"
}

func (this *Connect) CommandCompletion() bool {
	return false
}

func (this *Connect) MinArgs() int {
	return 0
}

func (this *Connect) MaxArgs() int {
	return 1
}

func (this *Connect) ParseCommand(queryurl []string) error {

	QUERYURL = queryurl[0]
	fmt.Println("\nCouchbase query shell connected to " + QUERYURL + " . Type Ctrl-D / \\exit / \\quit to exit.")
	return nil
}

func (this *Connect) PrintHelp() {
	fmt.Println("\\CONNECT <url>")
	fmt.Println()
	fmt.Println("Connect to the query service or cluster endpoint url.")
	fmt.Println("\n\t\t Default : http://localhost:8091 \n\t\t \\CONNECT https://my.secure.node.com:8093 ; \n\t\t Connects to query node at my.secure.node.com:8093 using secure https protocol.")
	fmt.Println(" For Example : \n\t \\CONNECT http://172.6.23.2:8091 ;")
	fmt.Println()
}
