package meth_exec

import (
	"reflect"
)

/**Информаци о параметре (поле одного входного параметра)**/
type TParamInfo struct {
	Name string
	Type reflect.Type
}
