package models

import (
    // "fmt"
    // "os"
    // "strings"
    // "errors"
)


/////////////
// STRUCTS //
/////////////

// Conceptual structs
type Commit struct {
    Hash      string
    Message   string
    Author    string
    Timestamp string
}

type FileData struct {
    Path      string
    Modified  bool
    Tracked   bool
}

type Tag struct {
    Name    string
    Version string
    Commit  string
}

// Storage structs
type Blob struct {
    Hash    string
    Content []byte
}

type TreeEntry struct {
    Name    string
    Mode    string
    Hash    string
    IsBlob  bool
}

type Tree struct {
    Hash    string
    Entries []TreeEntry
}
