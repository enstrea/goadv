package service

import (
	"github.com/pkg/errors"
	"goadv/two/dao"
	_err "goadv/two/libs/errors"
	"goadv/two/model"
)

var Employee *EmployeeService

func init() {
	Employee = NewEmployeeService(dao.Employee)
}

func NewEmployeeService(dao *dao.EmployeeDao) *EmployeeService {
	return &EmployeeService{
		employeeDao: dao,
	}
}

type EmployeeService struct {
	employeeDao *dao.EmployeeDao
}

func (e *EmployeeService) Entry(employee *model.Employee) error {
	// 项目中其他模块抛出的错误，原样返回
	if err := e.checkName(employee.Name); err != nil {
		return err
	}

	if err := e.employeeDao.Insert(employee); err != nil {
		return err
	}

	return nil
}

// 校验name是否合法
func (e *EmployeeService) checkName(name string) error {
	if name == "" {
		// 正常业务逻辑错误使用自定义error，方便客户端展示友好的提示信息
		return errors.Wrap(_err.New(_err.CodeEmployeeNameInvalid, "用户名非法"), "***[EmployeeService] checkName***")
	}

	return nil
}

func (e *EmployeeService) GetEmployeeInfo(id int64) (*model.Employee, error) {
	return e.employeeDao.FindById(id)
}
