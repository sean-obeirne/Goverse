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

    rootTree := models.Tree {}
    err := readFiles(BaseDir, &rootTree)
    if err != nil {
        return err
    }
    // root, err := os.Create(BaseDir + OBJECTS_DIR + rootTree.Hash)
    // if err != nil {
    //     return err
    // }
    // root.Close()
    return nil
}


func readFiles(path string, tree *models.Tree) (error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return err
    }
    
    for _, entry := range entries {
        if entry.Name() == GOVERSE {continue}


        thisPath := path + entry.Name() 
        if entry.IsDir() {
            thisPath += "/"
        }


        te, err := createTreeEntry(thisPath)
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
                Content: content,
            }
            blobHash, _ := hashBlob(b)
            fmt.Printf("entry hash: %s\n blob hash: %s\n", te.Hash, blobHash)
        } else {
            t := models.Tree {
                Entries: []models.TreeEntry{},
            }
            err := readFiles(thisPath, &t)
            if err != nil {
                return err
            }
            treeHash, _ := hashTree(t)
            fmt.Printf("entry hash: %s\n tree hash: %s\n", te.Hash, treeHash)
        }
    }
    return nil
}

func getHash(str string) (string, error) {
    hasher := sha1.New()

    _, err := hasher.Write([]byte(str))
    if err != nil{
        return "", err
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
        return "", err
    }
    var entryHash string
    for _, entry := range entries {
        if entry.IsDir() {
            entryHash, err = hashDir(path + entry.Name() + "/")
        } else {
            entryHash, err = hashFile(path + entry.Name())
        }
        if err != nil {
            return "", err
        }
        hash += entryHash
    }

    return getHash(hash)
}

func hashFile(path string) (string, error){
    content, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }

    return getHash(string(content))
}

func createTreeEntry(path string) (models.TreeEntry, error) {
    println("\tcreateTreeEntry path: " + path)

    fileInfo, err := os.Stat(path)
    if err != nil {
        return models.TreeEntry{}, err
    }

    var hash string
    if fileInfo.IsDir() {
        hash, err = hashDir(path)
    } else {
        hash, err = hashFile(path)
    }
    if err != nil {
        return models.TreeEntry{}, err
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
