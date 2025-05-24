package company

// Sector represents the industry sector a company belongs to.
type Sector int

// Defines the available sectors for a company.
const (
	UndefinedSector Sector = iota // Default or unknown sector
	Technology
	Healthcare
	Financials
	ConsumerDiscretionary
	ConsumerStaples
	Industrials
	Energy
	Utilities
	RealEstate
	Materials
	TelecommunicationServices
)

// String returns the string representation of a Sector.
func (s Sector) String() string {
	switch s {
	case Technology:
		return "Technology"
	case Healthcare:
		return "Healthcare"
	case Financials:
		return "Financials"
	case ConsumerDiscretionary:
		return "Consumer Discretionary"
	case ConsumerStaples:
		return "Consumer Staples"
	case Industrials:
		return "Industrials"
	case Energy:
		return "Energy"
	case Utilities:
		return "Utilities"
	case RealEstate:
		return "Real Estate"
	case Materials:
		return "Materials"
	case TelecommunicationServices:
		return "Telecommunication Services"
	default:
		return "UndefinedSector"
	}
}

// ParseSector converts a string to a Sector type.
// It returns UndefinedSector if the string does not match any known sector.
func ParseSector(s string) Sector {
	switch s {
	case "Technology":
		return Technology
	case "Healthcare":
		return Healthcare
	case "Financials":
		return Financials
	case "Consumer Discretionary":
		return ConsumerDiscretionary
	case "Consumer Staples":
		return ConsumerStaples
	case "Industrials":
		return Industrials
	case "Energy":
		return Energy
	case "Utilities":
		return Utilities
	case "Real Estate":
		return RealEstate
	case "Materials":
		return Materials
	case "Telecommunication Services":
		return TelecommunicationServices
	default:
		return UndefinedSector
	}
}
