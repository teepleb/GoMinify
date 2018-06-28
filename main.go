package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// load current working directory to save as default in directory flag pointer
	temp, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	// generate flags for parsing supported files
	directoryPathPtr := flag.String("", temp, "Enter the path of the directory you want to minify. Default is current working directory.")
	phpFlagPtr := flag.Bool("php", false, "PHP")
	cssFlagPtr := flag.Bool("css", false, "CSS")
	jsFlagPtr := flag.Bool("js", false, "JS")
	htmlFlagPtr := flag.Bool("html", false, "HTML")

	flag.Parse()

	// easy print introduction to program w/ help info
	printIntro()
	easyPrint("=> Your files are being indexed and located, please wait.")
	checkDirectoryAndFiles(*phpFlagPtr, *cssFlagPtr, *jsFlagPtr, *htmlFlagPtr, *directoryPathPtr)

}

func checkDirectoryAndFiles(php, css, js, html bool, directoryPath string) {
	// check if directory exists
	// start pulling file names on files currently in directory
	// recursively pull files?
	pauseProgram(2)
	filesCorrect := userVerification("=> Do these files look correct to you? (true / false)")
	fmt.Println(filesCorrect)
}

//easyPrint is just an easier way to print to console when needed, will expand later.
func easyPrint(s string) {
	fmt.Println(s)
}

//userVerification is meant to be a quick T/F questionnaire for the user for validation.
func userVerification(question string) bool {
	var userInput string

	for {
		fmt.Println(question)
		fmt.Print("=> ")
		fmt.Scanln(&userInput)
		if userInput == "false" {
			return false
		} else if userInput == "true" {
			return true
		} else if userInput != "false" {
			fmt.Println("You have entered an incorrect value, please try again.")
		}
	}

}

//printIntro is a simple CLI introduction to the software.
func printIntro() {
	fmt.Println()
	fmt.Println()
	fmt.Println("   ==============================")
	fmt.Println("        Welcome  to  GoMinify    ")
	fmt.Println("   ==============================")
	fmt.Println()
	fmt.Println()
	fmt.Println(" If you're having any problems with this program,\n please consult https://www.github.com/gophersion/gominify/")
	fmt.Println()
	fmt.Println()
}

//pauseProgram is used to slow down the program, if required. It is sometimes used for the user to read prompts as needed.
func pauseProgram(x int) {
	time.Sleep(time.Duration(x) * time.Second)
}

//clearTerminal is pretty self explanatory, it will clear the screen on MacOS, Linux Distros, and Windows.
func clearTerminal() {
	var c *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
	case "linux":
		c = exec.Command("clear")
	case "windows":
		c = exec.Command("cmd", "/c", "cls")
	}

	c.Stdout = os.Stdout
	c.Run()
}
