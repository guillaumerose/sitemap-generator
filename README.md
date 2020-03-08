# sitemap-generator

## Build

Please install go and run the following command:

```
make build
```

This will output the `sitemap-generator` binary in `bin/`

## Usage

Standalone crawler:

```
$ ./bin/sitemap-generator -h
Usage of ./bin/sitemap-generator:
  -d int
    	maximum depth to crawl (default 5)
  -p int
    	maximum number of concurrent requests (default 2)

$ ./bin/sitemap-generator -d 2 -p 2 https://www.redhat.com
INFO[0000] #0 Visiting /
INFO[0000] #2 Visiting /en
INFO[0000] #1 Visiting /en/search
...
INFO[0003] Found 145 URLs
- /
- /en
  - /about
    - /all-policies-guidelines
    - /around-the-world
    - /company
    - /development-model
    - /feedback
...
```

client/server:

```
$ ./bin/server
INFO[0000] Listening on :8080
```

```
$ ./bin/client -h
Usage of ./bin/client:
  -d int
    	maximum depth to crawl (default 5)
  -p int
    	maximum number of concurrent requests (default 2)
  -s string
    	crawler URL (default "http://127.0.0.1:8080")

$ ./bin/client https://kompose.io
INFO[0000] Crawling https://kompose.io (parallelism: 10, maxDepth: 5)
INFO[0000] 7 URLs found
- /
- /architecture
- /conversion
- /docs
  - /conversion.md
  - /maven-example.md
- /getting-started
- /installation
- /integrations
- /user-guide
```

## Deploy

Using Kubernetes:

```
$ kubectl apply -f config/deployment.yaml
$ kubectl apply -f config/service.yaml
$ kubectl get all
NAME                                     READY   STATUS    RESTARTS   AGE
pod/sitemap-generator-66b7488f6f-rzlg2   1/1     Running   0          4m11s

NAME                        TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/kubernetes          ClusterIP      10.96.0.1        <none>        443/TCP          6m9s
service/sitemap-generator   LoadBalancer   10.110.229.220   localhost     8080:30463/TCP   4m1s

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/sitemap-generator   1/1     1            1           4m11s

NAME                                           DESIRED   CURRENT   READY   AGE
replicaset.apps/sitemap-generator-66b7488f6f   1         1         1       4m11s
```
