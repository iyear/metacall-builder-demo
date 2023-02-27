## MetaCall Builder

Work in progress.

```bash
go build -o builder && ./builder deps python | ./hack/buildctl.sh build --output type=docker,name=metacalldemo:deps-py | docker load
```

Maybe staging package is not necessary, it can be done in the builder package.

In DEMO, `deps` = `dev-deps-base` + `core-repo` + `core-build` + `py-env` - `dev-deps-base` - `core-repo` = `core-build` + `py-env`
