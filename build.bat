echo compiling Go via Docker
docker run -it --rm -v //c/projects://app -e GOPATH="//app" -w "//app/src/github.com/rogeralsing/consulreaper" golang sh -c "CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags='-s' -o consulreaper"
echo Building Docker image

docker build -t rogeralsing/consulreaper .
docker push rogeralsing/consulreaper