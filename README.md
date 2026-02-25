# codeprovenance - Code Provenance Tracker

[![Go](https://img.shields.io/badge/Go-1.21-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

**Track code origins, verify build integrity, and ensure supply chain security.**

Maintain complete audit trails of code changes and verify build reproducibility.

## 🚀 Features

- **Artifact Tracking**: Register and track code artifacts with full metadata
- **Provenance Chains**: Maintain complete history of all code changes
- **Build Verification**: Verify build integrity and reproducibility
- **Hash Verification**: SHA-256 hash-based artifact verification
- **Unauthorized Change Detection**: Detect modifications by unauthorized actors
- **Integrity Scoring**: Calculate code integrity scores (0-100%)

## 📦 Installation

### Build from Source

```bash
git clone https://github.com/hallucinaut/codeprovenance.git
cd codeprovenance
go build -o codeprovenance ./cmd/codeprovenance
sudo mv codeprovenance /usr/local/bin/
```

### Install via Go

```bash
go install github.com/hallucinaut/codeprovenance/cmd/codeprovenance@latest
```

## 🎯 Usage

### Track Artifact

```bash
# Track new artifact
codeprovenance track myapp

# Register artifact with metadata
# codeprovenance track --name=myapp --version=1.0.0 --source=https://...
```

### Verify Artifact

```bash
# Verify artifact integrity
codeprovenance verify build-001

# Check provenance chain
codeprovenance chain build-001
```

### Check Build

```bash
# Verify build integrity
codeprovenance check build-info.json
```

### Programmatic Usage

```go
package main

import (
    "fmt"
    "github.com/hallucinaut/codeprovenance/pkg/provenance"
    "github.com/hallucinaut/codeprovenance/pkg/verify"
)

func main() {
    // Create tracker
    tracker := provenance.NewTracker()

    // Register artifact
    artifact := tracker.RegisterArtifact("myapp", "1.0.0", "https://github.com/example/repo", "dev@example.com")
    artifact.Hash = provenance.ComputeHash("source content")

    // Add provenance record
    record := tracker.AddProvenanceRecord(
        "myapp:1.0.0",
        "build",
        "ci@example.com",
        "abc123",
        "def456",
        []provenance.Change{{Type: "add", Path: "main.go", LineCount: 100}},
    )

    // Get provenance chain
    chain := tracker.GetProvenanceChain("myapp:1.0.0")
    fmt.Printf("Provenance chain: %d records\n", len(chain))

    // Verify integrity
    verif := tracker.VerifyArtifact("myapp:1.0.0", "expected_hash")
    fmt.Printf("Valid: %v\n", verif.Valid)

    // Check unauthorized changes
    unauthorized := tracker.DetectUnauthorizedChanges("myapp:1.0.0", "authorized@example.com")
    fmt.Printf("Unauthorized changes: %d\n", len(unauthorized))

    // Calculate integrity score
    score := provenance.CalculateIntegrityScore(chain)
    fmt.Printf("Integrity Score: %.0f%%\n", score)
}
```

## 🔍 Provenance Features

### Artifact Tracking

- Unique artifact identification
- Version tracking
- Hash-based integrity
- Source URL tracking
- Author attribution
- Timestamp recording

### Change Tracking

- Add/Delete/Modify operations
- Path tracking
- Line count analysis
- Hash-based diff
- Actor attribution

### Build Verification

- Hash verification
- Commit integrity check
- Timestamp validation
- Artifact completeness
- Build reproducibility

## 📊 Security Scores

| Score | Status | Meaning |
|-------|--------|---------|
| 90-100 | Excellent | Full provenance, verified |
| 70-89 | Good | Minor gaps in tracking |
| 50-69 | Fair | Significant gaps |
| <50 | Poor | Incomplete provenance |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./pkg/provenance -run TestVerifyArtifact
```

## 📋 Example Output

```
Tracking artifact: myapp

Artifact registered: myapp v1.0.0
Hash: 5d41402abc4b2a76...
Source: https://github.com/example/repo

Provenance record created: prov-20240225150405-abc123...
Action: build
Actor: ci@example.com

Provenance chain length: 1 records
Integrity Score: 95%

Verifying artifact: myapp:1.0.0

✓ Artifact integrity verified
Provenance records: 1
```

## 🔒 Security Use Cases

- **Supply Chain Security**: Track all code origins
- **Compliance Audits**: Maintain complete audit trails
- **Incident Response**: Trace affected artifacts
- **Build Reproducibility**: Verify reproducible builds
- **Unauthorized Access**: Detect tampering

## 🛡️ Best Practices

1. **Register all artifacts** in the provenance system
2. **Track every change** with detailed records
3. **Verify before deployment** using hash checks
4. **Maintain immutable logs** of provenance data
5. **Regular audits** of provenance chains

## 🏗️ Architecture

```
codeprovenance/
├── cmd/
│   └── codeprovenance/
│       └── main.go          # CLI entry point
├── pkg/
│   ├── provenance/
│   │   ├── provenance.go    # Provenance tracking
│   │   └── provenance_test.go # Unit tests
│   └── verify/
│       ├── verify.go        # Build verification
│       └── verify_test.go   # Unit tests
└── README.md
```

## 📄 License

MIT License

## 🙏 Acknowledgments

- Supply chain security community
- Build verification research
- Software bill of materials standards

## 🔗 Resources

- [CycloneDX](https://cyclonedx.org/)
- [SPDX Specification](https://spdx.dev/)
- [Supply Chain Levels for Software Artifacts (SLSA)](https://slsa.dev/)

---

**Built with GPU by [hallucinaut](https://github.com/hallucinaut)**