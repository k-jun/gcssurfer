VERSION := $$(make -s show-version)
MAIN = ./cmd/gcssurfer

.PHONY: show-version
show-version:
	@gobump show -r $(MAIN)

.PHONY: test
test:
	go test -v ./...

.PHONY: security
security:
	gosec ./...

.PHONY: release
release:
	gobump up -w "$(MAIN)"
	git add cmd/gcssurfer/main.go
	git commit -m "bump up version to $(VERSION)"
	git push origin master
	goxc
	ghr "v$(VERSION)" ./dist/snapshot/

	