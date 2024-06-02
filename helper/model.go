package helper

import (
	"fmt"
)

func GetUrlImageClassClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/class/", %s)`, fieldName)
}

func GetUrlImageSpellClause(fieldName string) string {
	return fmt.Sprintf(`CONCAT("/public/img/spell/", %s, ".svg")`, fieldName)
}
