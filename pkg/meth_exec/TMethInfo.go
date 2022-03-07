package meth_exec

import (
	"reflect"
)

type TMethInfo struct {
	pMeth    reflect.Value
	sName    string
	arInput  []TParamInfo
	arOutput []TParamInfo
}
