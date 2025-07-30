package dtos

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponseDTO struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code,omitempty"` // opsional: kode error internal
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationMeta berisi metadata untuk paginasi
type PaginationMeta struct {
	Page       int `json:"page"`        // halaman saat ini
	Limit      int `json:"limit"`       // jumlah data per halaman
	Total      int `json:"total"`       // total seluruh data
	TotalPages int `json:"total_pages"` // total halaman
}

// PaginatedSuccessResponse untuk response data yang dipaginasi
type PaginatedSuccessResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}
