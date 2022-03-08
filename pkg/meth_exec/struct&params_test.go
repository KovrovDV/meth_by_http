package meth_exec_test

import (
	"reflect"
	"testing"
	"time"

	"meth_by_http/pkg/meth_exec"

	"github.com/stretchr/testify/require"
)

type TestReq struct {
	Index int
	Name  string
	Data  time.Time
}

/**Проверка преобразования структуры в список параметров**/
func TestStructToParamInfo(t *testing.T) {

	pInfos := meth_exec.StructToParamInfo(reflect.TypeOf((*TestReq)(nil)).Elem())
	require.Truef(t, len(pInfos) == 3, "Кол-во параметров должно быть 3 а оно %d", len(pInfos))
	require.Truef(t, pInfos[0].Name == "Index", "Неправильное имя параметра %s, а должно быть Index", pInfos[0].Name)
	require.Truef(t, pInfos[1].Type == reflect.TypeOf(""),
		"Неправильный тип параметра %s - %s а должно быть %s", pInfos[1].Name, pInfos[1].Type, reflect.TypeOf(""))
}

/**Проверка преобразования структуры в список параметров**/
func TestParamsToStruct(t *testing.T) {
	arParams := []meth_exec.TParamInfo{
		{Name: "Index", Type: reflect.TypeOf((*int)(nil)).Elem()},
		{Name: "Name", Type: reflect.TypeOf("")},
		{Name: "Data", Type: reflect.TypeOf((*time.Time)(nil)).Elem()},
	}
	rDate := time.Date(2000, 01, 01, 10, 10, 10, 0, time.Local)
	arValues := map[string]interface{}{
		"Index": 1,
		"Name":  "test",
		"Data":  rDate,
	}
	pInfos, fOk, sError := meth_exec.ParamsToStruct(arParams, arValues)

	if !fOk {
		t.Fatalf("Ошибка при формировании структуры из списка значений %s", sError)
	}
	// Преобразование допускается только в неименованные объекты
	pReq := pInfos.(struct {
		Index int
		Name  string
		Data  time.Time
	})
	require.Truef(t, pReq.Index == 1, "Неправильно перенесено число %d, а должно быть 1", pReq.Index)
	require.Truef(t, pReq.Name == "test", "Неправильно перенесена строка %s, а должно быть test", pReq.Name)
	require.Truef(t, pReq.Data == rDate, "Неправильный перенесена дата %s а должно быть %s", pReq.Data, rDate)
}

/**Проверка преобразования структуры в список параметров со значениями**/
func TestStructToParams(t *testing.T) {
	pReq := TestReq{10, "test", time.Date(2000, 01, 01, 10, 10, 10, 0, time.Local)}
	pMap, fOk, sError := meth_exec.StructToParams(pReq)
	if !fOk {
		t.Fatalf("Ошибка при формировании списка значений из структуры %s", sError)
	}
	require.Truef(t, len(pMap) == 3, "Неправильно число параметров - %d, а должно быть 3", len(pMap))
	require.Truef(t, pMap["Index"] == 10, "Неправильно перенесено число %d, а должно быть 10", pMap["Index"])
	require.Truef(t, pMap["Name"] == "test", "Неправильно перенесена строка %s, а должно быть test", pMap["Name"])
	require.Truef(t, pMap["Data"] == pReq.Data, "Неправильный перенесена дата %s а должно быть %s", pMap["Data"], pReq.Data)
}
