package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	po "pokedex/internal/PokeAPImanager"
	"strings"
)

type cliCommand struct {
	name    string
	command func()
}

type config struct {
	NextUrl     string
	PreviousUrl string
	StepSize    int64
}

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *config) mapCommandFunc() {

	step := c.StepSize
	url := c.NextUrl + fmt.Sprint(step)

	res, err := po.GetLocations(url)

	if err != nil {
		fmt.Println("Couldn't read!")
	}
	fmt.Print(string(res))

	var locations Locations
	errLoc := json.Unmarshal(res, &locations)

	if errLoc != nil {
		fmt.Println("Error unmarshaling JSON:", errLoc)
	}

	for _, j := range locations.Results {
		fmt.Println(j.Name)
	}

	c.NextUrl = locations.Next
	c.PreviousUrl = locations.Previous

}

func helpCommandFunc() {
	fmt.Println(`Welcome to the Pokedex!
Usage:
help: Displays a help message
exit: Exit the Pokedex`)
}

func exitCommandFunc() {
	os.Exit(1)
}

func commandMap(cfg config) map[string]cliCommand {
	commandMap := make(map[string]cliCommand)

	commandMap["help"] = cliCommand{name: "help",
		command: helpCommandFunc}
	commandMap["exit"] = cliCommand{name: "exit",
		command: exitCommandFunc}
	commandMap["map"] = cliCommand{name: "map",
		command: cfg.mapCommandFunc}

	return commandMap
}

func main() {

	var cfg config

	cfg.NextUrl = "https://pokeapi.co/api/v2/location-area/?limit=20&offset=20"
	cfg.PreviousUrl = ""

	for {
		fmt.Print("Pokedex > ")
		mapOfFuncs := commandMap(cfg)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		commandTextSlice := strings.Split(text, ">")
		commandText := commandTextSlice[0]

		inputCommand, exists := mapOfFuncs[commandText]
		if exists {
			inputCommand.command()
		} else {
			fmt.Printf("Command %v doesn't exist!\n", text)
			continue
		}

	}

}
