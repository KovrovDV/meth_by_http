package meth_exec

import "reflect"

type IParamConverter interface {
	/**Преобразование поля из системы наружу
	    	in значение внутри метода, тип для выдачи наружу
		**/
	ConvertOut(_pInternal interface{}, _pType reflect.Type) (_pExternal interface{}, _fOk bool, _sError string)
	/**Преобразование поля с наружу для передачи внутрь системы
	   		in значение внутри метода, тип для выдачи наружу
		**/
	ConvertIn(_pInternal interface{}, _pType reflect.Type) (_pExternal interface{}, _fOk bool, _sError string)
}
