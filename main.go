package main

import (
	"fmt"
	"net/http"
	"reflect"
)

type TSrv struct {
	arMethods  []TMethInfo
	pProc      interface{}
	pProcTypes map[string]reflect.Type
}

type TParamInfo struct {
	Name string
	Type reflect.Type
}

type TMethInfo struct {
	pMeth    reflect.Value
	sName    string
	arInput  []TParamInfo
	arOutput []TParamInfo
}

type TInfoProc struct {
}

func (pProc *TInfoProc) GetName(_pReq struct{ _iIndex int }) {

}

func (pProc TInfoProc) GetNameWOPtr(_pReq struct {
	_iIndex int
	_sName  *string
}) (_pResp struct {
	_fOk    bool
	_sError string
}) {
	return struct {
		_fOk    bool
		_sError string
	}{true, ""}
}

func (pSrv TSrv) ProcRequest(page http.ResponseWriter, r *http.Request) {

}

func (pSrv *TSrv) ConvertStructToParams(_pType reflect.Type, _pInfo *[]TParamInfo) {
	*_pInfo = make([]TParamInfo, _pType.NumField())

}

func (pSrv *TSrv) Init() (_fOk bool, _sError string) {
	_sError = ""
	_fOk = true
	// 0) Определяем тип обработчика
	pProcVal := reflect.ValueOf(pSrv.pProc)
	pProcType := reflect.TypeOf(pSrv.pProc)

	if pProcType.Kind() != reflect.Ptr {
		return false, "Тип должен быть по ссылке"
	}
	// 1)  Определяем список методов
	pSrv.arMethods = make([]TMethInfo, pProcType.NumMethod())
	for i := 0; i < len(pSrv.arMethods); i++ {
		pMethElem := pProcType.Method(i)
		pMethValue := pProcVal.Method(i)
		pMessInfo := TMethInfo{pMeth: pMethValue, sName: pMethElem.Name}
		pSrv.arMethods[i] = pMessInfo
		// 1.1) Заполняем параметры (имена параметров не доступны - поэтому работаем через анонимные структуры)
		pMessInfo.arInput = make([]TParamInfo, pMethElem.Type.NumIn())
		pMessInfo.arOutput = make([]TParamInfo, pMethElem.Type.NumOut())
		// Первый параметр - сама структура и по одному - сообщение на вход и на выход
		if len(pMessInfo.arInput) > 2 || len(pMessInfo.arOutput) > 1 {
			return false, "Входной параметр должен быть структурой выходной тоже"
		}
		if len(pMessInfo.arInput) == 2 {
			pSrv.ConvertStructToParams(pMethElem.Type.In(1), &pMessInfo.arInput)
		}
		if len(pMessInfo.arOutput) == 1 {
			pSrv.ConvertStructToParams(pMethElem.Type.Out(0), &pMessInfo.arOutput)
		}
	}
}

func main() {

	pSrv := TSrv{pProc: new(TInfoProc)}

	pSrv.Init()
	fmt.Printf("Список методов %d", len(pSrv.arMethods))
	fmt.Println()

	//fmt.Println(pSrv.arMethods[0].sName)

	for _, pMethInfo := range pSrv.arMethods {
		fmt.Println(pMethInfo.sName)
	}

	/*
		http.HandleFunc("/", pSrv.ProcRequest)
		http.ListenAndServe(":8080", nil)
	*/

}
