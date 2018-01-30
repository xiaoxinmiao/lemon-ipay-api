FROM pangpanglabs/golang:jan AS builder
WORKDIR /go/src/lemon-ipay-api/
COPY ./ /go/src/lemon-ipay-api/
# disable cgo 
ENV CGO_ENABLED=0
# build steps
RUN echo ">>> 1: go version" && go version \
    && echo ">>> 2: go get" && go-wrapper download \
    && echo ">>> 3: go install" && go-wrapper install

# make application docker image use alpine
FROM  alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /go/bin
# copy config cert to image 
COPY ./tmp/ ./tmp/
RUN tar xf ./tmp/wxcert.tar.gz -C ./
# copy execute file to image
COPY --from=builder /go/bin/ ./
EXPOSE 5000
CMD ["./lemon-ipay-api"]