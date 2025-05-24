* Name: Company
* Description: Tracks financial metrics and calculates value investment scores
* Context: Investment Analysis
* Properties:
  - Ticker (string)
  - FinancialMetrics (struct)
  - CurrentScore (float64)
  - Sector (enum)
  - UpdatedAt (time.Time)
* Enforced Invariants:
  1. Metrics age ≤ 24h
  2. Score ∈ [0,100]
* Corrective Policies:
  - Refresh stale metrics automatically
  - Recalculate score on metric update
* Domain Events:
  - ScoreRecalculated
* Ways to access:
  - FindByTicker
  - SearchByScoreRange