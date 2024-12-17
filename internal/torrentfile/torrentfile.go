// internal/services/torrentfile/torrentfile.go
package torrentfile

import (
    "bytes"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    "github.com/anacrolix/torrent/bencode"
    "github.com/anacrolix/torrent/metainfo"
)

// Torrent represents a torrent file structure
type Torrent struct {
    info     *metainfo.Info
    metainfo *metainfo.MetaInfo
    filepath string
    errors   []error
}

// New creates a new torrent instance from file
func New(path string) (*Torrent, error) {
    mi, err := metainfo.LoadFromFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to load torrent: %w", err)
    }

    info, err := mi.UnmarshalInfo()
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal info: %w", err)
    }

    return &Torrent{
        info:     &info,
        metainfo: mi,
        filepath: path,
    }, nil
}

// NewFromBytes creates a torrent instance from bytes
func NewFromBytes(data []byte) (*Torrent, error) {
    var mi metainfo.MetaInfo
    if err := bencode.NewDecoder(bytes.NewReader(data)).Decode(&mi); err != nil {
        return nil, fmt.Errorf("failed to decode torrent: %w", err)
    }

    info, err := mi.UnmarshalInfo()
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal info: %w", err)
    }

    return &Torrent{
        info:     &info,
        metainfo: &mi,
    }, nil
}

// CreateFromPath creates a new torrent from a file or directory
func CreateFromPath(path string, opts CreateOptions) (*Torrent, error) {
    // Validate path
    info, err := os.Stat(path)
    if err != nil {
        return nil, fmt.Errorf("failed to stat path: %w", err)
    }

    // Create builder
    builder := metainfo.Builder{}
    if info.IsDir() {
        if err := builder.AddDir(path); err != nil {
            return nil, fmt.Errorf("failed to add directory: %w", err)
        }
    } else {
        if err := builder.AddFile(path); err != nil {
            return nil, fmt.Errorf("failed to add file: %w", err)
        }
    }

    // Set options
    builder.PieceLength = opts.PieceLength
    if opts.Private {
        builder.Private = true
    }

    // Build metainfo
    mi, err := builder.Submit()
    if err != nil {
        return nil, fmt.Errorf("failed to build torrent: %w", err)
    }

    // Set announce lists
    if len(opts.Trackers) > 0 {
        mi.AnnounceList = [][]string{opts.Trackers}
    }

    // Set creation date
    mi.CreationDate = time.Now().Unix()

    // Set created by
    mi.CreatedBy = opts.CreatedBy

    // Get info
    info, err := mi.UnmarshalInfo()
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal info: %w", err)
    }

    return &Torrent{
        info:     &info,
        metainfo: mi,
    }, nil
}

// CreateOptions holds options for torrent creation
type CreateOptions struct {
    PieceLength int64
    Private     bool
    Trackers    []string
    CreatedBy   string
}

// Save saves the torrent to a file
func (t *Torrent) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer f.Close()

    return bencode.NewEncoder(f).Encode(t.metainfo)
}

// GetName returns torrent name
func (t *Torrent) GetName() string {
    return t.info.Name
}

// GetSize returns total size
func (t *Torrent) GetSize() int64 {
    return t.info.TotalLength()
}

// GetInfoHash returns info hash
func (t *Torrent) GetInfoHash() string {
    return t.metainfo.HashInfoBytes().String()
}

// GetAnnounceList returns tracker list
func (t *Torrent) GetAnnounceList() [][]string {
    return t.metainfo.AnnounceList
}

// GetCreationDate returns creation date
func (t *Torrent) GetCreationDate() time.Time {
    return time.Unix(t.metainfo.CreationDate, 0)
}

// GetCreatedBy returns created by string
func (t *Torrent) GetCreatedBy() string {
    return t.metainfo.CreatedBy
}

// GetFiles returns file list
func (t *Torrent) GetFiles() []File {
    files := make([]File, 0)
    
    if !t.info.IsDir() {
        // Single file mode
        files = append(files, File{
            Path: t.info.Name,
            Size: t.info.Length,
        })
    } else {
        // Multi file mode
        for _, f := range t.info.Files {
            files = append(files, File{
                Path: filepath.Join(f.Path...),
                Size: f.Length,
            })
        }
    }

    return files
}

// File represents a file in the torrent
type File struct {
    Path string
    Size int64
}

// Errors returns any errors encountered
func (t *Torrent) Errors() []error {
    return t.errors
}

// AddTracker adds a tracker
func (t *Torrent) AddTracker(announce string) {
    if len(t.metainfo.AnnounceList) == 0 {
        t.metainfo.AnnounceList = [][]string{{announce}}
    } else {
        t.metainfo.AnnounceList = append(t.metainfo.AnnounceList, []string{announce})
    }
}

// SetPrivate sets private flag
func (t *Torrent) SetPrivate(private bool) {
    t.info.Private = &private
}

// IsPrivate returns private status
func (t *Torrent) IsPrivate() bool {
    return t.info.Private != nil && *t.info.Private
}

// CalculateHash calculates info hash
func CalculateHash(data []byte) string {
    hash := sha1.New()
    hash.Write(data)
    return hex.EncodeToString(hash.Sum(nil))
}