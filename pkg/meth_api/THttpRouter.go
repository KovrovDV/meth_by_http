package meth_api

import (
	"fmt"
	"meth_by_http/pkg/meth_exec"
	"net/http"
)

/** машрутизатор на обработку**/
type THttpRouter struct {
	pServer *http.Server
	pRoutes map[string]meth_exec.TProcCover
}

/**Конструктор**/
func NewHttpRouter() *THttpRouter {
	pRes := new(THttpRouter)
	pRes.pRoutes = make(map[string]meth_exec.TProcCover)
	pRes.pServer = nil
	return pRes
}

/**добавляем путь **/
func (pRouter *THttpRouter) AddRoute(_sRoute string, _pProcCover meth_exec.TProcCover) (_fOk bool, _sError string) {
	// Инициализируем обработку
	_fOk, _sError = _pProcCover.Init()
	if !_fOk {
		return _fOk, _sError
	}
	// Добавляем в список обработки
	pRouter.pRoutes[_sRoute] = _pProcCover
	return true, ""
}

/**обработка запроса **/
func (pRouter *THttpRouter) getProc() {
}

/**обработка запроса **/
func (pRouter *THttpRouter) ProcRequest(w http.ResponseWriter, r *http.Request) {

	// Инициализируем обработку
	_fOk, _sError = _pProcCover.Init()
	if !_fOk {
		return
	}
	// Добавляем в список обработки
	pRouter.pRoutes[_sRoute] = _pProcCover
	return
}

/** Запуск прослушивания адреса и порта**/
func (pRouter *THttpRouter) Listen(_sAddress string) (_fOk bool, _sError string) {
	if pRouter.pServer != nil {
		return false, "сервер уже запущен"
	}
	// Инициализиреу сервер
	pRouter.pServer = &http.Server{Addr: _sAddress, Handler: handler}
	// Запуск
	go func() {
		if err := pRouter.pServer.ListenAndServe(); err != nil {
			// handle err
		}
	}()

	http.HandleFunc("/", pRouter.ProcRequest)
	http.ListenAndServe(_sAddres, nil)
}

/** Запуск прослушивания адреса и порта**/
func (pRouter *THttpRouter) StopListen() {
	// Обработка паники при остановке (можно выдать и наружу)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Ошибка остановки сервера  %s ", r)
		}
	}()
	// Проверяем что сервер инициализирован
	if pRouter.pServer == nil {
		return
	}
	pRouter.pServer.close()
	pRouter.pServer = nil
}
