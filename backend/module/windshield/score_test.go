package score

import (
//	"context"
	"testing"
  "github.com/stretchr/testify/assert"
//	"github.com/lyft/clutch/backend/module/moduletest"
//	"github.com/lyft/clutch/backend/service"
)

type Finding_test_data struct {
    severity string
    confidence  string
}

func TestSingleFindingHigh(t *testing.T) {
  var findings []Finding_test_data
  findings = append(findings, Finding_test_data{"high", "high"})
  score := CalculateScore(findings)
  assert.Equal(t, score, 90, "Incorrectly calculated score")
}

func TestSingleFindingMedium(t *testing.T) {
  var findings []Finding_test_data
  findings = append(findings, Finding_test_data{"medium", "medium"})
  score := CalculateScore(findings)
  assert.Equal(t, score, 95, "Incorrectly calculated score")
}

func TestSingleFindingLow(t *testing.T) {
  var findings []Finding_test_data
  findings = append(findings, Finding_test_data{"low", "low"})
  score := CalculateScore(findings)
  assert.Equal(t, score, 98, "Incorrectly calculated score")
}

func TestMultipleFindingsHigh(t *testing.T) {
  var findings []Finding_test_data
  findings = append(findings, Finding_test_data{"high", "high"})
  findings = append(findings, Finding_test_data{"high", "high"})
  findings = append(findings, Finding_test_data{"high", "high"})
  score := CalculateScore(findings)
  assert.Equal(t, score, 72, "Incorrectly calculated score")
}

func TestMultipleFindingsMixed(t *testing.T) {
  var findings []Finding_test_data
  findings = append(findings, Finding_test_data{"high", "high"})
  findings = append(findings, Finding_test_data{"high", "high"})
  findings = append(findings, Finding_test_data{"high", "medium"})
  findings = append(findings, Finding_test_data{"medium", "high"})
  findings = append(findings, Finding_test_data{"medium", "medium"})
  findings = append(findings, Finding_test_data{"medium", "low"})
  findings = append(findings, Finding_test_data{"low", "low"})
  findings = append(findings, Finding_test_data{"low", "low"})
  score := CalculateScore(findings)
  assert.Equal(t, score, 62, "Incorrectly calculated score")
}
