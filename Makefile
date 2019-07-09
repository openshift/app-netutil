default:
	./hack/build.sh
image:
	./hack/build-image.sh
clean:
	rm -rf gopath/
	rm -rf bin/
