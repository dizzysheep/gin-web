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
