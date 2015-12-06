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

func (data *Tab) loop() {
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	chRowSel := func(dir int) func(termui.Event) {
		return func(termui.Event) {
			next := data.selected + dir
			_, h := termbox.Size()
			if next >= 0 && next <= h {
				data.selected = next
				if dir == 1 {
					next--
				}
				data.redrawTwoRows(next)
			}
		}
	}
	chView := func(dir int, field *int, canGoFwd func() bool) func(termui.Event) {
		return func(termui.Event) {
			next := *field + dir
			if (dir < 0 && next >= 0) || (dir > 0 && canGoFwd()) {
				*field = next
				data.redraw()
			}
		}
	}
	chUpSel, chDownSel := chRowSel(-1), chRowSel(+1)
	chLeft := chView(-1, &data.currentX, data.trimmed)
	chRight := chView(+1, &data.currentX, data.trimmed)
	chUp := chView(-1, &data.currentY, data.canGoDown)
	chDown := chView(+1, &data.currentY, data.canGoDown)
	termui.Handle("/sys/kbd/k", chUpSel)
	termui.Handle("/sys/kbd/j", chDownSel)
	termui.Handle("/sys/kbd/<down>", chDown)
	termui.Handle("/sys/kbd/<up>", chUp)
	termui.Handle("/sys/kbd/h", chLeft)
	termui.Handle("/sys/kbd/l", chRight)
	termui.Handle("/sys/kbd/<left>", chLeft)
	termui.Handle("/sys/kbd/<right>", chRight)
	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		data.redraw()
	})
	termui.Loop()
}

func (data *Tab) redraw() {
	termbox.Clear(FgColor, BgColor)
	defer termbox.Flush()
	width, height := termbox.Size()
	if height < 2 {
		drawString("wnd size too small", 0, 0, FgColor, BgColor)
		return
	}

	data.toSkip = data.currentX
	data.firstCol = 0
	for data.toSkip > 0 && data.firstCol+1 != data.cols() {
		if data.toSkip-data.maxLen[data.firstCol] >= 0 {
			data.toSkip -= data.maxLen[data.firstCol]
			data.firstCol++
		} else {
			break
		}
	}

	data.drawRow(width, 0, 0)
	for i := 1; i < height && i+data.currentY < data.rows(); i++ {
		data.drawRow(width, i+data.currentY, i)
	}
}

func (data *Tab) redrawTwoRows(firstIdx int) {
	width, _ := termbox.Size()
	data.drawRow(width, data.currentY+firstIdx, firstIdx)
	data.drawRow(width, data.currentY+firstIdx+1, firstIdx+1)
	termbox.Flush()
}

func (data *Tab) drawRow(width, sourceIdx, viewIdx int) {
	column, x := data.firstCol, 0
	fg, bg := FgColor, BgColor
	if data.selected == viewIdx {
		fg = SelColor
	} else if sourceIdx == 0 {
		fg, bg = bg, fg // header line
	}
	if column&1 == 1 {
		fg, bg = bg, fg // odd number of coulmns skipped
	}

	negShift := -data.toSkip
	for x < width && column < data.cols() {
		s := data.records[sourceIdx][column]
		if x+data.maxLen[column] >= width {
			// clap (maybe)
			s = runewidth.Truncate(s, width-x-negShift, "")
			if len(s) == 0 {
				break
			}
		}

		drawString(s, x+negShift, viewIdx, fg, bg)
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

func (data *Tab) trimmed() bool {
	w, _ := termbox.Size()
	return data.currentX+w < data.lenSum()
}

func (data *Tab) canGoDown() bool {
	_, h := termbox.Size()
	return data.currentY+h < data.rows()
}
