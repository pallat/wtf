package ticket

// swagger:route POST /tickets ticket buyticket
// responses:
//

// swagger:parameters buyticket
type swaggerBuyTicketParamsWrapper struct {
	// in:header
	// required:true
	// example: 7c129eb1-c479-47bb-9c73-d263e2673011
	XRefID string `json:"X-Ref-Id"`

	// in:body
	Body Ticket
}
