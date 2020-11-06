package score

var per_finding_impact = 0.1

var severity = map[string]float64{
    "high": 1.0,
    "medium": 0.6,
    "low": 0.2,
}

var confidence = map[string]float64{
    "high": 1.0,
    "medium": 0.8,
    "low": 0.6,
}

func CalculateScore(findings []Finding_test_data) int {
  // given a slice of findings iterate through them and keep a running count of the score
  var score = 100.0
  for _, finding := range findings {
    finding_impact := (per_finding_impact * severity[finding.severity] *
                      confidence[finding.confidence])
    score = score * (1.0 - finding_impact)
  }
  // score is a float for calculations which we cast to an int to return
  return int(score)


}
