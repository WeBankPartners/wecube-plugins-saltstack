# Salt-Stack Plugin Compile and Package Guide

## Before Compiling

1. Please prepare a linux host with the Internet connection. To speed up the compilation, we recommends the host with 4 cores CPU and  8GB RAM.

2. Ubuntu 16.04+ or CentOS 7.3+;

3. Git installed
    - Use yum to install
	```shell script
 	yum install -y git
 	```
	- Manual installation, please refer to [git installation documentation](https://github.com/WeBankPartners/we-cmdb/blob/master/cmdb-wiki/docs/install/git_install_guide_en.md)

4. Installed Docker 1.17.03.x+
    - Please refer to [docker installation documentation](https://github.com/WeBankPartners/we-cmdb/blob/master/cmdb-wiki/docs/install/docker_install_guide_en.md)

5. Please use `netstat` or `ss` command to confirm the host ports `8082`, `9090`, `4505`, `4506` are not occupied.

## Compiling and packaging

1. Clone the code through github

    Switch to the local repository directory and execute the command:

    ```shell script
    cd /data
	git clone https://github.com/WeBankPartners/wecube-plugins-saltstack.git
    ```

    Enter the Github account password as prompted, and you can pull the code to the local.

    After that, enter the `wecube-plugins-saltstack` directory and the structure is as follows:

    ![saltstack_dir](images/saltstack_dir.png)\

2. Compile and package the plugin

    - Get plugin binary package.

    ```shell script
	make build
	```
	![saltstack_build](images/saltstack_build.png)

    - Generate a docker image.

    ```shell script
	make image
	```
    ![saltstack_image](images/saltstack_image.png)

    - Make the plugin package.

    ```shell script
	make package PLUGIN_VERSION=v1.0
	```

    The variable **PLUGIN_VERSION** is the number of the plugin package version. After the compilation is completed, a zip plugin package will be generated.

    ![saltstack_zip](images/saltstack_zip.png)

