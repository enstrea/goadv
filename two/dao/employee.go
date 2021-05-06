package dao

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"goadv/two/db"
	_err "goadv/two/libs/errors"
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
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(_err.NotFound, fmt.Sprintf("id: %d, err: %v", id, err))
	} else if err != nil {
		return nil, errors.Wrap(_err.Internal, fmt.Sprintf("id: %d, err: %v", id, err))
	}
	return employee, nil
}
