// util/minifier/minifier.go
package minifier

import (
    "bytes"
    "fmt"
    "strings"
    "unicode"
)

// Options represents minifier options
type Options struct {
    FlaggedComments bool
}

// DefaultOptions returns default minifier options
var DefaultOptions = Options{
    FlaggedComments: true,
}

// Minifier represents a JavaScript minifier instance
type Minifier struct {
    input        string
    output       string
    index        int
    len          int
    a            string // last character
    b            string // current character
    c            string // next character
    lastCharType rune  // last non-whitespace character type
    options      Options
}

// keywords that can appear before a regex
var keywords = map[string]bool{
    "delete":      true,
    "do":          true,
    "for":         true,
    "in":          true,
    "instanceof":  true,
    "return":      true,
    "typeof":      true,
    "yield":       true,
}

// stringDelimiters is a set of quote characters
var stringDelimiters = map[rune]bool{
    '"':  true,
    '\'': true,
    '`':  true,
}

// New creates a new minifier instance
func New(opts Options) *Minifier {
    return &Minifier{
        options: opts,
    }
}

// Minify minifies JavaScript code
func (m *Minifier) Minify(js string) (string, error) {
    // Reset state
    m.input = js + "\n" // Add newline to handle final comments
    m.output = ""
    m.index = 0
    m.len = len(m.input)
    m.a = "\n"
    m.b = "\n"
    m.lastCharType = 0

    // Lock string patterns for later replacement
    js = m.lock(js)

    // Process input
    if err := m.loop(); err != nil {
        return "", err
    }

    // Unlock string patterns
    return m.unlock(m.output), nil
}

// lock preserves string patterns we don't want to modify
func (m *Minifier) lock(js string) string {
    // Implement pattern locking similar to JShrink
    // Add implementation here based on needs
    return js
}

// unlock restores preserved string patterns
func (m *Minifier) unlock(js string) string {
    // Implement pattern unlocking
    // Add implementation here based on needs
    return js
}

// loop processes the input character by character
func (m *Minifier) loop() error {
    for m.a != "" {
        if err := m.processChar(); err != nil {
            return err
        }
    }
    return nil
}

// processChar processes the current character
func (m *Minifier) processChar() error {
    switch m.a {
    case "/":
        if err := m.handlePotentialComment(); err != nil {
            return err
        }
    case "'", "\"", "`":
        if err := m.handleString(); err != nil {
            return err
        }
    case "\n":
        if err := m.handleNewline(); err != nil {
            return err
        }
    default:
        if err := m.handleDefault(); err != nil {
            return err
        }
    }
    
    return m.advance()
}

// handlePotentialComment handles potential comment starts
func (m *Minifier) handlePotentialComment() error {
    if m.b == "*" {
        return m.handleMultiLineComment()
    }
    if m.b == "/" {
        return m.handleSingleLineComment()
    }
    if m.shouldStartRegex() {
        return m.handleRegex()
    }
    
    m.output += m.a
    return nil
}

// shouldStartRegex determines if we're at a position where a regex can start
func (m *Minifier) shouldStartRegex() bool {
    switch m.lastCharType {
    case 0:
        return true
    case ')', ']':
        return false
    case '}':
        // Check for keywords before the block
        // TODO: Implement keyword checking
        return false
    default:
        // Check if previous token is a keyword
        return keywords[m.lastToken()]
    }
}

// handleMultiLineComment handles /* ... */ comments
func (m *Minifier) handleMultiLineComment() error {
    var comment bytes.Buffer
    
    // Skip /* characters
    m.advance()
    m.advance()
    
    for m.a != "" && !(m.a == "*" && m.b == "/") {
        comment.WriteString(m.a)
        if err := m.advance(); err != nil {
            return err
        }
    }
    
    // Keep important comments
    if m.options.FlaggedComments && strings.HasPrefix(comment.String(), "!") {
        m.output += "/*" + comment.String() + "*/"
    }
    
    // Skip closing */
    m.advance()
    m.advance()
    return nil
}

// handleSingleLineComment handles // comments
func (m *Minifier) handleSingleLineComment() error {
    // Skip to end of line
    for m.a != "" && m.a != "\n" {
        if err := m.advance(); err != nil {
            return err
        }
    }
    return nil
}

// handleString handles string literals
func (m *Minifier) handleString() error {
    delimiter := m.a
    m.output += delimiter
    
    if err := m.advance(); err != nil {
        return err
    }
    
    for m.a != "" && m.a != delimiter {
        if m.a == "\\" {
            m.output += m.a
            if err := m.advance(); err != nil {
                return err
            }
        }
        m.output += m.a
        if err := m.advance(); err != nil {
            return err
        }
    }
    
    m.output += delimiter
    return nil
}

// handleRegex handles regular expression literals
func (m *Minifier) handleRegex() error {
    m.output += "/"
    
    if err := m.advance(); err != nil {
        return err
    }
    
    for m.a != "" && m.a != "/" {
        if m.a == "\\" {
            m.output += m.a
            if err := m.advance(); err != nil {
                return err
            }
        }
        m.output += m.a
        if err := m.advance(); err != nil {
            return err
        }
    }
    
    m.output += "/"
    
    // Include regex flags
    for m.b != "" && isRegexFlag(rune(m.b[0])) {
        if err := m.advance(); err != nil {
            return err
        }
        m.output += m.a
    }
    
    return nil
}

// handleNewline handles newline characters
func (m *Minifier) handleNewline() error {
    // Skip consecutive newlines
    for m.a == "\n" {
        if err := m.advance(); err != nil {
            return err
        }
    }
    
    // Preserve newline only if necessary
    if m.b != "" && shouldPreserveNewline(rune(m.b[0])) {
        m.output += "\n"
    }
    return nil
}

// handleDefault handles standard characters
func (m *Minifier) handleDefault() error {
    if m.isWhitespace(m.a) {
        // Skip whitespace unless necessary
        if !m.isWhitespace(m.b) && isAlphanumeric(rune(m.b[0])) {
            m.output += " "
        }
        return nil
    }
    
    m.output += m.a
    m.lastCharType = rune(m.a[0])
    return nil
}

// advance moves to the next character
func (m *Minifier) advance() error {
    if m.index >= m.len {
        m.a = ""
        m.b = ""
        return nil
    }
    
    m.a = m.b
    m.b = string(m.input[m.index])
    m.index++
    return nil
}

// Helper functions

func (m *Minifier) isWhitespace(s string) bool {
    if s == "" {
        return false
    }
    return unicode.IsSpace(rune(s[0]))
}

func isAlphanumeric(r rune) bool {
    return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isRegexFlag(r rune) bool {
    return strings.ContainsRune("gimsuy", r)
}

func shouldPreserveNewline(r rune) bool {
    return strings.ContainsRune("([+-", r)
}

// lastToken returns the last meaningful token
func (m *Minifier) lastToken() string {
    // TODO: Implement token tracking
    return ""
}

// Package level convenience function
func Minify(js string) (string, error) {
    m := New(DefaultOptions)
    return m.Minify(js)
}