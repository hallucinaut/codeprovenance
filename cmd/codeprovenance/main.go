package main

import (
	"fmt"
	"os"

	"github.com/hallucinaut/codeprovenance/pkg/provenance"
	"github.com/hallucinaut/codeprovenance/pkg/verify"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "track":
		if len(os.Args) < 3 {
			fmt.Println("Error: artifact name required")
			printUsage()
			return
		}
		trackArtifact(os.Args[2])
	case "verify":
		if len(os.Args) < 3 {
			fmt.Println("Error: artifact ID required")
			printUsage()
			return
		}
		verifyArtifact(os.Args[2])
	case "check":
		if len(os.Args) < 3 {
			fmt.Println("Error: build info required")
			printUsage()
			return
		}
		checkBuild(os.Args[2])
	case "chain":
		if len(os.Args) < 3 {
			fmt.Println("Error: artifact ID required")
			printUsage()
			return
		}
		showChain(os.Args[2])
	case "version":
		fmt.Printf("codeprovenance version %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
	}
}

func printUsage() {
	fmt.Printf(`codeprovenance - Code Provenance Tracker

Usage:
  codeprovenance <command> [options]

Commands:
  track <name>    Track artifact provenance
  verify <id>     Verify artifact integrity
  check <file>    Check build integrity
  chain <id>      Show artifact provenance chain
  version         Show version information
  help            Show this help message

Examples:
  codeprovenance track myapp
  codeprovenance verify build-001
  codeprovenance check build-info.json
`, "codeprovenance")
}

func trackArtifact(name string) {
	tracker := provenance.NewTracker()

	fmt.Printf("Tracking artifact: %s\n", name)
	fmt.Println()

	// Register artifact
	artifact := tracker.RegisterArtifact(name, "1.0.0", "https://github.com/example/repo", "developer@example.com")
	artifact.Hash = provenance.ComputeHash("sample content")

	fmt.Printf("Artifact registered: %s v%s\n", artifact.Name, artifact.Version)
	fmt.Printf("Hash: %s\n", artifact.Hash[:16]+"...")
	fmt.Printf("Source: %s\n", artifact.SourceURL)
	fmt.Println()

	// Add provenance record
	record := tracker.AddProvenanceRecord(
		name+":1.0.0",
		"build",
		"ci@example.com",
		"abc123",
		"def456",
		[]provenance.Change{
			{Type: "add", Path: "main.go", LineCount: 100},
			{Type: "modify", Path: "config.yaml", LineCount: 5},
		},
	)

	fmt.Printf("Provenance record created: %s\n", record.ID[:20]+"...")
	fmt.Printf("Action: %s\n", record.Action)
	fmt.Printf("Actor: %s\n", record.Actor)
	fmt.Println()

	// Get provenance chain
	chain := tracker.GetProvenanceChain(name + ":1.0.0")
	fmt.Printf("Provenance chain length: %d records\n", len(chain))

	// Calculate integrity score
	score := provenance.CalculateIntegrityScore(chain)
	fmt.Printf("Integrity Score: %.0f%%\n", score)
}

func verifyArtifact(artifactID string) {
	tracker := provenance.NewTracker()

	// Register sample artifact
	tracker.RegisterArtifact(artifactID, "1.0.0", "https://github.com/example/repo", "dev@example.com")

	fmt.Printf("Verifying artifact: %s\n", artifactID)
	fmt.Println()

	// Verify
	verif := tracker.VerifyArtifact(artifactID, "expected_hash")

	if verif.Valid {
		fmt.Println("✓ Artifact integrity verified")
	} else {
		fmt.Println("⚠ Verification failed:")
		for _, err := range verif.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
	}

	// Get chain
	chain := tracker.GetProvenanceChain(artifactID)
	fmt.Printf("Provenance records: %d\n", len(chain))
}

func checkBuild(filepath string) {
	fmt.Printf("Checking build integrity: %s\n", filepath)
	fmt.Println()

	// In production: read and parse build info file
	// For demo: show verification template
	fmt.Println("Build verification template:")
	fmt.Println("1. Verify artifact hashes")
	fmt.Println("2. Check commit integrity")
	fmt.Println("3. Verify build timestamp")
	fmt.Println("4. Check artifact completeness")
	fmt.Println()

	// Example verification
	verifier := verify.NewVerifier()
	buildInfo := &verify.BuildInfo{
		BuildID:     "build-001",
		CommitHash:  "abc123def456",
		BuildTime:   time.Now(),
		Builder:     "ci-system",
		Artifacts:   []verify.ArtifactInfo{{Name: "app", Version: "1.0.0", Hash: "hash123"}},
		Environment: map[string]string{"OS": "linux", "Arch": "amd64"},
	}

	check := verifier.VerifyBuild(buildInfo)
	fmt.Println(verify.GenerateIntegrityReport(check))
}

func showChain(artifactID string) {
	tracker := provenance.NewTracker()

	fmt.Printf("Showing provenance chain for: %s\n", artifactID)
	fmt.Println()

	// Sample chain
	chain := tracker.GetProvenanceChain(artifactID)
	if len(chain) == 0 {
		fmt.Println("No provenance records found")
		return
	}

	fmt.Printf("Provenance Chain:\n")
	for i, record := range chain {
		fmt.Printf("\n[%d] Action: %s\n", i+1, record.Action)
		fmt.Printf("    Timestamp: %s\n", record.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Actor: %s\n", record.Actor)
		fmt.Printf("    Changes: %d\n", len(record.Changes))
	}
}