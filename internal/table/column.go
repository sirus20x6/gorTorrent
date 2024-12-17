// internal/table/column.go

package table

import (
    "encoding/json"
)

// ColumnType represents the data type of a column
type ColumnType int

const (
    TypeString ColumnType = iota
    TypeNumber
    TypeProgress
    TypeDate
    TypePeers
    TypeSeeds
    TypeStringNoCase
)

// ColumnAlign represents column text alignment
type ColumnAlign int

const (
    AlignLeft ColumnAlign = iota
    AlignCenter
    AlignRight
    AlignAuto
)

// Column represents a table column configuration
type Column struct {
    ID       string      // Unique identifier for the column
    Text     string      // Display text for column header
    Width    int         // Column width in pixels
    Type     ColumnType  // Data type of the column
    Align    ColumnAlign // Text alignment
    Enabled  bool        // Whether column is visible
    Sortable bool        // Whether column can be sorted
}

// ColumnState represents the persistent state of a column
type ColumnState struct {
    Width   int  `json:"width"`
    Enabled bool `json:"enabled"`
}

// ColumnStorage handles persistence of column states
type ColumnStorage interface {
    SaveColumnStates(tableID string, states map[string]ColumnState) error
    LoadColumnStates(tableID string) (map[string]ColumnState, error)
}

// NewColumn creates a new column with default settings
func NewColumn(id string, text string) Column {
    return Column{
        ID:       id,
        Text:     text,
        Width:    100, // Default width
        Type:     TypeString,
        Align:    AlignLeft,
        Enabled:  true,
        Sortable: true,
    }
}

// WithType sets the column type
func (c Column) WithType(t ColumnType) Column {
    c.Type = t
    return c
}

// WithWidth sets the column width
func (c Column) WithWidth(width int) Column {
    c.Width = width
    return c
}

// WithAlign sets the column alignment
func (c Column) WithAlign(align ColumnAlign) Column {
    c.Align = align
    return c
}

// WithSortable sets whether the column is sortable
func (c Column) WithSortable(sortable bool) Column {
    c.Sortable = sortable
    return c
}

// WithEnabled sets whether the column is enabled
func (c Column) WithEnabled(enabled bool) Column {
    c.Enabled = enabled
    return c
}

// SaveState saves the column's current state
func (c *Column) SaveState() ColumnState {
    return ColumnState{
        Width:   c.Width,
        Enabled: c.Enabled,
    }
}

// LoadState loads the column's state
func (c *Column) LoadState(state ColumnState) {
    c.Width = state.Width
    c.Enabled = state.Enabled
}

// DefaultColumnState returns the default column state
func DefaultColumnState() ColumnState {
    return ColumnState{
        Width:   100,
        Enabled: true,
    }
}

// Helper functions for column management

// FormatValue formats a value based on the column type
func (c *Column) FormatValue(val interface{}) string {
    if val == nil {
        return ""
    }

    switch c.Type {
    case TypeProgress:
        if v, ok := val.(float64); ok {
            return fmt.Sprintf("%.1f%%", v)
        }
    case TypeNumber:
        if v, ok := val.(float64); ok {
            return fmt.Sprintf("%.2f", v)
        }
    case TypeDate:
        if v, ok := val.(int64); ok {
            return time.Unix(v, 0).Format("2006-01-02 15:04:05")
        }
    case TypePeers, TypeSeeds:
        if v, ok := val.(int); ok {
            return fmt.Sprintf("%d", v)
        }
    case TypeStringNoCase:
        return strings.ToLower(fmt.Sprint(val))
    }

    return fmt.Sprint(val)
}

// GetAlignClass returns the CSS class for column alignment
func (c *Column) GetAlignClass() string {
    switch c.Align {
    case AlignLeft:
        return "text-left"
    case AlignCenter:
        return "text-center"
    case AlignRight:
        return "text-right"
    default:
        return "text-left"
    }
}

// GetSortIcon returns the appropriate sort icon class based on sort state
func (c *Column) GetSortIcon(sortCol string, sortReverse bool) string {
    if c.ID != sortCol {
        return ""
    }
    if sortReverse {
        return "sort-desc"
    }
    return "sort-asc"
}

// ValidateColumnWidth ensures the width is within acceptable bounds
func ValidateColumnWidth(width int) int {
    const (
        MinWidth = 50
        MaxWidth = 1000
    )
    
    if width < MinWidth {
        return MinWidth
    }
    if width > MaxWidth {
        return MaxWidth
    }
    return width
}

// ParseColumnState parses column state from JSON
func ParseColumnState(data []byte) (map[string]ColumnState, error) {
    var states map[string]ColumnState
    err := json.Unmarshal(data, &states)
    if err != nil {
        return nil, fmt.Errorf("invalid column state data: %w", err)
    }
    return states, nil
}

// SerializeColumnState serializes column state to JSON
func SerializeColumnState(states map[string]ColumnState) ([]byte, error) {
    data, err := json.Marshal(states)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize column state: %w", err)
    }
    return data, nil
}