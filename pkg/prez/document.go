package prez

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Document struct {
	Title  string
	Author string
	Lang   string
	Pages  [][]Block
}

type Block interface{ _guardBlock() }

type H1 string
type Text string
type Image string
type Code string
type ListItem string

func (H1) _guardBlock()       {}
func (Text) _guardBlock()     {}
func (Image) _guardBlock()    {}
func (Code) _guardBlock()     {}
func (ListItem) _guardBlock() {}

func ParseFile(fpath string) (doc *Document, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	defer f.Close()
	dirpath := filepath.Dir(fpath)

	doc = &Document{}
	r := bufio.NewReader(f)

	// Parse metadata.
	line, err := r.ReadString('\n') // Read first line ("---").
	if err != nil {
		return nil, fmt.Errorf("read metadata start line: %w", err)
	} else if strings.TrimSpace(line) != "---" {
		return nil, fmt.Errorf("unexpected metadata start line: %q", line)
	}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read metadata line: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "---" {
			break
		}
		k, v, _ := strings.Cut(line, ":")
		switch strings.ToLower(k) {
		default:
			continue // Ignore unknown metadata fields.
		case "title":
			doc.Title = strings.TrimSpace(v)
		case "author":
			doc.Author = strings.TrimSpace(v)
		case "lang":
			doc.Lang = strings.TrimSpace(v)
		}
	}

	// Parse pages.
	var page []Block
	for {
		line, err = r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("read line: %w", err)
		}
		line = strings.TrimSpace(line)
		switch {
		default:
			page = append(page, Text(line))
		case line == "---":
			doc.Pages = append(doc.Pages, page)
			page = nil
		case line == "":
			continue // Ignore empty lines.
		case strings.HasPrefix(line, "# "):
			page = append(page, H1(strings.TrimPrefix(line, "# ")))
		case strings.HasPrefix(line, "!["):
			line = strings.TrimPrefix(line, "![")
			line = strings.TrimSuffix(line, ")")
			_, fpath, _ := strings.Cut(line, "](")
			page = append(page, Image(filepath.Join(dirpath, fpath)))
		case strings.HasPrefix(line, "```"):
			code := ""
			for {
				line, err = r.ReadString('\n')
				if errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					return nil, fmt.Errorf("read line: %w", err)
				} else if line == "```\n" {
					break
				} else if strings.HasPrefix(line, "````") {
					line = line[1:] // Support escaping triple-backtick ("````" results in "```").
				}
				code += line
			}
			page = append(page, Code(code))
		case strings.HasPrefix(line, "- "):
			page = append(page, ListItem(strings.TrimPrefix(line, "- ")))
		}
	}
	if page != nil {
		doc.Pages = append(doc.Pages, page)
	}

	return doc, nil
}
