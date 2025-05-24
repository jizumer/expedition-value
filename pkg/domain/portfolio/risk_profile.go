package portfolio

// RiskProfile represents the investor's tolerance for risk.
type RiskProfile int

// Defines the available risk profiles.
const (
	UndefinedProfile RiskProfile = iota // Default or unknown profile
	Conservative
	Moderate
	Aggressive
)

// String returns the string representation of a RiskProfile.
func (rp RiskProfile) String() string {
	switch rp {
	case Conservative:
		return "Conservative"
	case Moderate:
		return "Moderate"
	case Aggressive:
		return "Aggressive"
	default:
		return "UndefinedProfile"
	}
}

// ParseRiskProfile converts a string to a RiskProfile type.
// It returns UndefinedProfile if the string does not match any known profile.
func ParseRiskProfile(s string) RiskProfile {
	switch s {
	case "Conservative":
		return Conservative
	case "Moderate":
		return Moderate
	case "Aggressive":
		return Aggressive
	default:
		return UndefinedProfile
	}
}
