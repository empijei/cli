package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/empijei/cli/lg"
)

// Commands is the list of all registered commands.
// Each command `cmd` is invoked via `mainprogram cmd`
var Commands []*Cmd

var DefaultCommand *Cmd

type Cmd struct {
	// Name is the name of the command. It's what comes after the main executable.
	Name string

	// Run is the command entrypoint.
	Run func(...string)

	// UsageLine is the header of what's printed by flag.PrintDefaults.
	UsageLine string

	// Short is a one-line description of what the command does.
	Short string

	// Long is the detailed description of what the command does.
	Long string

	// Flag is the set of flags accepted by the command. This should be initialized in
	// the command's module's `init` function. The parsing of these flags is issued
	// by the main entrypoint, so each command doesn't have to do it itself.
	Flag flag.FlagSet
}

// AddCommand allows packages to setup their own command. In order for them to
// be compiled, they must be imported by the main package with the "_" alias
func AddCommand(c *Cmd) {
	Commands = append(Commands, c)
}

func (c *Cmd) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	c.Flag.PrintDefaults()
	os.Exit(2)
}

// FindCommand takes a string and searches for a command whose name has that string as
// prefix. If more than 1 command name has that string as a prefix (and no command name
// equals that string), an error is returned. If no suitable command is found, an error
// is returned.
func FindCommand(name string) (command *Cmd, err error) {
	for _, cmd := range Commands {
		if cmd.Name == name {
			command = cmd
			// If there were several commands beginning with this string, but I
			// have an exact match, the error should not be returned.
			err = nil
			return
		}
		if strings.HasPrefix(cmd.Name, name) {
			if command != nil {
				err = fmt.Errorf("Ambiguous command: '%s'.", name)
			} else {
				command = cmd
			}
		}
	}
	if command == nil {
		err = fmt.Errorf("Command not found: '%s'.", name)
	}
	return
}

func callCommand(command *Cmd) {
	command.Flag.Usage = command.Usage
	err := command.Flag.Parse(os.Args[1:])
	if err != nil {
		lg.Error(err)
		return
	}
	command.Run(command.Flag.Args()...)
}
