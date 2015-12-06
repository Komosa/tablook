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

	data.redraw()
	data.loop()
}

func (data Tab) loop() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	chRow := func(dir int) func(termui.Event) {
		return func(termui.Event) {
			next := data.selected + dir
			if next != 0 && next < data.rows() {
				data.selected = next
				if dir == 1 {
					next--
				}
				data.redrawTwoRows(next)
			}
		}
	}
	chUp, chDown := chRow(-1), chRow(+1)
	termui.Handle("/sys/kbd/k", chUp)
	termui.Handle("/sys/kbd/j", chDown)
	termui.Handle("/sys/kbd/<down>", chDown)
	termui.Handle("/sys/kbd/<up>", chUp)
	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		data.redraw()
	})
	termui.Loop()
}

func (data Tab) redraw() {
	termbox.Sync()
	termbox.Clear(FgColor, BgColor)
	width, height := termbox.Size()
	if height < 2 {
		drawString("wnd size too small", 0, 0, FgColor, BgColor)
		return
	}

	for i := 0; i < height && i < data.rows(); i++ {
		data.drawRow(width, i, i)
	}
	termbox.Flush()
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
