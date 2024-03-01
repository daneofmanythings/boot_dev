package repl

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"gitlab.com/daneofmanythings/pokedex/repl/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

func newCommand(name, description string, callback func([]string) error) *cliCommand {
	return &cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

func (cc *cliCommand) Run(args []string) error {
	return cc.callback(args)
}

var help cliCommand = *newCommand(
	"help",
	"Displays the help message",
	commandHelp,
)

func commandHelp(args []string) error {
	output := "\n\nWelcome to the Pokedex!\nUsage:\n\n"
	for _, c := range commandSlice {
		output += c.name + ": " + c.description + "\n"
	}
	output += "\n\n"
	fmt.Print(output)
	return nil
}

var exit cliCommand = *newCommand(
	"exit",
	"Exit the pokedex",
	commandExit,
)

func commandExit(args []string) error {
	defer os.Exit(0)
	return nil
}

var map_ cliCommand = *newCommand(
	"map",
	"returns the next 20 locations to explore",
	commandMap,
)

func commandMap(args []string) error {
	response, err := internal.MapResponse(internal.NEXT)
	if err != nil {
		return err
	}
	fmt.Print("\n")
	for _, r := range response.Results {
		fmt.Println(r.Name)
	}
	fmt.Print("\n")
	return nil
}

var bmap_ cliCommand = *newCommand(
	"bmap",
	"returns the previous 20 locations to explore",
	commandBmap,
)

func commandBmap(args []string) error {
	response, err := internal.MapResponse(internal.PREV)
	if err != nil {
		return err
	}
	fmt.Println()
	for _, r := range response.Results {
		fmt.Println(r.Name)
	}
	fmt.Println()
	return nil
}

var explore cliCommand = *newCommand(
	"explore",
	"Lists all pokemon in a specified area",
	commandExplore,
)

func commandExplore(args []string) error {
	if len(args) < 2 {
		return errors.New("please input the area to explore: 'explore <area>'")
	}
	response, err := internal.ExploreResponse(args[1])
	if err != nil {
		return err
	}

	fmt.Println()
	for _, pe := range response.PokemonEncounters {
		fmt.Println(pe.Pokemon.Name)
	}
	fmt.Println()

	return nil
}

var catch cliCommand = *newCommand(
	"catch",
	"Attempts to catch the specified pokemon",
	commandCatch,
)

var pokedex map[string]internal.PokemonApiResponse = make(map[string]internal.PokemonApiResponse)

func commandCatch(args []string) error {
	if len(args) < 2 {
		return errors.New("please input the area to explore: 'explore <area>'")
	}
	response, err := internal.CatchResponse(args[1])
	if err != nil {
		return err
	}

	baseCatchChance := int(100 * (1 - (float64(response.BaseExperience) / 300.0))) // mew's base xp
	skew := rand.Intn(20)
	catch := rand.Intn(100)

	baseCatchChance += skew - 10
	if catch < baseCatchChance {
		fmt.Printf("%s was caught!", args[1])
		pokedex[args[1]] = response
	} else {
		fmt.Printf("%s escaped!", args[1])
	}
	fmt.Print("\n\n")

	return nil
}

var inspect cliCommand = *newCommand(
	"inspect",
	"Displays the name, height, weight, stats, and type of specified pokemon",
	commandInspect,
)

func commandInspect(args []string) error {
	if len(args) < 2 {
		return errors.New("please input the area to explore: 'explore <area>'")
	}
	response, ok := pokedex[args[1]]
	if !ok {
		return fmt.Errorf("%s does not have an entry in your pokedex. Catch one to get its information", args[1])
	}

	displayString := `
Name: %s
Height: %d
Weight: %d
Stats:
 - hp: %d
 - attack: %d
 - special-attack: %d
 - defense: %d
 - special-defense: %d
 - speed: %d
Types:
`
	displayString = fmt.Sprintf(
		displayString,
		response.Name,
		response.Height,
		response.Weight,
		response.Stats[0].BaseStat,
		response.Stats[1].BaseStat,
		response.Stats[2].BaseStat,
		response.Stats[3].BaseStat,
		response.Stats[4].BaseStat,
		response.Stats[5].BaseStat,
	)
	for _, typeEntry := range response.Types {
		displayString += fmt.Sprintf(" - %s\n", typeEntry.Type.Name)
	}

	fmt.Println(displayString)

	return nil
}

var displayPokedex = *newCommand(
	"pokedex",
	"Displays the names of all pokemon in pokedex",
	commandPokedex,
)

func commandPokedex(args []string) error {
	fmt.Println("\nYour Pokedex:")
	for e := range pokedex {
		fmt.Printf(" - %s\n", e)
	}
	fmt.Println()
	return nil
}

var commandSlice []cliCommand = []cliCommand{
	map_,
	bmap_,
	explore,
	catch,
	inspect,
	displayPokedex,
	exit,
}

func commands() map[string]cliCommand {
	commands := make(map[string]cliCommand)
	for _, command := range commandSlice {
		commands[command.name] = command
	}
	return commands
}

var COMMANDS map[string]cliCommand = commands()
