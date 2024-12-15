// xmlrpc/xmlrpc.go
package xmlrpc

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "time"
)

const (
    XMLRPC_MAX_I4 = 2147483647
    XMLRPC_MIN_I4 = -2147483648
    XMLRPC_MIN_I8 = -9.999999999999999E+15
    XMLRPC_MAX_I8 = 9.999999999999999E+15
)

// XMLRPCParam represents a parameter in an XML-RPC request/response
type XMLRPCParam struct {
    Type  string
    Value interface{}
}

// XMLRPCCommand represents an XML-RPC command
type XMLRPCCommand struct {
    Command string
    Params  []XMLRPCParam
}

// XMLRPCRequest represents an XML-RPC request batch
type XMLRPCRequest struct {
    Commands []XMLRPCCommand
    Content  string
    Results  struct {
        I8s     []int64
        Strings []string
        Values  []interface{}
    }
    Fault      bool
    Important  bool
}

// NewXMLRPCCommand creates a new XML-RPC command
func NewXMLRPCCommand(cmd string, args ...interface{}) XMLRPCCommand {
    command := XMLRPCCommand{
        Command: cmd,
        Params:  make([]XMLRPCParam, 0),
    }

    for _, arg := range args {
        command.AddParameter(arg)
    }

    return command
}

// AddParameter adds a parameter to the command
func (c *XMLRPCCommand) AddParameter(value interface{}, typ ...string) {
    paramType := ""
    if len(typ) > 0 {
        paramType = typ[0]
    } else {
        paramType = getParamType(value)
    }

    c.Params = append(c.Params, XMLRPCParam{
        Type:  paramType,
        Value: value,
    })
}

// getParamType determines the XML-RPC type for a value
func getParamType(value interface{}) string {
    switch v := value.(type) {
    case int:
        if v >= XMLRPC_MIN_I4 && v <= XMLRPC_MAX_I4 {
            return "i4"
        }
        return "i8"
    case float64:
        return "i8"
    default:
        return "string"
    }
}

// NewXMLRPCRequest creates a new XML-RPC request
func NewXMLRPCRequest(cmds ...XMLRPCCommand) *XMLRPCRequest {
    req := &XMLRPCRequest{
        Important: true,
    }
    if cmds != nil {
        req.AddCommands(cmds...)
    }
    return req
}

// AddCommand adds a command to the request
func (r *XMLRPCRequest) AddCommand(cmd XMLRPCCommand) {
    r.Commands = append(r.Commands, cmd)
}

// AddCommands adds multiple commands to the request
func (r *XMLRPCRequest) AddCommands(cmds ...XMLRPCCommand) {
    r.Commands = append(r.Commands, cmds...)
}

// makeNextCall prepares the next XML-RPC call
type methodCall struct {
    XMLName    xml.Name `xml:"methodCall"`
    MethodName string   `xml:"methodName"`
    Params     []param  `xml:"params>param"`
}

type param struct {
    Value value `xml:"value"`
}

type value struct {
    String string `xml:"string,omitempty"`
    Int    int64  `xml:"i8,omitempty"`
    Base64 string `xml:"base64,omitempty"`
}

func (r *XMLRPCRequest) makeNextCall() bool {
    r.Fault = false
    r.Content = ""

    if len(r.Commands) == 0 {
        return false
    }

    var call methodCall
    if len(r.Commands) == 1 {
        cmd := r.Commands[0]
        call.MethodName = cmd.Command
        call.Params = make([]param, len(cmd.Params))
        for i, p := range cmd.Params {
            call.Params[i] = paramToXML(p)
        }
    } else {
        call.MethodName = "system.multicall"
        calls := make([]map[string]interface{}, len(r.Commands))
        for i, cmd := range r.Commands {
            calls[i] = map[string]interface{}{
                "methodName": cmd.Command,
                "params":     cmd.Params,
            }
        }
        call.Params = []param{{
            Value: value{Base64: encodeStruct(calls)},
        }}
    }

    data, err := xml.Marshal(call)
    if err != nil {
        r.Fault = true
        return false
    }

    r.Content = string(data)
    return true
}

// paramToXML converts a parameter to XML format
func paramToXML(p XMLRPCParam) param {
    switch p.Type {
    case "i4", "i8":
        if val, ok := p.Value.(int64); ok {
            return param{Value: value{Int: val}}
        }
        if val, ok := p.Value.(int); ok {
            return param{Value: value{Int: int64(val)}}
        }
    default:
        return param{Value: value{String: fmt.Sprintf("%v", p.Value)}}
    }
    return param{Value: value{String: fmt.Sprintf("%v", p.Value)}}
}

// Send sends the XML-RPC request
func (r *XMLRPCRequest) Send(endpoint string) error {
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    for r.makeNextCall() {
        resp, err := client.Post(endpoint, "text/xml", bytes.NewReader([]byte(r.Content)))
        if err != nil {
            return fmt.Errorf("failed to send request: %w", err)
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("failed to read response: %w", err)
        }

        if err := r.parseResponse(body); err != nil {
            return err
        }
    }

    return nil
}

// parseResponse parses the XML-RPC response
func (r *XMLRPCRequest) parseResponse(data []byte) error {
    var resp struct {
        XMLName xml.Name `xml:"methodResponse"`
        Params  []param  `xml:"params>param"`
        Fault   *struct {
            Value value `xml:"value"`
        } `xml:"fault"`
    }

    if err := xml.Unmarshal(data, &resp); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }

    if resp.Fault != nil {
        r.Fault = true
        return fmt.Errorf("XML-RPC fault: %v", resp.Fault.Value)
    }

    for _, p := range resp.Params {
        if p.Value.String != "" {
            r.Results.Strings = append(r.Results.Strings, p.Value.String)
        }
        if p.Value.Int != 0 {
            r.Results.I8s = append(r.Results.I8s, p.Value.Int)
        }
        r.Results.Values = append(r.Results.Values, p.Value)
    }

    return nil
}

// encodeStruct encodes a struct as base64
func encodeStruct(v interface{}) string {
    data, _ := xml.Marshal(v)
    return string(data)
}