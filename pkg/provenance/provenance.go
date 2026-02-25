// Package provenance provides code provenance tracking and verification.
package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Artifact represents a code artifact with provenance.
type Artifact struct {
	Name        string
	Version     string
	Hash        string
	Algorithm   string
	Author      string
	CommittedAt time.Time
	SourceURL   string
	Metadata    map[string]string
}

// ProvenanceRecord tracks the origin and history of code.
type ProvenanceRecord struct {
	ID           string
	ArtifactID   string
	Timestamp    time.Time
	Action       string
	Actor        string
	SourceCommit string
	DestCommit   string
	Changes      []Change
	Verification *Verification
}

// Change represents a code change.
type Change struct {
	Type        string // add, delete, modify
	Path        string
	BeforeHash  string
	AfterHash   string
	LineCount   int
}

// Verification represents verification of provenance.
type Verification struct {
	Valid        bool
	Signed       bool
	Signature    string
	VerifiedBy   string
	VerifiedAt   time.Time
	Errors       []string
}

// Tracker tracks code provenance.
type Tracker struct {
	artifacts      map[string]*Artifact
	provenanceLogs []ProvenanceRecord
}

// NewTracker creates a new provenance tracker.
func NewTracker() *Tracker {
	return &Tracker{
		artifacts:      make(map[string]*Artifact),
		provenanceLogs: make([]ProvenanceRecord, 0),
	}
}

// RegisterArtifact registers a new artifact.
func (t *Tracker) RegisterArtifact(name, version, sourceURL, author string) *Artifact {
	artifact := &Artifact{
		Name:        name,
		Version:     version,
		SourceURL:   sourceURL,
		Author:      author,
		CommittedAt: time.Now(),
		Algorithm:   "SHA-256",
		Metadata:    make(map[string]string),
	}
	t.artifacts[name+":"+version] = artifact
	return artifact
}

// ComputeHash computes hash of content.
func ComputeHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// AddProvenanceRecord adds a provenance record.
func (t *Tracker) AddProvenanceRecord(artifactID, action, actor, sourceCommit, destCommit string, changes []Change) *ProvenanceRecord {
	record := &ProvenanceRecord{
		ID:           generateID(),
		ArtifactID:   artifactID,
		Timestamp:    time.Now(),
		Action:       action,
		Actor:        actor,
		SourceCommit: sourceCommit,
		DestCommit:   destCommit,
		Changes:      changes,
	}
	t.provenanceLogs = append(t.provenanceLogs, *record)
	return record
}

// GetProvenanceChain retrieves the full provenance chain for an artifact.
func (t *Tracker) GetProvenanceChain(artifactID string) []ProvenanceRecord {
	var chain []ProvenanceRecord
	for _, record := range t.provenanceLogs {
		if record.ArtifactID == artifactID {
			chain = append(chain, record)
		}
	}
	return chain
}

// VerifyArtifact verifies artifact provenance.
func (t *Tracker) VerifyArtifact(artifactID, expectedHash string) *Verification {
	verif := &Verification{
		Valid:      true,
		VerifiedBy: "codeprovenance",
		VerifiedAt: time.Now(),
		Errors:     make([]string, 0),
	}

	artifact := t.artifacts[artifactID]
	if artifact == nil {
		verif.Valid = false
		verif.Errors = append(verif.Errors, "Artifact not found")
		return verif
	}

	if artifact.Hash != expectedHash {
		verif.Valid = false
		verif.Errors = append(verif.Errors, "Hash mismatch")
	}

	// Check provenance chain exists
	chain := t.GetProvenanceChain(artifactID)
	if len(chain) == 0 {
		verif.Errors = append(verif.Errors, "No provenance records found")
	}

	return verif
}

// GenerateID generates a unique ID.
func generateID() string {
	return "prov-" + time.Now().Format("20060102150405") + "-" + ComputeHash(time.Now().String())[:8]
}

// CalculateIntegrityScore calculates code integrity score.
func CalculateIntegrityScore(records []ProvenanceRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}

	score := 100.0
	for _, record := range records {
		if record.Action == "modify" && len(record.Changes) > 100 {
			score -= 5
		}
		if record.Action == "delete" {
			score -= 2
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}

// GetBuildChain retrieves the build chain for an artifact.
func (t *Tracker) GetBuildChain(artifactID string) []ProvenanceRecord {
	var chain []ProvenanceRecord
	for _, record := range t.provenanceLogs {
		if record.ArtifactID == artifactID && record.Action == "build" {
			chain = append(chain, record)
		}
	}
	return chain
}

// DetectUnauthorizedChanges detects unauthorized modifications.
func (t *Tracker) DetectUnauthorizedChanges(artifactID, authorizedActor string) []ProvenanceRecord {
	var unauthorized []ProvenanceRecord
	for _, record := range t.provenanceLogs {
		if record.ArtifactID == artifactID && record.Actor != authorizedActor {
			unauthorized = append(unauthorized, record)
		}
	}
	return unauthorized
}

// GetArtifactHistory gets full history of an artifact.
func (t *Tracker) GetArtifactHistory(artifactID string) []ProvenanceRecord {
	return t.GetProvenanceChain(artifactID)
}

// CalculateChangeImpact calculates impact of changes.
func CalculateChangeImpact(changes []Change) string {
	totalLines := 0
	adds := 0
	deletes := 0

	for _, change := range changes {
		totalLines += change.LineCount
		if change.Type == "add" {
			adds++
		} else if change.Type == "delete" {
			deletes++
		}
	}

	if totalLines > 1000 {
		return "HIGH"
	} else if totalLines > 100 {
		return "MEDIUM"
	}
	return "LOW"
}

// GetProvenanceStats gets provenance statistics.
func GetProvenanceStats(records []ProvenanceRecord) map[string]interface{} {
	stats := map[string]interface{}{
		"total_records": len(records),
		"by_action":     make(map[string]int),
		"by_actor":      make(map[string]int),
	}

	for _, record := range records {
		stats["by_action"].(map[string]int)[record.Action]++
		stats["by_actor"].(map[string]int)[record.Actor]++
	}

	return stats
}

// VerifyBuildReproducibility checks if build is reproducible.
func VerifyBuildReproducibility(artifact *Artifact, buildRecords []ProvenanceRecord) bool {
	if len(buildRecords) == 0 {
		return false
	}

	// Check if all builds use same source commit
	sourceCommits := make(map[string]bool)
	for _, record := range buildRecords {
		sourceCommits[record.SourceCommit] = true
	}

	return len(sourceCommits) == 1
}

// GenerateReport generates provenance report.
func GenerateReport(artifactID string, records []ProvenanceRecord) string {
	var report string

	report += "=== Provenance Report ===\n"
	report += "Artifact: " + artifactID + "\n\n"
	report += "Total Provenance Records: " + string(rune(len(records)+48)) + "\n\n"

	for i, record := range records {
		report += "[" + string(rune(i+49)) + "] " + record.Action + "\n"
		report += "    Actor: " + record.Actor + "\n"
		report += "    Timestamp: " + record.Timestamp.Format("2006-01-02 15:04:05") + "\n"
		report += "    Source Commit: " + record.SourceCommit[:8] + "\n"
		report += "    Changes: " + string(rune(len(record.Changes)+48)) + "\n\n"
	}

	return report
}