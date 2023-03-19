module github.com/celso-patiri/go-micro/broker

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
)

require (
	github.com/celso-patiri/go-micro/helpers v0.0.0-20230220231013-264a87f2dcf9
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/rabbitmq/amqp091-go v1.7.0
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

replace github.com/celso-patiri/go-micro/helpers => ../helpers/
