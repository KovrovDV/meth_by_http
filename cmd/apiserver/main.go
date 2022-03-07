package main

import (
	"fmt"
	"meth_by_http/cmd/apiserver/test_proc"
	"meth_by_http/pkg/meth_exec"
)

func main() {

	pSrv := meth_exec.NewProcCoverSmp(new(test_proc.TestHandler))

	pSrv.Init()
	/*	fmt.Printf("Список методов %d", len(pSrv.arMethods))
		fmt.Println()

	*/
	//fmt.Println(pSrv.arMethods[0].sName)

	/*	for _, pMethInfo := range pSrv.arMethods {
			fmt.Println(pMethInfo.sName)
		}
	*/
	fmt.Println(" Exec GetName")
	//	fmt.Println(pSrv.Exec("GetName", struct{ Index int }{1}))

	//	fmt.Println(pSrv.ExecOut("GetName", map[string]interface{}{"Index": "1"}))
	fmt.Println(pSrv.ExecOut("GetNameWOPtr", map[string]interface{}{"Index": "1"}))

	/*
		http.HandleFunc("/", pSrv.ProcRequest)
		http.ListenAndServe(":8080", nil)
	*/

}
