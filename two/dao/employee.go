package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"goadv/two/db"
	"goadv/two/model"
)

var Employee *EmployeeDao

func init() {
	Employee = NewEmployeeDao()
}

func NewEmployeeDao() *EmployeeDao {
	return &EmployeeDao{}
}

type EmployeeDao struct {
}

// dao层属于Application，可以用Wrap包装错误，获得更为详细的错误信息

func (e *EmployeeDao) Insert(employee *model.Employee) error {
	if err := db.Insert(employee); err != nil {
		return errors.Wrap(err, fmt.Sprintf("***[EmployeeDao] insert fail, data: %+v***", employee))
	}
	return nil
}

func (e *EmployeeDao) FindById(id int64) (*model.Employee, error) {
	employee, err := db.FindById(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("***[EmployeeDao] data not found, id: %d***", id))
	}
	return employee, nil
}
