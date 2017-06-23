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

	"github.com/jroimartin/gocui"
	"github.com/ttacon/chalk"
)

// Section represents a section in a display. Each section has a header text and
// several strings as the body.
type Section struct {
	Header  string
	Body    []string
	display *Display
}

func newSection(display *Display, header string) *Section {

	return &Section{
		Header:  header,
		display: display,
	}

}

// Writer returns io.WriteCloser to write messages into the section. Users have
// to close the returned writer.
func (s *Section) Writer() io.WriteCloser {

	reader, writer := io.Pipe()
	go func() {
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			s.Body = append(s.Body, scanner.Text())
			s.display.gui.Execute(func(g *gocui.Gui) (err error) {
				v, err := g.View(s.Header)
				if err != nil {
					return
				}
				v.Clear()
				for _, line := range s.Body {
					fmt.Fprintln(v, line)
				}
				return
			})
		}

	}()

	return writer

}

// String returns a string representing this section.
func (s *Section) String() string {

	return fmt.Sprintf(
		"%v\n%v",
		chalk.Bold.TextStyle(s.Header),
		strings.Join(s.Body, "\n"))

}

// Display represents a display which consists of several sections.
type Display struct {
	mutex    sync.Mutex
	sections []*Section
	closed   bool
	done     chan error
	gui      *gocui.Gui
}

// NewDisplay creates a new display. The new display hooks any key inputs
// including SIGINT, if it receives that signal, the given cancel function
// will be called.
func NewDisplay(cancel context.CancelFunc) (display *Display, err error) {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	display = &Display{
		gui:  g,
		done: make(chan error),
	}
	g.SetManager(display)

	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		cancel()
		return gocui.ErrQuit
	})
	if err != nil {
		g.Close()
		return
	}

	go func() {
		err := g.MainLoop()
		if err == gocui.ErrQuit {
			err = nil
		}
		display.done <- err
	}()

	return

}

// Layout is called every time the GUI is redrawn, it must contain the
// base views and its initializations.
func (d *Display) Layout(g *gocui.Gui) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if len(d.sections) == 0 {
		return nil
	}
	if d.closed {
		return fmt.Errorf("Display has been closed already")
	}

	width, height := g.Size()
	hOffset := height / len(d.sections)
	// hExtra := height - hOffset*len(d.sections)

	for i, s := range d.sections {

		v, err := g.SetView(s.Header, 0, i*hOffset, width-1, (i+1)*hOffset-1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = s.Header
			v.Autoscroll = true
		}

	}

	return nil

}

// Close closes this display.
func (d *Display) Close() {

	if !d.closed {
		d.gui.Close()
		d.closed = true
	}

}

// AddSection adds a new section to this display.
func (d *Display) AddSection(header string) *Section {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s := newSection(d, header)
	d.sections = append(d.sections, s)
	sort.Slice(d.sections, func(i int, j int) bool {
		return d.sections[i].Header < d.sections[j].Header
	})

	d.gui.Execute(d.Layout)
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

	d.gui.Execute(func(g *gocui.Gui) error {
		return g.DeleteView(sec.Header)
	})

}

// Wait blocks until the given context will be canceled.
func (d *Display) Wait() error {
	return <-d.done
}
