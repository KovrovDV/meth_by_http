package meth_exec

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

const (
	// ============ строки ошибок ===========================
	S_OUT_TYPE_ONLY_STR_ERR   = "конвертер поддерживает только строки во внешнем формате"
	S_NO_PTR_IN_PARAMS        = "Параметр не может быть передан по ссылке"
	S_NUMBER_COVERT_ERR       = "Ошибка конвертации числа %s - %s"
	S_OUT_OF_RANGE_NUMBER_ERR = "Число не попадает в органичение конвертации числа %s - %s"
	S_NUMBER_FP_COVERT_ERR    = "Ошибка преобразования числа с пл. точкой %s - %s"
	S_DATETIME_CONVERT_ERR    = "Ошибка преобразования даты %s - %s"
	S_BYTEA_CONVERT_ERR       = "Ошибка преобразования массива байт base64 %s - %s"
)

type TStrParamConverter struct {
	DateFormat string
}

func NewStrParamConverter() *TStrParamConverter {
	return &TStrParamConverter{
		DateFormat: S_STD_DATE_FORMAT,
	}
}

/**Преобразование поля из системы наружу
   	in значение внутри метода, тип для выдачи наружу
**/
func (pConv TStrParamConverter) ConvertOut(_pInternal interface{}, _pType reflect.Type) (_pExternal interface{}, _fOk bool, _sError string) {
	// Не строка на выходе - выходим
	if _pType != reflect.TypeOf("") {
		return nil, false, S_OUT_TYPE_ONLY_STR_ERR
	}
	switch _pInternal := _pInternal.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", _pInternal), true, ""
	case float32, float64:
		return fmt.Sprintf("%f", _pInternal), true, ""
	case string, []rune:
		return _pInternal.(string), true, ""
	case time.Time:
		return _pInternal.Format(pConv.DateFormat), true, ""
	case []byte:
		return base64.StdEncoding.EncodeToString(_pInternal), true, ""
	default:
		return fmt.Sprint(_pInternal), true, ""
	}
}

/**Преобразование поля с наружи для передачи внутрь системы
	in значение внутри метода, тип для выдачи наружу
**/
func (pConv TStrParamConverter) ConvertIn(_pExternal interface{}, _pType reflect.Type) (_pInternal interface{}, _fOk bool, _sError string) {
	// Не строка на входе - выходим
	if reflect.TypeOf(_pExternal) != reflect.TypeOf("") {
		return nil, false, S_OUT_TYPE_ONLY_STR_ERR
	}
	sValue := _pExternal.(string)
	pElem := reflect.New(_pType).Elem().Interface()
	if _pType.Kind() == reflect.Ptr {
		return nil, false, S_NO_PTR_IN_PARAMS
	}

	switch pElem.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		iVal, err := strconv.Atoi(sValue)
		if err != nil {
			return nil, false, fmt.Sprintf(S_NUMBER_COVERT_ERR, sValue, err)
		}
		// TODO переделать на массив/map с ограничениями
		switch pElem.(type) {
		case int:
			return iVal, true, ""
		case int8:
			if iVal >= math.MinInt8 && iVal <= math.MaxInt8 {
				return int8(iVal), true, ""
			}
		case int16:
			if iVal >= math.MinInt16 && iVal <= math.MaxInt16 {
				return int16(iVal), true, ""
			}
		case int32:
			if iVal >= math.MinInt32 && iVal <= math.MaxInt32 {
				return int32(iVal), true, ""
			}
		case int64:
			return int64(iVal), true, ""
		case uint:
			if iVal >= 0 {
				return uint(iVal), true, ""
			}
		case uint8:
			if iVal >= 0 && iVal <= math.MaxUint8 {
				return uint8(iVal), true, ""
			}
		case uint16:
			if iVal >= 0 && iVal <= math.MaxUint16 {
				return uint16(iVal), true, ""
			}
		case uint32:
			if iVal >= 0 && iVal <= math.MaxUint32 {
				return uint32(iVal), true, ""
			}
		case uint64:
			if iVal >= 0 {
				return uint64(iVal), true, ""
			}
		}
		return nil, false, fmt.Sprintf(S_OUT_OF_RANGE_NUMBER_ERR, sValue, err)
	case float32, float64:
		pRes, err := strconv.ParseFloat(sValue, 64)
		if err != nil {
			return nil, false, fmt.Sprintf(S_NUMBER_FP_COVERT_ERR, sValue, err)
		}
		switch pElem.(type) {
		case float64:
			return pRes, true, ""
		case float32:
			if pRes <= math.MaxFloat32 {
				return float32(pRes), true, ""
			}
		}
		return nil, false, fmt.Sprintf(S_OUT_OF_RANGE_NUMBER_ERR, sValue, err)
	// TODO проверить преобразование из rune[]
	case string:
		return sValue, true, ""
	case []rune:
		return []rune(sValue), true, ""
	case time.Time:
		pRes, err := time.Parse(pConv.DateFormat, sValue)
		if err != nil {
			return nil, false, fmt.Sprintf(S_DATETIME_CONVERT_ERR, sValue, err)
		}
		return pRes, true, ""
	case []byte:
		pRes, err := base64.StdEncoding.DecodeString(sValue)
		if err != nil {
			return nil, false, fmt.Sprintf(S_BYTEA_CONVERT_ERR, sValue, err)
		}
		return pRes, true, ""
	default:
		return fmt.Sprint(_pInternal), true, ""
	}
}

/* преоразование внутреннего  типа во внешний*/
func (pConv TStrParamConverter) InTypeToOut(_pInternal reflect.Type) reflect.Type {
	return reflect.TypeOf("")
}
