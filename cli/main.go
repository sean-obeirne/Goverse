package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "io"
	// "flag"
	// "strings"
	"errors"
)


//////////////////////////////
// PRINT / HELPER FUNCTIONS //
//////////////////////////////

const MAKE_YELLOW = "\033[33m"
const MAKE_GREEN = "\033[32m"
const MAKE_RED = "\033[31m"
const CLEAR_COLOR = "\033[0m"
const GRAY_VALUE = 192
const CAUTION = '⚠'

func printGray(str string, bold bool) {
    if bold {
        fmt.Printf("\033[1m")
    }
    fmt.Printf("\033[38;2;%d;%d;%dm", GRAY_VALUE, GRAY_VALUE, GRAY_VALUE)
    fmt.Printf(str)
    fmt.Printf(CLEAR_COLOR)
}

func printError(str string) {
    fmt.Fprintf(os.Stderr,"%s%c  %s%s\n", MAKE_RED, CAUTION, str, CLEAR_COLOR)
}
func printErr(err error) {
    fmt.Fprintf(os.Stderr,"%s%c  %s%s\n", MAKE_RED, CAUTION, err, CLEAR_COLOR)
}

///////////////////
// COMMAND BLOCK //
///////////////////

// Paths
var baseDir string
const (
    GOVERSE_DIR = ".goverse/"
    OBJECTS_DIR = GOVERSE_DIR + "objects/"
    TAGS_DIR    = GOVERSE_DIR + "tags/"
    CONFIG_FILE = GOVERSE_DIR + "config"
    HEAD_FILE   = GOVERSE_DIR + "head"
)

func initGoverse() {
    dirs := []string{ GOVERSE_DIR, OBJECTS_DIR, TAGS_DIR }
    files := []string{ CONFIG_FILE, HEAD_FILE }
    for _, dir := range dirs {
        err := os.MkdirAll(baseDir + dir, 0755)
        if err != nil {
            printErr(err)
        }
    }
    for _, file := range files {
        newFile, err := os.Create(baseDir + file)
        if err != nil {
            printErr(err)
        }
        newFile.Close()
    }
}

func status(){
    ls, err := os.ReadDir(baseDir + GOVERSE_DIR)
	if err != nil {
        printErr(err)
    }
    for _, entry := range ls {
        fmt.Println(entry)
    }
}

func add() {
    
}

func diff() {
    
}

func tag() {
    
}

func commit() {
    
}

func log() {
    
}
    
func flush() {
    err := os.RemoveAll(baseDir + GOVERSE_DIR)
    if err != nil {
        printErr(err)
    }
}

func printHelp() {
    printGray("valid commands:\n", true)
    printGray("  i   init\tInit cwd as a repository\n", false)
    printGray("  s   status\tCheck repository status\n", false)
    printGray("  a   add\tAdd file to next commit\n", false)
    printGray("  d   diff\tIdentify changes\n", false)
    printGray("  t   tag\tTag this commit with version\n", false)
    printGray("  c   commit\tSend code to remote\n", false)
    printGray("  l   log\tShow history log\n", false)
    printGray("  f   flush\tDelete all goverse files\n", false)
    printGray("  h   help\tDisplay this message\n", false)
    printGray("  q   quit\tTerminate this interactive application\n", false)
}

func interactive() {
    running := true
    for running {
        reader := bufio.NewReader(os.Stdin)

        fmt.Print(MAKE_YELLOW)
        fmt.Print("goverse command > ")
        fmt.Print(CLEAR_COLOR)

        fmt.Print(MAKE_GREEN)
        cmd, _ := reader.ReadString('\n')  // Reads until newline
        fmt.Print(CLEAR_COLOR)

        cmd = strings.Trim(cmd, "\n")
        // fmt.Println("You entered:", strings.TrimSpace(input))

        switch(cmd) {
        case "i", "init":
            initGoverse()
        case "s", "status":
            status()
        case "a", "add":
        case "d", "diff":
        case "t", "tag":
        case "c", "commit":
        case "l", "log":
        case "f", "flush":
            flush()
        case "h", "help":
            printHelp()
        case "q", "quit":
            printGray("    quitting...\n", false)
            running = false
        default:
            printError("command invalid")
            printHelp()
        }
    }
}




////////////////
// MAIN CHUNK //
////////////////

// really only necessary for testing purposes to keep output clean
var exit bool

func parseArgs(args []string) ([]string, error) {
    exit = args[0] == "./goverse"
    if len(args) < 2 {
        return args, errors.New("usage error: Not enough arguments")
    } 
    if len(args) > 2 {
        return args, errors.New("usage error: Too many arguments")
    }
    if args[1][len(args[1])-1] != '/' {
        entry, err := os.Stat(args[1])
        if err != nil {
            return args, errors.New("usage error: First argument is not a valid directory")
        }

        if entry.IsDir() {
            args[1] = args[1] + "/"
        } else {
            return args, errors.New("usage error: First argument is not a valid directory")
        }
    }
    return args, nil
}

func main() {
    args, err := parseArgs(os.Args)
    if err != nil {
        printErr(err)
        // fmt.Fprintf(os.Stderr, "%s\n", err)
        if exit {os.Exit(1)}
        return
    }
    baseDir = args[1]
    fmt.Println("Base dir: ", baseDir)
    interactive()

    // switch args[1]{
    // case "i", "interactive":
    //     interactive()
    // default:
    //     printError("usage error: invalid argument, did you mean 'interactive'?")
    //     if exit {os.Exit(1)}
    //     return
    // }
}
