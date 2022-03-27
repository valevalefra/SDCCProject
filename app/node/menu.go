package main

import (
	"SDCCProject/app/utility"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

type CommandOption struct {
	Command, Description string
	Function             func(args ...string) error
}

type MenuOptions struct {
	Prompt     string
	MenuLength int
}

// Menu struct encapsulates Commands and Options
type Menu struct {
	Title    string
	Commands []CommandOption
	Options  MenuOptions
}

func menu() {

	commandOptions := []CommandOption{
		{Command: "lmp", Description: "digit lmp for choose lamport's algorithm", Function: lamport},
		{Command: "ra", Description: "digit ra for choose Ricart-Agrawala's algorithm", Function: ricart},
		{Command: "quit or exit", Description: "Close application", Function: nil},
	}
	menuOptions := NewMenuOptions("Insert command > ", 0)
	newMenu := NewMenu("...MAIN MENU...", commandOptions, menuOptions)
	newMenu.Start()

}

func ricart(args ...string) error {

	algorithmChoosen = 1
	commandOptions := []CommandOption{
		{Command: "send", Description: "Send message use: send arg1 arg2 ...", Function: sendMessages},
		{Command: "quit or exit", Description: "Close application", Function: nil},
	}
	menuOptions := NewMenuOptions("Insert command > ", 0)
	newMenu := NewMenu("...MAIN MENU...", commandOptions, menuOptions)
	newMenu.Start()
	return nil
}

func lamport(args ...string) error {

	algorithmChoosen = 0
	commandOptions := []CommandOption{
		{Command: "send", Description: "Send message use: send arg1 arg2 ...", Function: sendMessages},
		{Command: "quit or exit", Description: "Close application", Function: nil},
	}
	menuOptions := NewMenuOptions("Insert command > ", 0)
	newMenu := NewMenu("...MAIN MENU...", commandOptions, menuOptions)
	newMenu.Start()
	return nil
}

func NewMenuOptions(prompt string, length int) MenuOptions {
	if prompt == "" {
		prompt = "> "
	}

	if length == 0 {
		length = 100
	}

	return MenuOptions{prompt, length}
}

func (m *Menu) Start() {
	m.start(os.Stdin)
}

// Creates a new menu with options
func NewMenu(title string, cmds []CommandOption, options MenuOptions) *Menu {
	return &Menu{title, cmds, options}
}

func (m *Menu) menu() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	layoutMenu(w, m.Commands, m.Options.MenuLength, m.Title)
}

// Handle building menu layout
func layoutMenu(w *tabwriter.Writer, cmds []CommandOption, width int, title string) {
	fmt.Fprintln(w, title)
	fmt.Fprintln(w, "*\tCommand\tDescription\t")
	for i := range cmds {
		// Write command
		fmt.Fprintf(w, "*\t%s\t", cmds[i].Command)

		// Check description length
		description_length := len(cmds[i].Description)

		if description_length <= width {
			fmt.Fprintf(w, "%s\t\n", cmds[i].Description)
			continue
		}

		if description_length > width {
			layoutLongDescription(w, cmds[i].Description, width)
		}

	}
	fmt.Fprintln(w)
	w.Flush()
}

func cleanCommand(cmd string) ([]string, error) {
	cmd_args := strings.Split(strings.Trim(cmd, " \n"), " ")
	return cmd_args, nil
}

// Main loop
func (m *Menu) start(reader io.Reader) {
	m.menu()
MainLoop:
	for {
		input := bufio.NewReader(reader)
		// Prompt for input
		m.prompt()

		inputString, err := input.ReadString('\n')
		if err != nil {
			// If we didn't receive anything from ReadString
			// we shouldn't continue because we're not blocking
			// anymore but we also don't have any data
			break MainLoop
		}

		cmd, _ := cleanCommand(inputString)
		if len(cmd) < 1 {
			break MainLoop
		}
		// Route the first index of the cmd slice to the appropriate case
	Route:
		switch cmd[0] {
		case "exit", "quit":
			fmt.Println("Exiting from the actual menu ... ")
			break MainLoop

		case "menu":
			m.menu()
			break

		default:
			// Loop through commands and find the right one
			// Probably a more efficient way to do this, but unless we have
			// tons of commands, it probably doesn't matter
			for i := range m.Commands {
				if m.Commands[i].Command == cmd[0] {
					err := m.Commands[i].Function(cmd[1:]...)
					if err != nil {
						panic(err)
					}

					break Route
				}
			}
			// Shouldn't get here if we found a command
			fmt.Println("Unknown command print list")
			for e := scalarMsgQueue.Front(); e != nil; e = e.Next() {
				item := e.Value.(utility.Message)
				log.Printf("MESSAGE IN QUEUE:send id %d:: text %s:tipo %d", item.SendID, item.Text, item.Type)
			}
		}
	}
}

func (m *Menu) prompt() {
	fmt.Print(m.Options.Prompt)
}

// Return tokens up cumulative maxsize
func getDescriptionRange(tokens []string, start int, maxsize int) ([]string, int) {
	total := 0
	token_part := tokens[start:]
	for i := range token_part {
		length := len(token_part[i])
		if total+length > maxsize {
			return token_part[0 : i-1], start + i
		}
		total = total + length
	}
	return token_part[0:], -1
}

func layoutLongDescription(w *tabwriter.Writer, d string, width int) {

	// Tokenize description
	tokens := strings.Fields(d)

	// Get description for range
	description, lastIndex := getDescriptionRange(tokens, 0, width)

	// Write first MAX_LENGTH of description
	fmt.Fprintf(w, "%s\t\n", strings.Join(description, " "))

	for {
		if lastIndex == -1 {
			break
		}

		description, lastIndex = getDescriptionRange(tokens, lastIndex, width)
		fmt.Fprintf(w, "*\t\t%s\t\n", strings.Join(description, " "))
	}
}
