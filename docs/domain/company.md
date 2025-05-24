**Name**: Company  
**Description**: Tracks financial metrics and scoring  
**Context**: Investment Analysis  
**Properties**:
- Ticker (Value Object)
- FinancialMetrics (Entity)
- HistoricalScores (Collection<ScoreSnapshot>)
- CurrentScore (Value Object)
- SectorClassification (Enum)

**Enforced Invariants**:
1. Metrics must be ≤ 24hrs stale[3][10]
2. Score calculation atomic transaction[9][19]

**Corrective Policies**:
- Automated data refresh on stale metrics[10]
- Score recalc on metric update[20]

**Domain Events**:
- FinancialMetricsUpdated
- ScoreRecalculated
- SectorReclassified
