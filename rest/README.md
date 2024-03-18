# running local dynamodb
docker run --rm -d -p 8000:8000 amazon/dynamodb-local
# swaggo
write comments in main.go
go install github.com/swaggo/swag/cmd/swag@latest
"swag init" to generate docs files
more info references: https://github.com/swaggo/swag?tab=readme-ov-file#how-to-use-it-with-gin
access http://localhost:8080/docs/index.html
# format, lint and test
go fmt ./...
golangci-lint run ./... -v
go mod tidy
go test ./... -v -count=1
# live reloading
go get -u github.com/cosmtrek/air
create .air.toml file
run air
# go build
git add .
git commit -m 'module vx.x.x'
git tag vx.x.x
git push origin vx.x.x
# go module
$env:GOPROXY="goproxy.cn"
go list -m github.com/clarkhao/go-docker@vx.x.x
