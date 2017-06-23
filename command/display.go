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

const (
	headerHeight = 6
)

// Header represents a header space in a display.
type Header struct {
	body    []string
	mutex   sync.Mutex
	display *Display
}

// Println prints a new line to this header space.
func (h *Header) Println(msg string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.body = append(h.body, msg)
	h.display.gui.Execute(func(g *gocui.Gui) (err error) {
		v, err := g.View("header")
		if err != nil {
			return
		}
		v.Clear()
		for _, line := range h.body {
			fmt.Fprintln(v, line)
		}
		return
	})

}

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
	MaxSection int

	mutex    sync.Mutex
	header   *Header
	sections []*Section
	closed   bool
	done     chan error
	gui      *gocui.Gui
}

// NewDisplay creates a new display.
func NewDisplay(ctx context.Context, maxSection int) (display *Display, nctx context.Context, err error) {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	display = &Display{
		MaxSection: maxSection,
		gui:        g,
		done:       make(chan error),
	}
	display.header = &Header{
		display: display,
	}
	g.SetManager(display)

	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		g.Close()
		return
	}

	nctx, cancel := context.WithCancel(ctx)
	go func() {
		err := g.MainLoop()
		if err == gocui.ErrQuit {
			err = nil
		}
		cancel()
		display.done <- err
	}()

	return

}

// Layout is called every time the GUI is redrawn, it must contain the
// base views and its initializations.
func (d *Display) Layout(g *gocui.Gui) error {
	if d.closed {
		return fmt.Errorf("Display has been closed already")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	width, height := g.Size()
	if v, err := g.SetView("root", 0, 0, width-1, height-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Loci v0.x.x"
		v.Frame = true
	}

	if v, err := g.SetView("header", 0, 0, width-2, headerHeight+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Autoscroll = true
	}

	sectionHeight := (height - headerHeight - 2) / d.MaxSection
	for i, s := range d.sections {

		v, err := g.SetView(s.Header, 1, i*sectionHeight+headerHeight+1, width-2, (i+1)*sectionHeight+headerHeight)
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
