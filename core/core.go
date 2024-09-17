package core

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "goverse/internal/models"
)

///////////////////
// COMMAND BLOCK //
///////////////////

// Paths
var BaseDir string
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


func InitGoverse() (error) {

    // Create necessary dirs and files for goverse VCS
    dirs := []string{ GOVERSE_DIR, OBJECTS_DIR, TAGS_DIR }
    // dirs = append(dirs, TD1, TD2, TD3, TD4, TD5)
    files := []string{ CONFIG_FILE, HEAD_FILE }
    for _, dir := range dirs {
        err := os.MkdirAll(BaseDir + dir, 0755)
        if err != nil {
            return err
        }
    }
    for _, file := range files {
        newFile, err := os.Create(BaseDir + file)
        if err != nil {
            return err
        }
        newFile.Close()
    }

    rootTree := models.Tree {
        Hash: "0",
    }
    readFiles(BaseDir, rootTree)
    root, err := os.Create(BaseDir + OBJECTS_DIR + rootTree.Hash)
    if err != nil {
        return err
    }
    root.Close()
    return nil
}


func readFiles(path string, tree models.Tree) (error) {
    // Establish FileData for each file in project
    files, err := os.ReadDir(path)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if file.Name() == GOVERSE {continue}


        thisPath := path + file.Name() 
        if file.IsDir() {
            thisPath += "/"
        }
        


        hashString, te, err := store(thisPath)
        if err != nil {
            return err
        }

        tree.Entries = append(tree.Entries, te)



        if te.IsBlob {
            content, err := os.ReadFile(thisPath)
            if err != nil {
                return err
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
    return nil
}


func store(path string) (string, models.TreeEntry, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", models.TreeEntry{}, err
    }
    defer file.Close()

    hasher := sha1.New()

    fileInfo, err := file.Stat()
    if err != nil {
        return "", models.TreeEntry{}, err
    }
    if fileInfo.IsDir() {
        // fill hasher with dir path
        _, err := hasher.Write([]byte(path))
        if err != nil {
            return "", models.TreeEntry{}, err
        }
    } else {
        // fill hasher with file content
        _, err := io.Copy(hasher, file)
        if err != nil {
            return "", models.TreeEntry{}, err
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
    
    return hashString, te, nil
}


func Add() {
    
}

func Status() {

}

func Diff() {
    
}

func Tag() {
    
}

func Commit() {
    
}

func Log() {
    
}
    
func Flush() (error) {
    err := os.RemoveAll(BaseDir + GOVERSE_DIR)
    if err != nil {
        return err
    }
    return nil
}
