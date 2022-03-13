package main

import (
	"flag"
	"fmt"
	"meth_by_http/cmd/apiserver/test_proc"
	"meth_by_http/pkg/meth_api"
	"os"
)

func main() {
	// Формируем привязку методов test_proc.TestHandler с путем http
	pSrv := meth_api.NewHttpRouter()
	fOk, sErr := pSrv.AddRoute("test", new(test_proc.TestHandler))
	if !fOk {
		fmt.Printf("Ошибка запуска %s\r\n", sErr)
		os.Exit(-1)
	}
	// Получаем адрес из параметров и запускаем сервер
	sAddress := flag.String("address", "localhost:8080", "Адрес сервера host:port")
	fOk, sErr = pSrv.Listen(*sAddress, "/")
	if !fOk {
		fmt.Printf("Ошибка запуска %s\r\n", sErr)
		os.Exit(-1)
	}
	// Завершение
	defer func() {
		pSrv.UserStop <- true
	}()
	// Ждем получение  строки - и после него останавливаем сервер
	var input string
	fmt.Scanln(&input)

	// Прямой вызов методов объекта
	/*
		pSrv := meth_exec.NewProcCoverSmp(new(test_proc.TestHandler))
		pSrv.Init()
		fmt.Println(pSrv.Exec("GetName", struct{ Index int }{1}))
		fmt.Println(pSrv.ExecOut("GetName", map[string]interface{}{"Index": "1"}))
		fmt.Println(pSrv.ExecOut("GetNameWOPtr", map[string]interface{}{"Index": "1"}))
	*/

}
