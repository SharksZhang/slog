package model

import (
	"bytes"
	"strconv"
)

type Log struct {
	Type        string `json:"type"`
	FileSizeKb  int    `json:"file_size_KB"`
	FilterLevel string `json:"filter_level"`

	Action string `json:"action"`
	Name   string `json:"name"`

	Level         string `json:"level"`
	Time          string `json:"time"`
	Source        string `json:"source"`
	File          string `json:"file"`
	Line          int    `json:"line"`
	Function      string `json:"function"`
	Content       string `json:"content"`
}

func (l *Log) transferToString() string {
	var buf bytes.Buffer
	buf.WriteString(l.Time)
	buf.WriteString("\t")
	buf.WriteString(l.Level)
	buf.WriteString("\t")
	buf.WriteString(l.Source)
	buf.WriteString("\t")
	buf.WriteString(l.File)
	buf.WriteString("\t")
	buf.WriteString(strconv.Itoa(l.Line))
	buf.WriteString("\t")
	buf.WriteString(l.Function)
	buf.WriteString("\t")
	buf.WriteString(l.Content)
	buf.WriteString("\n")
	return buf.String()
}
