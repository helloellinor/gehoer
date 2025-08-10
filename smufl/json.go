package smufl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Core types matching SMuFL JSON files:

type Glyph struct {
	Codepoint          string `json:"codepoint"`
	Description        string `json:"description"`
	AlternateCodepoint string `json:"alternateCodepoint,omitempty"`
}

type Range struct {
	Description string   `json:"description"`
	Glyphs      []string `json:"glyphs"`
	RangeStart  string   `json:"range_start"`
	RangeEnd    string   `json:"range_end"`
}

type Class []string

type Metadata struct {
	Glyphs  map[string]Glyph `json:"glyphs"`
	Ranges  map[string]Range `json:"ranges"`
	Classes map[string]Class `json:"classes"`
}

// LoadMetadata loads the SMuFL metadata from a repo path
func LoadMetadata(smuflRepoPath string) (*Metadata, error) {
	meta := &Metadata{
		Glyphs:  make(map[string]Glyph),
		Ranges:  make(map[string]Range),
		Classes: make(map[string]Class),
	}

	if err := loadJSON(filepath.Join(smuflRepoPath, "metadata", "glyphnames.json"), &meta.Glyphs); err != nil {
		return nil, fmt.Errorf("glyphnames.json: %w", err)
	}
	if err := loadJSON(filepath.Join(smuflRepoPath, "metadata", "ranges.json"), &meta.Ranges); err != nil {
		return nil, fmt.Errorf("ranges.json: %w", err)
	}
	if err := loadJSON(filepath.Join(smuflRepoPath, "metadata", "classes.json"), &meta.Classes); err != nil {
		return nil, fmt.Errorf("classes.json: %w", err)
	}

	return meta, nil
}

func loadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// GetGlyphRune returns the Unicode rune for a glyph name
func (m *Metadata) GetGlyphRune(name string) (rune, error) {
	glyph, ok := m.Glyphs[name]
	if !ok {
		return 0, fmt.Errorf("glyph %s not found", name)
	}
	codeStr := strings.TrimPrefix(glyph.Codepoint, "U+")
	code, err := strconv.ParseInt(codeStr, 16, 32)
	if err != nil {
		return 0, err
	}
	return rune(code), nil
}

// Utility functions to get glyph lists by range or list ranges/classes

func (m *Metadata) GetGlyphsByRange(name string) ([]string, error) {
	r, ok := m.Ranges[name]
	if !ok {
		return nil, fmt.Errorf("range %s not found", name)
	}
	return r.Glyphs, nil
}

func (m *Metadata) ListRanges() []string {
	var keys []string
	for k := range m.Ranges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *Metadata) ListClasses() []string {
	var keys []string
	for k := range m.Classes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
