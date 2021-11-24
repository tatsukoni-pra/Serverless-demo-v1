build-ThumbnailFunction:
	GOARCH=arm64 GOOS=linux go build -o thumbnail handler.go
	cp thumbnail $(ARTIFACTS_DIR)
	chmod 755 bootstrap
	cp bootstrap $(ARTIFACTS_DIR)