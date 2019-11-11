current_dir=$(shell pwd)
version=$(PLUGIN_VERSION)
project_name=$(shell basename "${current_dir}")

APP_HOME=src/github.com/WeBankPartners/wecube-plugins-saltstack
PORT_BINDING={{host_port}}:8081

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
	./build/register.xml.tpl > ./register.xml
	sed -i 's/{{PLUGIN_VERSION}}/$(version)/' ./register.xml
	sed -i 's/{{IMAGENAME}}/$(project_name):$(version)/' ./register.xml
	sed -i 's/{{PORTBINDING}}/$(PORT_BINDING)/' ./register.xml 
	docker save -o  image.tar $(project_name):$(version)
	zip  $(project_name)_$(verson).zip image.tar register.xml
	rm -rf $(project_name)
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)	
	


	
