package meth_exec_test

import (
	"meth_by_http/pkg/meth_exec"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/**Проверка преобразования внутрь**/
func Test_TStrParamConverter_ConvertIn(t *testing.T) {
	pConv := meth_exec.NewStrParamConverter()
	// Число
	pRes, fOk, sError := pConv.ConvertIn("123", reflect.TypeOf((*int)(nil)).Elem())
	require.Truef(t, fOk, "Ошибка пребразования в число %s", sError)
	require.Truef(t, pRes == 123, "Ошибка пребразования в число %d, a должно быть 123", pRes)
	_, fOk, _ = pConv.ConvertIn("test", reflect.TypeOf((*int)(nil)).Elem())
	require.Truef(t, !fOk, "Отсуствие ошибки пребразования строки в число")
	// Дата
	pRes, fOk, sError = pConv.ConvertIn("2010-01-01T18:03:02", reflect.TypeOf((*time.Time)(nil)).Elem())
	require.Truef(t, fOk, "Ошибка пребразования в дату %s", sError)
	require.Truef(t, pRes.(time.Time).Day() == 1, "Ошибка пребразования в дату - день %d, a должно быть 1", pRes.(time.Time).Day())
	_, fOk, _ = pConv.ConvertIn("2010ывывыв01-01T18:03:02", reflect.TypeOf((*time.Time)(nil)).Elem())
	require.Truef(t, !fOk, "Отсуствие ошибки пребразования строки в дату")
	// Строка
	pRes, fOk, sError = pConv.ConvertIn("test", reflect.TypeOf(""))
	require.Truef(t, fOk, "Ошибка пребразования в строки в строку %s", sError)
	require.Truef(t, pRes == "test", "Ошибка пребразования в строку %s, a должно быть test", pRes)
	// []rune
	pRes, fOk, sError = pConv.ConvertIn("test", reflect.TypeOf((*[]rune)(nil)).Elem())
	require.Truef(t, fOk, "Ошибка пребразования в строки в срез рун %s", sError)
	require.Truef(t, string(pRes.([]rune)) == "test", "Ошибка пребразования в срез рун %s, a должно быть test", pRes)

	// массив байт 54686524
	pRes, fOk, sError = pConv.ConvertIn("VGhlJA==", reflect.TypeOf((*[]byte)(nil)).Elem())
	require.Truef(t, fOk, "Ошибка пребразования в массив байт %s", sError)
	require.Truef(t, len(pRes.([]byte)) == 4, "Ошибка пребразования массива байт длина %d, a должно быть 4", len(pRes.([]byte)))
	require.Truef(t, pRes.([]byte)[1] == 0x68, "Ошибка пребразования массива байт второй байт %d, a должно быть 0x68", pRes.([]byte)[1])
	_, fOk, _ = pConv.ConvertIn("xcxcfssdsd", reflect.TypeOf((*[]byte)(nil)).Elem())
	require.Truef(t, !fOk, "Отсуствие ошибка при пребразования в массив байт")

}

/**Проверка преобразования наружу**/

func Test_TStrParamConverter_ConvertOut(t *testing.T) {
	pConv := meth_exec.NewStrParamConverter()
	pStrType := reflect.TypeOf("")

	// Число
	pRes, fOk, sError := pConv.ConvertOut(123, pStrType)
	require.Truef(t, fOk, "Ошибка пребразования в числа %s", sError)
	require.Truef(t, pRes == "123", "Ошибка пребразования в число %s, a должно быть 123", pRes)
	// Дата
	pRes, fOk, sError = pConv.ConvertOut(time.Date(2010, 1, 1, 18, 3, 2, 0, time.Local), pStrType)
	require.Truef(t, fOk, "Ошибка пребразования даты %s", sError)
	require.Truef(t, pRes == "2010-01-01T18:03:02", "Ошибка пребразования даты %s, a должно быть 2010-01-01T18:03:02", pRes)
	// Строка
	pRes, fOk, sError = pConv.ConvertOut("test", pStrType)
	require.Truef(t, fOk, "Ошибка пребразования в строки в строку %s", sError)
	require.Truef(t, pRes == "test", "Ошибка пребразования в строку %s, a должно быть test", pRes)
	// массив байт 54686524

	pRes, fOk, sError = pConv.ConvertOut([]byte{0x54, 0x68, 0x65, 0x24}, pStrType)
	require.Truef(t, fOk, "Ошибка пребразования массива байт %s", sError)
	require.Truef(t, pRes == "VGhlJA==", "Ошибка пребразования массива байт  %s, a должно быть VGhlJA==", pRes)
}

/**Проверка выдачи внешнего типа**/
func Test_TStrParamConverter_InTypeToOut(t *testing.T) {
	pConv := meth_exec.NewStrParamConverter()
	pStrType := reflect.TypeOf("")
	require.True(t, pConv.InTypeToOut(reflect.TypeOf((*time.Time)(nil)).Elem()) == pStrType, "Ошибка опредления выходного типа даты")
	require.True(t, pConv.InTypeToOut(reflect.TypeOf((*int)(nil)).Elem()) == pStrType, "Ошибка опредления выходного типа числа")
	require.True(t, pConv.InTypeToOut(pStrType) == pStrType, "Ошибка опредления выходного типа строки")
	require.True(t, pConv.InTypeToOut(reflect.TypeOf((*[]byte)(nil)).Elem()) == pStrType, "Ошибка опредления выходного типа массива")
	require.True(t, pConv.InTypeToOut(reflect.TypeOf((*float32)(nil)).Elem()) == pStrType, "Ошибка опредления выходного типа числа с пл. точкой")
}
