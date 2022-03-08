package meth_exec_test

import (
	"fmt"
	"meth_by_http/pkg/meth_exec"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestHandler struct {
}

func (pProc *TestHandler) GetName(_pReq struct{ Index int }) {
	fmt.Printf("Get Name, Index %d\r\n", _pReq.Index)
}

func (pProc TestHandler) GetNameWOPtr(_pReq struct {
	Index int
	Name  string
}) (Resp struct {
	Ok    bool
	Error string
}) {
	return struct {
		Ok    bool
		Error string
	}{true, "test"}
}

/**Тестирование создания обертки**/
func Test_ProcCover_New(t *testing.T) {
	pCover := meth_exec.NewProcCoverSmp(new(TestHandler))
	require.Truef(t, pCover != nil, "Не создана обрертка")
	/*
		require.Truef(t, reflect.TypeOf(pCover.pParamsConverter) != reflect.TypeOf((*meth_exec.TStrParamConverter)(nil)).Elem(),
			"Не правильный конвертер по умолчанию")
	*/
	// TODO Custom Converter
}

/**Тестирование инициализации метода**/
func Test_ProcCover_Init(t *testing.T) {
	pCover := meth_exec.NewProcCoverSmp(new(TestHandler))
	fOk, sError := pCover.Init()
	if !fOk {
		t.Fatalf("Ошибка при инициализации списка методов %s", sError)
	}
	/*
		require.Truef(t, len(pCover.ArMethods) == 2, "Неправильное число методов - %d, а должно быть 2", len(pCover.ArMethods))

		require.Truef(t, len(pCover.ArMethods[0].arInput) == 1,
			"Неправильное число параметров первого метода- %d, а должно быть 1", len(pCover.ArMethods[0].arInput))
		require.Truef(t, len(pCover.ArMethods[0].arOutput) == 0,
			"Неправильное число выходныхпараметров первого метода- %d, а должно быть 0", len(pCover.ArMethods[0].arOutput))
		require.Truef(t, len(pCover.ArMethods[1].arOutput) == 1,
			"Неправильное число выходныхпараметров второго метода- %d, а должно быть 1", len(pCover.ArMethods[1].arOutput))
	*/
}

/**Тестирование вызова метода**/
func Test_ProcCover_Exec(t *testing.T) {
	pCover := meth_exec.NewProcCoverSmp(new(TestHandler))
	fOk, sError := pCover.Init()
	if !fOk {
		t.Fatalf("Ошибка при инициализации списка методов %s", sError)
	}
	// Простой метод структурой
	pRes, fOk, sError := pCover.Exec("GetName", struct{ Index int }{1})
	require.Truef(t, fOk, "Ошибка выполнения метода структурой %s - %s ", "GetName", sError)
	require.Truef(t, pRes == nil, "Не пустой вывод у метода %s структурой", "GetName")
	// Метод c возвратом
	_, fOk, _ = pCover.Exec("GetNameWOPtr", struct{ Index int }{1})
	require.Truef(t, !fOk, "Отсуствие ошибки при недостающих элементах выполнения метода структурой %s ", "GetNameWOPtr")
	pRes, fOk, sError = pCover.Exec("GetNameWOPtr", struct {
		Index int
		Name  string
	}{1, "test"})
	require.Truef(t, fOk, "Ошибка выполнения метода структурой %s - %s ", "GetNameWOPtr", sError)
	require.Truef(t, pRes != nil, "Пустой вывод у метода %s структурой", "GetNameWOPtr")
	pErr := pRes.(struct {
		Ok    bool
		Error string
	})
	require.Truef(t, pErr.Error == "test", "Неправильный ответ у метода %s структурой - %s", "GetNameWOPtr", pErr.Error)

	// Простой метод набором параметров
	pResMap, fOk, sError := pCover.ExecOut("GetName", map[string]interface{}{"Index": "1"})
	require.Truef(t, fOk, "Ошибка выполнения метода списком параметров %s - %s ", "GetName", sError)
	require.Truef(t, len(pResMap) == 0, "Не пустой вывод у метода %s списком параметров", "GetName")
	// метод с возвратом набором параметров
	pResMap, fOk, sError = pCover.ExecOut("GetNameWOPtr", map[string]interface{}{"Index": "1"})
	require.Truef(t, fOk, "Ошибка выполнения метода списком параметров %s - %s ", "GetNameWOPtr", sError)
	require.Truef(t, len(pResMap) == 2, "Вывод у метода %s списком параметров %d, а должно быть 2", "GetNameWOPtr", len(pResMap))
	require.Truef(t, pResMap["Error"] == "test", "Неправильный ответ у метода %s списком параметров - %s", "GetNameWOPtr", pErr.Error)
}
