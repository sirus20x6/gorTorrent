// internal/webui/webui.go

package webui

import (
    "your-project/internal/rtorrent"
    "your-project/internal/table"
    "sync"
)

type WebUI struct {
    Version string
    RTorrent *rtorrent.Client
    Tables   map[string]*table.Table
    Settings map[string]interface{}
    mu       sync.RWMutex
}

// Settings keys
const (
    SettingShowCats           = "webui.show_cats"
    SettingShowDets           = "webui.show_dets"
    SettingUpdateInterval     = "webui.update_interval"
    SettingConfirmDelete      = "webui.confirm_when_deleting"
    SettingAlternateColor     = "webui.alternate_color"
    SettingShowSearchLabel    = "webui.show_searchlabelsize"
    SettingLabelTextOverflow  = "webui.show_label_text_overflow"
)

// New creates a new WebUI instance
func New(rtorrentURL string) *WebUI {
    return &WebUI{
        Version:  "1.0.0",
        RTorrent: rtorrent.New(rtorrentURL),
        Tables:   make(map[string]*table.Table),
        Settings: map[string]interface{}{
            SettingShowCats:          true,
            SettingShowDets:          true,
            SettingUpdateInterval:    2500,
            SettingConfirmDelete:     true,
            SettingAlternateColor:    false,
            SettingShowSearchLabel:   false,
            SettingLabelTextOverflow: true,
        },
    }
}