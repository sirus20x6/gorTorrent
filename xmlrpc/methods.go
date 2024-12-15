// xmlrpc/methods.go
package xmlrpc

import (
    "fmt"
    "strings"
)

// MethodAlias represents an XML-RPC method alias
type MethodAlias struct {
    Name        string
    ActualName  string
    ParamCount  int
}

// MethodManager handles XML-RPC method aliases
type MethodManager struct {
    aliases map[string]MethodAlias
}

// NewMethodManager creates a new method manager
func NewMethodManager() *MethodManager {
    return &MethodManager{
        aliases: make(map[string]MethodAlias),
    }
}

// AddAlias adds a method alias
func (m *MethodManager) AddAlias(name, actualName string, paramCount int) {
    m.aliases[name] = MethodAlias{
        Name:       name,
        ActualName: actualName,
        ParamCount: paramCount,
    }
}

// GetMethod returns the actual method name and parameter count
func (m *MethodManager) GetMethod(name string) (string, int) {
    if alias, ok := m.aliases[name]; ok {
        return alias.ActualName, alias.ParamCount
    }
    return name, 0
}

// GetCommand gets a command with proper method name and parameters
func (m *MethodManager) GetCommand(name string, args ...interface{}) XMLRPCCommand {
    methodName, paramCount := m.GetMethod(name)
    cmd := NewXMLRPCCommand(methodName)

    // Add provided arguments
    for i, arg := range args {
        if i >= paramCount {
            break
        }
        cmd.AddParameter(arg)
    }

    // Add empty strings for missing parameters
    for i := len(args); i < paramCount; i++ {
        cmd.AddParameter("")
    }

    return cmd
}

// RegisterDefaultAliases registers the default rTorrent method aliases
func (m *MethodManager) RegisterDefaultAliases() {
    // Core methods
    m.AddAlias("d.get_hash", "d.hash", 0)
    m.AddAlias("d.get_name", "d.name", 0)
    m.AddAlias("d.get_state", "d.state", 0)
    m.AddAlias("d.get_size_bytes", "d.size_bytes", 0)
    m.AddAlias("d.get_bytes_done", "d.bytes_done", 0)
    m.AddAlias("d.get_up_rate", "d.up.rate", 0)
    m.AddAlias("d.get_down_rate", "d.down.rate", 0)
    m.AddAlias("d.get_chunk_size", "d.chunk_size", 0)
    m.AddAlias("d.get_custom1", "d.custom1", 0)
    m.AddAlias("d.get_peers_accounted", "d.peers_accounted", 0)
    m.AddAlias("d.get_peers_complete", "d.peers_complete", 0)
    m.AddAlias("d.get_creation_date", "d.creation_date", 0)

    // System methods
    m.AddAlias("system.get_cwd", "system.cwd", 0)
    m.AddAlias("get_directory", "directory.default", 0)
    m.AddAlias("get_session", "session.path", 0)

    // Setting methods
    m.AddAlias("d.set_custom1", "d.custom1.set", 1)
    m.AddAlias("d.set_directory", "d.directory.set", 1)
    m.AddAlias("d.set_peer_exchange", "d.peer_exchange.set", 1)
}

// TransformMulticall transforms a multicall command based on aliases
func (m *MethodManager) TransformMulticall(multiCmd XMLRPCCommand) XMLRPCCommand {
    for i, param := range multiCmd.Params {
        if strings.HasPrefix(param.Value.(string), "d.get_") {
            actualName, _ := m.GetMethod(param.Value.(string))
            multiCmd.Params[i].Value = actualName
        }
    }
    return multiCmd
}