# thrift-gen-docs

Step 1: Package it into an executable file
```
go build .
```

Step 2: Place the executable file in the bin directory under gopath

Step 3: Use docs plugin by :

```
thriftgo -g go -p docs xxx.thrift
```
