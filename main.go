package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	// load file extension flags
	phpPtr := flag.Bool("php", false, "Minify files with the .php extension.")
	cssPtr := flag.Bool("css", false, "Minify files with the .css extension.")
	htmlPtr := flag.Bool("html", false, "Minify files with the .html extension.")
	jsPtr := flag.Bool("js", false, "Minify files with the .js extension.")
	allPtr := flag.Bool("all", false, "Add this flag if you want to go through all sub folders inside the current working directory.")

	flag.Parse()

	// save current working directory as start point
	dirPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// print introduction to program w/ help info
	printIntro()

	// start of program process
	checkDirectoryAndFiles(*phpPtr, *cssPtr, *jsPtr, *htmlPtr, *allPtr, dirPath)

}

func checkDirectoryAndFiles(php, css, js, html, allDirectories bool, directoryPath string) {
	var fileExtToValidate []string
	var fileNames []string

	if php == true {
		fileExtToValidate = append(fileExtToValidate, ".php")
	}

	if css == true {
		fileExtToValidate = append(fileExtToValidate, ".css")
	}

	if js == true {
		fileExtToValidate = append(fileExtToValidate, ".js")
	}

	if html == true {
		fileExtToValidate = append(fileExtToValidate, ".html")
	}

	if allDirectories == true {

		clearTerminal()
		fmt.Println("")

		var subDirectoryList []string
		err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			if info.IsDir() {
				fmt.Println(path)
				subDirectoryList = append(subDirectoryList, path)
			}
			return nil
		})
		checkError(err)

		valid := userVerification("\n\n=> Do these directories look correct to you? (true / false)")

		if valid {

			clearTerminal()

			for _, element := range subDirectoryList {
				err := filepath.Walk(element, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						panic(err)
					}
					if info.Mode().IsRegular() {
						temp := strings.Index(path, ".")
						if temp > 0 {
							hasExt := checkExtension(path[temp:], fileExtToValidate)
							if hasExt == true {
								fmt.Println(path)
								fileNames = append(fileNames, path)
							}
						}
					}
					return nil
				})
				checkError(err)
			}
		} else {
			fmt.Println("\n\n=> The program has encountered a problem finding all of your sub directories. Please let the developer know, thank you.")
			os.Exit(0)
		}
	} else {
		err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				//todo - don't want to just end program, get user input
				panic(err)
			}
			if info.Mode().IsRegular() {
				temp := strings.Index(path, ".")
				if temp > 0 {
					hasExt := checkExtension(path[temp:], fileExtToValidate)
					if hasExt == true {
						fmt.Println(path)
						fileNames = append(fileNames, path)
					}
				}
			}
			return nil
		})
		checkError(err)
	}

	// start pulling file names on files currently in directory
	// recursively pull files?

	filesCorrect := userVerification("\n\n=> Do these files look correct to you? (true / false)")

	if filesCorrect {
		clearTerminal()
		// start minify
		fmt.Println("=> Starting to minify your files. Please wait.")
	} else {
		fmt.Println("It seems the program wasn't accurately able to find your files, please tell the developer. Thank you.")
		os.Exit(0)
	}
}

func checkExtension(path string, extensions []string) bool {
	for i := 0; i < len(extensions); i++ {
		if path == extensions[i] {
			return true
		}
	}
	return false
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
		} else if userInput != "true" {
			fmt.Println("You have enetered an incorrect value, please try again.")
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
	fmt.Println(" If you're having any problems with this program,\n please consult https://www.github.com/teepleb/gominify/")
	fmt.Println()
	fmt.Println()
	fmt.Println("=> Your files are being indexed and located, please wait.")
	fmt.Println()
	fmt.Println()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
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
