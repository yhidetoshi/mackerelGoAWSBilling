# mackerelGoAWSBilling

Blog: https://yhidetoshi.hatenablog.com/entry/2019/09/20/190937

■ デプロイ
```
export MKRKEY=XXX

curl -X POST https://api.mackerelio.com/api/v0/services \
    -H "X-Api-Key: ${MKRKEY}" \
    -H "Content-Type: application/json" \
    -d '{"name": "AWS", "memo": "aws cost"}'

make build
sls deploy --aws-profile <PROFILE> --mkrkey ${MKRKEY}
```


`$ make help`
```
build:             Build binaries
build-deps:        Setup build
deps:              Install dependencies
devel-deps:        Setup development
help:              Show help
lint:              Lint
```
