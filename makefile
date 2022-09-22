.PHONY: test-unit test-coverage mock-api


test-unit:
	go test ./... -v


test-coverage:
	go test -cover `go list ./... | grep -v -e /vendor/ -e /mock/`
	go test `go list ./... | grep -v -e /vendor/ -e /mock/` -coverprofile=cover.out
	go tool cover -func=cover.out
	# remove file
	@unlink cover.out

mock-api:
	docker run -it --rm -p 9999:9999 -v $$(pwd)/.tests/openmock:/data/templates checkr/openmock 