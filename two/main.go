package main

import (
	"fmt"
	"github.com/pkg/errors"
	_err "goadv/two/libs/errors"
	"goadv/two/model"
	"goadv/two/service"
	"log"
)

func main() {
	insertOne()
	insertOne()
	insertTwo()

	getOne(5)
	getOne(7)
}

func insertOne() {
	err := service.Employee.Entry(&model.Employee{
		Id:   6,
		Name: "HEE",
	})

	if err != nil {
		handleErr(err)
	}
}

func insertTwo() {
	err := service.Employee.Entry(&model.Employee{
		Id:   7,
		Name: "",
	})

	if err != nil {
		handleErr(err)
	}
}

func getOne(id int64) {
	employee, err := service.Employee.GetEmployeeInfo(id)
	if err != nil {
		if errors.Is(err, _err.NotFound) {
			// 做些处理，比如插入新数据，或者打印错误
			fmt.Println(fmt.Sprintf("%+v", err))
		} else {
			handleErr(err)
		}

		return
	}
	fmt.Println(employee)
}

func handleErr(err error) {
	code, err := getErrCode(err)

	if err != nil {
		log.Printf("%+v\n\n", err)
	} else {
		fmt.Println("err code: ", code)
	}
}

func getErrCode(err error) (int32, error) {
	target := &_err.AppError{}
	if errors.As(err, &target) {
		err2, ok := errors.Cause(err).(*_err.AppError)
		if ok {
			return err2.Code, nil
		}
	}
	return 0, err
}