package migrations

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UpdateForums_20171214_192507 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UpdateForums_20171214_192507{}
	m.Created = "20171214_192507"

	migration.Register("UpdateForums_20171214_192507", m)
}

// Run the migrations
func (m *UpdateForums_20171214_192507) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("alter table t_forums add is_deleted tinyint(2)")

}

// Reverse the migrations
func (m *UpdateForums_20171214_192507) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("alter table t_forums drop is_deleted")
}
