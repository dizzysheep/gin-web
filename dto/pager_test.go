package dto

import "testing"

func TestPagerReqToDTO(t *testing.T) {
	cases := []struct {
		pageNo     int
		pageSize   int
		wantOffset int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{5, 20, 80},
		{0, 10, 0},  // pageNo<=0 时回退为第1页
		{-3, 10, 0}, // 负数页码回退为第1页
	}

	for _, c := range cases {
		p := PagerReqToDTO(c.pageNo, c.pageSize)
		if p.Offset != c.wantOffset {
			t.Errorf("PagerReqToDTO(%d, %d) offset = %d, want %d",
				c.pageNo, c.pageSize, p.Offset, c.wantOffset)
		}
	}
}
