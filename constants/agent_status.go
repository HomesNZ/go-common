package constants

const (
	// AgentStatusListingAndSale represents an agent who was both the listing and selling
	// agent for a given sale
	AgentStatusListingAndSale = iota
	// AgentStatusListing represents an agent who was/is the listing agent for a given sale
	AgentStatusListing
	// AgentStatusSale represents an agent who was/is the selling agent for a given sale
	AgentStatusSale
)
