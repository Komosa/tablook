package tablook

import (
	"errors"

	"github.com/gizak/termui" // for events
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ErrTermboxAlreadyInitialized = errors.New("tablook: termbox (or termui) already initialized")

const (
	FgColor  = termbox.ColorWhite
	BgColor  = termbox.ColorBlack
	SelColor = termbox.ColorBlue
)

func initTermui() error {
	if termbox.IsInit {
		return ErrTermboxAlreadyInitialized
	}
	return termui.Init()
}

func closeTermui() {
	termui.Close()
}

func (data Tab) Show() {
	if err := initTermui(); err != nil {
		panic(err)
	}
	defer closeTermui()

	width, height := termbox.Size()
	if height < 2 {
		return
	}

	for i := 0; i < height && i < data.rows(); i++ {
		data.drawRow(width, i, i)
	}
	termbox.Flush()

	data.loop()
}

func (data Tab) loop() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/k", func(termui.Event) {
		if data.selected > 1 {
			data.selected--
			data.redrawTwoRows(data.selected)
		}
	})
	termui.Handle("/sys/kbd/j", func(termui.Event) {
		if data.selected+1 < data.rows() {
			data.selected++
			data.redrawTwoRows(data.selected - 1)
		}
	})
	termui.Loop()
}

func (data Tab) redrawTwoRows(firstIdx int) {
	width, _ := termbox.Size()
	data.drawRow(width, firstIdx, firstIdx)
	data.drawRow(width, firstIdx+1, firstIdx+1)
	termbox.Flush()
}

func (data Tab) drawRow(width, sourceIdx, viewIdx int) {
	column, x := 0, 0
	fg, bg := FgColor, BgColor
	if data.selected == sourceIdx {
		fg = SelColor
	} else if sourceIdx == 0 {
		fg, bg = bg, fg // header
	}

	for x < width && column < data.cols() {
		s := data.records[sourceIdx][column]
		if x+data.maxLen[column] >= width {
			// clap
			s = runewidth.Truncate(s, width-x, "")
			if len(s) == 0 {
				break
			}
		}

		drawString(s, x, viewIdx, fg, bg)
		fg, bg = bg, fg
		x += data.maxLen[column]
		column++
	}
}

func drawString(s string, x, y int, fg, bg termbox.Attribute) {
	for _, ch := range s {
		termbox.SetCell(x, y, ch, fg, bg)
		x += runewidth.RuneWidth(ch)
	}
}
