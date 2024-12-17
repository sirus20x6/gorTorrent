// internal/table/table.go

package table

import (
    "fmt"
    "sort"
    "strings"
    "sync"
)

// Table represents a table instance with its configuration and state
type Table struct {
    ID            string
    Columns       []Column
    Rows          map[string]*Row
    RowIDs        []string         // Maintains row order
    SortCol       string          // Primary sort column
    SortReverse   bool            // Primary sort direction
    SecondSortCol string          // Secondary sort column 
    SecondReverse bool            // Secondary sort direction
    PageSize      int             // Number of rows per page
    MaxRows       int             // Maximum rows to display
    Filter        string          // Current filter/search text
    SelectedRows  map[string]bool // Track selected rows
    mu            sync.RWMutex    // Protects concurrent access
}

// NewTable creates a new table instance
func NewTable(id string) *Table {
    return &Table{
        ID:           id,
        Rows:         make(map[string]*Row),
        RowIDs:       make([]string, 0),
        SelectedRows: make(map[string]bool),
        PageSize:     50,    // Default page size
        MaxRows:      1000,  // Default maximum rows
    }
}

// AddColumn adds a column definition to the table
func (t *Table) AddColumn(col Column) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.Columns = append(t.Columns, col)
}

// SetRow adds or updates a row in the table
func (t *Table) SetRow(id string, data map[string]interface{}) {
    t.mu.Lock()
    defer t.mu.Unlock()

    row := &Row{
        ID:      id,
        Data:    data,
        Enabled: true,
    }

    if _, exists := t.Rows[id]; !exists {
        t.RowIDs = append(t.RowIDs, id)
    }
    t.Rows[id] = row

    // Resort if a sort column is set
    if t.SortCol != "" {
        t.sort()
    }
}

// RemoveRow removes a row from the table
func (t *Table) RemoveRow(id string) {
    t.mu.Lock()
    defer t.mu.Unlock()

    delete(t.Rows, id)
    for i, rowID := range t.RowIDs {
        if rowID == id {
            t.RowIDs = append(t.RowIDs[:i], t.RowIDs[i+1:]...)
            break
        }
    }
    delete(t.SelectedRows, id)
}

// Sort sorts the table data
func (t *Table) Sort(col string, secondary bool) {
    t.mu.Lock()
    defer t.mu.Unlock()

    if secondary {
        if t.SecondSortCol == col {
            t.SecondReverse = !t.SecondReverse
        } else {
            t.SecondSortCol = col
            t.SecondReverse = false
        }
    } else {
        if t.SortCol == col {
            t.SortReverse = !t.SortReverse
        } else {
            t.SortCol = col
            t.SortReverse = false
        }
    }

    t.sort()
}

// sort performs the actual sorting (internal method)
func (t *Table) sort() {
    sort.SliceStable(t.RowIDs, func(i, j int) bool {
        row1 := t.Rows[t.RowIDs[i]]
        row2 := t.Rows[t.RowIDs[j]]

        // Primary sort
        cmp := t.compareValues(row1.Data[t.SortCol], row2.Data[t.SortCol])
        if t.SortReverse {
            cmp = -cmp
        }
        if cmp != 0 {
            return cmp < 0
        }

        // Secondary sort
        if t.SecondSortCol != "" {
            cmp = t.compareValues(row1.Data[t.SecondSortCol], row2.Data[t.SecondSortCol])
            if t.SecondReverse {
                cmp = -cmp
            }
            if cmp != 0 {
                return cmp < 0
            }
        }

        // Stable sort - use row ID as final tiebreaker
        return row1.ID < row2.ID
    })
}

// compareValues compares two values based on their types
func (t *Table) compareValues(v1, v2 interface{}) int {
    // Handle nil values
    if v1 == nil && v2 == nil {
        return 0
    }
    if v1 == nil {
        return -1
    }
    if v2 == nil {
        return 1
    }

    switch v1.(type) {
    case float64:
        f1 := v1.(float64)
        f2, ok := v2.(float64)
        if !ok {
            return -1
        }
        if f1 < f2 {
            return -1
        }
        if f1 > f2 {
            return 1
        }
        return 0

    case string:
        s1 := v1.(string)
        s2, ok := v2.(string)
        if !ok {
            return -1
        }
        return strings.Compare(s1, s2)
    }

    // Default string comparison
    return strings.Compare(fmt.Sprint(v1), fmt.Sprint(v2))
}

// GetVisibleRows returns the visible rows for the current page
func (t *Table) GetVisibleRows(start, limit int) []*Row {
    t.mu.RLock()
    defer t.mu.RUnlock()

    if start >= len(t.RowIDs) {
        return []*Row{}
    }

    end := start + limit
    if end > len(t.RowIDs) {
        end = len(t.RowIDs)
    }

    visible := make([]*Row, 0, end-start)
    for _, id := range t.RowIDs[start:end] {
        if row := t.Rows[id]; row.Enabled {
            visible = append(visible, row)
        }
    }

    return visible
}

// ApplyFilter filters the table rows based on search text
func (t *Table) ApplyFilter(filter string) {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.Filter = strings.ToLower(filter)
    for _, row := range t.Rows {
        row.Enabled = t.rowMatchesFilter(row)
    }
}

// rowMatchesFilter checks if a row matches the current filter
func (t *Table) rowMatchesFilter(row *Row) bool {
    if t.Filter == "" {
        return true
    }

    for _, col := range t.Columns {
        if val, ok := row.Data[col.ID]; ok {
            if strings.Contains(
                strings.ToLower(fmt.Sprint(val)),
                t.Filter,
            ) {
                return true
            }
        }
    }
    return false
}

// GetColumnByID returns a column by its ID
func (t *Table) GetColumnByID(id string) *Column {
    t.mu.RLock()
    defer t.mu.RUnlock()

    for i, col := range t.Columns {
        if col.ID == id {
            return &t.Columns[i]
        }
    }
    return nil
}

// SetSelection sets the selection state for a row
func (t *Table) SetSelection(id string, selected bool) {
    t.mu.Lock()
    defer t.mu.Unlock()

    if selected {
        t.SelectedRows[id] = true
    } else {
        delete(t.SelectedRows, id)
    }
}

// GetSelectedRows returns the currently selected row IDs
func (t *Table) GetSelectedRows() []string {
    t.mu.RLock()
    defer t.mu.RUnlock()

    selected := make([]string, 0, len(t.SelectedRows))
    for id := range t.SelectedRows {
        selected = append(selected, id)
    }
    return selected
}