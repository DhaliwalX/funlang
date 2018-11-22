package lex

import (
    "io/ioutil"
    "os"
    "unicode/utf8"
)

// Source represents the source code of a program
type Source interface {
    ReadChar() rune
    Pos() Position
}

// FileSource represents the source code present in a file
type FileSource struct {
    path string
    pos Position
    source string
    currentPos int
}

func NewStringSource(source string) *FileSource {
    return &FileSource{ pos: Position{}, source: source}
}

func NewFileSource(path string) *FileSource {
    f, err := os.Open(path)
    if err != nil {
        panic(path + " does not exist")
    }

    defer f.Close()
    bytes, err := ioutil.ReadAll(f)
    if err != nil {
        panic(path + ": unable to read file")
    }

    source := &FileSource{ path: path, pos: Position{}, source: string(bytes) }
    return source
}

func (source *FileSource) ReadChar() rune {
    if source.EOF() {
        return -1
    }
    runeValue, width := utf8.DecodeRuneInString(source.source[source.currentPos:])
    source.currentPos += width
    return runeValue
}

func (source *FileSource) EOF() bool {
    return len(source.source) <= source.currentPos
}

func (source *FileSource) Pos() Position {
    return source.pos
}
