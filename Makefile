define build_env
	go build -o ./bin/bycrod_center_$(1) ./main.go
endef

all: clean prod

clean: clean_test clean_prod

test:
	$(call build_env,test)

prod:
	$(call build_env,prod)
