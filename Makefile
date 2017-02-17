define build_env
	eGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/bycrodata_grab_$(1) ./main.go
	tar zcvf ./bycrodata_grab_$(1).tar.gz ./bin/bycrodata_grab_$(1) ./conf/conf.$(1).toml ./start_$(1).sh
endef

all: clean prod

clean: clean_test clean_prod

clean_test:
	rm -rf ./bycrodata_grab_test.tar.gz

clean_prod:
	rm -rf ./bycrodata_grab_prod.tar.gz

test:
	$(call build_env,test)

prod:
	$(call build_env,prod)
