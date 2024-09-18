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

type Tag struct {
    Name    string
    Version string
    Commit  string
}

// Storage structs
type Blob struct {
    Content []byte
}

type TreeEntry struct {
    Name      string
    Mode      string
    Hash      string
    IsBlob    bool
    Modified  bool
    Tracked   bool
}

type Tree struct {
    Entries []TreeEntry
}
