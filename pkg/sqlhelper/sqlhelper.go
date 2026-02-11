package sqlhelper

import (
	"fmt"
	"strings"
)

// URLImageClassClause returns a SQL CONCAT clause for building class image URLs.
func URLImageClassClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/class/", %s)`, fieldName)
}

// URLImageSpellClause returns a SQL CONCAT clause for building spell image URLs.
func URLImageSpellClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/spell/", %s, ".svg")`, fieldName)
}

// InClause returns a parameterized IN clause and the corresponding args slice.
func InClause(ids []int) (string, []interface{}) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	return "(" + strings.Join(placeholders, ",") + ")", args
}
