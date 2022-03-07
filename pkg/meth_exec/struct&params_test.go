package meth_exec

import (
	"reflect"
	"testing"
	"time"

	//	"github.com/stretchr/testify"
	"github.com/stretchr/testify/require"
)

type TestReq struct {
	Index int
	Name  string
	Data  time.Time
}

/**Проверка преобразования структуры в список параметров**/
func TestStructToParamInfo(t *testing.T) {

	pInfos := StructToParamInfo(reflect.TypeOf((*TestReq)(nil)).Elem())
	require.Truef(t, len(pInfos) == 3, "Кол-во параметров должно быть 3 а оно %d", len(pInfos))
	require.Truef(t, pInfos[0].Name == "Index", "Неправильное имя параметра %s, а должно быть Index", pInfos[0].Name)
	require.Truef(t, pInfos[1].Type == reflect.TypeOf(""),
		"Неправильный тип параметра %s - %s а должно быть %s", pInfos[1].Name, pInfos[1].Type, reflect.TypeOf(""))
}

/**Проверка преобразования структуры в список параметров**/
func TestParamsToStruct(t *testing.T) {
	//pInfos := StructToParamInfo(reflect.TypeOf((*TestReq)(nil)).Elem())
	/*

		pInfos := ParamsToStruct(reflect.TypeOf((*TestReq)(nil)).Elem())

		require.Truef(t, len(pInfos) == 3, "Кол-во параметров должно быть 3 а оно %d", len(pInfos))
		require.Truef(t, pInfos[0].Name == "Index", "Неправильное имя параметра %s, а должно быть Index", pInfos[0].Name)
		require.Truef(t, pInfos[1].Type == reflect.TypeOf(""),
			"Неправильный тип параметра %s - %s а должно быть %s", pInfos[1].Name, pInfos[1].Type, reflect.TypeOf(""))
	*/
}
