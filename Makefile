.PHONY: build

include Makefile.dev

# WE DON'T WANT TO FORGET ABOUT THE PROD VERSION, SO WE ALWAYS COMPILE FOR PROD
build_dev:
	GOPATH=${GOPATH}:`pwd`/vendor go build -o build/main main.go utils.go erroneouspage.go statusmachine.go

build:
	GOOS=linux GOARCH=amd64 GOPATH=${GOPATH}:`pwd`/vendor go build -v -o build/main_prod main.go utils.go erroneouspage.go statusmachine.go

tests:
	GOPATH=${GOPATH}:`pwd`/vendor go test

loop:
	env CRAWLER_USER_AGENT_STRING="Mozilla/5.0 (compatible; StatusMachineBot/1.0; +http://www.statusmachine.com/crawler)" \
		LINKSCRAWLER_TMP_DIR="/Users/juliendesrosiers/Sites/statusmachine/tmp" \
		DATABASE_URL="postgres://zlmnohxlwzyisd:yvT1ZaYudMeP4XN7Uxz7OcFtdY@ec2-54-204-43-200.compute-1.amazonaws.com:5432/d2k6mco30hv8js" \
		go run main.go utils.go erroneouspage.go statusmachine.go -islooping

# test:
#   env CRAWLER_USER_AGENT_STRING="Mozilla/5.0 (compatible; StatusMachineBot/1.0; +http://www.statusmachine.com/crawler)" \
#     DATABASE_URL="postgres://zlmnohxlwzyisd:yvT1ZaYudMeP4XN7Uxz7OcFtdY@ec2-54-204-43-200.compute-1.amazonaws.com:5432/d2k6mco30hv8js" \
#     ./main -islooping

test:
	env CRAWLER_USER_AGENT_STRING="Mozilla/5.0 (compatible; StatusMachineBot/1.0; +http://www.statusmachine.com/crawler)" \
		LINKSCRAWLER_TMP_DIR="./linkscrawler_tmp_dir" \
		./build/main 'http://www.juliendesrosiers.ca/sites/'

dev_test:
	env CRAWLER_USER_AGENT_STRING="Mozilla/5.0 (compatible; StatusMachineBot/1.0; +http://www.statusmachine.com/crawler)" \
		DATABASE_URL="postgres://zlmnohxlwzyisd:yvT1ZaYudMeP4XN7Uxz7OcFtdY@ec2-54-204-43-200.compute-1.amazonaws.com:5432/d2k6mco30hv8js" \
		./build/main2 -baseurl='http://test.juliendesrosiers.ca'

run_test:
	GOPATH=${GOPATH}:`pwd`/vendor go run testing_stuff.go
