package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"goverse/internal/models"
	// "hash"
	// "io"
	"os"
	// "strings"
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
const TD1 = OBJECTS_DIR + "another/file/smd/lol/lmao"
const TD2 = OBJECTS_DIR + "another/file/smd/lol/haha"
const TD3 = OBJECTS_DIR + "another/file/to/fuck/with"
const TD4 = OBJECTS_DIR + "another/yep"
const TD5 = OBJECTS_DIR + "another/yep/ha"


func InitGoverse() (error) {

    // Create necessary dirs and files for goverse VCS
    dirs := []string{ GOVERSE_DIR, OBJECTS_DIR, TAGS_DIR }
    // dirs = append(dirs, TD1, TD2, TD3, TD4, TD5)
    files := []string{ CONFIG_FILE, HEAD_FILE }
    for _, dir := range dirs {
        err := os.MkdirAll(BaseDir + dir, 0755)
        if err != nil {
            return fmt.Errorf("Unable to create dir at \"%s\"\n%w\n", BaseDir + dir, err)
        }
    }
    for _, file := range files {
        newFile, err := os.Create(BaseDir + file)
        if err != nil {
            return fmt.Errorf("Unable to create file at \"%s\"\n%w\n", BaseDir + file, err)
        }
        newFile.Close()
    }

    rootTree := models.Tree {}
    err := readFiles(BaseDir, &rootTree)
    if err != nil {
        return fmt.Errorf("Unable to read files at BaseDir \"%s\" into rootTree\n%w\n", BaseDir, err)
    }
    // println("Entries: ")
    // root, err := os.Create(BaseDir + OBJECTS_DIR + rootTree.Hash)
    // if err != nil {
    //     return err
    // }
    // root.Close()
    return nil
}

/*
func printTree(t models.Tree) {
    for _, entry := range t.Entries {
        println(entry.Name)
        if !entry.IsBlob {
            printTree(entry)
        }
    }
}
*/

func readFiles(path string, tree *models.Tree) (error) {
    // read directory
    entries, err := os.ReadDir(path)
    if err != nil {
        return fmt.Errorf("Unable to read files at BaseDir \"%s\"\n%w\n", BaseDir, err)
    }
    
    // loop through directory
    for _, entry := range entries {

        // skip meta goverse directory
        if entry.Name() == GOVERSE {continue}

        // get new entry's path
        thisPath := path + entry.Name() 
        if entry.IsDir() {
            thisPath += "/"
        }

        // create tree entry for this entry
        te, err := createTreeEntry(thisPath)
        if err != nil {
            return fmt.Errorf("Unable to create TreeEntry at path: %s\n%w\n", thisPath, err)
        }
        
        // add new tree entry to tree
        tree.Entries = append(tree.Entries, te)


        // store blob content, or tree with tree entrys populated
        if te.IsBlob {
            content, err := os.ReadFile(thisPath)
            if err != nil {
                return fmt.Errorf("Unable to read Blob at path: %s\n%w\n", thisPath, err)
            }
            b := models.Blob {
                Content: content,
            }
            err = storeBlob(b)
            if err != nil {
                return fmt.Errorf("Unable to store Blob at path: %s\n%w\n", thisPath, err)
            }

            // newBlobHash, _ := hashBlob(b)
            // newContent, _ := getContent(newBlobHash)
            // println(string(newContent))
        } else {
            t := models.Tree {
                Entries: []models.TreeEntry{},
            }
            err := readFiles(thisPath, &t)
            if err != nil {
                return fmt.Errorf("Unable to read Tree at path: %s\n%w\n", thisPath, err)
            }
            treeHash, err := hashTree(t)
            if err != nil {
                return fmt.Errorf("Unable to hash Tree, got hash %s\n%w\n", truncHash(treeHash), err)
            }
            fmt.Printf("entry hash: %s\n tree hash: %s\n", truncHash(te.Hash), truncHash(treeHash))
        }
    }
    return nil
}

func storeBlob(b models.Blob) (error) {
    hashString, err := hashBlob(b)
    if err != nil {
        return fmt.Errorf("Unable to hash Blob, got hash %s\n%w\n", truncHash(hashString), err)
    }
    err = os.WriteFile(BaseDir + OBJECTS_DIR + hashString, b.Content, 0755)
    if err != nil {
        return fmt.Errorf("Unable to store Blob at %s\n%w\n", BaseDir + OBJECTS_DIR + truncHash(hashString), err)
    }

    return nil
}

func getContent(hash string) ([]byte, error) {
    files, err := os.ReadDir(BaseDir + OBJECTS_DIR)
    if err != nil {
        return nil, fmt.Errorf("Unable to read content at path: %s\n%w\n", BaseDir + OBJECTS_DIR, err)
    }
    for _, file := range files {
        if file.Name() == hash {
            return os.ReadFile(BaseDir + OBJECTS_DIR + hash)
        }
    }
    return nil, fmt.Errorf("Unable to find hash content at path: %s\n%w\n", BaseDir + OBJECTS_DIR + truncHash(hash), err)
}

func truncHash(hash string) (string) {
    return hash[:4] + "..."
}

func getHash(str string) (string, error) {
    hasher := sha1.New()

    _, err := hasher.Write([]byte(str))
    if err != nil {
        return "", fmt.Errorf("Unable to hash %s\n%w", str, err)
    }

    hashBytes := hasher.Sum(nil)
    hashString := hex.EncodeToString(hashBytes)

    return hashString, nil
}

func hashBlob(b models.Blob) (string, error) {
    return getHash(string(b.Content))
}

func hashTreeEntry(te models.TreeEntry) (string, error) {
    return getHash(string(te.Name + te.Mode + te.Hash))
}

func hashTree(t models.Tree) (string, error) {
    contents := ""
    for _, te := range t.Entries {
        contents += te.Hash
    }
    return getHash(contents)
}

func hashDir(path string) (string, error) {
    hash := ""
    entries, err := os.ReadDir(path)  // read all files and subdirs
    if err != nil {
        return "", fmt.Errorf("Unable to read files at path: %s\n%w\n", path, err)
    }
    var entryHash string
    for _, entry := range entries {
        if entry.IsDir() {
            entryHash, err = hashDir(path + entry.Name() + "/")
        } else {
            entryHash, err = hashFile(path + entry.Name())
        }
        if err != nil {
            return "", fmt.Errorf("Unable to hash %s\n%w", path + entry.Name(), err)
        }
        hash += entryHash
    }

    return getHash(hash)
}

func hashFile(path string) (string, error){
    content, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("Unable to read file at path: %s\n%w\n", path, err)
    }
    return getHash(string(content))
}

func createTreeEntry(path string) (models.TreeEntry, error) {
    // println("\tcreateTreeEntry path: " + path)

    fileInfo, err := os.Stat(path)
    if err != nil {
        return models.TreeEntry{}, fmt.Errorf("Unable to find file at path: %s, cannot create TreeEntry\n%w\n", path, err)
    }

    var hash string
    if fileInfo.IsDir() {
        hash, err = hashDir(path)
        if err != nil {
            return models.TreeEntry{}, fmt.Errorf("Unable to hash dir for TreeEntry at path: %s\n%w\n", path, err)
        }
    } else {
        hash, err = hashFile(path)
        if err != nil {
            return models.TreeEntry{}, fmt.Errorf("Unable to hash file for TreeEntry at path: %s\n%w\n", path, err)
        }
    }

    te := models.TreeEntry {
        Name: fileInfo.Name(),
        Mode: fmt.Sprintf("%o", fileInfo.Mode()),
        Hash: hash,
        IsBlob: !fileInfo.IsDir(),
        Modified: false,
        Tracked: false,
    }

    return te, nil
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
