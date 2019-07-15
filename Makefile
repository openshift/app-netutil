default:
	./hack/build.sh
image:
	./hack/build-image.sh
c_sample:
	@cd c_api && go build -o ../bin/libnetutil_api.so -buildmode=c-shared -v
	@cd samples/c_app && gcc -I$(CURDIR)/bin -L$(CURDIR)/bin -Wall -o $(CURDIR)/bin/c_sample app_sample.c -lnetutil_api
clean:
	rm -rf gopath/
	rm -rf bin/
