package main

import (
	"io"
	"strings"
)

type CSVParserImpl struct {
	line   []byte
	fields []string

	lastByte byte
}

func (p *CSVParserImpl) ReadLine(r io.Reader) (string, error) {
	p.line = make([]byte, 0)
	p.fields = make([]string, 0)

	inQuotes := false
	firstQ := false
	newField := 0

	// Check for CR
	if p.lastByte != 0 {
		if p.lastByte == ',' {
			newField = 1
			p.fields = append(p.fields, string(p.line[0:len(p.line)]))
		}
		p.line = append(p.line, p.lastByte)
		p.lastByte = 0
	}

	for {
		buffer := make([]byte, 1)
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {

				if len(p.line) > 0 && !inQuotes {
					if p.line[len(p.line)-1] == '"' {
						p.line = p.line[:len(p.line)-1]
					}
					p.fields = append(p.fields, string(p.line[newField:len(p.line)]))
					if len(p.fields[len(p.fields)-1]) == 0 {
						return "", ErrLast
					}
					return string(p.line), nil
				} else if inQuotes {
					return "", ErrQuote
				}
				return "", io.EOF
			}
		}

		if n == 0 {
			return "", io.EOF // No data read
		}

		char := buffer[0]
		switch char {
		case '\r', '\n':
			if inQuotes {
				p.line = append(p.line, char)
				continue
			} else {
				if char == '\r' {
					peekBuf := make([]byte, 1)
					n, err := r.Read(peekBuf)
					if err != nil {
						return "", err
					}
					if n == 0 {
						return "", nil
					}
					if peekBuf[0] != '\n' {
						p.lastByte = peekBuf[0]
					}
				}
				if len(p.line) > 0 {
					if p.line[len(p.line)-1] == '"' {
						p.line = p.line[:len(p.line)-1]
					}
					p.fields = append(p.fields, string(p.line[newField:len(p.line)]))
					if len(p.fields[len(p.fields)-1]) == 0 {
						return "", ErrLast
					}
				}
				return string(p.line), nil
			}
		case '"': // previous symbol is not comma or quote if quote is over

			if !inQuotes && len(p.line) > 0 && p.line[len(p.line)-1] != '"' && p.line[len(p.line)-1] != ',' {
				return "", ErrQuote
			}
			if !inQuotes && (len(p.line) == 0 || p.line[len(p.line)-1] == ',') {
				firstQ = true
			}
			inQuotes = !inQuotes

		case ',':
			if !inQuotes {
				if len(p.line) > 0 && p.line[len(p.line)-1] == '"' {
					p.line = p.line[:len(p.line)-1]
				}
				p.fields = append(p.fields, string(p.line[newField:len(p.line)]))
				p.line = append(p.line, char)
				newField = len(p.line)
				continue
			}
		}
		if !firstQ {
			p.line = append(p.line, char)
		}
		firstQ = false
	}
}

func (p *CSVParserImpl) GetField(n int) (string, error) {
	if p.fields == nil {
		return "", io.EOF
	}
	if n < 0 || n >= len(p.fields) || len(p.fields) == 0 {
		return "", ErrFieldCount
	}
	p.fields[n] = strings.ReplaceAll(p.fields[n], "\"\"", "\"")

	return p.fields[n], nil
	/*
		Fields are separated by commas
		Fields may be surrounded by "..."; such quotes are removed
		There can be an arbitrary number of fields of any length
	*/
}

func (p *CSVParserImpl) GetNumberOfFields() int {
	if len(p.fields) > 0 {
		return len(p.fields)
	}
	return -1
}
