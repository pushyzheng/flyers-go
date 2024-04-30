package flyers

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Scanner interface {
	Scan(fn func(line string) bool) error
	GetAll() ([]string, error)
}

type reader struct {
	r io.Reader
}

func NewScanner(r io.Reader) Scanner {
	return &reader{r: r}
}

func NewStdinScanner() Scanner {
	return &reader{r: os.Stdin}
}

func (s *reader) Scan(fn func(line string) bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		continueFlag := fn(line)
		if !continueFlag {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}
	return nil
}

func (s *reader) GetAll() ([]string, error) {
	var res []string
	err := s.Scan(func(line string) bool {
		res = append(res, line)
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
