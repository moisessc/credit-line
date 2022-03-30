package model

const (
	// Approved identify the approved credit request status
	Approved CreditStatus = "APPROVED"
	// Declined identify the declined credit request status
	Declined CreditStatus = "DECLINED"
)

// CreditStatus type to specify the credit request status
type CreditStatus string

// CreditLine struct for the credit line entity
type CreditLine struct {
	foundingType        string
	cashBalance         float64
	monthlyRevenue      float64
	requestedCreditLine float64
	requestedDate       string
}

// CreditLineResponse struct that represents the credit line response
type CreditLineResponse struct {
	CreditStatus         CreditStatus `json:"creditStatus"`
	CreditLineAuthorized string       `json:"creditLineAuthorized"`
}

// NewCreditLine creates a new pointer of CreditLine struct
func NewCreditLine(foundingType, requestedDate string, cashBalance, monthlyRevenue, requestedCreditLine float64) *CreditLine {
	return &CreditLine{
		foundingType:        foundingType,
		cashBalance:         cashBalance,
		monthlyRevenue:      monthlyRevenue,
		requestedCreditLine: requestedCreditLine,
		requestedDate:       requestedDate,
	}
}

// FoundingType getter for the foundingType attribute
func (cl *CreditLine) FoundingType() string { return cl.foundingType }

// CashBalance getter for the cashBalance attribute
func (cl *CreditLine) CashBalance() float64 { return cl.cashBalance }

// MonthlyRevenue getter for the monthlyRevenue attribute
func (cl *CreditLine) MonthlyRevenue() float64 { return cl.monthlyRevenue }

// RequestedCreditLine getter for the requestedCreditLine attribute
func (cl *CreditLine) RequestedCreditLine() float64 { return cl.requestedCreditLine }

// RequestedDate getter for the requestedDate attribute
func (cl *CreditLine) RequestedDate() string { return cl.requestedDate }

// NewCreditLineResponse creates a new pointer of CreditLineResponse
func NewCreditLineResponse(creditStatus CreditStatus, creditLineAuthorized string) *CreditLineResponse {
	return &CreditLineResponse{
		CreditStatus:         creditStatus,
		CreditLineAuthorized: creditLineAuthorized,
	}
}
