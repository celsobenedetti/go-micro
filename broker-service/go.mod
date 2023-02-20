module github.com/celso-patiri/go-micro/broker

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
)

require github.com/celso-patiri/go-micro/helpers v0.0.0-20230220231013-264a87f2dcf9 // indirect

replace github.com/celso-patiri/go-micro/helpers => ../helpers/
