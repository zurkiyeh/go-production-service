# go-production-service

# Notes
- Go mantra: accept interfaces and return structs 
# Why context?
This design helps us to retrieve values passed on with the request. Also helps when we log out/... Better for tracibility.
For example, at a certain layer we may:
> ctx = context.WithValue(ctx, "request_id", "id_value")

We can retrieve that in another layer:
> fmt.Println(ctx.Value("request_id")) // This will return "id_value"

Contexts are very important when doing timeouts and making sure they're honored while passing through the different layers.
