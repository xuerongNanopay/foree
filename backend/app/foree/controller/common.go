package foree_controller

import (
	"fmt"
	"net/http"

	"xue.io/go-pay/app/foree/service"
	"xue.io/go-pay/server/transport"
)

func afterLogger[Q any](req service.LoginReq, resp Q, hErr transport.HError) {
	if v, is := hErr.(*transport.InteralServerError); is {
		// use logger.
		fmt.Print(v.OriginalError.Error())
	} else {
		fmt.Println(hErr.Error())
	}
}

func emptyBeforeResponse[Q any](w http.ResponseWriter, resp Q) http.ResponseWriter {
	return w
}
