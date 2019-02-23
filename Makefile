default: build

build:
	docker run --rm -i \
		-e CGO_ENABLED=1 \
		-v "$(CURDIR):/go/src/github.com/Luzifer/auth-supercookie" \
		-w /go/src/github.com/Luzifer/auth-supercookie \
		golang:alpine \
		sh -c 'apk add --update build-base git && go get -d && go build -buildmode=plugin -o auth-supercookie.so'
