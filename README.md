# Using session affinity with external services in Kubernetes clusters

This project shows how to connect pods running in a Kubernetes cluster to
external services, using session affinity. When session affinity is enabled,
a pod connecting to a service (backed by several endpoints) will always use
the same instance, depending on the pod client IP.

## How to use it?

You need two Kubernetes clusters to use this demo:
 - cluster A is hosting backend services
 (2 endpoints with public IPs)
 - cluster B is hosting a frontend app, connected to
 the backend service

Switch to cluster A, and deploy the backend service:
```bash
$ kubectl apply -f k8s-backend.yml
```

Backend endpoints are exposed using a `LoadBalancer` object
with a public IP. Get these addresses:
```bash
$ kubectl -n backend get svc
NAME          TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
backend1-lb   LoadBalancer   10.110.59.187   34.76.187.102   80:31359/TCP   48s
backend2-lb   LoadBalancer   10.99.40.128    34.76.93.247    80:31281/TCP   48s
```

Now, you need to update the frontend app descriptor to include
the backend endpoints IP: edit `k8s-frontend.yml`.

```yaml
apiVersion: v1
kind: Endpoints
metadata:
  name: backend
  namespace: frontend
subsets:
  - addresses:
    # Replace these IPs with your public backend endpoint IPs.
    # Get endpoint IPs by running: $ kubectl -n backend get svc
    # Warning: your backend endpoints must run in a different K8s cluster.
    - ip: 34.76.187.102
    - ip: 34.76.93.247
    ports:
    - port: 80
```

Switch to cluster B, and deploy the frontend app:
```bash
$ kubectl apply -f k8s-frontend.yml
```

The frontend app is exposed using a `LoadBalancer` object.
Using the allocated IP address for this load balancer,
open a connection to the frontend app:
```bash
$ watch curl http://1.2.3.4
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
```

Note how the same backend instance is always used
by the frontend app.

You can disable session affinity in `k8s-frontend.yml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend
  namespace: frontend
spec:
  # Try disabling sessionAffinity (comment this property),
  # and see how different backend hosts are targeted.
  sessionAffinity: ClientIP
  ports:
  - protocol: TCP
    port: 9000
    targetPort: 9000
```

If you deploy the frontend app again, you can see
this time each connection is load balanced with
different backend endpoints:
```bash
$ watch curl http://1.2.3.4
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend2-74ddf7564f-dfrrx]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend1-54467975d9-sp2h2]
frontend[frontend-8b5d86d8f-29j49] -> backend[backend2-74ddf7564f-dfrrx]
```

## Contribute

Contributions are always welcome!

Feel free to open issues & send PR.

## License

Copyright &copy; 2019 [Pivotal Software, Inc](https://pivotal.io).

This project is licensed under the [Apache Software License version 2.0](https://www.apache.org/licenses/LICENSE-2.0).
