package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitors = 3
const delay = 5

func main() {
	showIntroduction()

	for {
		showOptions()
		command := runCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Showing logs...")
			showLogs()
		case 0:
			fmt.Println("Exiting")
			os.Exit(0)
		default:
			fmt.Println("Invalid command!")
			os.Exit(-1)
		}
	}
}

func showIntroduction() {
	myName := "Lucas"
	appVersion := 1.1

	fmt.Println("Hello, sr", myName)
	fmt.Println("Current version of this software is", appVersion)
}

func showOptions() {
	fmt.Println("1 - Starting the application...")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Exit")
}

func runCommand() int {
	var commandRead int
	fmt.Scan(&commandRead)
	fmt.Println("Command: ", commandRead)

	return commandRead
}

func startMonitoring() {
	fmt.Println("Starting the monitoring...")

	websites := readFileWebsite()
	fmt.Println(websites)

	for i := 0; i < monitors; i++ {
		for i, website := range websites {
			fmt.Println("Testing website: ", i, " - ", website)
			testWebsite(website)
		}

		time.Sleep(delay * time.Minute)
	}
}

func testWebsite(website string) {
	resp, err := http.Get(website)

	if err != nil {
		fmt.Println("Occurred an error during the website test: ", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site: ", website, "is up!")
		registerLog(website, true)
	} else {
		fmt.Println("Site: ", website, "is down! Status code: ", resp.StatusCode)
		registerLog(website, false)
	}
}

func readFileWebsite() []string {
	var websites []string
	file, err := os.Open("websites.txt")

	if err != nil {
		fmt.Println("Occurred an error during read file website: ", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		websites = append(websites, line)

		if err == io.EOF {
			break
		}
	}

	file.Close()

	return websites
}

func registerLog(website string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Occurred an error during register log: ", err)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + website + " - online: " + strconv.FormatBool(status) + "\n")
	file.Close()
}

func showLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Occurred an error during get logs: ", err)
	}

	fmt.Println(string(file))
}
