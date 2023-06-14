# gRPC-Gateway-Demo
A simple repository to showcase gRPC gateway

## Steps to start the application

steps to run the application is as follows

`` 1. git clone https://github.com/sujaysamanta/grpc-gateway-demo.git``

`` 2. cd grpc-gateway-demo``

`` 3. go mod tidy``

`` 4. go run server/main.go``

open a new terminal window and then go to the same folder `grpc-gateway-demo`

`` 5. go run proxy/main.go``


## Interact with the application 

``curl -s -H 'x-api-version:2.0.0' -d '{"name": "unruffled-galileo"}' 'http://localhost:8081/v1/sayHello' | jq .``

you should see a response similar to this 

```json
{
  "message": "hello unruffled-galileo",
  "apiVersion": "2.0.0"
}
```

``curl -s -d '{"name": "unruffled-galileo"}' 'http://localhost:8081/v1/sayHello' | jq .``

```json
{
  "message": "hello unruffled-galileo",
  "apiVersion": "1.0.0"
}
```


✅ Note:

   you will see from the above response that

   when the header `x-api-version` is omitted it defaults `"1.0.0"`

©️ [️Hashicorp](https://www.hashicorp.com/)

📧 [sujay.samanta@hashicorp.com](mailto:sujay.samanta@hashicorp.com)
