package constants

const (
	// AuthorityUnknown is a zero value Authority that is used when parsing an unknown auhtority string.
	AuthorityUnknown = iota
	// AuthorityAuction is when the property is offered for sale via auction at a particular date and time.
	AuthorityAuction
	// AuthorityExclusive is when the property is "For Sale"
	AuthorityExclusive
	// AuthorityMultilist is when the property is "For Sale" and listed with multiple agents
	AuthorityMultilist
	// AuthorityConjunctional is when the property is "For Sale" and listed with multiple agents working together - sales
	// commission is shared amongst agents
	AuthorityConjunctional
	// AuthorityOpen is when more than one real estate agent may be employed to sell the property. The owner pays a commission only
	// to the agent who finds the buyer
	AuthorityOpen
	// AuthoritySale is the most common sale type - The property is offered for sale by negotiation.
	AuthoritySale
	// AuthoritySetSale is when the property is offered for sale up to the set sale date.
	AuthoritySetSale
)
