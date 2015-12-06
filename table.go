package tablook

import (
	"errors"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ErrTooFewRecords = errors.New("tablook: at least one row must be given in addition to header row")

const (
	FgColor = termbox.ColorWhite
	BgColor = termbox.ColorBlack
)

type Tab struct {
	records [][]string
	maxLen  []int
}

func New(records [][]string) (Tab, error) {
	if len(records) < 2 {
		return Tab{}, ErrTooFewRecords
	}

	var data Tab
	data.records = records
	maxLen := make([]int, data.cols())
	for _, rcrd := range records {
		for i, col := range rcrd {
			l := runewidth.StringWidth(col)
			if maxLen[i] < l {
				maxLen[i] = l
			}
		}
	}
	data.maxLen = maxLen
	return data, nil
}

func (data Tab) Draw() {
	width, height := termbox.Size()
	if height < 2 {
		return
	}

	///termbox.Clear(FgColor, BgColor)
	for i := 0; i < height && i < len(data.records); i++ {
		data.drawRow(width, i, i, i == 0)
	}
	termbox.Flush()
}

func (data Tab) drawRow(width, sourceIdx, viewIdx int, inverse bool) {
	column, x := 0, 0
	fg, bg := FgColor, BgColor
	if inverse {
		fg, bg = bg, fg
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

func (data Tab) lenSum() int {
	s := 0
	for _, x := range data.maxLen {
		s += x
	}
	return s
}

func (data Tab) cols() int { return len(data.records[0]) }
