current_dir=$(shell pwd)
version=$(PLUGIN_VERSION)
project_name=$(shell basename "${current_dir}")

APP_HOME=src/github.com/WeBankPartners/wecube-plugins-saltstack
PORT_BINDING={{ALLOCATE_PORT}}:8082

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
	sed 's/{{PLUGIN_VERSION}}/$(version)/' ./build/register.xml.tpl > ./register.xml
	sed -i 's/{{IMAGENAME}}/$(project_name):$(version)/g' ./register.xml
	sed -i 's/{{CONTAINERNAME}}/$(project_name)-$(version)/g' ./register.xml
	sed -i 's/{{PORTBINDING}}/$(PORT_BINDING)/' ./register.xml 
	docker save -o  image.tar $(project_name):$(version)
	zip  $(project_name)-$(version).zip image.tar register.xml
	rm -rf $(project_name)
	rm -f register.xml
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)	

upload: package
	$(eval container_id:=$(shell docker run -v $(current_dir):/package -itd --entrypoint=/bin/sh minio/mc))
	docker exec $(container_id) mc config host add wecubeS3 $(s3_server_url) $(s3_access_key) $(s3_secret_key) wecubeS3
	docker exec $(container_id) mc cp /package/$(project_name)-$(version).zip wecubeS3/wecube-plugin-package-bucket
	docker rm -f $(container_id)
	rm -rf $(project_name)-$(version).zip

push: image
	docker login -u $(dockerhub_user) -p $(dockerhub_pass) $(dockerhub_server)
	docker tag $(project_name):$(version) wecube-statstack:$(version)
	docekr push $(dockerhub_server)/$(dockerhub_path)/wecube-statstack:$(version)

run_container: push
	pwd=$(shell pwd)
	docker run -d -H $(server_addr) -p 9099:80 -p 9090:8080 -p 4505:4505 -p 4506:4506 -p $(server_port):8082 -v /etc/localtime:/etc/localtime -v $(base_mount_path)/data/minions_pki:/etc/salt/pki/master/minions -v $(base_mount_path)/saltstack/logs:/home/app/saltstack/logs -v $(base_mount_path)/data:/home/app/data  -e minion_master_ip=$(minion_master_ip) -e minion_passwd=$(minion_passwd) -e minion_port=$(minion_port)
