## MetaCall Builder

Work in progress.

```bash
go build -o builder && ./builder dev py | ./hack/buildctl.sh build --output type=docker,name=metacalldemo:dev | docker load
```
