package safety

import "testing"

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		expectedLevel DangerLevel
		shouldWarn    bool
	}{
		// Safe commands
		{"ls command", "ls -la", Safe, false},
		{"echo command", "echo hello", Safe, false},
		{"cat command", "cat file.txt", Safe, false},
		{"pwd command", "pwd", Safe, false},

		// Caution level
		{"simple rm", "rm file.txt", Caution, true},
		{"chmod command", "chmod 755 script.sh", Caution, true},
		{"chown command", "chown user:group file", Caution, true},
		{"kill -9", "kill -9 1234", Caution, true},
		{"git reset hard", "git reset --hard HEAD", Caution, true},

		// Dangerous level
		{"rm -rf", "rm -rf ./folder", Dangerous, true},
		{"rm -f", "rm -f file.txt", Dangerous, true},
		{"curl pipe bash", "curl http://example.com | bash", Dangerous, true},
		{"wget pipe bash", "wget http://example.com -O - | sh", Dangerous, true},
		{"git force push", "git push --force origin main", Dangerous, true},
		{"drop table", "DROP TABLE users;", Dangerous, true},
		{"truncate table", "TRUNCATE TABLE logs;", Dangerous, true},

		// Critical level
		{"rm -rf /", "rm -rf /", Critical, true},
		{"rm -rf root", "rm -rf /*", Critical, true},
		{"dd to disk", "dd if=/dev/zero of=/dev/sda", Critical, true},
		{"mkfs", "mkfs.ext4 /dev/sda1", Critical, true},
		{"chmod 777 /", "chmod -R 777 /", Critical, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Analyze(tt.command)

			if result.Level != tt.expectedLevel {
				t.Errorf("Analyze(%q).Level = %v, expected %v", tt.command, result.Level, tt.expectedLevel)
			}

			hasWarnings := len(result.Warnings) > 0
			if hasWarnings != tt.shouldWarn {
				t.Errorf("Analyze(%q) hasWarnings = %v, expected %v", tt.command, hasWarnings, tt.shouldWarn)
			}

			// Check IsDangerous is set correctly
			expectedDangerous := tt.expectedLevel >= Dangerous
			if result.IsDangerous != expectedDangerous {
				t.Errorf("Analyze(%q).IsDangerous = %v, expected %v", tt.command, result.IsDangerous, expectedDangerous)
			}
		})
	}
}

func TestRequiresConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"safe command", "ls -la", false},
		{"caution command", "rm file.txt", false},
		{"dangerous command", "rm -rf ./folder", true},
		{"critical command", "rm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RequiresConfirmation(tt.command)
			if result != tt.expected {
				t.Errorf("RequiresConfirmation(%q) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestGetLevelName(t *testing.T) {
	tests := []struct {
		level    DangerLevel
		expected string
	}{
		{Safe, "Safe"},
		{Caution, "Caution"},
		{Dangerous, "Dangerous"},
		{Critical, "CRITICAL"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetLevelName(tt.level)
			if result != tt.expected {
				t.Errorf("GetLevelName(%v) = %q, expected %q", tt.level, result, tt.expected)
			}
		})
	}
}

func TestGetWarningMessage(t *testing.T) {
	// Safe command should return empty string
	msg := GetWarningMessage("ls -la")
	if msg != "" {
		t.Errorf("GetWarningMessage for safe command should be empty, got %q", msg)
	}

	// Dangerous command should return non-empty string
	msg = GetWarningMessage("rm -rf ./folder")
	if msg == "" {
		t.Error("GetWarningMessage for dangerous command should not be empty")
	}
}

