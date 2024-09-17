package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"errors"
    "runtime"

    "goverse/internal/models"
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
    // Capture the file and line number of where the error occurred
    _, file, line, ok := runtime.Caller(1)
    if ok {
        fmt.Fprintf(os.Stderr, "%s%c  %s (at %s:%d)%s\n", MAKE_RED, CAUTION, err, file, line, CLEAR_COLOR)
    } else {
        // Fallback to just printing the error if runtime.Caller fails
fmt.Fprintf(os.Stderr, "%s%c  %s%s\n", MAKE_RED, CAUTION, err, CLEAR_COLOR)
    }
}

///////////////////
// COMMAND BLOCK //
///////////////////

// Paths
var baseDir string
const GOVERSE = ".goverse"
const (
    GOVERSE_DIR = ".goverse/"
    OBJECTS_DIR = GOVERSE_DIR + "objects/"
    TAGS_DIR    = GOVERSE_DIR + "tags/"
    CONFIG_FILE = GOVERSE_DIR + "config"
    HEAD_FILE   = GOVERSE_DIR + "head"
)
// test dirs
const TD4 = OBJECTS_DIR + "another/yep"
const TD1 = OBJECTS_DIR + "another/file/smd/lol/lmao"
const TD2 = OBJECTS_DIR + "another/file/smd/lol/haha"
const TD3 = OBJECTS_DIR + "another/file/to/fuck/with"
const TD5 = OBJECTS_DIR + "another/yep/ha"

func initGoverse() {

    // Create necessary dirs and files for goverse VCS
    dirs := []string{ GOVERSE_DIR, OBJECTS_DIR, TAGS_DIR }
    // dirs = append(dirs, TD1, TD2, TD3, TD4, TD5)
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

    rootTree := models.Tree {
        Hash: "0",
    }
    readFiles(baseDir, rootTree)
    root, err := os.Create(baseDir + OBJECTS_DIR + rootTree.Hash)
    if err != nil {
        printErr(err)
    }
    root.Close()
}

func readFiles(path string, tree models.Tree) {
    // Establish FileData for each file in project
    files, err := os.ReadDir(path)
    if err != nil {
        printErr(err)
    }
    
    for _, file := range files {
        if file.Name() == GOVERSE {continue}


        thisPath := path + file.Name() 
        if file.IsDir() {
            thisPath += "/"
        }
        


        hashString, te := store(thisPath)

        tree.Entries = append(tree.Entries, te)



        if te.IsBlob {
            content, err := os.ReadFile(thisPath)
            if err != nil {
                printErr(err)
            }
            b := models.Blob {
                Modified: false,
                Tracked: false,
                Hash: hashString,
                Content: content,
            }
            fmt.Printf("entry hash: %s\n blob hash: %s\n", te.Hash, b.Hash)
        } else {
            t := models.Tree {
                Hash: hashString,
                Entries: []models.TreeEntry{},
            }
            fmt.Printf("entry hash: %s\n tree hash: %s\n", te.Hash, t.Hash)
            if file.IsDir() {
                readFiles(thisPath + "/", t)
            }
        }

    }
}

func store(path string) (string, models.TreeEntry) {
    file, err := os.Open(path)
    if err != nil {
        printErr(err)
    }
    defer file.Close()

    hasher := sha1.New()

    fileInfo, err := file.Stat()
    if err != nil {
        printErr(err)
    }
    if fileInfo.IsDir() {
        // fill hasher with dir path
        _, err := hasher.Write([]byte(path))
        if err != nil {
            printErr(err)
        }
    } else {
        // fill hasher with file content
        _, err := io.Copy(hasher, file)
        if err != nil {
            printErr(err)
        }
    }


    // finalize hash
    hashBytes := hasher.Sum(nil)
    hashString := hex.EncodeToString(hashBytes)

    perms := fmt.Sprintf("%o", fileInfo.Mode())

    te := models.TreeEntry {
        Name: fileInfo.Name(),
        Mode: perms,
        Hash: hashString,
        IsBlob: !fileInfo.IsDir(),
    }



    // println(te.Name)
    // println(te.IsBlob)

    // println("String: " + path + "\nHash: " + hashString)
    
    return hashString, te
}


func getEntryString(entry os.DirEntry, depth int, lines []bool, siblings bool, first bool, last bool) (string) {
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
        entryString += MAKE_GREEN
    }
    entryString += entry.Name()

    entryString += CLEAR_COLOR
    entryString += "\n"

    return entryString
}


func getEntryString2(entry os.DirEntry, depth int, lines []int, siblings bool, first bool, last bool) (string) {
    entryString := ""
    entryString += MAKE_MEDIUM_GRAY
    if depth > 0 {
        entryString += " │"
        if len(lines) > 0 && siblings {
            for _, line := range lines {
                entryString += strings.Repeat(" ", line + 1)
                entryString += "│"
                entryString += strings.Repeat(" ", 0)
            }
        } else {
            entryString += strings.Repeat(" ", depth * 2)
        }

        // entryString += MAKE_BOLD
        if siblings {
            entryString += "├" 
        } else { 
            entryString += "└"
        }
        if entry.IsDir() {

        }
        // entryString += "─┐ "
        entryString += "─┬ "
        // entryString += "  ┗━ "
        // entryString += "┖─ "
        entryString += CLEAR_COLOR
    } else {
        if first {
            entryString += " ┌─ "
        } else if last { 
            entryString += " └─ "
        } else {
            entryString += " ├─ "
        }
        // entryString += " "
    }
    if entry.IsDir(){
        entryString += MAKE_BLUE
        entryString += MAKE_BOLD
    } else{
        entryString += MAKE_GREEN
    }
    // if entry.IsDir() {
        // entryString += "d "
    // } else {
        // entryString += "f "
    // }
    entryString += entry.Name()
    entryString += "\n"
    entryString += CLEAR_COLOR
    // if entry.Name() == "0" {
    // }
    return entryString
}

func checkHealth(path string, depth int, lines []bool) {
    ls, err := os.ReadDir(path)
    if err != nil {
        printErr(err)
    }
    first := true && depth == 0
    for i, entry := range ls {
        siblings := i != len(ls) - 1
        if depth >= len(lines) {
			lines = append(lines, siblings)
		} else {
			lines[depth] = siblings
		}
        // print(lines[1])
        // print(i)
        // siblings = false
        fmt.Print(getEntryString(entry, depth, lines, !siblings, first, false))
        first = false
        if entry.IsDir() {
            checkHealth(path + entry.Name() + "/", depth + 1, lines)
        }
    }
}

func add() {
    
}

func status() {

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
    printGray("  e   check\tCheck .goverse/ for entries\n", false)
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
        case "e", "checkHealth":
            checkHealth(baseDir + GOVERSE_DIR, 0, []bool{true})
        case "a", "add":
        case "s", "status":
            status()
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

}
