FROM  ccr.ccs.tencentyun.com/webankpartners/wecube-saltstack:v1.7

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
    && chmod +x $APP_HOME/scripts/rsautil.sh

COPY wecube-plugins-saltstack $APP_HOME/

ENTRYPOINT [ "/bin/bash","-c","/start.sh" ]

