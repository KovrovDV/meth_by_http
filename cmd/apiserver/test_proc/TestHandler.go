package test_proc

import "fmt"

type TestHandler struct {
}

func (pProc *TestHandler) GetName(_pReq struct{ Index int }) {
	fmt.Printf("Get Name, Index %d\r\n", _pReq.Index)
}

func (pProc TestHandler) GetNameWOPtr(_pReq struct {
	Index int
	Name  *string
}) (Resp struct {
	Ok    bool
	Error string
}) {
	return struct {
		Ok    bool
		Error string
	}{true, ""}
}
