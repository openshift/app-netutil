go_sample:
	./hack/build.sh
default:
	./hack/build.sh
c_sample:
	./hack/build-c.sh
image:
	./hack/build-image.sh
clean:
	rm -rf gopath/
	rm -rf bin/
