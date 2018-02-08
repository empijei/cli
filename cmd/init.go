package cmd

import (
	"os"

	"github.com/empijei/cli/lg"
)

// Printbanner will be called if the output is not to a pipe
var Printbanner func()

func Init() {
	stderrinfo, err := os.Stderr.Stat()
	if err == nil && stderrinfo.Mode()&os.ModeCharDevice == 0 {
		// Output is a pipe, turn off colors
		lg.Color = false
	} else {
		// Output is to terminal, print banner
		if Printbanner != nil {
			Printbanner()
		}
	}

	if len(os.Args) > 1 {
		//read the first argument
		directive := os.Args[1]
		if len(os.Args) > 2 {
			//shift parameters left, but keep argv[0]
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			os.Args = os.Args[:1]
		}
		command, err := FindCommand(directive)
		if err == nil {
			callCommand(command)
		} else {
			lg.Error(err.Error())
			lg.Error("Available commands are:\n")
			for _, cmd := range Commands {
				lg.Error("\t" + cmd.Name + "\n\t\t" + cmd.Short)
			}
			lg.Error("\nDefault command is: ", DefaultCommand.Name)
		}
	} else {
		callCommand(DefaultCommand)
	}
}
