// rtorrent/client.go
package rtorrent

import (
    "fmt"
    "time"
    "bytes"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/xml"
)

type Client struct {
    endpoint string
    client   *http.Client
}

func New(endpoint string) *Client {
    return &Client{
        endpoint: endpoint,
        client:   &http.Client{Timeout: 10 * time.Second},
    }
}

type methodCall struct {
    XMLName    xml.Name `xml:"methodCall"`
    MethodName string   `xml:"methodName"`
    Params     []param  `xml:"params>param"`
}

type param struct {
    Value value `xml:"value"`
}

type value struct {
    String  string  `xml:"string,omitempty"`
    Int     int64   `xml:"i8,omitempty"`
    Boolean int     `xml:"boolean,omitempty"`
}

type methodResponse struct {
    XMLName xml.Name `xml:"methodResponse"`
    Params  []param  `xml:"params>param"`
    Fault   *struct {
        Value struct {
            Struct struct {
                Members []struct {
                    Name  string `xml:"name"`
                    Value value  `xml:"value"`
                } `xml:"member"`
            } `xml:"struct"`
        } `xml:"value"`
    } `xml:"fault"`
}

func (c *Client) call(method string, args ...interface{}) (*methodResponse, error) {
    // Construct XML-RPC request
    call := methodCall{
        MethodName: method,
        Params:     make([]param, len(args)),
    }

    // Convert arguments to params
    for i, arg := range args {
        switch v := arg.(type) {
        case string:
            call.Params[i] = param{Value: value{String: v}}
        case int, int64:
            call.Params[i] = param{Value: value{Int: int64(v.(int))}}
        case bool:
            if v {
                call.Params[i] = param{Value: value{Boolean: 1}}
            } else {
                call.Params[i] = param{Value: value{Boolean: 0}}
            }
        default:
            return nil, fmt.Errorf("unsupported argument type: %T", arg)
        }
    }

    // Marshal request to XML
    xmlData, err := xml.MarshalIndent(call, "", "  ")
    if err != nil {
        return nil, err
    }

    // Make HTTP request
    resp, err := c.client.Post(c.endpoint, "text/xml", bytes.NewBuffer(xmlData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Parse response
    var response methodResponse
    if err := xml.Unmarshal(body, &response); err != nil {
        return nil, err
    }

    // Check for fault
    if response.Fault != nil {
        return nil, fmt.Errorf("XML-RPC fault: %v", response.Fault)
    }

    return &response, nil
}

type TorrentInfo struct {
    Hash       string
    Name       string
    Size       int64
    Progress   float64
    Status     string
    Seeds      int
    Peers      int
    DownSpeed  int64
    UpSpeed    int64
    AddedTime  time.Time
}

func (c *Client) GetTorrents() ([]TorrentInfo, error) {
    // Make multicall to get all torrent information
    resp, err := c.call("d.multicall2", "", "main",
        "d.hash=",
        "d.name=",
        "d.size_bytes=",
        "d.completed_bytes=",
        "d.up.rate=",
        "d.down.rate=",
        "d.connection_seed=",
        "d.peers_complete=",
        "d.peers_accounted=",
        "d.creation_date=")
    if err != nil {
        return nil, err
    }

    // Parse response into TorrentInfo structs
    var torrents []TorrentInfo
    // TODO: Parse the response data
    
    return torrents, nil
}

func (c *Client) AddTorrent(data []byte, start bool) error {
    method := "load.raw"
    if start {
        method = "load.raw_start"
    }
    _, err := c.call(method, string(data))
    return err
}

func (c *Client) AddMagnet(uri string, start bool) error {
    method := "load.start"
    if !start {
        method = "load"
    }
    _, err := c.call(method, uri)
    return err
}

func (c *Client) StartTorrent(hash string) error {
    _, err := c.call("d.start", hash)
    return err
}

func (c *Client) PauseTorrent(hash string) error {
    _, err := c.call("d.stop", hash)
    return err
}

func (c *Client) DeleteTorrent(hash string) error {
    _, err := c.call("d.erase", hash)
    return err
}

func (c *Client) GetTorrentDetails(hash string) (*TorrentInfo, error) {
    // TODO: Implement detailed torrent info retrieval
    return nil, nil
}

// Helper functions for formatting
func FormatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatSpeed(bytesPerSecond int64) string {
    return FormatBytes(bytesPerSecond) + "/s"
}
