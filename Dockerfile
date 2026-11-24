FROM  ccr.ccs.tencentyun.com/webankpartners/wecube-saltstack:v1.10

ENV LANG=en_US.utf8
ENV APP_HOME=/home/app/wecube-plugins-saltstack
ENV DEFAULT_S3_KEY=access_key
ENV DEFAULT_S3_PASSWORD=secret_key

RUN export LOG_PATH=$APP_HOME/logs \
    && mkdir -p $APP_HOME $LOG_PATH $APP_HOME/minio-conf /run/httpd && chown -R root:apache /run/httpd
RUN mkdir -p /var/www/html/salt-minion && mkdir -p /var/www/html/salt-minion/conf && mkdir -p /var/www/html/tmp
COPY scripts/salt/install/* /var/www/html/salt-minion/
RUN chmod +x /var/www/html/salt-minion/minion_install.sh && chmod +x /var/www/html/salt-minion/minion_uninstall.sh

COPY static  $APP_HOME/static
COPY conf  $APP_HOME/conf

COPY scripts/salt/minions /srv/salt/minions
COPY conf/s3conf /conf/s3conf
COPY template  /conf/template

COPY scripts/salt/rsautil.sh $APP_HOME/scripts/
COPY scripts/salt/rsa_gen.sh $APP_HOME/scripts/
COPY scripts/salt/user_manage.sh /srv/salt/base/user_manage.sh
COPY scripts/salt/formatAndMountDisk.py /srv/salt/base/formatAndMountDisk.py
COPY scripts/salt/getUnformatedDisk.py /srv/salt/base/getUnformatedDisk.py

COPY build/start.sh /start.sh
COPY scripts/salt/install_minion.sh $APP_HOME/scripts/salt/install_minion.sh
COPY scripts/salt/uninstall_minion.sh $APP_HOME/scripts/salt/uninstall_minion.sh
COPY scripts/salt/remove_master_unused_key.sh $APP_HOME/scripts/salt/remove_master_unused_key.sh

RUN chmod +x  /start.sh \
    && chmod +x $APP_HOME/scripts/salt/install_minion.sh \
    && chmod +x $APP_HOME/scripts/salt/remove_master_unused_key.sh \
    && chmod +x $APP_HOME/scripts/salt/uninstall_minion.sh \
    && chmod +x $APP_HOME/scripts/rsautil.sh \
    && chmod +x $APP_HOME/scripts/rsa_gen.sh \
    && curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.tencentyun.com/repo/centos7_base.repo \
    && rm -f /etc/yum.repos.d/epel.repo /etc/yum.repos.d/epel-testing.repo \
    && yum clean all \
    && yum makecache \
    && yum install -y wget xz libaio ncurses-compat-libs \
    && wget https://dev.mysql.com/get/Downloads/MySQL-8.4/mysql-8.4.0-linux-glibc2.17-x86_64-minimal.tar.xz \
    && tar -xvf mysql-8.4.0-linux-glibc2.17-x86_64-minimal.tar.xz -C /opt  \
    && mv /usr/bin/mysql /usr/bin/mysql5.6  \
    && ln -s /opt/mysql-8.4.0-linux-glibc2.17-x86_64-minimal/bin/mysql /usr/bin/mysql

COPY wecube-plugins-saltstack $APP_HOME/

ENTRYPOINT [ "/bin/bash","-c","/start.sh" ]

