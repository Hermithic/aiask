package update

import "testing"

func TestComparePrerelease(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{"equal versions", "beta.1", "beta.1", 0},
		{"simple numeric greater", "beta.2", "beta.1", 1},
		{"simple numeric lesser", "beta.1", "beta.2", -1},
		{"multi-digit numeric", "beta.10", "beta.2", 1},
		{"multi-digit numeric reverse", "beta.2", "beta.10", -1},
		{"alpha vs beta", "beta.1", "alpha.1", 1},
		{"rc vs beta", "rc.1", "beta.1", 1},
		{"different length same prefix", "beta.1.1", "beta.1", 1},
		{"different length same prefix reverse", "beta.1", "beta.1.1", -1},
		{"numeric vs string", "1", "beta", -1},
		{"string vs numeric", "beta", "1", 1},
		{"complex version", "beta.2.1", "beta.2.0", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := comparePrerelease(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("comparePrerelease(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareVersionCore(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{"equal versions", "1.0.0", "1.0.0", 0},
		{"major version greater", "2.0.0", "1.0.0", 1},
		{"minor version greater", "1.1.0", "1.0.0", 1},
		{"patch version greater", "1.0.1", "1.0.0", 1},
		{"major version lesser", "1.0.0", "2.0.0", -1},
		{"different lengths", "1.0", "1.0.0", 0},
		{"different lengths with diff", "1.0.1", "1.0", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersionCore(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareVersionCore(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		// Dev version
		{"dev version", "2.0.0", "dev", false},
		{"empty current", "2.0.0", "", false},

		// Core version comparisons
		{"newer major", "2.0.0", "1.0.0", true},
		{"same version", "1.0.0", "1.0.0", false},
		{"older major", "1.0.0", "2.0.0", false},
		{"newer minor", "1.1.0", "1.0.0", true},
		{"newer patch", "1.0.1", "1.0.0", true},

		// Pre-release comparisons
		{"release vs prerelease same core", "2.0.0", "2.0.0-beta.1", true},
		{"prerelease vs release same core", "2.0.0-beta.1", "2.0.0", false},
		{"newer prerelease", "2.0.0-beta.2", "2.0.0-beta.1", true},
		{"older prerelease", "2.0.0-beta.1", "2.0.0-beta.2", false},
		{"same prerelease", "2.0.0-beta.1", "2.0.0-beta.1", false},
		{"rc vs beta", "2.0.0-rc.1", "2.0.0-beta.1", true},
		{"beta.10 vs beta.2", "2.0.0-beta.10", "2.0.0-beta.2", true},

		// Different core with prerelease
		{"newer core with prerelease", "2.1.0-beta.1", "2.0.0", true},
		{"older core with prerelease", "2.0.0-beta.1", "2.1.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNewer(tt.latest, tt.current)
			if result != tt.expected {
				t.Errorf("isNewer(%q, %q) = %v, expected %v", tt.latest, tt.current, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name               string
		version            string
		expectedCore       string
		expectedPrerelease string
	}{
		{"simple version", "1.0.0", "1.0.0", ""},
		{"with prerelease", "2.0.0-beta.1", "2.0.0", "beta.1"},
		{"with rc", "1.0.0-rc.1", "1.0.0", "rc.1"},
		{"complex prerelease", "1.0.0-alpha.1.2", "1.0.0", "alpha.1.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, prerelease := parseVersion(tt.version)
			if core != tt.expectedCore {
				t.Errorf("parseVersion(%q) core = %q, expected %q", tt.version, core, tt.expectedCore)
			}
			if prerelease != tt.expectedPrerelease {
				t.Errorf("parseVersion(%q) prerelease = %q, expected %q", tt.version, prerelease, tt.expectedPrerelease)
			}
		})
	}
}

