package sql_generators

import (
	"strings"
)

type joinInfo struct {
	buildInfo *SqlGenerator
	joinType string
	thisRelationField string
	relationCondition string
	thatRelationField string
}

func (this *joinInfo) JoinOn(thatRelationField, relationCondition, thisRelationField string) {
	this.relationCondition = relationCondition
	this.thatRelationField = thatRelationField
	this.thisRelationField = thisRelationField
}

func (this *joinInfo) RelationFieldPrefix(table string) {
	if strings.Count(this.thatRelationField, ".") == 0 {
		this.thatRelationField = AddCharForArray([]string{table, this.thatRelationField})
	}
}