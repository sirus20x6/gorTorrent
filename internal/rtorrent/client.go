// internal/rtorrent/client.go

package rtorrent

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    endpoint string
    client   *http.Client
}

func New(endpoint string) *Client {
    return &Client{
        endpoint: endpoint,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// Request/Response types
type XMLRPCRequest struct {
    XMLName    xml.Name `xml:"methodCall"`  
    MethodName string   `xml:"methodName"`
    Params     []param  `xml:"params>param"`
}

type param struct {
    Value value `xml:"value"`
}

type value struct { 
    String  string `xml:"string,omitempty"`
    Int     int64  `xml:"i8,omitempty"`
    Boolean int    `xml:"boolean,omitempty"`
    Base64  string `xml:"base64,omitempty"`
}

type XMLRPCResponse struct {
    XMLName xml.Name   `xml:"methodResponse"`
    Params  []param    `xml:"params>param,omitempty"`
    Fault   *faultNode `xml:"fault,omitempty"` 
}

type faultNode struct {
    Value value `xml:"value"`
}

// Call makes an XML-RPC request to rTorrent
func (c *Client) Call(method string, args ...interface{}) (*XMLRPCResponse, error) {
    req := XMLRPCRequest{
        MethodName: method,
        Params:     make([]param, len(args)),
    }

    // Convert args to params
    for i, arg := range args {
        switch v := arg.(type) {
        case string:
            req.Params[i] = param{Value: value{String: v}}
        case int64:
            req.Params[i] = param{Value: value{Int: v}} 
        case bool:
            if v {
                req.Params[i] = param{Value: value{Boolean: 1}}
            } else {
                req.Params[i] = param{Value: value{Boolean: 0}}
            }
        default:
            return nil, fmt.Errorf("unsupported argument type: %T", arg)
        }
    }

    body, err := xml.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("error marshaling request: %w", err) 
    }

    httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "text/xml")

    resp, err := c.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    var xmlResp XMLRPCResponse
    if err := xml.Unmarshal(respBody, &xmlResp); err != nil {
        return nil, fmt.Errorf("error parsing response: %w", err)
    }

    if xmlResp.Fault != nil {
        return nil, fmt.Errorf("rTorrent error: %v", xmlResp.Fault.Value)
    }

    return &xmlResp, nil
}

// Helper types for torrent info
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

// Common rTorrent commands
func (c *Client) GetDownloadList() ([]string, error) {
    resp, err := c.Call("download_list")
    if err != nil {
        return nil, err
    }

    var hashes []string
    for _, param := range resp.Params {
        hashes = append(hashes, param.Value.String)
    }
    return hashes, nil
}

func (c *Client) GetTorrentInfo(hash string) (map[string]interface{}, error) {
    resp, err := c.Call("d.multicall2", "", hash,
        "d.name=",
        "d.size_bytes=", 
        "d.completed_bytes=",
        "d.down.rate=",
        "d.up.rate=",
        "d.state=",
        "d.peers_complete=",
        "d.peers_connected=",
    )
    if err != nil {
        return nil, err
    }

    info := make(map[string]interface{})
    if len(resp.Params) > 0 {
        info["name"] = resp.Params[0].Value.String
        info["size"] = resp.Params[1].Value.Int
        info["completed"] = resp.Params[2].Value.Int
        info["download_rate"] = resp.Params[3].Value.Int
        info["upload_rate"] = resp.Params[4].Value.Int 
        info["state"] = resp.Params[5].Value.Int
        info["seeders"] = resp.Params[6].Value.Int
        info["peers"] = resp.Params[7].Value.Int
    }

    return info, nil
}

func (c *Client) AddTorrent(data []byte, start bool) error {
    method := "load.raw"
    if start {
        method = "load.raw_start"
    }
    _, err := c.Call(method, string(data))
    return err
}

func (c *Client) AddMagnet(uri string, start bool) error {
    method := "load.start"
    if !start {
        method = "load"
    }
    _, err := c.Call(method, uri)
    return err
}

func (c *Client) StartTorrent(hash string) error {
    _, err := c.Call("d.start", hash)
    return err
}

func (c *Client) PauseTorrent(hash string) error {
    _, err := c.Call("d.stop", hash) 
    return err
}

func (c *Client) DeleteTorrent(hash string) error {
    _, err := c.Call("d.erase", hash)
    return err
}