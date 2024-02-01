## saltstack基础软件镜像制作

1、 镜像基于centos标准模板镜像，安装了一些基础工具 (注意点: 在拉取centos最新镜像版本时请先去saltstack repo上查看支不支持该最新版)
```bash
docker run --name salt-master -itd --privileged=true centos:7.6.1810 /usr/sbin/init
docker exec -it salt-master /bin/bash
yum install -y epel-release net-tools iproute telnet zip unzip tcpdump
```

2、 安装saltstack服务
```bash
rpm --import https://repo.saltstack.com/yum/redhat/6/x86_64/latest/SALTSTACK-GPG-KEY.pub
## 编辑yum repo
vi /etc/yum.repos.d/saltstack.repo
```
粘贴内容
```text
[saltstack-repo]
name=SaltStack repo for RHEL/CentOS $releasever
baseurl=https://repo.saltstack.com/yum/redhat/$releasever/$basearch/latest
enabled=1
gpgcheck=1
gpgkey=https://repo.saltstack.com/yum/redhat/$releasever/$basearch/latest/$releaseverSALTSTACK-GPG-KEY.pub
```
然后执行
```bash
yum install -y salt-master salt-ssh salt-api pyOpenSSL make openssl python-pip
```

3、 修改salt-master配置
```bash
sed -i 's/#auto_accept: False/auto_accept: True/g' /etc/salt/master
sed -i '/#default_include/s/#default/default/g' /etc/salt/master
```

4、 配置salt-api
```bash
cd /etc/pki/tls/certs
make testcert (with phrase=saltapi)
cd /etc/pki/tls/private
openssl rsa -in localhost.key -out localhost_nopass.key (with phrase=saltapi)
chmod 755 /etc/pki/tls/certs/localhost.crt
chmod 755 /etc/pki/tls/private/localhost.key
chmod 755 /etc/pki/tls/private/localhost_nopass.key
useradd -M -s /sbin/nologin saltapi
passwd saltapi
```
```bash
echo 'rest_cherrypy:'>>/etc/salt/master.d/api.conf
echo ' port: 8080'>>/etc/salt/master.d/api.conf
echo ' ssl_crt: /etc/pki/tls/certs/localhost.crt'>>/etc/salt/master.d/api.conf
echo ' ssl_key: /etc/pki/tls/private/localhost_nopass.key'>>/etc/salt/master.d/api.conf
```
```bash
echo 'external_auth:'>>/etc/salt/master.d/eauth.conf
echo '  pam:'>>/etc/salt/master.d/eauth.conf
echo '    saltapi:'>>/etc/salt/master.d/eauth.conf
echo '      - .*'>>/etc/salt/master.d/eauth.conf
```

5、 测试salt-master和salt-api服务是否正常  
```bash
systemctl start salt-master
systemctl start salt-api
systemctl enable salt-master
systemctl enable salt-api
```

6、 增加s3cmd  
```bash
yum install -y mysql git  python-dateutil
git clone https://github.com/s3tools/s3cmd.git /opt/s3cmd
ln -s /opt/s3cmd/s3cmd /usr/bin/s3cmd
```

7、 安装httpd  
```bash
yum install httpd -y
#测试启动httpd
systemctl start httpd
```

8、 如果是升级，请把上一版本基础镜像中的salt-master key拷过来，在/etc/salt/pki/master/下的master.pem和master.pub替换到新镜像中

9、 保存镜像和上传
```bash
exit;
#保存镜像
docker export -o salt-master-3000.3.tar salt-master
docker import salt-master-3000.3.tar salt-base:v1.2
```
修改项目Dockerfile中的FROM引入新基础镜像
