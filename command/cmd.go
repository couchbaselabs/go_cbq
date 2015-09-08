package command

const (
	SHELL_VERSION  = "1.0"
	SERVER_VERSION = "4.0"
)

var QUERYURL = ""
var DISCONNECT = false

/*
	Interface to be implemented for shell commands.
*/
type ShellCommand interface {
	/* Name of the comand */
	Name() string
	/* Return true if included in shell command completion */
	CommandCompletion() bool
	/* Returns the Minimum number of input arguments required by the function */
	MinArgs() int
	/* Returns the Maximum number of input arguments allowed by the function */
	MaxArgs() int
	/* Method that implements the functionality*/
	ParseCommand() error
	/* Print Help information for command and its usage with an example */
	PrintHelp()
}
