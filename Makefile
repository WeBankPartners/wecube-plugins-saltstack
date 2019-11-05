export GOPATH=$(PWD)
current_dir=$(shell pwd)
version=$(shell ./build/version.sh)
project_name=$(shell basename "${current_dir}" )

APP_HOME=src/github.com/WeBankPartners/wecube-plugins-saltstack

archive:
	tar cvfz source.tar.gz *
	rm -rf src
	mkdir -p $(APP_HOME)
	rm -rf target
	mkdir target
	tar zxvf source.tar.gz -C $(APP_HOME)
	rm -rf source.tar.gz
	cd $(APP_HOME) && CGO_ENABLED=0 GOOS=linux go build
	cp start.sh stop.sh docker_run.sh docker_stop.sh makefile dockerfile register.xml target
	cp -R scripts target
	cp -R conf    target
	cd target && chmod -R 755 *.sh
	cp $(APP_HOME)/wecube-plugins-saltstack target
	cd target && tar cvfz $(PKG_NAME) *

fmt:
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) -w /go/src/github.com/WeBankPartners/$(project_name)/  golang:1.12.5 go fmt ./...

clean:
	rm -rf $(project_name)
	rm -rf  ./*.tar
	rm -rf ./*.zip

build: clean
	chmod +x ./build/*.sh
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) golang:1.12.5 /bin/bash /go/src/github.com/WeBankPartners/$(project_name)/build/build.sh 

image: build
	docker build -t $(project_name):$(version) .
     
package: image 
	sed 's/{{IMAGE_TAG}}/$(version)/' ./build/register.xml.tpl > ./register.xml
	sed -i 's/{{PLUGIN_VERSION}}/$(PLUGIN_VERSION)/' ./register.xml 
	docker save -o  $(project_name).tar $(project_name):$(version)
	zip  $(project_name)_$(PLUGIN_VERSION).zip $(project_name).tar register.xml
	rm -rf $(project_name)
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)	
	


	
