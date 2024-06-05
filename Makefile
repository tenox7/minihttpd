docker:
	docker buildx build --platform linux/amd64,linux/arm64 -t tenox7/minihttpd:latest --load .

clean:
	docker buildx prune -a -f