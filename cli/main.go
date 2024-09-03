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

const MAKE_YELLOW = "\033[33m"
const MAKE_GREEN = "\033[32m"
const MAKE_RED = "\033[31m"
const CLEAR_COLOR = "\033[0m"
const GRAY_VALUE = 192

var exit bool

func printGray(str string, bold bool) {
    if bold {
        fmt.Printf("\033[1m")
    }
    fmt.Printf("\033[38;2;%d;%d;%dm", GRAY_VALUE, GRAY_VALUE, GRAY_VALUE)
    fmt.Printf(str)
    fmt.Printf(CLEAR_COLOR)
}

func printError(str string) {
    var caution rune = '⚠'
    fmt.Fprintf(os.Stderr, MAKE_RED)
    fmt.Fprintf(os.Stderr,"%c  %s\n", caution, str)
    fmt.Fprintf(os.Stderr, CLEAR_COLOR)
}
func printErr(err error) {
    var caution rune = '⚠'
    fmt.Fprintf(os.Stderr, MAKE_RED)
    fmt.Fprintf(os.Stderr,"%c  %s\n", caution, err)
    fmt.Fprintf(os.Stderr, CLEAR_COLOR)
}
    


func printHelp() {
    printGray("valid commands:\n", true)
    printGray("  i   init\tInit cwd as a repository\n", false)
    printGray("  a   add\tAdd file to next commit\n", false)
    printGray("  d   diff\tIdentify changes\n", false)
    printGray("  t   tag\tTag this commit with version\n", false)
    printGray("  c   commit\tSend code to remote\n", false)
    printGray("  l   log\tShow history log\n", false)
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
        case "i":
        case "a":
        case "d":
        case "t":
        case "c":
        case "l":
        case "h":
            printHelp()
        case "q":
            fmt.Println("    quitting...")
            running = false
        default:
            printError("command invalid")
            printHelp()
        }
        
    }
}

func parseArgs(args []string) ([]string, error) {
    exit = args[0] == "./goverse"
    if len(args) < 2 {
        return args, errors.New("usage error: Not enough arguments")
    } 
    if len(args) > 2 {
        return args, errors.New("usage error: Too many arguments")
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
    switch args[1]{
    case "i", "interactive":
        // fmt.Println("Interactive start")
        interactive()
    default: // invalid command, print to stderr
        printError("usage error: invalid argument")
        if exit {os.Exit(1)}
        return
    }
}
