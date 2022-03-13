package meth_exec

import (
	"fmt"
	"reflect"
)

const (
	// ============================== Строки ошибок ==================================================
	S_PARAM_TO_STRUCT_ERR = "Ошибка преобразования параметров в структуру %s"
	S_STRUCT_TO_PARAM_ERR = "Ошибка преобразования структуры в параметры %s"
)

/***
     Преобразование типа структуры в описание параметров
**/
func StructToParamInfo(_pType reflect.Type) (_pInfo []TParamInfo) {
	_pInfo = make([]TParamInfo, _pType.NumField())
	for i := 0; i < _pType.NumField(); i++ {
		pFInfo := _pType.Field(i)
		_pInfo[i].Name = pFInfo.Name
		_pInfo[i].Type = pFInfo.Type
	}
	return _pInfo
}

/***
     Преобразование параметров в структуру
**/
func ParamsToStruct(_pInfo []TParamInfo, _pValues map[string]interface{}) (_pRes interface{}, _fOk bool, _sError string) {

	// Обработка паники внутри метода - чтобы сервер не упал из-за ошибки внутри обработки
	defer func() {
		if r := recover(); r != nil {
			_pRes = nil
			_fOk = false
			_sError = fmt.Sprintf(S_PARAM_TO_STRUCT_ERR, r)
		}
	}()
	// Формируем саму струкутуру
	pFields := make([]reflect.StructField, len(_pInfo))
	for iIndex, pInfo := range _pInfo {
		pFields[iIndex] = reflect.StructField{Name: pInfo.Name, Type: pInfo.Type}
	}
	pVal := reflect.New(reflect.StructOf(pFields)).Elem()
	for iIndex, pInfo := range _pInfo {
		pVal.Field(iIndex).Set(reflect.ValueOf(_pValues[pInfo.Name]))
	}
	return pVal.Interface(), true, ""
}

/***
     Преобразование структуры с данными в набор параметров
**/
func StructToParams(_pParams interface{}) (_pRes map[string]interface{}, _fOk bool, _sError string) {
	// Обработка паники внутри метода - чтобы сервер не упал из-за ошибки внутри обработки
	defer func() {
		if r := recover(); r != nil {
			_pRes = nil
			_fOk = false
			_sError = fmt.Sprintf(S_STRUCT_TO_PARAM_ERR, r)
		}
	}()

	if _pParams != nil {
		pVal := reflect.ValueOf(_pParams)
		pType := reflect.TypeOf(_pParams)
		_pRes = make(map[string]interface{}, pVal.NumField())
		for i := 0; i < pVal.NumField(); i++ {
			_pRes[pType.Field((i)).Name] = pVal.Field(i).Interface()
		}
	} else {
		_pRes = make(map[string]interface{}, 0)
	}
	return _pRes, true, ""
}
