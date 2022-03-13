package meth_api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meth_by_http/pkg/meth_exec"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

const (
	/**Строка информации**/
	S_SRV_STARTED_MSG = "Запуск сервера %s\r\n"
	// ================ Строки ошибок ====================
	S_SRV_START_ERR           = "Ошибка запуска сервера %s\r\n"
	S_SRV_STOP_ERR            = "Ошибка остановки сервера %s\r\n"
	S_SRV_ALREADY_STARTED_ERR = "сервер уже запущен"
	S_PROC_ERR                = "Ошибка обработки %s"
	S_MUX_ISNT_STARTED_ERR    = "Mux еще не запущен"
)

/** машрутизатор на обработку**/
type THttpRouter struct {
	pServer  *http.Server
	pMux     *http.ServeMux
	UserStop chan bool
	sysStop  chan os.Signal
	pRoutes  map[string]meth_exec.TProcCover
}

/**Конструктор**/
func NewHttpRouter() *THttpRouter {
	pRes := new(THttpRouter)
	pRes.pRoutes = make(map[string]meth_exec.TProcCover)
	pRes.pMux = nil
	pRes.pServer = nil
	pRes.UserStop = make(chan bool, 1)
	pRes.sysStop = make(chan os.Signal, 1)
	return pRes
}

/**добавляем путь **/
func (pRouter *THttpRouter) AddRoute(_sRoute string, _pProcObj interface{}) (_fOk bool, _sError string) {
	pProcCover := meth_exec.NewProcCoverSmp(_pProcObj)
	// Инициализируем обработку
	_fOk, _sError = pProcCover.Init()
	if !_fOk {
		return _fOk, _sError
	}
	// Добавляем в список обработки
	pRouter.pRoutes[_sRoute] = *pProcCover
	return true, ""
}

/** поиск подходящего обработкичика **/
func (pRouter *THttpRouter) getProc(_sUrl string) (_pCover *meth_exec.TProcCover, _sMethName string, _fExists bool) {
	for key, val := range pRouter.pRoutes {
		if strings.HasPrefix(_sUrl, key) {
			// Сделать более тонкую обработку + поиск метода()
			sMethName := _sUrl[len(key)+1:]
			// TODO Default метод попытаться вытащить через tags+reflections
			return &val, sMethName, true
		}
	}
	return nil, "", false
}

/**обработка запроса **/
func (pRouter *THttpRouter) procRequest(w http.ResponseWriter, r *http.Request) {

	// Пустой вход
	if len(r.URL.Path) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Ищем обработкчик
	pProc, sMethName, fExists := pRouter.getProc(r.URL.Path[1:])
	if !fExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Формируем список значений для параметров
	pParams := make(map[string]interface{})
	// Строка
	/*
		// Не нужно - все попадает в форму, хотя должно только из POST
		for key, val := range r.URL.Query() {
			// Пока берем первое (можно позже сменить на срез)
			if len(val) > 0 {
				pParams[key] = val[0]
			}
		}
	*/
	// Форма
	if err := r.ParseForm(); err == nil {
		for key, val := range r.Form {
			// Пока берем первое (можно позже сменить на срез)
			if len(val) > 0 {
				pParams[key] = val[0]
			}
		}
	}
	// Весь POST
	if r.ContentLength > 0 {
		arPost, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, S_PROC_ERR, err)
			return
		}
		pParams[meth_exec.S_AR_POST_PARAM] = base64.StdEncoding.EncodeToString(arPost)
	}
	// Вызываем Орбаботку
	pRes, fOk, sError := pProc.ExecOut(sMethName, pParams)
	if !fOk {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, S_PROC_ERR, sError)
		return
	}
	// Пока - иначе это обрабоать глуюже
	delete(pRes, meth_exec.S_OK_FLAG_PARAM)
	delete(pRes, meth_exec.S_ERROR_PARAM)

	if len(pRes) == 1 {
		for _, v := range pRes {
			fmt.Fprintf(w, v.(string))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		sJSON, err := json.Marshal(pRes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, S_PROC_ERR, err)
		}
		w.Write(sJSON)
	}
	//return
}

/* Запуск прослушивания адреса и порта*/
func (pRouter *THttpRouter) ListenLocal(_sSubPath string) (_fOk bool, _sError string) {
	if pRouter.pServer != nil || pRouter.pMux != nil {
		return false, S_SRV_ALREADY_STARTED_ERR
	}
	// Инициализиреу сервер (+ TODO добаить hanler для ssl, авторизации)
	pRouter.pMux = http.NewServeMux()
	pRouter.pMux.HandleFunc(_sSubPath, pRouter.procRequest)
	return true, ""
}

/* Запуск прослушивания адреса и порта*/
func (pRouter *THttpRouter) Listen(_sAddress string, _sSubPath string) (_fOk bool, _sError string) {
	_fOk, _sError = pRouter.ListenLocal(_sSubPath)
	if !_fOk {
		return _fOk, _sError
	}
	pRouter.pServer = &http.Server{Addr: _sAddress, Handler: pRouter.pMux}

	// Запуск реального адреса
	go func() {
		fmt.Printf(S_SRV_STARTED_MSG, _sAddress)
		if err := pRouter.pServer.ListenAndServe(); err != nil {
			// handle err
			fmt.Printf(S_SRV_START_ERR, err)
		}
	}()
	go pRouter.stopListen()
	return true, ""
}

/** Запуск прослушивания адреса и порта**/
func (pRouter *THttpRouter) stopListen() {
	// Обработка паники при остановке (можно выдать и наружу)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(S_SRV_STOP_ERR, r)
		}
	}()
	// Ожидаем завершение приложения или вызов остановки
	select {
	case <-pRouter.UserStop:
	case <-pRouter.sysStop:
	}
	// Проверяем что сервер инициализирован
	if pRouter.pServer == nil {
		return
	}
	pRouter.pServer.Shutdown(context.Background())
	pRouter.pServer = nil
}

/* Обработка запроса локальная - без запуска сервера - прямо через mux*/
func (pRouter *THttpRouter) ProcRequestInternal(_pReq *http.Request) (_pResp *httptest.ResponseRecorder, _fOk bool, _sError string) {
	if pRouter.pMux == nil {
		return nil, false, S_MUX_ISNT_STARTED_ERR
	}
	_pResp = httptest.NewRecorder()
	pRouter.pMux.ServeHTTP(_pResp, _pReq)
	return _pResp, true, ""
}
