package web

func MiddlewareCORS(ctx Context) {
	ctx.CORS()
}
