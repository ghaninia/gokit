#### validation with translation handlers for gin:

```go
func (h handler) handler(ctx gin.Context) {

	var request struct {
		Name string `json:"name" binding:"required"`
		Age  int    `json:"age" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		
		// h.translation is a instance of the translation pkg
		// ctx is a instance of the gin.Context
		// err is a instance of the gin binding error
		// http.StatusUnprocessableEntity is a http status code
		response.NewResponse(h.translation).
			Validation(err).
			WithStatusCode(http.StatusUnprocessableEntity)
			Echo(ctx)
		
	    return 
	}
}
```

#### for the response has payload you can use the following code:
```go
func (h handler) handler(ctx gin.Context) {
    // payload is a instance of the struct that you want to return
    // http.StatusOK is a http status code
    response.NewResponse(h.translation).
        Payload(payload).
		WithStatusCode(http.StatusOK)
        Echo(ctx)
}
```

#### for the response has error you can use the following code:
```go
func (h handler) handler(ctx gin.Context) {
	
    // err is a instance of the response.Error
    // http.StatusForbidden is a http status code
    response.NewResponse(h.translation).
        WithError(err).
        WithStatusCode(http.StatusForbidden).
        Echo()
	
    // if you need mapping status code to the error override the StatusCodeMapping
    response.NewResponse(h.translation, StatusCodeMapping).
        WithError(err).
		Echo(ctx)
	
}
```