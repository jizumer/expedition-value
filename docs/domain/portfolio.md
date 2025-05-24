**Name**: Portfolio  
**Description**: Manages investment positions and rebalancing decisions  
**Context**: Portfolio Management  
**Properties**:
- ID (UUID)
- CashBalance (Money Value Object)
- Holdings (Collection<CompanyInvestment>)
- RiskProfile (Enum)
- LastRebalanceDate (DateTime)

**Enforced Invariants**:
1. Total investments + cash ≤ portfolio value[3][9]
2. No duplicate company positions[5][10]
3. Rebalance triggers when score delta ≥ threshold[19]

**Corrective Policies**:
- Auto-liquidation when maintenance margin breached[10]
- Position capping at 5% of total value[20]

**Domain Events**:
- CompanyAddedToPortfolio
- PositionAdjusted
- RebalanceTriggered
- MarginCallOccurred

**Handled Commands**:
- AddCompanyPosition
- RemoveCompanyPosition
- ExecuteRebalance
- AdjustCashBalance