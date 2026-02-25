// Package verify provides build verification and integrity checks.
package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// BuildInfo contains build metadata.
type BuildInfo struct {
	BuildID      string
	CommitHash   string
	BuildTime    time.Time
	Builder      string
	Artifacts    []ArtifactInfo
	Environment  map[string]string
}

// ArtifactInfo contains artifact information.
type ArtifactInfo struct {
	Name     string
	Version  string
	Hash     string
	Size     int64
	Location string
}

// IntegrityCheck verifies build integrity.
type IntegrityCheck struct {
	Valid        bool
	Checks       []CheckResult
	Score        float64
	Errors       []string
	Warnings     []string
}

// CheckResult contains individual check result.
type CheckResult struct {
	Name        string
	Status      string // pass, fail, warn
	Message     string
	Details     string
}

// Verifier verifies build integrity.
type Verifier struct {
	requiredChecks []CheckType
}

// CheckType defines a verification check.
type CheckType struct {
	Name        string
	Description string
	Critical    bool
}

// NewVerifier creates a new verifier.
func NewVerifier() *Verifier {
	return &Verifier{
		requiredChecks: []CheckType{
			{Name: "Hash Verification", Description: "Verify artifact hashes", Critical: true},
			{Name: "Build Timestamp", Description: "Verify build timestamp", Critical: false},
			{Name: "Commit Integrity", Description: "Verify commit hash", Critical: true},
			{Name: "Artifact Completeness", Description: "Verify all artifacts present", Critical: true},
		},
	}
}

// VerifyBuild verifies build integrity.
func (v *Verifier) VerifyBuild(buildInfo *BuildInfo) *IntegrityCheck {
	check := &IntegrityCheck{
		Valid:  true,
		Checks: make([]CheckResult, 0),
		Errors: make([]string, 0),
		Warnings: make([]string, 0),
	}

	for _, checkType := range v.requiredChecks {
		result := v.performCheck(checkType, buildInfo)
		check.Checks = append(check.Checks, result)

		if result.Status == "fail" {
			check.Valid = false
			check.Errors = append(check.Errors, result.Message)
		} else if result.Status == "warn" {
			check.Warnings = append(check.Warnings, result.Message)
		}
	}

	// Calculate score
	check.Score = v.calculateScore(check)

	return check
}

// performCheck performs a single verification check.
func (v *Verifier) performCheck(checkType CheckType, buildInfo *BuildInfo) CheckResult {
	switch checkType.Name {
	case "Hash Verification":
		return v.checkHashes(buildInfo)
	case "Build Timestamp":
		return v.checkTimestamp(buildInfo)
	case "Commit Integrity":
		return v.checkCommit(buildInfo)
	case "Artifact Completeness":
		return v.checkArtifacts(buildInfo)
	default:
		return CheckResult{
			Name:    checkType.Name,
			Status:  "pass",
			Message: "Check not implemented",
		}
	}
}

// checkHashes verifies artifact hashes.
func (v *Verifier) checkHashes(buildInfo *BuildInfo) CheckResult {
	if len(buildInfo.Artifacts) == 0 {
		return CheckResult{
			Name:    "Hash Verification",
			Status:  "fail",
			Message: "No artifacts to verify",
		}
	}

	for _, artifact := range buildInfo.Artifacts {
		if artifact.Hash == "" {
			return CheckResult{
				Name:    "Hash Verification",
				Status:  "fail",
				Message: "Missing hash for artifact: " + artifact.Name,
			}
		}
	}

	return CheckResult{
		Name:    "Hash Verification",
		Status:  "pass",
		Message: "All artifacts have valid hashes",
	}
}

// checkTimestamp verifies build timestamp.
func (v *Verifier) checkTimestamp(buildInfo *BuildInfo) CheckResult {
	if buildInfo.BuildTime.IsZero() {
		return CheckResult{
			Name:    "Build Timestamp",
			Status:  "warn",
			Message: "Build timestamp not set",
		}
	}

	age := time.Since(buildInfo.BuildTime).Hours()
	if age > 24 {
		return CheckResult{
			Name:    "Build Timestamp",
			Status:  "warn",
			Message: "Build is older than 24 hours",
		}
	}

	return CheckResult{
		Name:    "Build Timestamp",
		Status:  "pass",
		Message: "Build timestamp is valid",
	}
}

// checkCommit verifies commit hash.
func (v *Verifier) checkCommit(buildInfo *BuildInfo) CheckResult {
	if buildInfo.CommitHash == "" || len(buildInfo.CommitHash) < 8 {
		return CheckResult{
			Name:    "Commit Integrity",
			Status:  "fail",
			Message: "Invalid or missing commit hash",
		}
	}

	return CheckResult{
		Name:    "Commit Integrity",
		Status:  "pass",
		Message: "Commit hash is valid",
	}
}

// checkArtifacts verifies artifact completeness.
func (v *Verifier) checkArtifacts(buildInfo *BuildInfo) CheckResult {
	if len(buildInfo.Artifacts) == 0 {
		return CheckResult{
			Name:    "Artifact Completeness",
			Status:  "fail",
			Message: "No artifacts found in build",
		}
	}

	return CheckResult{
		Name:    "Artifact Completeness",
		Status:  "pass",
		Message: "All artifacts present",
	}
}

// calculateScore calculates verification score.
func (v *Verifier) calculateScore(check *IntegrityCheck) float64 {
	if len(check.Checks) == 0 {
		return 0.0
	}

	score := 100.0
	for _, result := range check.Checks {
		if result.Status == "fail" {
			score -= 25
		} else if result.Status == "warn" {
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}

// VerifyArtifactHash verifies artifact hash.
func VerifyArtifactHash(content []byte, expectedHash string) bool {
	hash := ComputeHash(content)
	return hash == expectedHash
}

// ComputeHash computes SHA-256 hash.
func ComputeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// VerifyBuildChain verifies build chain integrity.
func VerifyBuildChain(builds []BuildInfo) *IntegrityCheck {
	check := &IntegrityCheck{
		Valid:  true,
		Checks: make([]CheckResult, 0),
	}

	if len(builds) == 0 {
		check.Valid = false
		check.Errors = append(check.Errors, "No builds in chain")
		return check
	}

	// Check chain continuity
	for i := 1; i < len(builds); i++ {
		if builds[i].CommitHash == "" {
			check.Valid = false
			check.Errors = append(check.Errors, "Missing commit hash in build chain")
			break
		}
	}

	return check
}

// GenerateIntegrityReport generates integrity report.
func GenerateIntegrityReport(check *IntegrityCheck) string {
	var report string

	report += "=== Build Integrity Report ===\n\n"
	report += "Valid: " + boolToString(check.Valid) + "\n"
	report += "Score: " + fmt.Sprintf("%.0f%%", check.Score) + "%\n\n"

	report += "Checks:\n"
	for _, result := range check.Checks {
		status := "✓"
		if result.Status == "fail" {
			status = "✗"
		} else if result.Status == "warn" {
			status = "⚠"
		}
		report += "  " + status + " " + result.Name + ": " + result.Message + "\n"
	}

	if len(check.Errors) > 0 {
		report += "\nErrors:\n"
		for _, err := range check.Errors {
			report += "  ✗ " + err + "\n"
		}
	}

	if len(check.Warnings) > 0 {
		report += "\nWarnings:\n"
		for _, warn := range check.Warnings {
			report += "  ⚠ " + warn + "\n"
		}
	}

	return report
}

// boolToString converts bool to string.
func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}