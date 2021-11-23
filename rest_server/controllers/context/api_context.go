package context

const (
	Wallet_type_metamask = "metamask"
)

// page info
type PageInfo struct {
	PageOffset int64 `query:"page_offset" validate:"required"`
	PageSize   int64 `query:"page_size" validate:"required"`
}

// page response
type PageInfoResponse struct {
	PageOffset int64 `json:"page_offset"`
	PageSize   int64 `json:"page_size"`
	TotalSize  int64 `json:"total_size"`
}
