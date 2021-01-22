package query_test

import (
	"testing"

	"github.com/frozentech/query"
	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	ID        string `db:"bin_pk" json:"-" `
	Bin       string `db:"bin" json:"bin"`
	Brand     string `db:"brand" json:"brand"`
	Scheme    string `db:"scheme" json:"scheme"`
	CreatedAt string `db:"created_at" json:"-"`
	Ignore    string `db:"-" json:"-"`
}

func Test_Contain(t *testing.T) {
	assert := assert.New(t)
	array := []string{"A", "B", "C"}

	found := query.Contains(array, "A")
	assert.True(found)

	found = query.Contains(array, "Z")
	assert.False(found)
}

func Test_UpdateBuilder(t *testing.T) {
	object := &TestObject{}
	assert := assert.New(t)
	sql, _ := query.UpdateBuilder(object, "TestTable")
	assert.Equal("UPDATE `TestTable` SET  bin_pk = ? ,  bin = ? ,  brand = ? ,  scheme = ? ,  created_at = ?  ", sql)
}

func Test_InsertBuilder(t *testing.T) {
	object := &TestObject{}
	assert := assert.New(t)
	sql, _ := query.InsertBuilder(object, "TestTable", true)
	assert.Equal("INSERT INTO `TestTable` ( `TestTable`.`bin_pk`, `TestTable`.`bin`, `TestTable`.`brand`, `TestTable`.`scheme`, `TestTable`.`created_at` ) VALUES ( ?, ?, ?, ?, ? )", sql)
}

func Test_InsertBuilder_Ignore_ID(t *testing.T) {
	object := &TestObject{}
	assert := assert.New(t)
	sql, _ := query.InsertBuilder(object, "TestTable", false, "bin_pk")
	assert.Equal("INSERT INTO `TestTable` ( `bin`, `brand`, `scheme`, `created_at` ) VALUES ( ?, ?, ?, ? )", sql)
}

func Test_SelectBuilder(t *testing.T) {
	object := &TestObject{}
	assert := assert.New(t)
	sql := query.SelectBuilder(object, "TestTable", false)
	assert.Equal("SELECT `bin_pk`, `bin`, `brand`, `scheme`, `created_at` FROM `TestTable`", sql)
}
