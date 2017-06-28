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
	body   []string
	Logger io.Writer
}

func newHeader(update DisplayUpdateFunc) (header *Header) {

	header = new(Header)
	reader, writer := io.Pipe()
	go func() {
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {

			header.body = append(header.body, scanner.Text())
			update(func(writer io.Writer) {
				for _, line := range header.body {
					fmt.Fprintln(writer, line)
				}
			})

		}
	}()

	header.Logger = writer
	return

}

// Section represents a section in a display. Each section has a header text and
// several strings as the body.
type Section struct {
	Header string
	Body   []string
	update DisplayUpdateFunc
}

func newSection(header string, update DisplayUpdateFunc) *Section {

	return &Section{
		Header: header,
		update: update,
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
			s.update(func(writer io.Writer) {
				for _, line := range s.Body {
					fmt.Fprintln(writer, line)
				}
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
	Title      string
	Header     *Header
	mutex      sync.Mutex
	sections   []*Section
	closed     bool
	done       chan error
	gui        *gocui.Gui
}

// DisplayUpdateHandler defines a handler function to update section body.
type DisplayUpdateHandler func(writer io.Writer)

// DisplayUpdateFunc is a function which a section calls to update the section
// body.
type DisplayUpdateFunc func(handler DisplayUpdateHandler)

// NewDisplay creates a new display.
func NewDisplay(ctx context.Context, title string, maxSection int) (display *Display, nctx context.Context, err error) {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return
	}

	display = &Display{
		MaxSection: maxSection,
		Title:      title,
		gui:        g,
		done:       make(chan error),
		Header: newHeader(func(handler DisplayUpdateHandler) {
			g.Execute(func(g *gocui.Gui) (err error) {
				v, err := g.View("header")
				if err != nil {
					return
				}
				v.Clear()
				handler(v)
				return
			})
		}),
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
		v.Title = d.Title
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
func (d *Display) Close() (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.closed {
		d.gui.Execute(func(g *gocui.Gui) error {
			return gocui.ErrQuit
		})
		err = <-d.done
		d.gui.Close()
		d.closed = true
	}
	return

}

// AddSection adds a new section to this display.
func (d *Display) AddSection(header string) *Section {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	s := newSection(header, func(handler DisplayUpdateHandler) {
		d.gui.Execute(func(g *gocui.Gui) (err error) {
			v, err := g.View(header)
			if err != nil {
				return
			}
			v.Clear()
			handler(v)
			return
		})
	})

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
