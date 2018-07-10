package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	var filePaths []string

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

	// generate a map of the extensions for ease of use
	var extMap = make(map[string]bool)
	extMap[".php"] = *phpPtr
	extMap[".css"] = *cssPtr
	extMap[".js"] = *jsPtr
	extMap[".html"] = *htmlPtr

	// walk through the CWD to get all the file paths needed to load the files
	filePaths = getFilePaths(extMap, *allPtr, dirPath)

	if filePaths == nil {
		log.Fatal("It seems the program wasn't accurately able to find your files, please tell the developer. Thank you.")
	} else {
		// start minifying files
		minifyFiles(filePaths)
	}

}

func minifyFiles(filePaths []string) {
	clearTerminal()
	fmt.Println("=> Starting to minify your files. Please wait.")

	for _, element := range filePaths {
		if filepath.Ext(element) == ".css" {
			minifyCSS(element)
		}
		//if filepath.Ext(element) == ".js" {
		//	minifyJS(element)
		//}
	}

}

// minifyCSS will minify your CSS files with respect to comments
func minifyCSS(path string) {
	fileData := loadResource(path)
	cssString := ""
	fileName := filepath.Base(path[0:len(path)-len(filepath.Ext(path))]) + ".min.css"
	isComment := false

	for _, element := range fileData {
		if strings.HasPrefix(element, "/*") && strings.HasSuffix(element, "*/") {
			isComment = false
			continue
		} else if strings.HasPrefix(element, "/*") {
			isComment = true
			continue
		} else if strings.HasSuffix(element, "*/") {
			isComment = false
			continue
		}

		if !isComment {
			cssString += strings.TrimSpace(element)
		}

	}

	saveResource(cssString, fileName)

}

func minifyJS(path string) {
	// todo
	// - handle comments
	// - check for semi colon, if no semi colon then new line else place together
	// - only check for semi colon if preceeding is not {
	// - if next char is } then remove semi-colon before
	fileData := loadResource(path)
	jsString := ""
	fileName := filepath.Base(path[0:len(path)-len(filepath.Ext(path))]) + ".min.js"

	for _, element := range fileData {
		if !strings.HasPrefix(element, "//") {
			jsString += strings.TrimSpace(element)
		}

	}

	saveResource(jsString, fileName)
}

func saveResource(data, fileName string) {
	if _, err := os.Stat(fileName); os.IsExist(err) {
		err := os.Remove(fileName)
		if err != nil {
			log.Fatal("There was a problem removing the file from the directory.")
		}
	}

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatal("Error creating minified file.")
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	w.WriteString(data)

	w.Flush()
}

func loadResource(path string) []string {
	var tempData []string
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tempData = append(tempData, scanner.Text())
	}

	return tempData
}

func getFilePaths(extensions map[string]bool, allDirectories bool, directoryPath string) []string {
	var fileNames []string
	var directories []string

	// if the user wants all directories to be walked, do so, otherwise only walk the CWD
	if allDirectories {
		directories = getSubDirectories(directoryPath)
	} else {
		directories = append(directories, directoryPath)
	}

	// walk each directory found/given to snag all the file paths desired
	for _, element := range directories {
		err := filepath.Walk(element, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal("There was a problem walking the directory: " + path)
			}

			if info.Mode().IsRegular() {
				// need to make sure the path is clean with only the extension being grabbed
				fileExtension := filepath.Ext(path)

				if len(directories) == 1 {
					if filepath.Dir(path) == element {
						if fileExtension != "" {
							hasExt := checkExtension(fileExtension, extensions)
							if hasExt == true {
								fmt.Println(path)
								fileNames = append(fileNames, path)
							}
						}
					}
				} else {
					if fileExtension != "" {
						hasExt := checkExtension(fileExtension, extensions)
						if hasExt == true {
							if !checkFileName(path, fileNames) {
								fmt.Println(path)
								fileNames = append(fileNames, path)
							}
						}
					}
				}
			}
			return nil
		})

		if err != nil {
			log.Fatal("There was a problem going through the directory looking for files: " + element)
		}
	}

	if len(fileNames) == 0 {
		log.Fatal("We couldn't find any files with the specified extensions, please navigate to the proper directory or use -all on the command call to search all sub directories as well.")
	}

	// verify with the user that all files are correct before minifying (last check as a save, just in case)
	filesCorrect := userVerification("\n\n=> Do these files look correct to you? (true / false)")

	if filesCorrect {
		return fileNames
	}

	return nil
}

func getSubDirectories(currentDirectory string) []string {
	var subDirectoryList []string

	err := filepath.Walk(currentDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			subDirectoryList = append(subDirectoryList, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal("GoMinify had a problem getting the subdirectories in your current working directory.")
	}

	return subDirectoryList
}

func checkFileName(path string, filePaths []string) bool {
	for _, element := range filePaths {
		if element == path {
			return true
		}
	}
	return false
}

func checkExtension(currentFileExt string, extensions map[string]bool) bool {
	for index, element := range extensions {
		if currentFileExt == index {
			return element
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
		log.Fatal(err)
	}
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
