# tablook
less for tabular data, with csvless command as example

# library
you can look on your `[][]string` data by:

    tbl, err := tablook.New(some_data)
    // remember to handle errors
    err = tablook.Show()

full example in `example/simple.go`

# standalone command for csv files
just `go get github.com/komosa/tablook/cmd/csvless`
and then you can:
* pipe csv data to `csvless` (eg. from John Kerl's  [mlr tool](https://github.com/johnkerl/miller))
* open your csv files with `csvless filename.csv`

# keybindings
key | behaviour
--- | ---
**q** | quit app / returns from `Tab.Show()`
**arrows** | scroll view
**j/k** | move highlighted line down/up
**h/l** | change selected column to one on left/right side of current one
**d** | deletes selected column (or start over with all visible instead of removing the last one)

# bonus
* works with characters wider than one cell
* handles terminal resizes smoothly
