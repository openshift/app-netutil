go_sample:
	./hack/build.sh
c_sample:
	./hack/build-c.sh
dpdk_app:
	./hack/build-dpdkapp.sh
testpod:
	./hack/build-testpod.sh
# 'make' and 'make image' are left around for legacy, but not
# documented anywhere.
default:
	./hack/build.sh
image:
	./hack/build-testpod.sh
clean:
	rm -rf gopath/
	rm -rf bin/

.PHONY: install.tools lint gofmt

install.tools:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ${GOPATH}/bin

lint:
	@GOPATH=${GOPATH} ./hack/lint.sh

gofmt:
	@./hack/verify-gofmt.sh
