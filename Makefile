

.PHONY: gen-api
gen-api:
	goctl api go --api ./desc/website.api --dir ./


.PHONY: gen-model
gen-model:
	goctl model mysql ddl -src="./doc/website-v1.0.sql" -dir="./model" -c


.PHONY: gen-dockerfile
gen-dockerfile:
	goctl docker -go ./waitlist.go



.PHONY: website-api
website-api-quick:
	docker-compose stop website
	docker-compose build website
	docker-compose up website

.PHONY: restart
restart:
	docker-compose stop && docker-compose rm -f && docker-compose build && docker-compose up -d