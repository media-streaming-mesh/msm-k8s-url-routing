# msm-k8s-svc-helper

### Build image
`make docker-build`

### Deploy to cluster
`make deploy`

### Delete from cluster
`make clean-deploy`

### Http API Docs


| API               | Description                                  |
| ----------------- | -----------------------------------------    |
| `GET` /apiv1/url-routing?url=    | resolve to internal host      |

ex:  
req: `http://localhost:9898/apiv1/url-routing?url=http://kubernetes:443`  
resp: `["http://192.168.65.4:6443"]`
