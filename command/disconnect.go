package command

import (
	"fmt"
)

/* Disconnect Command */
type Disconnect struct {
	ShellCommand
}

func (this *Disconnect) Name() string {
	return "CONNECT"
}

func (this *Disconnect) CommandCompletion() bool {
	return false
}

func (this *Disconnect) MinArgs() int {
	return 0
}

func (this *Disconnect) MaxArgs() int {
	return 1
}

func (this *Disconnect) ParseCommand(queryurl []string) error {

	DISCONNECT = true
	fmt.Println("\nCouchbase query shell not connected to any endpoint. Use \\CONNECT command to connect.  ")
	return nil
}

func (this *Disconnect) PrintHelp() {
	fmt.Println("\\DISCONNECT;")
	fmt.Println()
	fmt.Println("Disconnect from the query service or cluster endpoint url.")
	fmt.Println(" For Example : \n\t \\DISCONNECT;")
	fmt.Println()
}
