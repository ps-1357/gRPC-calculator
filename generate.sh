#!bin/zshrc

protoc -I=greet/greetpb greet/greetpb/greet.proto --go-grpc_out=greet/

protoc -I=greet/greetpb --go_out=greet --go-grpc_out=greet greet.proto

protoc -I=calculator/calculatorpb --go_out=calculator --go-grpc_out=calculator calculator.proto

protoc -I=blog/blogpb --go_out=blog --go-grpc_out=blog blog.proto

