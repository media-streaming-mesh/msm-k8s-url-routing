# msm-k8s-url-routing

## Build image
`make docker-build`

## API transport modes
The the url routing service can be deployed with either an http or  
gRPC api. This can be set in the deployment container args section  
within the deployment yaml:  
`msm-k8s-svc-helper/deployment/deployment.yaml`

## Deploying and destroying
`make deploy`  
`make clean-deploy`

## HTTP-mode api docs


| API               | Description                                  |
| ----------------- | -----------------------------------------    |
| `GET` /apiv1/url-routing?url=    | resolve to internal host      |

ex:  
req: `http://localhost:9898/apiv1/url-routing?url=http://kubernetes:443`  
resp: `["http://192.168.65.4:6443"]`

## gRPC-mode api docs

golang gRPC client instrumenting:  
```
url := "http://kubernetes:443"

conn, err := grpc.Dial(grpcHost, grpc.WithInsecure(), grpc.WithBlock())  
client := NewGetEndpointClient(conn)

resp, err := client.Send(context.TODO(), &EndpointRequest{Req: url})
```

See `msm-k8s-svc-helper/pkg/api/gRPC/server_test.go` for more details. 


