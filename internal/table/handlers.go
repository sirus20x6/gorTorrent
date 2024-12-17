// internal/table/row.go

package table

import (
    "fmt"
    "strings"
    "sync"
)

// Row represents a single table row
type Row struct {
    ID       string                  // Unique identifier for the row
    Data     map[string]interface{} // Column values indexed by column ID
    Icon     string                 // Icon class or path
    Selected bool                   // Selection state
    Enabled  bool                   // Whether row is visible/enabled
    Attr     map[string]string      // Additional HTML attributes
    mu       sync.RWMutex          // Protects concurrent access to row data
}

// NewRow creates a new row instance
func NewRow(id string) *Row {
    return &Row{
        ID:      id,
        Data:    make(map[string]interface{}),
        Attr:    make(map[string]string),
        Enabled: true,
    }
}

// SetValue sets a value for a specific column
func (r *Row) SetValue(columnID string, value interface{}) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.Data[columnID] = value
}

// GetValue retrieves a value for a specific column
func (r *Row) GetValue(columnID string) interface{} {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.Data[columnID]
}

// SetIcon sets the row's icon
func (r *Row) SetIcon(icon string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.Icon = icon
}

// SetAttribute sets an HTML attribute for the row
func (r *Row) SetAttribute(key, value string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.Attr[key] = value
}

// GetAttribute gets an HTML attribute value
func (r *Row) GetAttribute(key string) string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.Attr[key]
}

// Enable enables/disables the row
func (r *Row) Enable(enabled bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.Enabled = enabled
}

// Select sets the row's selection state
func (r *Row) Select(selected bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.Selected = selected
}

// Clone creates a deep copy of the row
func (r *Row) Clone() *Row {
    r.mu.RLock()
    defer r.mu.RUnlock()

    newRow := NewRow(r.ID)
    newRow.Icon = r.Icon
    newRow.Selected = r.Selected
    newRow.Enabled = r.Enabled

    // Deep copy data
    for k, v := range r.Data {
        newRow.Data[k] = v
    }

    // Deep copy attributes
    for k, v := range r.Attr {
        newRow.Attr[k] = v
    }

    return newRow
}

// MatchesFilter checks if the row matches a filter string
func (r *Row) MatchesFilter(filter string, columns []Column) bool {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if filter == "" {
        return true
    }

    filter = strings.ToLower(filter)
    for _, col := range columns {
        if val, ok := r.Data[col.ID]; ok {
            strVal := strings.ToLower(fmt.Sprint(val))
            if strings.Contains(strVal, filter) {
                return true
            }
        }
    }
    return false
}

// GetCellClass returns the CSS class for a cell based on column type
func (r *Row) GetCellClass(col Column) string {
    classes := []string{col.GetAlignClass()}
    
    if r.Selected {
        classes = append(classes, "selected")
    }

    switch col.Type {
    case TypeProgress:
        classes = append(classes, "progress-cell")
    case TypeNumber:
        classes = append(classes, "numeric-cell")
    case TypeDate:
        classes = append(classes, "date-cell")
    }

    return strings.Join(classes, " ")
}

// GetRowClass returns the CSS class for the entire row
func (r *Row) GetRowClass() string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    classes := []string{"table-row"}
    
    if r.Selected {
        classes = append(classes, "selected")
    }
    
    if !r.Enabled {
        classes = append(classes, "disabled")
    }

    return strings.Join(classes, " ")
}

// FormatCell formats a cell value based on column type
func (r *Row) FormatCell(col Column) string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    val, ok := r.Data[col.ID]
    if !ok || val == nil {
        return ""
    }

    return col.FormatValue(val)
}

// GetAttributes returns HTML attributes as a string
func (r *Row) GetAttributes() string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if len(r.Attr) == 0 {
        return ""
    }

    var attrs []string
    for k, v := range r.Attr {
        attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
    }
    return strings.Join(attrs, " ")
}

// IsEmpty checks if the row has any data
func (r *Row) IsEmpty() bool {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return len(r.Data) == 0
}

// Clear removes all data from the row
func (r *Row) Clear() {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.Data = make(map[string]interface{})
    r.Attr = make(map[string]string)
    r.Icon = ""
    r.Selected = false
}

// Update updates multiple values at once
func (r *Row) Update(values map[string]interface{}) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for k, v := range values {
        r.Data[k] = v
    }
}