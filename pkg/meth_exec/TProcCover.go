package meth_exec

import (
	"fmt"
	"reflect"
)

type TProcCover struct {
	arMethods        []TMethInfo
	pProcHander      interface{}
	pParamsConverter IParamConverter
	//	pProcTypes       map[string]reflect.Type
}

const (

	// ================== Строки ошибок =========================
	S_MAIN_TYPE_ONLY_PTR_ERR = "Тип должен быть по ссылке"
	S_PARAM_ISNT_STRUCT_ERR  = "Входной параметр должен быть структурой, выходной тоже"
	S_INTERNAL_METH_ERR      = "Ошибка внутри метода  %s - %s"
	S_NO_SUCH_METHOD         = "Нет такого метода, %s"
	S_METH_PARAM_FORM_ERR    = "Ошибка при формировании параметра %s метода %s - %s"
)

/**Конструктор**/
func NewProcCover(_pProc interface{}, _pConverter IParamConverter) *TProcCover {
	pRes := new(TProcCover)
	pRes.pProcHander = _pProc
	pRes.pParamsConverter = _pConverter
	return pRes
}

/**Конструктор**/
func NewProcCoverSmp(_pProc interface{}) *TProcCover {
	return NewProcCover(_pProc, NewStrParamConverter())
}

func (pCover *TProcCover) Init() (_fOk bool, _sError string) {
	_sError = ""
	_fOk = true
	// 0) Определяем тип обработчика
	pProcVal := reflect.ValueOf(pCover.pProcHander)
	pProcType := reflect.TypeOf(pCover.pProcHander)

	if pProcType.Kind() != reflect.Ptr {
		return false, S_MAIN_TYPE_ONLY_PTR_ERR
	}
	// 1)  Определяем список методов
	pCover.arMethods = make([]TMethInfo, pProcType.NumMethod())
	for i := 0; i < len(pCover.arMethods); i++ {
		pMethElem := pProcType.Method(i)
		pMethValue := pProcVal.Method(i)
		pMessInfo := TMethInfo{pMeth: pMethValue, sName: pMethElem.Name}
		// 1.1) Заполняем параметры (имена параметров не доступны - поэтому работаем через анонимные структуры)
		pMessInfo.arInput = make([]TParamInfo, pMethElem.Type.NumIn())
		pMessInfo.arOutput = make([]TParamInfo, pMethElem.Type.NumOut())
		// Первый параметр - сама структура и по одному - сообщение на вход и на выход
		switch {
		case len(pMessInfo.arInput) > 2 || len(pMessInfo.arOutput) > 1:
			return false, S_PARAM_ISNT_STRUCT_ERR
		case len(pMessInfo.arInput) == 2:
			pMessInfo.arInput = StructToParamInfo(pMethElem.Type.In(1))
		case len(pMessInfo.arInput) == 1:
			pMessInfo.arInput = StructToParamInfo(pMethElem.Type.In(0))
		}
		if len(pMessInfo.arOutput) == 1 {
			pMessInfo.arOutput = StructToParamInfo(pMethElem.Type.Out(0))
		} else {
			pMessInfo.arOutput = make([]TParamInfo, 0)
		}

		pCover.arMethods[i] = pMessInfo
	}
	return true, ""
}

/*
   Выполнение метода по имени и входящим параметрам
   in имя метода,  струкура с входящими параметрами
   out выходные параметры, флаг успеха, строка  с ошибкой
**/
func (pCover *TProcCover) exec(_pInfo TMethInfo, _pParams interface{}) (_pOutParam interface{}, _fOk bool, _sError string) {

	// Обработка паники внутри метода - чтобы сервер не упал из-за ошибки внутри обработки
	defer func() {
		if r := recover(); r != nil {
			_pOutParam = nil
			_fOk = false
			_sError = fmt.Sprintf(S_INTERNAL_METH_ERR, _pInfo.sName, r)
		}
	}()
	// Ищем нужный метод и запускаем
	arOutParams := _pInfo.pMeth.Call([]reflect.Value{reflect.ValueOf(_pParams)})
	// Не ответа
	if len(arOutParams) == 0 {
		return nil, true, ""
	}
	// Проверяем ответ внутри
	pVal := arOutParams[0].Interface()
	pResp, pRespVal := reflect.TypeOf(pVal), reflect.Indirect(reflect.ValueOf(pVal))
	_, fOk := pResp.FieldByName(S_OK_FLAG_PARAM)
	_, fError := pResp.FieldByName(S_ERROR_PARAM)
	if fOk && fError {
		return pVal, pRespVal.FieldByName(S_OK_FLAG_PARAM).Bool(), pRespVal.FieldByName(S_ERROR_PARAM).String()
	} else {
		return pVal, true, ""
	}
}

/*
   Выполнение метода по имени и входящим параметрам
   in имя метода,  струкура с входящими параметрами
   out выходные параметры, флаг успеха, строка  с ошибкой
**/
func (pCover *TProcCover) Exec(_sMethName string, _pParams interface{}) (_pOutParam interface{}, _fOk bool, _sError string) {
	// Ищем нужный метод и запускаем
	for _, pMethInfo := range pCover.arMethods {
		// Case sensitive или ограничиваем в Init
		if pMethInfo.sName == _sMethName {
			return pCover.exec(pMethInfo, _pParams)
		}
	}
	return nil, false, fmt.Sprintf(S_NO_SUCH_METHOD, _sMethName)
}

/*
   Выполнение метода по имени и входящим параметрам c преобразованием
   in имя метода,  струкура с входящими параметрами
   out выходные параметры, флаг успеха, строка  с ошибкой
**/
func (pCover *TProcCover) ExecOut(_sMethName string, _pParams map[string]interface{}) (_pOutParam map[string]interface{}, _fOk bool, _sError string) {
	// Ищем нужный метод и запускаем
	for _, pMethInfo := range pCover.arMethods {
		// Case sensitive или ограничиваем в Init

		if pMethInfo.sName == _sMethName {
			// Формирум список параметров
			pParams := make(map[string]interface{})
			for _, pParamInfo := range pMethInfo.arInput {
				pVal, found := _pParams[pParamInfo.Name]
				// Значение по умодчанию
				if !found {
					pParams[pParamInfo.Name] = reflect.Zero(pParamInfo.Type).Interface()
				} else {
					pRes, fOk, sErr := pCover.pParamsConverter.ConvertIn(pVal, pParamInfo.Type)
					if !fOk {
						return nil, fOk, fmt.Sprintf(S_METH_PARAM_FORM_ERR, pParamInfo.Name, _sMethName, sErr)
					}
					pParams[pParamInfo.Name] = pRes
				}
			}
			pParamsVal, fOk, sError := ParamsToStruct(pMethInfo.arInput, pParams)
			if !fOk {
				return nil, fOk, sError
			}
			pOut, fOk, sError := pCover.exec(pMethInfo, pParamsVal)
			if !fOk {
				return nil, fOk, sError
			}
			// Преобразуем ответ в список параметров
			pInternalOutParams, fOk, sError := StructToParams(pOut)
			// Преобразуем для выдачи наружу
			pOutParams := make(map[string]interface{})
			for _, pParamInfo := range pMethInfo.arOutput {
				pValue, fExists := pInternalOutParams[pParamInfo.Name]
				if !fExists {
					continue
				}
				pResp, fOk, sErr := pCover.pParamsConverter.ConvertOut(pValue, pCover.pParamsConverter.InTypeToOut(pParamInfo.Type))
				if !fOk {
					return nil, fOk, fmt.Sprintf(S_METH_PARAM_FORM_ERR, pParamInfo.Name, _sMethName, sErr)
				}
				pOutParams[pParamInfo.Name] = pResp
			}
			return pOutParams, true, ""
		}

	}
	return nil, false, fmt.Sprintf(S_NO_SUCH_METHOD, _sMethName)
}
