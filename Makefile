docker:
	docker build -t cmdls .
	docker run --network host -it cmdls cmdls