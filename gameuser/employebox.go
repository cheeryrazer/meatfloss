package gameuser

import "meatfloss/common"

// EmployeeBox ...
type EmployeeBox struct {
	EmployeesInfo  []*common.EmployeesInfo // 雇员
	EmployeesToken string                  // 检测变化的token值
}

// NewEmployeeBox ...
func NewEmployeeBox() (employeesinfo *EmployeeBox) {
	employeesinfo = &EmployeeBox{}
	employeesinfo.EmployeesInfo = make([]*common.EmployeesInfo, 0)
	return
}
