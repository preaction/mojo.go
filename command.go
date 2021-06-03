package mojo

import (
	"fmt"
	"net"
	"os"
	"path"
	"sort"
)

// Command is an interface for application commands. Types that
// implement this interface can be registered as Application commands.
type Command interface {
	Description() string
	Usage() string
	Run(args []string) error
}

// HelpCommand is a command to show help for other commands.
type HelpCommand struct {
	App *Application
}

// Description returns the short description of the command's function
func (cmd *HelpCommand) Description() string {
	return "Display help for commands"
}

// Usage returns the full documentation for the command
func (cmd *HelpCommand) Usage() string {
	return `[COMMAND]

ARGUMENTS
  COMMAND                          Display help for the given command
`
}

// Run runs the command. If run with no arguments, will list all
// commands in the application. If run with an argument, it is the name
// of a command to display full documentation for.
func (cmd *HelpCommand) Run(args []string) error {
	// XXX: Look at Go's built-in `flag` library to see what kind of
	// help it can auto-generate
	if len(args) == 0 {
		fmt.Printf("Usage: %s COMMAND [OPTIONS]\n\n", path.Base(os.Args[0]))
		// Display all commands and Description
		fmt.Printf("Commands:\n")
		names := make([]string, 0, len(cmd.App.Commands))
		for name, _ := range cmd.App.Commands {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			cmd := cmd.App.Commands[name]
			fmt.Printf(" %-10s- %s\n", name, cmd.Description())
		}
	} else {
		// Display one command's Usage
		name := args[0]
		cmd, ok := cmd.App.Commands[name]
		if !ok {
			fmt.Printf("Command not found: %s", name)
			os.Exit(1)
		}
		fmt.Printf("Usage: %s %s %s\n", path.Base(os.Args[0]), args[0], cmd.Usage())
	}
	return nil
}

// DaemonCommand starts the web application server command
type DaemonCommand struct {
	App *Application
}

// Description returns the short description of the command's function
func (cmd *DaemonCommand) Description() string {
	return "Start the web application server"
}

// Usage returns the full documentation for the command
func (cmd *DaemonCommand) Usage() string {
	return `[OPTIONS]

OPTIONS
  -l <listen>                      The host/port to listen on. Defaults to
                                   "http://*:3000"
`
}

// Run runs the web application
func (cmd *DaemonCommand) Run(args []string) error {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on http://%s\n", l.Addr())
	srv := Server{
		App: cmd.App,
	}
	srv.Serve(l)
	return nil
}

// VersionCommand shows the current version of Mojolicious
type VersionCommand struct {
	App *Application
}

// Description returns the short description of the command's function
func (cmd *VersionCommand) Description() string {
	return "Show the current Mojolicious version"
}

// Usage returns the full documentation for the command
func (cmd *VersionCommand) Usage() string {
	return ""
}

// Run displays the version and other useful information
func (cmd *VersionCommand) Run(args []string) error {
	fmt.Printf("The current version")
	return nil
}
