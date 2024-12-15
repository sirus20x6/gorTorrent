// network/client.go
package network

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)

// Client represents the HTTP client with extra features
type Client struct {
    // HTTP client settings
    host          string
    port          int
    proxyHost     string
    proxyPort     int
    proxyProto    string
    agent         string
    referer       string
    cookies       map[string]string
    rawHeaders    map[string]string
    
    // Configuration
    maxRedirs     int
    connectTimeout time.Duration
    readTimeout    time.Duration
    useGzip       bool
    useIPv4       bool
    bindIP        string

    // Internal state
    status        int
    responseCode  string
    headers       []string
    lastRedirect  string
    client        *http.Client
}

// NewClient creates a new HTTP client instance
type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
    c := &Client{
        maxRedirs:      5,
        connectTimeout: 30 * time.Second,
        readTimeout:    30 * time.Second,
        useGzip:        true,
        cookies:        make(map[string]string),
        rawHeaders:     make(map[string]string),
    }

    // Apply options
    for _, opt := range opts {
        opt(c)
    }

    // Create transport with settings
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   c.connectTimeout,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        MaxIdleConns:          100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }

    // Configure proxy if set
    if c.proxyHost != "" {
        proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", 
            c.proxyProto, c.proxyHost, c.proxyPort))
        if err == nil {
            transport.Proxy = http.ProxyURL(proxyURL)
        }
    }

    // Configure IP binding if set
    if c.bindIP != "" {
        transport.DialContext = (&net.Dialer{
            Timeout:   c.connectTimeout,
            LocalAddr: &net.TCPAddr{IP: net.ParseIP(c.bindIP)},
        }).DialContext
    }

    // Create HTTP client
    c.client = &http.Client{
        Transport: transport,
        Timeout:   c.readTimeout,
    }

    // Configure redirect handling
    c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if len(via) >= c.maxRedirs {
            return fmt.Errorf("stopped after %d redirects", c.maxRedirs)
        }
        c.lastRedirect = req.URL.String()
        return nil
    }

    return c
}

// WithProxy sets proxy settings
func WithProxy(proto, host string, port int) ClientOption {
    return func(c *Client) {
        c.proxyProto = proto
        c.proxyHost = host
        c.proxyPort = port
    }
}

// WithTimeout sets timeouts
func WithTimeout(connect, read time.Duration) ClientOption {
    return func(c *Client) {
        c.connectTimeout = connect
        c.readTimeout = read
    }
}

// WithBindIP sets the local IP to bind to
func WithBindIP(ip string) ClientOption {
    return func(c *Client) {
        c.bindIP = ip
    }
}

// SetCookie sets a cookie
func (c *Client) SetCookie(name, value string) {
    c.cookies[name] = value
}

// SetHeader sets a raw header
func (c *Client) SetHeader(name, value string) {
    c.rawHeaders[name] = value
}

// Fetch makes an HTTP request
func (c *Client) Fetch(method, urlStr string, body io.Reader) ([]byte, error) {
    // Create request
    req, err := http.NewRequest(method, urlStr, body)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Set headers
    if c.agent != "" {
        req.Header.Set("User-Agent", c.agent)
    }
    if c.referer != "" {
        req.Header.Set("Referer", c.referer)
    }
    if c.useGzip {
        req.Header.Set("Accept-Encoding", "gzip")
    }

    // Add cookies
    var cookies []string
    for name, value := range c.cookies {
        cookies = append(cookies, fmt.Sprintf("%s=%s", name, value))
    }
    if len(cookies) > 0 {
        req.Header.Set("Cookie", strings.Join(cookies, "; "))
    }

    // Add raw headers
    for name, value := range c.rawHeaders {
        req.Header.Set(name, value)
    }

    // Make request
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    // Update state
    c.status = resp.StatusCode
    c.responseCode = resp.Status
    for name, values := range resp.Header {
        c.headers = append(c.headers, fmt.Sprintf("%s: %s", name, strings.Join(values, ", ")))
    }

    // Save cookies from response
    for _, cookie := range resp.Cookies() {
        if cookie.Value == "deleted" {
            delete(c.cookies, cookie.Name)
        } else {
            c.cookies[cookie.Name] = cookie.Value
        }
    }

    // Read response body
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    return data, nil
}

// FetchComplex makes a request with optional content type and body
func (c *Client) FetchComplex(urlStr, method, contentType string, body string) ([]byte, error) {
    var bodyReader io.Reader
    if body != "" {
        bodyReader = bytes.NewBufferString(body)
        if contentType != "" {
            c.SetHeader("Content-Type", contentType)
            c.SetHeader("Content-Length", strconv.Itoa(len(body)))
        }
    }

    return c.Fetch(method, urlStr, bodyReader)
}

// Status returns the last response status code
func (c *Client) Status() int {
    return c.status
}

// Headers returns the last response headers
func (c *Client) Headers() []string {
    return c.headers
}

// LastRedirect returns the URL of the last redirect
func (c *Client) LastRedirect() string {
    return c.lastRedirect
}

// GetCookie gets a cookie value by name
func (c *Client) GetCookie(name string) string {
    return c.cookies[name]
}