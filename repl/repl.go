package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
	"n128.xyz/n128/envo/lexer"
	"n128.xyz/n128/envo/object"
	"n128.xyz/n128/envo/parser"
	"n128.xyz/n128/envo/runtime"
)

const PROMPT = "$> "

func Start(in io.Reader, out io.Writer) {
	env := runtime.NewEnvironment()
	fmt.Fprintln(out, "Envo REPL - Write 'exit' or press Ctrl+C to exit.")

	history := make([]string, 0, 50)
	historyIndex := 0

	for {
		line, ok := readLine(in, out, PROMPT, &history, &historyIndex)
		if !ok {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" {
			fmt.Fprintln(out, "bye!")
			break
		}

		l := lexer.New(strings.NewReader(line))
		p := parser.New(l)

		program, parseErrors := p.Parse()
		if len(parseErrors) > 0 {
			printParserErrors(out, parseErrors)
			continue
		}

		evaluated := runtime.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NIL_OBJ {
			fmt.Fprintln(out, evaluated.String())
		}
	}
}

func readLine(in io.Reader, out io.Writer, prompt string, history *[]string, historyIndex *int) (string, bool) {
	if file, ok := in.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		return readLineTerminal(file, out, prompt, history, historyIndex)
	}

	scanner := bufio.NewScanner(in)
	fmt.Fprint(out, prompt)
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}

func readLineTerminal(in *os.File, out io.Writer, prompt string, history *[]string, historyIndex *int) (string, bool) {
	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return "", false
	}
	defer term.Restore(int(in.Fd()), oldState)

	fmt.Fprint(out, prompt)
	reader := bufio.NewReader(in)
	buffer := []rune{}
	cursorPos := 0
	currentIndex := len(*history)
	*historyIndex = currentIndex

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", false
		}

		switch b {
		case '\r', '\n':
			fmt.Fprint(out, "\r\n")
			line := string(buffer)
			if strings.TrimSpace(line) != "" {
				if len(*history) == 0 || (*history)[len(*history)-1] != line {
					*history = append(*history, line)
				}
			}
			*historyIndex = len(*history)
			return line, true
		case 0x7f, 0x08:
			if cursorPos > 0 {
				buffer = append(buffer[:cursorPos-1], buffer[cursorPos:]...)
				cursorPos--
				refreshLine(out, prompt, buffer, cursorPos)
			}
		case 0x03: // CTRL+C
			fmt.Fprint(out, "\r\n")
			return "", false
		case 0x1b: //
			seq, err := reader.Peek(2)
			if err != nil {
				continue
			}
			r1, r2 := seq[0], seq[1]
			if r1 == '[' {
				_, _ = reader.Discard(2)
				switch r2 {
				case 'A':
					if len(*history) == 0 {
						continue
					}
					if *historyIndex > 0 {
						*historyIndex--
						buffer = []rune((*history)[*historyIndex])
						cursorPos = len(buffer)
						refreshLine(out, prompt, buffer, cursorPos)
					}
				case 'B':
					if *historyIndex < len(*history) {
						*historyIndex++
					}
					if *historyIndex < len(*history) {
						buffer = []rune((*history)[*historyIndex])
					} else {
						buffer = []rune{}
					}
					cursorPos = len(buffer)
					refreshLine(out, prompt, buffer, cursorPos)
				case 'C':
					if cursorPos < len(buffer) {
						cursorPos++
						refreshLine(out, prompt, buffer, cursorPos)
					}
				case 'D':
					if cursorPos > 0 {
						cursorPos--
						refreshLine(out, prompt, buffer, cursorPos)
					}
				}
			}
		default:
			if b >= 0x20 && b <= 0x7e {
				if cursorPos < len(buffer) {
					buffer = append(buffer[:cursorPos], append([]rune{rune(b)}, buffer[cursorPos:]...)...)
				} else {
					buffer = append(buffer, rune(b))
				}
				cursorPos++
				refreshLine(out, prompt, buffer, cursorPos)
			}
		}
	}
}

func refreshLine(out io.Writer, prompt string, buffer []rune, cursorPos int) {
	fmt.Fprintf(out, "\r%s%s", prompt, string(buffer))
	fmt.Fprint(out, "\x1b[0K")
	fmt.Fprintf(out, "\r%s", prompt)
	for i := 0; i < cursorPos; i++ {
		fmt.Fprint(out, string(buffer[i]))
	}
}

func printParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, err := range errors {
		fmt.Fprintf(out, "  [ syntax error ]: %s\n", err)
	}
}
