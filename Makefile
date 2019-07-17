default:
	./hack/build.sh
image:
	./hack/build-image.sh
c_sample:
	./hack/build-c.sh
clean:
	rm -rf gopath/
	rm -rf bin/
