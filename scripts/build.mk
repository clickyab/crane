godep:
	$(GO) install -v github.com/tools/godep

grpc:
	$(GO) get  -v google.golang.org/grpc
	$(GO) get  -v github.com/golang/protobuf/protoc-gen-go
	protoc -I $(ROOT)/openrtb $(ROOT)/openrtb/*.proto --go_out=plugins=grpc:$(GOPATH)/src
	sed -i -r 's/"IAB([[:digit:]]+)S([[:digit:]]+)"/"IAB\1-\2"/' $(ROOT)/openrtb/v2.5/ortb.pb.go