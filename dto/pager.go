package dto

type Pager struct {
	PageNo   int   `json:"page_no"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Offset   int   `json:"-"`
}

func PagerReqToDTO(PageNo, pageSize int) *Pager {
	if PageNo <= 0 {
		PageNo = 1
	}
	return &Pager{
		PageNo:   PageNo,
		PageSize: pageSize,
		Offset:   (PageNo - 1) * pageSize,
	}
}
