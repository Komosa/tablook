package tablook

import (
	"errors"

	"github.com/gizak/termui" // for events
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ErrTermboxAlreadyInitialized = errors.New("tablook: termbox (or termui) already initialized")

var (
	FgColor        = termbox.ColorWhite
	BgColor        = termbox.ColorBlack
	SelColor       = termbox.ColorBlue
	SelColumnColor = termbox.ColorGreen
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
			if next >= -1 && next <= h {
				data.selected = next
				if dir == 1 {
					next--
				} else if next < 1 {
					next = 1
				}
				data.redrawTwoRows(next)
			}
		}
	}
	chColSel := func(dir int) func(termui.Event) {
		return func(termui.Event) {
			next := data.selColumn + dir
			for next >= 0 && next < data.cols() && data.isColDel[next] {
				next += dir
			}
			if next >= -1 && next <= data.cols() {
				data.selColumn = next
				data.redraw()
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
	chLeftSel, chRightSel := chColSel(-1), chColSel(+1)
	chLeft := chView(-1, &data.currentX, data.trimmed)
	chRight := chView(+1, &data.currentX, data.trimmed)
	chUp := chView(-1, &data.currentY, data.canGoDown)
	chDown := chView(+1, &data.currentY, data.canGoDown)
	termui.Handle("/sys/kbd/k", chUpSel)
	termui.Handle("/sys/kbd/j", chDownSel)
	termui.Handle("/sys/kbd/h", chLeftSel)
	termui.Handle("/sys/kbd/l", chRightSel)
	termui.Handle("/sys/kbd/<down>", chDown)
	termui.Handle("/sys/kbd/<up>", chUp)
	termui.Handle("/sys/kbd/<left>", chLeft)
	termui.Handle("/sys/kbd/<right>", chRight)
	termui.Handle("/sys/kbd/d", func(termui.Event) {
		if data.selColumn >= 0 && data.selColumn < data.cols() {
			data.isColDel[data.selColumn] = true
			if data.lenSum() == 0 {
				// all columns removed, lets start over again
				for i := 0; i < data.cols(); i++ {
					data.isColDel[i] = false
				}
			}
			data.redraw()
		}
	})
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
		if data.isColDel[data.firstCol] {
			data.firstCol++
			continue
		}
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
	if firstIdx < 2 {
		data.redraw()
		return
	}
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
	}
	if sourceIdx == 0 {
		fg, bg = bg, fg // header line
	}
	columnIsOdd := column&1 == 1
	if columnIsOdd {
		fg, bg = bg, fg // odd number of coulmns skipped
	}

	negShift := -data.toSkip
	for ; x < width && column < data.cols(); column++ {
		if data.isColDel[column] {
			continue
		}
		s := data.records[sourceIdx][column]
		if x+data.maxLen[column] >= width {
			// clap (maybe)
			s = runewidth.Truncate(s, width-x-negShift, "")
			if len(s) == 0 {
				break
			}
		}

		fgcol, bgcol := fg, bg
		if column == data.selColumn {
			if column&1 == 0 {
				fgcol = SelColumnColor
			} else {
				bgcol = SelColumnColor
			}
		}
		drawString(s, x+negShift, viewIdx, fgcol, bgcol)

		fg, bg = bg, fg
		x += data.maxLen[column]
		columnIsOdd = !columnIsOdd
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
