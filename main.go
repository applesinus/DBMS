package main

import (
	"DBMS/database"
	"DBMS/task0+3"
	"DBMS/task789"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func executeCommand(db *task0.Database, settings map[string]string, command string) string {
	words := strings.Fields(command)
	switch len(words) {
	// Empty line
	case 0:
		return "empty"

	// Single command (application side)
	case 1:
		switch strings.ToLower(words[0]) {
		case "stop", "exit", "-1":
			return "exit"
		case "help":
			help()
			return "ok"
		case "pools":
			for _, poolName := range *db.ListPools(settings) {
				fmt.Print(poolName + " ")
			}
			fmt.Println()
			return "ok"
		default:
			_, err := os.Open(words[0])
			if err == nil {
				ret := executeFile(db, settings, words[0])
				if ret == "error" {
					fmt.Println("Could not execute file")
				} else {
					return ret
				}
			}
			fmt.Printf("%s\n", command)
			return "error"
		}
	default:
		res := database.ExecuteCommand(command)
		if res != "ok" {
			fmt.Printf("%s\n", res)
		}
		return res
	}

}

func executeFile(db *task0.Database, settings map[string]string, filePath string) string {
	file, err := os.Open(filePath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			command := scanner.Text()
			ret := executeCommand(db, settings, command)
			fmt.Printf("command: %v, ret: %v\n", command, ret)
			if ret == "exit" {
				return "exit"
			} else if ret == "error" {
				return "error"
			}
		}
	} else {
		fmt.Printf("Could not open file. %v\n", err)
		return "error"
	}

	return "ok"
}

func help() {
	fmt.Println("Commands:")
	fmt.Println("\t'stop', 'exit' or '-1' to stop the program")
	fmt.Println("\t'help' to show this help")
}

func main() {
	server := task789.StartHTTPServer(8080)

	filepath := flag.String("f", "", "initial filepath")
	persistant := flag.Bool("p", false, "persistant mode")
	flag.Parse()

	settings := make(map[string]string)
	if *persistant {
		settings["persistant"] = "on"
	}
	database.SetSettings(settings)

	db := database.Database()

	// if there's a filename
	if *filepath != "" {
		ret := executeFile(db, settings, *filepath)
		if ret == "error" {
			fmt.Println("Could not execute file")
		} else if ret == "exit" {
			task789.StopHTTPServer(server)
			return
		}
	}

	help()
	line := ""

	in := bufio.NewReader(os.Stdin)
	for line != "exit" {
		fmt.Println("You can type your command any time")
		line, err := in.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		ret := executeCommand(db, settings, line)

		if ret == "exit" {
			task789.StopHTTPServer(server)
			return
		}
	}
}
