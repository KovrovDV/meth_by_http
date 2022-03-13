package meth_api_test

import (
	"bytes"
	"io/ioutil"
	"meth_by_http/pkg/meth_api"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestHandler struct {
}

func (pProc *TestHandler) GetName(_pReq struct{ Index int }) {
}

func (pProc TestHandler) GetNameWOPtr(_pReq struct {
	Index int
	Name  string
}) (Resp struct {
	Ok    bool
	Error string
	Info  string
}) {
	return struct {
		Ok    bool
		Error string
		Info  string
	}{true, "", "Meth, Name - " + _pReq.Name}
}
func (pProc TestHandler) GetBytesLen(_pReq struct{ ArPost []byte }) (Resp struct {
	Ok  bool
	Len int
}) {
	return struct {
		Ok  bool
		Len int
	}{true, len(_pReq.ArPost)}
}

/**Запуск сервера**/
func startSrv(t *testing.T, _sAddress string, _fFull bool) meth_api.THttpRouter {
	pSrv := meth_api.NewHttpRouter()
	fOk, sErr := pSrv.AddRoute("test", new(TestHandler))
	if !fOk {
		require.Fail(t, "Ошибка создания пути для обработки %s", sErr)
	}
	// Получаем адрес
	if _fFull {
		fOk, sErr = pSrv.Listen(_sAddress, "/")
	} else {
		fOk, sErr = pSrv.ListenLocal("/")
	}
	if !fOk {
		require.Fail(t, "Ошибка запуска сервера пути для обработки %s", sErr)
	}
	return *pSrv
}

/**Тестирование запуска сервера**/
func TestRealSrvStart(t *testing.T) {
	pSrv := startSrv(t, "localhost:8082", true)
	defer func() {
		pSrv.UserStop <- true
	}()
}

/**Тестирование обработки запроса через http**/
func TestRealSrvRequest(t *testing.T) {
	pSrv := startSrv(t, "localhost:8081", true)
	defer func() {
		pSrv.UserStop <- true
	}()
	// Простой метод
	resp, err := http.Get("http://localhost:8081/test/GetName?Index=2")
	require.Truef(t, err == nil, "Ошибка запроса простого метода %s", err)
	require.Truef(t, resp.StatusCode == 200, "Ошибка ответа простого метода - код %d", resp.StatusCode)
	// Простой метод с ответом
	resp, err = http.Get("http://localhost:8081/test/GetNameWOPtr?Index=2&Name=Temp")
	require.Truef(t, err == nil, "Ошибка запроса простого метода c ответом %s", err)
	require.Truef(t, resp.StatusCode == 200, "Ошибка ответа простого метода ответом- код %d", resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.Truef(t, err == nil, "Ошибка чтения ответа %s", err)
	require.Truef(t, string(body) == "Meth, Name - Temp", "Ошибка ответа метода %s а должно быть Meth, Name - Temp", string(body))
	// Форма  с ответом
	pParam := url.Values{"Index": {"2"}, "Name": {"Temp"}}
	resp, err = http.PostForm("http://localhost:8081/test/GetNameWOPtr", pParam)
	require.Truef(t, err == nil, "Ошибка запроса простого метода c ответом %s", err)
	require.Truef(t, resp.StatusCode == 200, "Ошибка ответа простого метода ответом- код %d", resp.StatusCode)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	require.Truef(t, err == nil, "Ошибка чтения ответа %s", err)
	require.Truef(t, string(body) == "Meth, Name - Temp", "Ошибка ответа метода %s а должно быть Meth, Name - Temp", string(body))
}

/**Тестирование обработки запроса напрямую через mux**/
func TestSmpMethInternal(t *testing.T) {
	pSrv := startSrv(t, "", false)
	// Простой метод
	req, _ := http.NewRequest("GET", "/test/GetName?Index=2", nil)
	resp, ok, err := pSrv.ProcRequestInternal(req)
	require.Truef(t, ok, "Ошибка запроса простого метода %s", err)
	require.Truef(t, resp.Code == 200, "Ошибка ответа простого метода - код %d", resp.Code)
}

/**Тестирование обработки запроса с ответом напрямую через mux**/
func TestSmpMethWithRespInternal(t *testing.T) {
	pSrv := startSrv(t, "", false)
	// Простой метод с ответом
	req, _ := http.NewRequest("GET", "/test/GetNameWOPtr?Index=2&Name=Temp", nil)
	resp, ok, err := pSrv.ProcRequestInternal(req)
	require.Truef(t, ok, "Ошибка запроса простого метода c ответом %s", err)
	require.Truef(t, resp.Code == 200, "Ошибка ответа простого метода ответом- код %d", resp.Code)
	require.Truef(t, resp.Body.String() == "Meth, Name - Temp", "Ошибка ответа метода %s а должно быть Meth, Name - Temp", resp.Body.String())
}

/**Тестирование обработки запроса формой напрямую через mux**/
func TestFormMethWithRespInternal(t *testing.T) {
	pSrv := startSrv(t, "", false)
	pParam := url.Values{"Index": {"2"}, "Name": {"Temp"}}
	req, _ := http.NewRequest("POST", "/test/GetNameWOPtr", strings.NewReader(pParam.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, ok, err := pSrv.ProcRequestInternal(req)
	require.Truef(t, ok, "Ошибка запроса формой метода c ответом %s", err)
	require.Truef(t, resp.Code == 200, "Ошибка ответа формой метода ответом- код %d", resp.Code)
	require.Truef(t, resp.Body.String() == "Meth, Name - Temp", "Ошибка ответа формой метода %s а должно быть Meth, Name - Temp", resp.Body.String())
}

/**Тестирование обраотки массива POST через mux**/
func TestPostArInternal(t *testing.T) {
	pSrv := startSrv(t, "", false)
	arBytes := []byte{1, 2, 3, 5, 6, 7}
	req, _ := http.NewRequest("POST", "/test/GetBytesLen", bytes.NewReader(arBytes))
	resp, ok, err := pSrv.ProcRequestInternal(req)
	require.Truef(t, ok, "Ошибка запроса  метода массива c ответом %s", err)
	require.Truef(t, resp.Code == 200, "Ошибка ответа метода массива ответом- код %d", resp.Code)
	iLen, pErr := strconv.ParseInt(resp.Body.String(), 10, 32)
	require.Truef(t, pErr == nil, "Ошибка чтения ответа метода массива %s", pErr)
	require.Truef(t, int(iLen) == len(arBytes), "Ошибка ответа метода массива %d а должно быть %d", iLen, len(arBytes))
}
