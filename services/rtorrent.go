// services/rtorrent.go
package services

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
)

type RTorrentClient struct {
    endpoint string
}

func NewRTorrentClient(endpoint string) *RTorrentClient {
    return &RTorrentClient{
        endpoint: endpoint,
    }
}

// XMLRPCRequest represents an XML-RPC request
type XMLRPCRequest struct {
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
    Base64  string  `xml:"base64,omitempty"`
}

// XMLRPCResponse represents an XML-RPC response
type XMLRPCResponse struct {
    XMLName xml.Name    `xml:"methodResponse"`
    Params  []param     `xml:"params>param,omitempty"`
    Fault   *faultNode `xml:"fault,omitempty"`
}

type faultNode struct {
    Value value `xml:"value"`
}

func (c *RTorrentClient) Call(method string, args ...interface{}) (*XMLRPCResponse, error) {
    // Create request body
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

    // Marshal request to XML
    body, err := xml.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("error marshaling request: %w", err)
    }

    // Create HTTP request
    httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "text/xml")

    // Send request
    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    // Read response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    // Parse response
    var xmlResp XMLRPCResponse
    if err := xml.Unmarshal(respBody, &xmlResp); err != nil {
        return nil, fmt.Errorf("error parsing response: %w", err)
    }

    // Check for fault
    if xmlResp.Fault != nil {
        return nil, fmt.Errorf("rTorrent error: %v", xmlResp.Fault.Value)
    }

    return &xmlResp, nil
}

// Helper methods for common rTorrent commands
func (c *RTorrentClient) GetDownloadList() ([]string, error) {
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

func (c *RTorrentClient) GetTorrentInfo(hash string) (map[string]interface{}, error) {
    // Get multiple torrent attributes in one call
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

    // Parse response into map
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