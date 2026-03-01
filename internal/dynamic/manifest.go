package dynamic

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func DiscoverManifestPath(explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}
	candidates := []string{"prepare.yaml", "prepare.yml", "prepare.json"}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", os.ErrNotExist
}

func LoadManifest(path string) (Manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "" {
		return Manifest{APIVersion: "v1"}, nil
	}
	if strings.HasPrefix(trimmed, "{") {
		var m Manifest
		if err := json.Unmarshal([]byte(trimmed), &m); err != nil {
			return Manifest{}, fmt.Errorf("invalid JSON manifest: %w", err)
		}
		if m.APIVersion == "" {
			m.APIVersion = "v1"
		}
		if m.Profiles == nil {
			m.Profiles = map[string]Profile{}
		}
		return m, nil
	}
	return parseYAMLManifest(trimmed)
}

func parseYAMLManifest(content string) (Manifest, error) {
	m := Manifest{APIVersion: "v1", Profiles: map[string]Profile{}}
	s := bufio.NewScanner(strings.NewReader(content))
	var section string
	var inProfile string
	var profileField string
	lineNo := 0

	for s.Scan() {
		lineNo++
		raw := s.Text()
		trimmed := strings.TrimSpace(stripComment(raw))
		if trimmed == "" {
			continue
		}
		indent := leadingSpaces(raw)

		switch {
		case indent == 0 && strings.HasPrefix(trimmed, "apiVersion:"):
			m.APIVersion = strings.TrimSpace(strings.TrimPrefix(trimmed, "apiVersion:"))
			section = ""
			inProfile = ""
		case indent == 0 && strings.HasPrefix(trimmed, "profile:"):
			m.Profile = strings.TrimSpace(strings.TrimPrefix(trimmed, "profile:"))
			section = ""
			inProfile = ""
		case indent == 0 && trimmed == "tools:":
			section = "tools"
			inProfile = ""
		case indent == 0 && trimmed == "profiles:":
			section = "profiles"
			inProfile = ""
		case section == "tools" && indent == 2 && strings.HasPrefix(trimmed, "- "):
			m.Tools = append(m.Tools, strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
		case section == "profiles" && indent == 2 && strings.HasSuffix(trimmed, ":"):
			inProfile = strings.TrimSuffix(trimmed, ":")
			m.Profiles[inProfile] = m.Profiles[inProfile]
			profileField = ""
		case section == "profiles" && inProfile != "" && indent == 4 && strings.HasPrefix(trimmed, "tools:"):
			profileField = "tools"
		case section == "profiles" && inProfile != "" && indent == 4 && strings.HasPrefix(trimmed, "extends:"):
			profileField = "extends"
		case section == "profiles" && inProfile != "" && indent == 6 && strings.HasPrefix(trimmed, "- "):
			p := m.Profiles[inProfile]
			item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if profileField == "tools" {
				p.Tools = append(p.Tools, item)
			} else if profileField == "extends" {
				p.Extends = append(p.Extends, item)
			} else {
				return Manifest{}, fmt.Errorf("line %d: list item outside profile tools/extends", lineNo)
			}
			m.Profiles[inProfile] = p
		default:
			return Manifest{}, fmt.Errorf("line %d: unsupported manifest syntax", lineNo)
		}
	}
	if err := s.Err(); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func stripComment(s string) string {
	if idx := strings.Index(s, "#"); idx >= 0 {
		return s[:idx]
	}
	return s
}

func leadingSpaces(s string) int {
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	return i
}

func ValidateManifest(m Manifest, catalog map[string]ToolSpec, builtinProfiles map[string]Profile) error {
	if m.APIVersion == "" {
		return errors.New("apiVersion is required")
	}
	if m.APIVersion != "v1" {
		return fmt.Errorf("unsupported apiVersion %q", m.APIVersion)
	}
	for _, tool := range m.Tools {
		if _, ok := catalog[tool]; !ok {
			return fmt.Errorf("unknown tool %q", tool)
		}
	}

	profiles := mergeProfiles(builtinProfiles, m.Profiles)
	for profileName, profile := range profiles {
		for _, base := range profile.Extends {
			if _, ok := profiles[base]; !ok {
				return fmt.Errorf("profile %q extends unknown profile %q", profileName, base)
			}
		}
		for _, tool := range profile.Tools {
			if _, ok := catalog[tool]; !ok {
				return fmt.Errorf("profile %q references unknown tool %q", profileName, tool)
			}
		}
	}

	for name := range profiles {
		if _, err := resolveProfile(name, profiles, nil, nil); err != nil {
			return err
		}
	}
	return nil
}

func ResolveTools(m Manifest, selectedProfile string, catalog map[string]ToolSpec, builtinProfiles map[string]Profile) ([]string, error) {
	profiles := mergeProfiles(builtinProfiles, m.Profiles)
	tools := append([]string{}, m.Tools...)

	if selectedProfile == "" {
		selectedProfile = m.Profile
	}
	if selectedProfile == "" {
		selectedProfile = "fullstack"
	}
	if selectedProfile != "" {
		resolved, err := resolveProfile(selectedProfile, profiles, nil, nil)
		if err != nil {
			return nil, err
		}
		tools = append(tools, resolved...)
	}

	tools = unique(tools)
	for _, tool := range tools {
		if _, ok := catalog[tool]; !ok {
			return nil, fmt.Errorf("unknown tool %q", tool)
		}
	}
	return tools, nil
}

func resolveProfile(name string, profiles map[string]Profile, visiting map[string]bool, visited map[string]bool) ([]string, error) {
	if visiting == nil {
		visiting = map[string]bool{}
	}
	if visited == nil {
		visited = map[string]bool{}
	}
	if visiting[name] {
		return nil, fmt.Errorf("profile inheritance cycle detected at %q", name)
	}
	if visited[name] {
		return nil, nil
	}
	profile, ok := profiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown profile %q", name)
	}

	visiting[name] = true
	resolved := []string{}
	for _, base := range profile.Extends {
		baseTools, err := resolveProfile(base, profiles, visiting, visited)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, baseTools...)
	}
	resolved = append(resolved, profile.Tools...)
	visiting[name] = false
	visited[name] = true

	return unique(resolved), nil
}

func mergeProfiles(builtin map[string]Profile, user map[string]Profile) map[string]Profile {
	all := map[string]Profile{}
	for k, v := range builtin {
		all[k] = Profile{Extends: append([]string{}, v.Extends...), Tools: append([]string{}, v.Tools...)}
	}
	for k, v := range user {
		all[k] = v
	}
	return all
}

func unique(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

func FormatExampleManifestPath() string {
	return filepath.Join(".", "prepare.yaml")
}

func SupportedProfiles(profiles map[string]Profile) []string {
	out := make([]string, 0, len(profiles))
	for k := range profiles {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
