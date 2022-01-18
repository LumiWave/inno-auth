package context

type Company struct {
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`
}

func NewCompany() *Company {
	return new(Company)
}
