# Description of IP List services
Run the following command from the root directory to generate the client and service stubs based on the proto file:

*protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/iplist.proto*

## Running the server
From the **server** directory run the following command:
*go run main.go*

## Running client for IP List services
From the **ven_client** directory run the following command:
*go run main.go*

## Running the bidirectional client/server for IP List services
From the **bidirectional_ven_client** directory run the following command:
*go run main.go*


