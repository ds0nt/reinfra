
regenerate:
	protoc --go_out=plugins=grpc:. test/test.proto

test:
	ginkgo test -v .

PHONY: test-travis regenerate test