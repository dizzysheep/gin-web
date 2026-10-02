package common

import (
	"strings"
	"testing"
)

func TestEqCond(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "username", Value: "admin"},
	}, nil)

	if !strings.Contains(sql, "username") || !strings.Contains(sql, "= ?") {
		t.Fatalf("expect `username = ?` in sql, got: %s", sql)
	}
	if !containsVar(vars, "admin") {
		t.Fatalf("expect admin in vars, got: %v", vars)
	}
}

func TestEqCondNilPtr(t *testing.T) {
	var state *int8
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "state", Value: state},
	}, nil)

	if strings.Contains(sql, "WHERE") {
		t.Fatalf("nil ptr should not add WHERE, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestEqCondPtr(t *testing.T) {
	state := int8(1)
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "state", Value: &state},
	}, nil)

	if !strings.Contains(sql, "state") || !strings.Contains(sql, "= ?") {
		t.Fatalf("expect `state = ?` in sql, got: %s", sql)
	}
	if !containsVar(vars, int8(1)) {
		t.Fatalf("expect 1 in vars, got: %v", vars)
	}
}

func TestEqCondZeroValue(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "is_draft", Value: 0},
	}, nil)

	if !strings.Contains(sql, "is_draft") {
		t.Fatalf("zero value should still filter, got: %s", sql)
	}
	if !containsVar(vars, 0) {
		t.Fatalf("expect 0 in vars, got: %v", vars)
	}
}

func TestGteCondNilPtr(t *testing.T) {
	var publishedFrom *int64
	sql, vars := buildSQL(t, GormConditions{
		&GteCond{Field: "published_on", Value: publishedFrom},
	}, nil)

	if strings.Contains(sql, "WHERE") {
		t.Fatalf("nil ptr should not add WHERE, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestLtCondNilPtr(t *testing.T) {
	var publishedTo *int64
	sql, vars := buildSQL(t, GormConditions{
		&LtCond{Field: "published_on", Value: publishedTo},
	}, nil)

	if strings.Contains(sql, "WHERE") {
		t.Fatalf("nil ptr should not add WHERE, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestGteLtCondPtr(t *testing.T) {
	from, to := int64(100), int64(200)
	sql, vars := buildSQL(t, GormConditions{
		&GteCond{Field: "published_on", Value: &from},
		&LtCond{Field: "published_on", Value: &to},
	}, nil)

	if strings.Count(sql, "published_on") != 2 {
		t.Fatalf("expect both publication bounds, got: %s", sql)
	}
	if len(vars) != 2 || !containsVar(vars, from) || !containsVar(vars, to) {
		t.Fatalf("expect range vars %d and %d, got: %v", from, to, vars)
	}
}

func TestLikeCond(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&LikeCond{Field: "name", Value: "go"},
	}, nil)

	if !strings.Contains(sql, "LIKE ?") {
		t.Fatalf("expect `name LIKE ?` in sql, got: %s", sql)
	}
	if !containsVar(vars, "%go%") {
		t.Fatalf("expect %%go%% in vars, got: %v", vars)
	}
}

func TestLikeCondEmptyValue(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&LikeCond{Field: "name", Value: ""},
	}, nil)

	if strings.Contains(sql, "LIKE") {
		t.Fatalf("empty value should not add LIKE cond, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestOrConditions(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&OrConditions{GormCond: []GormCond{
			&EqCond{Field: "username", Value: "admin"},
			&EqCond{Field: "state", Value: 1},
		}},
	}, nil)

	if !strings.Contains(sql, "OR") || !strings.Contains(sql, "username") || !strings.Contains(sql, "state") {
		t.Fatalf("expect OR cond in sql, got: %s", sql)
	}
	if !containsVar(vars, "admin") || !containsVar(vars, 1) {
		t.Fatalf("expect admin/1 in vars, got: %v", vars)
	}
}

func TestOrConditionsEmpty(t *testing.T) {
	sql, _ := buildSQL(t, GormConditions{
		&OrConditions{GormCond: []GormCond{}},
	}, nil)

	if strings.Contains(sql, "WHERE") {
		t.Fatalf("empty OR conds should not add WHERE, got: %s", sql)
	}
}

func TestOrConditionsNoOpChildren(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&OrConditions{GormCond: []GormCond{
			&LikeCond{Field: "title", Value: ""},
			&LikeCond{Field: "desc", Value: ""},
		}},
	}, nil)

	if strings.Contains(sql, "WHERE") {
		t.Fatalf("all no-op children should not add WHERE, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestOrConditionsSkipsNoOpChildren(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&OrConditions{GormCond: []GormCond{
			&LikeCond{Field: "title", Value: ""},
			&LikeCond{Field: "desc", Value: "go"},
		}},
	}, nil)

	if !strings.Contains(sql, "LIKE") || !containsVar(vars, "%go%") {
		t.Fatalf("expect only non-no-op child in sql, got: %s, vars: %v", sql, vars)
	}
}

func TestOrConditionsDoesNotRepeatParentConditions(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "state", Value: int8(1)},
		&OrConditions{GormCond: []GormCond{
			&LikeCond{Field: "title", Value: "go"},
			&LikeCond{Field: "desc", Value: "go"},
		}},
	}, nil)

	if strings.Count(sql, "state") != 1 {
		t.Fatalf("parent condition should appear once, got: %s", sql)
	}
	if strings.Count(sql, "LIKE") != 2 {
		t.Fatalf("expect two OR branches, got: %s", sql)
	}
	if len(vars) != 3 {
		t.Fatalf("expect one parent var and two OR vars, got: %v", vars)
	}
}

func TestOrConditionsEmptyChildrenPreserveParentConditions(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&EqCond{Field: "state", Value: int8(1)},
		&OrConditions{GormCond: []GormCond{
			&LikeCond{Field: "title", Value: ""},
			&LikeCond{Field: "desc", Value: ""},
		}},
	}, nil)

	if strings.Count(sql, "state") != 1 || strings.Contains(sql, " OR ") {
		t.Fatalf("empty OR children should preserve only parent conditions, got: %s", sql)
	}
	if len(vars) != 1 || !containsVar(vars, int8(1)) {
		t.Fatalf("expect only parent var, got: %v", vars)
	}
}

func TestDefaultArticleFiltersSkipEmptyPublicationRangeAndKeyword(t *testing.T) {
	var publishedFrom, publishedTo *int64
	sql, vars := buildSQL(t, NotDeleted(GormConditions{
		&GteCond{Field: "published_on", Value: publishedFrom},
		&LtCond{Field: "published_on", Value: publishedTo},
		&OrConditions{GormCond: []GormCond{
			&LikeCond{Field: "title", Value: ""},
			&LikeCond{Field: "desc", Value: ""},
			&LikeCond{Field: "content_md", Value: ""},
		}},
	}), nil)

	if strings.Contains(sql, "NULL") || strings.Contains(sql, " OR ") {
		t.Fatalf("empty optional filters should not add SQL conditions, got: %s", sql)
	}
	if strings.Count(sql, "deleted_on") != 1 {
		t.Fatalf("expect only the soft-delete filter, got: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("expect no vars, got: %v", vars)
	}
}

func TestLikeCondEscape(t *testing.T) {
	sql, vars := buildSQL(t, GormConditions{
		&LikeCond{Field: "name", Value: `100%_`},
	}, nil)

	if !strings.Contains(sql, "LIKE ?") {
		t.Fatalf("expect LIKE in sql, got: %s", sql)
	}
	if !containsVar(vars, `%100\%\_%`) {
		t.Fatalf("expect escaped var, got: %v", vars)
	}
}
