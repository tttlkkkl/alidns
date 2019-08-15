test:
	TEST_ZONE_NAME=lihuasheng.cn go test .
build:
	docker build -t tttlkkkl/cert-manage-alidns .
	docker push tttlkkkl/cert-manage-alidns