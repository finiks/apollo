package main

import (
    "github.com/nsf/termbox-go"
    "strconv"
    "sort"
)

type EntriesTab struct {
    a *Apollo
    entries *[]Entry
    slice []*Entry
    name string
    status string
    sortField string
    view string
    offset int
    cursor int
    ratings bool
    search []Entry
    additionalField string
}

func (t *EntriesTab) Name() string {
    return t.name
}

func (t *EntriesTab) Status() string {
    return t.status
}

func (t *EntriesTab) changeView(view string) {
    t.view = view
    t.cursor = 0
    t.offset = 0
    t.refreshSlice()
}

func (t *EntriesTab) toggleSort() {
    if t.sortField == "title" {
        t.sortField = "year"
    } else if t.sortField == "year" {
        t.sortField = "rating"
    } else if t.sortField == "rating" && t.additionalField != "" {
        t.sortField = t.additionalField
    } else {
        t.sortField = "title"
    }

    t.cursor = 0
    t.offset = 0
    t.sort()
}

func (t *EntriesTab) handleKeyEvent(ev *termbox.Event) bool {
    if t.view == "edit" {
        switch ev.Ch {
        case '0':
            t.a.inputActive = true
            t.a.input = []rune(":t " + t.slice[t.cursor].Title)
            t.a.inputCursor = len(t.a.input)
            return true
        case '1':
            t.a.inputActive = true
            t.a.input = []rune(":y " + t.slice[t.cursor].Year)
            t.a.inputCursor = len(t.a.input)
            return true
        }
    }

    if t.view == "tag" {
        indexes := map[rune]int{'0': 0, '1': 1, '2': 2, '3': 3,
                                '4': 4, '5': 5, '6': 6, '7': 7,
                                '8': 8, '9': 9,}
        if i, exist := indexes[ev.Ch]; exist {
            if i < len(t.search) {
                t.slice[t.cursor].Title = t.search[i].Title
                t.slice[t.cursor].Year = t.search[i].Year
                t.slice[t.cursor].TagID = t.search[i].TagID

                t.view = "passive"
                t.refreshSlice()

                for j := range(t.slice) {
                    if t.slice[j].TagID == t.search[i].TagID {
                        t.cursor = j
                        t.offset = 0
                        if t.cursor > t.a.height - 4 {
                            t.offset = t.cursor - (t.a.height - 4)
                        }
                    }
                }

                t.search = t.search[:0]
                t.a.d.save()
            }
        }

        return true
    }

    switch ev.Ch {
    case '1':
        t.changeView("passive")
    case '2':
        t.changeView("active")
    case '3':
        t.changeView("inactive")
    case '4':
        t.changeView("all")
    case 's':
        t.toggleSort()
    case 'D':
        if len(t.slice) > 0 {
            for i := range *t.entries {
                if (*t.entries)[i] == *t.slice[t.cursor] {
                    t.a.logDebug("deleted id." + strconv.Itoa(i))
                    *t.entries = append((*t.entries)[:i], (*t.entries)[i+1:]...)
                    break
                }
            }
            t.a.d.save()
            t.refreshSlice()
        }
    case 'e':
        if len(t.slice) > 0 {
            if t.view != "edit" {
                t.view = "edit"
            } else {
                t.view = "passive"
                t.refreshSlice()
            }
        }
    case 'r':
        t.ratings = !t.ratings
    case 'a':
        if len(t.slice) > 0 {
            if t.slice[t.cursor].State == "passive" {
                t.slice[t.cursor].State = "inactive"
            } else if t.slice[t.cursor].State == "inactive" {
                t.slice[t.cursor].State = "active"
            } else {
                t.slice[t.cursor].State = "passive"
            }
            t.slice[t.cursor].Rating = 0
            t.a.d.save()
        }
    case 'z':
        if len(t.slice) > 0 {
            if t.slice[t.cursor].Rating > 0 {
                t.slice[t.cursor].Rating--
                t.a.d.save()
            }
        }
    case 'x':
        if len(t.slice) > 0 {
            if t.slice[t.cursor].Rating < 6 {
                t.slice[t.cursor].Rating++
                t.a.d.save()
            }
        }
    }

    switch ev.Key {
    case termbox.KeyArrowUp:
        t.cursor--
        if t.cursor < 0 {
            t.cursor = 0
        } else if t.cursor - t.offset < 0 {
            t.offset--
        }
    case termbox.KeyPgup:
        t.cursor -= 5
        if t.cursor < 0 {
            t.cursor = 0
            t.offset = 0
        } else if t.cursor - t.offset < 0 {
            t.offset -= 5
            if t.offset < 0 {
                t.offset = 0
            }
        }
    case termbox.KeyArrowDown:
        t.cursor++
        if t.cursor > len(t.slice) - 1 {
            t.cursor--
        } else if t.cursor - t.offset > t.a.height - 4 {
            t.offset++
        }
    case termbox.KeyPgdn:
        t.cursor += 5
        if t.cursor > len(t.slice) - 1 {
            t.cursor = len(t.slice) - 1
            if len(t.slice) > t.a.height - 3 {
                t.offset = len(t.slice) - (t.a.height - 3)
            }
        } else if t.cursor - t.offset > t.a.height - 4 {
            t.offset += 5
            if t.cursor - t.offset < t.a.height - 3 {
                t.offset = t.cursor - (t.a.height - 4)
            }
        }
    default:
        return false
    }

    t.a.logDebug("cursor: " + strconv.Itoa(t.cursor) + " offset: " + strconv.Itoa(t.offset) + 
    " len:" + strconv.Itoa(len(t.slice)) + " height:" + strconv.Itoa(t.a.height))

    return true
}

func (t *EntriesTab) drawEditView() {
    runes := []rune("0. " + t.slice[t.cursor].Title)
    for i := 0; i < len(runes); i++ {
        termbox.SetCell(i, 1, runes[i], colors['d'], colors['d'])
    }

    runes = []rune("1. " + t.slice[t.cursor].Year)
    for i := 0; i < len(runes); i++ {
        termbox.SetCell(i, 2, runes[i], colors['d'], colors['d'])
    }
}

func (t *EntriesTab) drawTagView() {
    for j := 0; j < len(t.search); j++ {
        runes := []rune(strconv.Itoa(j) + ". [" + t.search[j].Year + "] " + t.search[j].Title)
        for i := 0; i < len(runes); i++ {
            termbox.SetCell(i, j + 1, runes[i], colors['d'], colors['d'])
        }
    }
}

func (t *EntriesTab) drawEntries() {
    for j := 0; j < t.a.height - 3; j++ {
        if j < len(t.slice) {
            if t.ratings {
                for i := 0; i < t.slice[j + t.offset].Rating; i++ {
                    if t.slice[j + t.offset].State == "passive" {
                        termbox.SetCell(i + 3, j + 1, '*', colors['y'], colors['d'])
                    } else if t.slice[j + t.offset].State == "inactive" {
                        termbox.SetCell(i + 3, j + 1, '*', colors['B'], colors['d'])
                    }
                }
            }

            year := t.slice[j + t.offset].Year
            if year == "" {
                year = "    "
            }
            title := t.slice[j + t.offset].Title
            runes := []rune(year + " " + title)
            for i := 0; i < len(runes); i++ {
                fg := colors['d']
                if i < 4 {
                    if t.slice[j + t.offset].State == "passive" {
                        fg = colors['g']
                    } else if t.slice[j + t.offset].State == "inactive" {
                        fg = colors['b']
                    } else {
                        fg = colors['Y']
                    }
                }

                if t.ratings {
                    termbox.SetCell(i + 10, j + 1, runes[i], fg, colors['d'])
                } else {
                    termbox.SetCell(i + 3, j + 1, runes[i], fg, colors['d'])
                }
            }
        }
    }

    termbox.SetCell(1, t.cursor - t.offset + 1, '*', colors['d'], colors['d'])
}

func (t *EntriesTab) Draw() {
    if t.view == "edit" {
        t.drawEditView()
    } else if t.view == "tag" {
        t.drawTagView()
    } else {
        t.drawEntries()
    }
}

func (t *EntriesTab) appendEntry(e Entry) {
    *t.entries = append(*t.entries, e)
    t.a.d.save()
    t.refreshSlice()

    t.cursor = len(t.slice) - 1
    if t.cursor > t.a.height - 4 {
        t.offset = t.cursor - (t.a.height - 4)
    }
}

func (t *EntriesTab) editCurrentEntry(field rune, value string) {
    switch field {
    case 't':
        t.slice[t.cursor].Title = value
    case 'y':
        t.slice[t.cursor].Year = value
    }

    t.a.inputActive = false
}

func (t *EntriesTab) query(query string) {
    if query[0] != ':' {
        t.appendEntry(Entry{Title: query, State: "passive"})
    } else {
        t.editCurrentEntry(rune(query[1]), query[3:])
    }
}

type By func(e1, e2 *Entry) bool

func (by By) Sort(entries []*Entry) {
    es := &entrySorter{
        entries: entries,
        by: by,
    }
    sort.Sort(es)
}

type entrySorter struct {
    entries []*Entry
    by func(e1, e2 *Entry) bool
}

func (s *entrySorter) Len() int {
    return len(s.entries)
}

func (s *entrySorter) Swap(i, j int) {
    s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func (s *entrySorter) Less(i, j int) bool {
    return s.by(s.entries[i], s.entries[j])
}

func (t *EntriesTab) sort() {

    title := func(e1, e2 *Entry) bool {
        return e1.Title < e2.Title
    }

    year := func(e1, e2 *Entry) bool {
        return e1.Year < e2.Year
    }

    rating := func(e1, e2 *Entry) bool {
        return e1.Rating < e2.Rating
    }

    info1 := func(e1, e2 *Entry) bool {
        return e1.Info1 < e2.Info1
    }

    By(title).Sort(t.slice)

    if t.sortField == "year" {
        By(year).Sort(t.slice)
    } else if t.sortField == "rating" {
        By(rating).Sort(t.slice)
    } else if t.sortField != "title" {
        By(info1).Sort(t.slice)
    }

    t.status = t.name + " - " + t.view + " (" + strconv.Itoa(len(t.slice)) +
        " entries) sorted by " + t.sortField
}

func (t *EntriesTab) refreshSlice() {
    t.slice = t.slice[:0]
    for i := range *t.entries {
        if (*t.entries)[i].State == "passive" {
            if t.view == "passive" || t.view == "all" {
                t.slice = append(t.slice, &(*t.entries)[i])
            }
        } else if (*t.entries)[i].State == "inactive" {
            if t.view == "inactive" || t.view == "all" {
                t.slice = append(t.slice, &(*t.entries)[i])
            }
        } else {
            if t.view == "active" || t.view == "all" {
                t.slice = append(t.slice, &(*t.entries)[i])
            }
        }
    }
    t.sort()

    if t.cursor > len(t.slice) - 1 {
        t.cursor--
    }

    if t.offset > 0 {
        if len(t.slice) - t.offset < t.a.height - 3 {
            t.offset--
        }
    }
}
