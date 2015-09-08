package command

import (
	"fmt"
	"math"
)

/* Help Command */
type Help struct {
	ShellCommand
}

func (this *Help) Name() string {
	return "HELP"
}

func (this *Help) CommandCompletion() bool {
	return false
}

func (this *Help) MinArgs() int {
	return 1
}

func (this *Help) MaxArgs() int {
	return math.MaxInt16
}

func (this *Help) ParseCommand(v []string) error {
	/*for i, vals := range v {
		if strings.Contains(vals, "*") {

		} else {

		}
		vals.PrintHelp()
	} */
	return nil
}

func (this *Help) PrintHelp() {
	fmt.Println("\\HELP [<args>]")
	fmt.Println()
	fmt.Println("The input arguments are shell commands. If a * is input then the command displays HELP information for all input shell commands")
	fmt.Println()
	fmt.Println(" For Example : \n\t \\HELP VERSION; \n\t \\HELP EXIT DISCONNECT VERSION; \n\t \\HELP *;")
	fmt.Println()
}
