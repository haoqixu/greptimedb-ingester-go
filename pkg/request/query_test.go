package request

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func mockRowsToSqlRows(mockRows *sqlmock.Rows) *sql.Rows {
	db, mock, _ := sqlmock.New()
	mock.ExpectQuery("select").WillReturnRows(mockRows)
	rows, _ := db.Query("select")
	return rows
}

type Person struct {
	Name string
	Age  int
}

func TestFillStructSliceFromRows(t *testing.T) {
	expected := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}

	// Set up a mock rows object
	rows := sqlmock.NewRows([]string{"name", "age"}).
		AddRow("Alice", 25).
		AddRow("Bob", 30).
		AddRow("Charlie", 35)

	// Call the function and check the result
	slice := []Person{}
	err := fillStructSlice(&slice, mockRowsToSqlRows(rows))
	assert.Nil(t, err)
	assert.Equal(t, slice, expected)
}

func TestFillStructSliceWithInvalidRowData(t *testing.T) {
	rows := sqlmock.NewRows([]string{"name", "age"}).
		AddRow("test", 123.345)
	slice := []Person{}
	err := fillStructSlice(&slice, mockRowsToSqlRows(rows))
	// fmt.Printf("rows: %+v", slice)
	assert.NotNil(t, err)
}

func TestFillStructSliceWithIncorrectNumberOfColumns(t *testing.T) {
	// TODO(vinland-avalon): a better way to mock, since sqlmock need column number to match
	// rows := sqlmock.NewRows([]string{"id", "name"}).
	// 	AddRow(1).
	// 	AddRow(2, "test2")
	// var slice []struct {
	// 	Id   int
	// 	Name string
	// }
	// err := fillStructSlice(&slice, mockRowsToSqlRows(rows))
}

func TestIsStructSliceSettableWithNilSlicePointer(t *testing.T) {
	err := isStructSliceSettable(nil)
	assert.NotNil(t, err)
	assert.Equal(t, "dest must be a pointer to a slice", err.Error())
}

func TestIsStructSliceSettableWithNonPointerSlice(t *testing.T) {
	slice := make([]Person, 0)
	err := isStructSliceSettable(slice)
	assert.NotNil(t, err)
	assert.Equal(t, "dest must be a pointer to a slice", err.Error())
}

func TestIsStructSliceSettableWithFieldCanNotSet(t *testing.T) {
	type NonSettableStruct struct {
		ptr int
	}

	slice := []NonSettableStruct{}
	err := isStructSliceSettable(&slice)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not settable")
}
