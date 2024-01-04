package paths

import (
	"fmt"
	"path/filepath"
)

// PathManager path manager
type PathManager struct {
	resourcePath string
}

// NewPathManager returns a path manager object.
func NewPathManager(resourcePath string) *PathManager {
	return &PathManager{resourcePath: resourcePath}
}

// ServiceFilePath returns a file path of domain.
func (p *PathManager) ServiceFilePath(domain string) string {
	file := fmt.Sprintf("%s.json", domain)
	return filepath.Join(p.resourcePath, "services", file)
}

// TemporaryPath returns a temporary path.
func (p *PathManager) TemporaryPath() string {
	return filepath.Join(p.resourcePath, ".tmp")
}
