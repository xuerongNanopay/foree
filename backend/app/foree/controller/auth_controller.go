package controller

import "context"

func GetWrapper[P any, Q any](f func(context.Context, P) (Q, error)) {
}
