package foree_router

import (
	"context"
	"fmt"
	"net/http"

	foree_service "xue.io/go-pay/app/foree/app_service"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
	http_util "xue.io/go-pay/util/http"
	reflect_util "xue.io/go-pay/util/reflect"
)

func sessionGetWrapper[P transport.ForeeSessionRequest, Q any](
	serviceName string,
	permission string,
	authService *foree_service.AuthService,
	handler func(context.Context, P) (Q, transport.HError),
) func(http.ResponseWriter, *http.Request) {
	validatePayloadAndPermissionBeforeProcess := func(ctx context.Context, r *http.Request, req P) (context.Context, transport.HError) {
		err := req.Validate()
		if err != nil {
			foree_logger.Logger.Warn(
				"Payload_Validation_Error",
				"ip", http_util.LoadRealIp(r),
				"sessionId", r.Header.Get("SESSION_ID"),
				"requestType", fmt.Sprintf("%T", req),
				"cause", req.Validate().Error(),
			)
			return ctx, err
		}
		session, sErr := authService.Authorize(ctx, req.GetSessionId(), permission)
		if sErr != nil {
			var userId int64
			if session != nil {
				userId = session.UserId
			}
			// Normal error when the token expired
			foree_logger.Logger.Info(fmt.Sprintf("%v_Fail", serviceName), "ip", http_util.LoadRealIpFromContext(ctx), "userId", userId, "sessionId", req.GetSessionId(), "cause", sErr.Error())
			return nil, sErr
		}
		return ctx, nil
	}

	return restful_wrapper.RestGetWrapper(
		handler,
		validatePayloadAndPermissionBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}

func sessionPostWrapper[P transport.ForeeSessionRequest, Q any](
	serviceName string,
	permission string,
	authService *foree_service.AuthService,
	handler func(context.Context, P) (Q, transport.HError),
) func(http.ResponseWriter, *http.Request) {
	validatePayloadAndPermissionBeforeProcess := func(ctx context.Context, r *http.Request, req P) (context.Context, transport.HError) {
		err := req.Validate()
		if err != nil {
			foree_logger.Logger.Warn(
				"Payload_Validation_Error",
				"ip", http_util.LoadRealIp(r),
				"sessionId", r.Header.Get("SESSION_ID"),
				"requestType", fmt.Sprintf("%T", req),
				"cause", req.Validate().Error(),
			)
			return ctx, err
		}
		session, sErr := authService.Authorize(ctx, req.GetSessionId(), permission)
		if sErr != nil {
			var userId int64
			if session != nil {
				userId = session.UserId
			}
			// Normal error when the token expired
			foree_logger.Logger.Info(fmt.Sprintf("%v_Fail", serviceName), "ip", http_util.LoadRealIpFromContext(ctx), "userId", userId, "sessionId", req.GetSessionId(), "cause", sErr.Error())
			return nil, sErr
		}
		return ctx, nil
	}

	return restful_wrapper.RestPostWrapper(
		handler,
		validatePayloadAndPermissionBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}

func commonEndFunc[P any, Q any](ctx context.Context, req P, resp Q, hErr transport.HError) {
	if reflect_util.IsNil(hErr) {
		return
	}

	if v, is := hErr.(*transport.InteralServerError); is {
		//TODO: alert
		foree_logger.Logger.Error("System Error", "cause", v.OriginalError.Error())
	}
}

func validatePayloadBeforeProcess[P transport.ForeeRequest](ctx context.Context, r *http.Request, req P) (context.Context, transport.HError) {
	err := req.Validate()
	if err != nil {
		foree_logger.Logger.Warn(
			"Payload_Validation_Error",
			"ip", http_util.LoadRealIp(r),
			"sessionId", r.Header.Get("SESSION_ID"),
			"requestType", fmt.Sprintf("%T", req),
			"cause", req.Validate().Error())
	}
	return ctx, err
}

func emptyAfterProcess[Q any](ctx context.Context, w http.ResponseWriter, resp Q, hErr transport.HError) (context.Context, http.ResponseWriter) {
	return ctx, w
}

func simpleGetWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestGetWrapper(
		handler,
		validatePayloadBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}

func simplePostWrapper[P transport.ForeeRequest, Q any](handler func(context.Context, P) (Q, transport.HError)) func(http.ResponseWriter, *http.Request) {
	return restful_wrapper.RestPostWrapper(
		handler,
		validatePayloadBeforeProcess,
		emptyAfterProcess,
		commonEndFunc,
		true,
	)
}
