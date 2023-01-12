package tools

import "context"

func OperationID(ctx context.Context) string {
	s, _ := ctx.Value("operationID").(string)
	return s
}

func OpUserID(ctx context.Context) string {
	s, _ := ctx.Value("opUserID").(string)
	return s
}
