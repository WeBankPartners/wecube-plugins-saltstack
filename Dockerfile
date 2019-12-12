FROM  webankpartners/salt-master-base:v1

ENV APP_HOME=/home/app/wecube-plugins-saltstack

RUN export LOG_PATH=$APP_HOME/logs \
    && mkdir -p $APP_HOME $LOG_PATH

COPY static  $APP_HOME/static

COPY scripts/salt/minions /srv/salt/minions
COPY conf/s3conf /conf/s3conf
COPY template  /conf/template

COPY scripts/salt/user_manage.sh /srv/salt/base/user_manage.sh
COPY scripts/salt/formatAndMountDisk.py /srv/salt/base/formatAndMountDisk.py
COPY scripts/salt/getUnformatedDisk.py /srv/salt/base/getUnformatedDisk.py

COPY build/start.sh /start.sh
COPY scripts/salt/install_minion.sh $APP_HOME/scripts/salt/install_minion.sh
COPY scripts/salt/remove_master_unused_key.sh $APP_HOME/scripts/salt/remove_master_unused_key.sh

RUN chmod +x  /start.sh \
    && chmod +x $APP_HOME/scripts/salt/install_minion.sh \
    && chmod +x $APP_HOME/scripts/salt/remove_master_unused_key.sh

COPY wecube-plugins-saltstack $APP_HOME/

ENTRYPOINT [ "/bin/bash","-c","/start.sh" ]

