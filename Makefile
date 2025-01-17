build:
	@go build -o /bin/loadbalancer main.go

run: build
	@ ./bin/loadbalancer

clean:
	@rm -rf /bin/loadbalancer