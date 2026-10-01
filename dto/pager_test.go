package dto

import "testing"

func TestPagerReqToDTO(t *testing.T) {
	cases := []struct {
		pageNo     int
		pageSize   int
		wantOffset int
		wantPage   int
		wantSize   int
	}{
		{1, 10, 0, 1, 10},
		{2, 10, 10, 2, 10},
		{5, 20, 80, 5, 20},
		{0, 0, 0, 1, 10},     // 缺省页码和大小使用1/10
		{-3, 101, 0, 1, 100}, // 非法页码回退，大小限制为100
	}

	for _, c := range cases {
		p := PagerReqToDTO(c.pageNo, c.pageSize)
		if p.Offset != c.wantOffset {
			t.Errorf("PagerReqToDTO(%d, %d) offset = %d, want %d",
				c.pageNo, c.pageSize, p.Offset, c.wantOffset)
		}
		if p.Page != c.wantPage || p.PageSize != c.wantSize {
			t.Errorf("PagerReqToDTO(%d, %d) = page %d size %d, want page %d size %d",
				c.pageNo, c.pageSize, p.Page, p.PageSize, c.wantPage, c.wantSize)
		}
	}
}
