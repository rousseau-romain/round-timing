package sqlhelper

import "fmt"

// URLImageClassClause returns a SQL CONCAT clause for building class image URLs.
func URLImageClassClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/class/", %s)`, fieldName)
}

// URLImageSpellClause returns a SQL CONCAT clause for building spell image URLs.
func URLImageSpellClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/spell/", %s, ".svg")`, fieldName)
}
