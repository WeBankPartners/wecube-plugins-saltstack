current_dir=$(shell pwd)
version=$(PLUGIN_VERSION)
project_name=$(shell basename "${current_dir}")

APP_HOME=src/github.com/WeBankPartners/wecube-plugins-saltstack
PORT_BINDING={{host_port}}:8081

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
	chmod +x ./build/register.xml.tpl
	sed -i 's/{{PLUGIN_VERSION}}/$(version)/' ./build/register.xml.tpl > ./register.xml
	sed -i 's/{{IMAGENAME}}/$(project_name):$(version)/' ./register.xml
	sed -i 's/{{PORTBINDING}}/$(PORT_BINDING)/' ./register.xml 
	docker save -o  image.tar $(project_name):$(version)
	zip  $(project_name)-$(version).zip image.tar register.xml
	rm -rf $(project_name)
	rm -f register.xml
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)	
