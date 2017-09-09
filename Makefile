

regenerate:
	protoc --go_out=plugins=grpc:. test/test.proto

test: regenerate
	ginkgo test -v .

PHONY: regenerate test