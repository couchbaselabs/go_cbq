package command

import (
	"fmt"
)

/* Help Command */
type Version struct {
	ShellCommand
}

func (this *Version) Name() string {
	return "VERSION"
}

func (this *Version) CommandCompletion() bool {
	return false
}

func (this *Version) MinArgs() int {
	return 0
}

func (this *Version) MaxArgs() int {
	return 0
}

func (this *Version) ParseCommand() error {
	fmt.Println("SHELL VERSION : " + SHELL_VERSION)
	fmt.Println("SERVER VERSION : " + SERVER_VERSION)
	return nil
}

func (this *Version) PrintHelp() {
	fmt.Println("\\VERSION")
	fmt.Println()
	fmt.Println("Print the Shell Version")
	fmt.Println()
	fmt.Println(" For Example : \n\t \\VERSION;")
	fmt.Println()
}
