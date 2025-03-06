// Stack represents a go stacktrace.
package stack

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

type Stack struct {
	R   []*Routine // List of deduplicated goroutines
	All []*Routine // List of all goroutines
}

type Routine struct {
	N       int           // Goroutine number
	State   string        // State
	Blocked time.Duration // Time blocked
	CS      []*CallStack  // Callstack

	Ns           []int           // List of goroutine IDs
	States       []string        // List of states if deduplicated
	BlockedRange []time.Duration // List of times blocked if deduplicated
	I            []int           // Indexes into the parent All array that make up this deduplicated Goroutine
}

type CallStack struct {
	F        string // Function
	Filename string // Filename
	Line     int    // Line number in Filename
}

// Decode and return a Stack object from an io.Reader. Decode will run until
// EOF or an error is reached.
func Decode(in io.Reader) (*Stack, error) {
	s := new(Stack)
	var err error
	var r *Routine
	scanner := bufio.NewScanner(in)
	for {
		r, err = decodeRoutine(scanner)
		if err != nil {
			break
		}
		s.All = append(s.All, r)
	}
	if err != io.EOF {
		return nil, err
	}

	s.dedup()
	return s, nil
}

func decodeRoutine(s *bufio.Scanner) (*Routine, error) {
	return stateEntry(s)
}

func stateEntry(s *bufio.Scanner) (*Routine, error) {
	r := new(Routine)

	for s.Scan() {
		// goroutine DDDD [state]:
		re := regexp.MustCompile(`goroutine (\d+) [^\[]+\[([^\[]+)\]:`)
		m := re.FindStringSubmatch(s.Text())
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
		r.N = n
		r.State = m[2]
		return stateFunction(r, s)
	}
	if s.Err() == nil {
		return r, io.EOF
	}
	return r, s.Err()
}

func stateFunction(r *Routine, s *bufio.Scanner) (*Routine, error) {
	if !s.Scan() {
		return r, nil
	} else if s.Text() == "" {
		return r, nil
	}

	re := regexp.MustCompile(`(.+)\([^\)]*\)?`)
	m := re.FindStringSubmatch(s.Text())
	data := s.Text()
	if m != nil {
		data = m[1]
	}
	c := new(CallStack)
	c.F = data
	return stateFile(c, r, s)
}

func stateFile(c *CallStack, r *Routine, s *bufio.Scanner) (*Routine, error) {
	if !s.Scan() || s.Text() == "" {
		return r, fmt.Errorf("could not read file line")
	}

	re := regexp.MustCompile(`([^:]+):(\d+)`)
	m := re.FindStringSubmatch(s.Text())
	if m == nil {
		return nil, fmt.Errorf("invalid input %v", s.Text())
	}
	c.Filename = m[1]
	l, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, err
	}
	c.Line = l
	r.CS = append(r.CS, c)
	return stateFunction(r, s)
}

func (s *Stack) dedup() {
	for i, v := range s.All {
		var found bool
		for _, w := range s.R {
			if w.CSString() == v.CSString() {
				w.Ns = append(w.Ns, v.N)

				var foundState bool
				for _, x := range w.States {
					if x == v.State {
						foundState = true
						break
					}
				}
				if !foundState {
					w.States = append(w.States, v.State)
				}
				var foundBlocked bool
				for _, x := range w.BlockedRange {
					if x == v.Blocked {
						foundBlocked = true
						break
					}
				}
				if !foundBlocked {
					w.BlockedRange = append(w.BlockedRange, v.Blocked)
				}
				w.I = append(w.I, i)
				found = true
				break
			}
		}
		if !found {
			s.R = append(s.R, &Routine{
				CS:           v.CS,
				Ns:           []int{v.N},
				States:       []string{v.State},
				BlockedRange: []time.Duration{v.Blocked},
				I:            []int{i},
			})
		}
	}
}

func (s *Stack) String() string {
	var z string
	for _, v := range s.R {
		z += v.String()
	}
	return z
}

func (r *Routine) String() string {
	z := fmt.Sprintf("goroutine %v %v %v\n", r.Ns, r.States, r.BlockedRange)
	for _, v := range r.CS {
		z += v.String()
	}
	z += "\n"
	return z
}

func (r *Routine) CSString() string {
	var z string
	for _, v := range r.CS {
		z += v.String()
	}
	z += "\n"
	return z
}

func (cs *CallStack) String() string {
	return fmt.Sprintf("%v\n%v:%v\n", cs.F, cs.Filename, cs.Line)
}
