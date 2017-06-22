//
// command/display.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	termbox "github.com/nsf/termbox-go"
	"github.com/ttacon/chalk"
)

// Section represents a section in a display. Each section has a header text and
// several strings as the body.
type Section struct {
	Header  string
	Body    []string
	display *Display
}

// Writer returns io.WriteCloser to write messages into the section. Users have
// to close the returned writer.
func (s *Section) Writer() io.WriteCloser {

	reader, writer := io.Pipe()
	go func(s *Section) {
		defer reader.Close()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			s.Body = append(s.Body, scanner.Text())
			if err := s.display.Refresh(); err != nil {
				return
			}
		}
	}(s)
	return writer

}

// String returns a string representing this section.
func (s *Section) String() string {

	return fmt.Sprintf(
		"─ %v %v\n%v",
		chalk.Bold.TextStyle(s.Header),
		strings.Repeat("─", s.display.Width-len(s.Header)-3),
		strings.Join(s.Body, "\n"))

}

// Display represents a display which consists of several sections.
type Display struct {
	HeaderForeground termbox.Attribute
	HeaderBackground termbox.Attribute
	BodyForeground   termbox.Attribute
	BodyBackground   termbox.Attribute
	Width            int
	Height           int
	mutex            sync.Mutex
	sections         []*Section
	closed           bool
	done             chan error
}

// NewDisplay creates a new display. The new display hooks any key inputs
// including SIGINT, if it receives that signal, the given cancel function
// will be called.
func NewDisplay(cencel context.CancelFunc) (display *Display, err error) {

	err = termbox.Init()
	if err != nil {
		return
	}

	width, height := termbox.Size()
	display = &Display{
		HeaderForeground: termbox.ColorDefault,
		HeaderBackground: termbox.ColorDefault,
		BodyForeground:   termbox.ColorDefault,
		BodyBackground:   termbox.ColorDefault,
		Width:            width,
		Height:           height,
		done:             make(chan error),
	}

	go func() {
		for {
			e := termbox.PollEvent()
			switch e.Type {
			case termbox.EventError:
				display.done <- e.Err
				return
			case termbox.EventKey:
				if e.Key == termbox.KeyCtrlC {
					cencel()
					display.done <- fmt.Errorf("Canceled")
					return
				}
			case termbox.EventResize:
				display.Width = e.Width
				display.Height = e.Height
				display.Refresh()
			}
		}
	}()

	return

}

// Close closes this display.
func (d *Display) Close() {

	if !d.closed {
		termbox.Close()
		d.closed = true
	}

}

// AddSection adds a new section to this display.
func (d *Display) AddSection(header string) *Section {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s := &Section{
		Header:  header,
		display: d,
	}
	d.sections = append(d.sections, s)
	sort.Slice(d.sections, func(i int, j int) bool {
		return d.sections[i].Header < d.sections[j].Header
	})

	return s
}

// DeleteSection deletes the given section from this display.
func (d *Display) DeleteSection(sec *Section) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	old := d.sections
	d.sections = make([]*Section, 0, len(d.sections)-1)
	for _, s := range old {
		if s != sec {
			d.sections = append(d.sections, s)
		}
	}

}

// Refresh refreshes this display.
func (d *Display) Refresh() error {

	if len(d.sections) == 0 {
		return nil
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.closed {
		return fmt.Errorf("Display has been closed already")
	}

	hOffset := d.Height / len(d.sections)
	hExtra := d.Height - hOffset*len(d.sections)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i, s := range d.sections {

		// Print header.
		x := 0
		for _, c := range []rune("─ ") {
			termbox.SetCell(x, i*hOffset, c, d.HeaderForeground, d.HeaderBackground)
			x++
		}
		if s == nil {
			panic(fmt.Sprint(d.sections))
		}
		for _, c := range []rune(s.Header) {
			termbox.SetCell(x, i*hOffset, c, d.HeaderForeground|termbox.AttrBold, d.HeaderBackground)
			x++
		}
		for _, c := range []rune(fmt.Sprintf(" %v", strings.Repeat("─", d.Width-x-1))) {
			termbox.SetCell(x, i*hOffset, c, d.HeaderForeground, d.HeaderBackground)
			x++
		}

		// Print body.
		if i < len(d.sections)-1 && len(s.Body) > hOffset-1 {
			s.Body = s.Body[len(s.Body)-hOffset+1:]
		} else if len(s.Body) > hOffset+hExtra-1 {
			s.Body = s.Body[len(s.Body)-hOffset-hExtra+1:]
		}
		for u, line := range s.Body {
			for v, c := range []rune(line) {
				termbox.SetCell(v, u+i*hOffset+1, c, termbox.ColorDefault, termbox.ColorDefault)
			}
		}

	}
	return termbox.Flush()

}

// Wait blocks until the given context will be canceled.
func (d *Display) Wait() error {

	return <-d.done

}
