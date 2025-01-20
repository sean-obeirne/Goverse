package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    // "strconv"
    "errors"
    // "runtime"

    "goverse/core"
)

//////////////////////////////
// PRINT / HELPER FUNCTIONS //
//////////////////////////////

const (
    MAKE_BLUE = "\033[34m"
    MAKE_YELLOW = "\033[33m"
    MAKE_GREEN = "\033[32m"
    MAKE_RED = "\033[31m"
    MAKE_BOLD = "\033[1m"
    CLEAR_COLOR = "\033[0m"
    LIGHT_GRAY_VALUE = "192"
    MAKE_LIGHT_GRAY = "\033[38;2;" + LIGHT_GRAY_VALUE + ";" + LIGHT_GRAY_VALUE + ";" + LIGHT_GRAY_VALUE + "m"
    MEDIUM_GRAY_VALUE = "140"
    MAKE_MEDIUM_GRAY = "\033[38;2;" + MEDIUM_GRAY_VALUE + ";" + MEDIUM_GRAY_VALUE + ";" + MEDIUM_GRAY_VALUE + "m"
    DARK_GRAY_VALUE = "100"
    MAKE_DARK_GRAY = "\033[38;2;" + DARK_GRAY_VALUE + ";" + DARK_GRAY_VALUE + ";" + DARK_GRAY_VALUE + "m"
    CAUTION = '⚠'
)

func printGray(str string, bold bool) {
    if bold {
        fmt.Printf("\033[1m")
    }
    fmt.Printf(MAKE_LIGHT_GRAY)
    fmt.Printf(str)
    fmt.Printf(CLEAR_COLOR)
}

func printError(str string) {
    fmt.Fprintf(os.Stderr,"%s%c  %s%s\n", MAKE_RED, CAUTION, str, CLEAR_COLOR)
}

func printErr(err error) {
    // TODO: make this more simple by eliminating new random line problems
    errStrings := strings.Split(string(err.Error()), "\n")
    trimmed := []string{}
    for _, er := range errStrings {
        if strings.TrimSpace(er) != "" {
            trimmed = append(trimmed, er)
            // trimmed = append(trimmed, strings.TrimSpace(er))
        }
    }
    for i, er := range trimmed {
        if i == 0 {
            trimmed[i] = fmt.Sprintf(" %s%c  %s%s\n", MAKE_RED, CAUTION, er, CLEAR_COLOR)
        } else if i == len(trimmed) - 1 {
            trimmed[i] = fmt.Sprintf("      %s%s%s%s\n", MAKE_RED, MAKE_BOLD, er, CLEAR_COLOR)
        } else {
            trimmed[i] = fmt.Sprintf(" %s%c  %s%s\n", MAKE_RED, CAUTION, er, CLEAR_COLOR)
        }
    }
    fmt.Fprintf(os.Stderr, "%s", strings.Join(trimmed, ""))

    // // Capture the file and line number of where the error occurred
    // _, file, line, ok := runtime.Caller(1)
    // if ok {
    //     fmt.Fprintf(os.Stderr, "%s%c  %s (at %s:%d)%s\n", MAKE_RED, CAUTION, err, file, line, CLEAR_COLOR)
    // } else {
    //  // Fallback to just printing the error if runtime.Caller fails
    // }
}


func getEntryString(entry os.DirEntry, isChanged bool, depth int, lines []bool, siblings bool, first bool, last bool, inGoverse bool) (string) {
    entryString := ""
    entryString += MAKE_MEDIUM_GRAY

    for i := 0; i < depth; i++ {
        if lines[i] == true {
            entryString += "│  "
        } else {
            entryString += "   "
        }
    }

    // prefix
    if first {
        entryString += "┌─ "
    } else if last {
        entryString += "└─ "
    } else if !siblings {
        entryString += "├─ "
    } else {
        entryString += "└─ "
    }


    // name
    if entry.IsDir() {
        entryString += MAKE_BLUE + MAKE_BOLD
    } else {
        if !isChanged || inGoverse{
            entryString += MAKE_GREEN
        } else {
            entryString += MAKE_RED
        }
    }
    entryString += entry.Name()

    entryString += CLEAR_COLOR
    entryString += "\n"

    return entryString
}



func printGoverse(path string, depth int, lines []bool, showGoverse bool, inGoverse bool) {
    ls, err := os.ReadDir(path)
    if err != nil {
        printErr(err)
    }
    first := true && depth == 0
    for i, entry := range ls {
        if entry.Name() == core.GOVERSE && !showGoverse { 
            continue
        }

        // goverse dir should not be checked for changes
        if path == core.BaseDir + core.GOVERSE_DIR {
            inGoverse = true
        }
        // if entry.Name() == core.GOVERSE {
            // println(inGoverse)
            // inGoverse = true
            // println(entry.Name())
        // }
        siblings := i != len(ls) - 1
        if depth >= len(lines) {
			lines = append(lines, siblings)
		} else {
			lines[depth] = siblings
		}
        // println("yep path: " + path + entry.Name())
        changed, _ := core.CheckChanged(path + entry.Name())
        fmt.Print(getEntryString(entry, changed, depth, lines, !siblings, first, false, inGoverse))
        first = false
        if entry.IsDir() {
            printGoverse(path + entry.Name() + "/", depth + 1, lines, showGoverse, inGoverse)
        }
    }
}

func getFile(reader *bufio.Reader) (string) {
    printGoverse(core.BaseDir, 0, []bool{true}, false, false)
    printGray("Add file: ", false)
    fmt.Print(MAKE_GREEN)
    file, _ := reader.ReadString('\n')
    fmt.Print(CLEAR_COLOR)
    return file
}


func printHelp() {
    printGray("valid commands:\n", true)
    printGray("  i   init\tInit cwd as a repository\n", false)
    printGray("  p   print\tPrint your entire project directory\n", false)
    printGray("  pr  print\tPrint your repository\n", false)
    printGray("  pg  print\tPrint your .goverse directory\n", false)
    printGray("  a   add\tAdd file to next commit\n", false)
    printGray("  s   status\tCheck repository status\n", false)
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
    reader := bufio.NewReader(os.Stdin)
    for running {

        fmt.Print(MAKE_YELLOW)
        fmt.Print("goverse command > ")
        fmt.Print(CLEAR_COLOR)

        fmt.Print(MAKE_GREEN)
        cmd, _ := reader.ReadString('\n')
        fmt.Print(CLEAR_COLOR)

        cmd = strings.Trim(cmd, "\n")

        switch(cmd) {
        case "i", "init":
            err := core.InitGoverse()
            if err != nil {
                printErr(err)
            }
        case "p", "print":
            printGoverse(core.BaseDir, 0, []bool{true}, true, false)
        case "pg":
            printGoverse(core.BaseDir + core.GOVERSE_DIR, 0, []bool{true}, true, false)
        case "pr":
            printGoverse(core.BaseDir, 0, []bool{false}, false, false)
        case "a", "add":
            err := core.Add(getFile(reader))
            if err != nil {
                printErr(err)
            }
        case "s", "status":
            err := core.Status()
            if err != nil {
                printErr(err)
            }
        case "d", "diff":
        case "t", "tag":
        case "c", "commit":
        case "l", "log":
        case "f", "flush":
            core.Flush()
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
        if exit {os.Exit(1)}
        return
    }
    core.BaseDir = args[1]
    fmt.Println("Base dir: ", core.BaseDir)
    interactive()

}
