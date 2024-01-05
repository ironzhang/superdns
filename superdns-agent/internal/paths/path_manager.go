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

// TemporaryPath returns a temporary path.
func (p *PathManager) TemporaryPath() string {
	return filepath.Join(p.resourcePath, ".tmp")
}

// ServiceModelPath returns a service model file path.
func (p *PathManager) ServiceModelPath(domain string) string {
	file := fmt.Sprintf("%s.json", domain)
	return filepath.Join(p.resourcePath, "services", file)
}

// RouteFilePath returns a route model file path.
func (p *PathManager) RouteModelPath(domain string) string {
	file := fmt.Sprintf("%s.json", domain)
	return filepath.Join(p.resourcePath, "routes", file)
}

// RouteScriptPath returns a route script file path.
func (p *PathManager) RouteScriptPath(domain string) string {
	file := fmt.Sprintf("%s.lua", domain)
	return filepath.Join(p.resourcePath, "routes", file)
}
