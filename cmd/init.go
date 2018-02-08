package cmd

import (
	"fmt"
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			fmt.Fprintln(os.Stderr, "Available commands are:\n")
			for _, cmd := range Commands {
				fmt.Fprintln(os.Stderr, "\t"+cmd.Name+"\n\t\t"+cmd.Short)
			}
			fmt.Fprintln(os.Stderr, "\nDefault command is: "+DefaultCommand.Name)
			return
		}
		if command == nil {
			if DefaultCommand != nil {
				callCommand(DefaultCommand)
				return
			}
			lg.Error("Command not found, default not specified")
			return
		}
		callCommand(command)
	}
	if DefaultCommand != nil {
		callCommand(DefaultCommand)
		return
	}
	lg.Error("No command specified, default not configured")
}
