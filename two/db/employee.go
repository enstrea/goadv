package db

import (
	"database/sql"
	"errors"
	"goadv/two/model"
)

// 模拟数据库

var datas map[int64]*model.Employee

func init() {
	datas = make(map[int64]*model.Employee)

	var names = []string{"a", "b", "c", "d", "e"}
	for i, name := range names {
		datas[int64(i+1)] = newEmployee(int64(i+1), name)
	}
}

func newEmployee(id int64, name string) *model.Employee {
	return &model.Employee{
		Id:   id,
		Name: name,
	}
}

func Insert(employee *model.Employee) error {
	if _, ok := datas[employee.Id]; ok {
		return errors.New("data already exist")
	}

	datas[employee.Id] = employee
	return nil
}

func FindById(id int64) (*model.Employee, error) {
	if employee, ok := datas[id]; ok {
		return employee, nil
	}

	return nil, sql.ErrNoRows
}
