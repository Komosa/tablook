package tablook

import (
	"errors"

	"github.com/mattn/go-runewidth"
)

var ErrTooFewRecords = errors.New("tablook: at least one row must be given in addition to header row")

type Tab struct {
	records   [][]string
	maxLen    []int
	selected  int // in view
	selColumn int
	currentY  int // current y-axis shift
	currentX  int // current x-axis shift
	toSkip    int // how many cells we should skip when printing firstCol
	firstCol  int // from which column we should start printing
	isColDel  []bool
}

func New(records [][]string) (Tab, error) {
	if len(records) < 2 {
		return Tab{}, ErrTooFewRecords
	}

	data := Tab{selected: 1, selColumn: -1}
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
	data.isColDel = make([]bool, data.cols())
	return data, nil
}

func (data Tab) lenSum() int {
	s := 0
	for i, x := range data.maxLen {
		if !data.isColDel[i] {
			s += x
		}
	}
	return s
}

func (data Tab) cols() int { return len(data.records[0]) }
func (data Tab) rows() int { return len(data.records) }
