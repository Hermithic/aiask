package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/Hermithic/aiask/internal/config"
	"github.com/Hermithic/aiask/internal/fileutil"
	"gopkg.in/yaml.v3"
)

// templateNameRegex validates template names (alphanumeric, dashes, underscores)
var templateNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Template represents a saved prompt template
type Template struct {
	Name        string    `yaml:"name"`
	Prompt      string    `yaml:"prompt"`
	Description string    `yaml:"description,omitempty"`
	CreatedAt   time.Time `yaml:"created_at"`
	UsageCount  int       `yaml:"usage_count"`
}

// Templates represents the collection of saved templates
type Templates struct {
	Items []Template `yaml:"templates"`
}

// GetTemplatesPath returns the path to the templates file
func GetTemplatesPath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "templates.yaml"), nil
}

// Load loads the templates from the templates file
func Load() (*Templates, error) {
	templatesPath, err := GetTemplatesPath()
	if err != nil {
		return nil, err
	}

	// Return empty templates if file doesn't exist
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		return &Templates{Items: []Template{}}, nil
	}

	data, err := os.ReadFile(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates file: %w", err)
	}

	templates := &Templates{}
	if err := yaml.Unmarshal(data, templates); err != nil {
		return nil, fmt.Errorf("failed to parse templates file: %w", err)
	}

	return templates, nil
}

// Save saves the templates to the templates file atomically to prevent corruption
func (t *Templates) Save() error {
	templatesPath, err := GetTemplatesPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal templates: %w", err)
	}

	if err := fileutil.AtomicWriteFile(templatesPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write templates file: %w", err)
	}

	return nil
}

// ValidateTemplateName validates a template name
func ValidateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("template name cannot exceed 50 characters")
	}
	if !templateNameRegex.MatchString(name) {
		return fmt.Errorf("template name must start with a letter and contain only letters, numbers, dashes, and underscores")
	}
	return nil
}

// Add adds a new template
func (t *Templates) Add(name, prompt, description string) error {
	// Validate the template name
	if err := ValidateTemplateName(name); err != nil {
		return err
	}

	// Check if template with same name exists
	for _, tmpl := range t.Items {
		if tmpl.Name == name {
			return fmt.Errorf("template '%s' already exists", name)
		}
	}

	// Validate prompt is not empty
	if prompt == "" {
		return fmt.Errorf("template prompt cannot be empty")
	}

	t.Items = append(t.Items, Template{
		Name:        name,
		Prompt:      prompt,
		Description: description,
		CreatedAt:   time.Now(),
		UsageCount:  0,
	})

	return nil
}

// Get gets a template by name
func (t *Templates) Get(name string) (*Template, error) {
	for i := range t.Items {
		if t.Items[i].Name == name {
			return &t.Items[i], nil
		}
	}
	return nil, fmt.Errorf("template '%s' not found", name)
}

// Remove removes a template by name
func (t *Templates) Remove(name string) error {
	for i, tmpl := range t.Items {
		if tmpl.Name == name {
			t.Items = append(t.Items[:i], t.Items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("template '%s' not found", name)
}

// IncrementUsage increments the usage count for a template
func (t *Templates) IncrementUsage(name string) {
	for i := range t.Items {
		if t.Items[i].Name == name {
			t.Items[i].UsageCount++
			return
		}
	}
}

// List returns all templates sorted by name
func (t *Templates) List() []Template {
	sorted := make([]Template, len(t.Items))
	copy(sorted, t.Items)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

// ListByUsage returns all templates sorted by usage count (descending)
func (t *Templates) ListByUsage() []Template {
	sorted := make([]Template, len(t.Items))
	copy(sorted, t.Items)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].UsageCount > sorted[j].UsageCount
	})
	return sorted
}

