FROM  ccr.ccs.tencentyun.com/wecube/salt-master-base:v1
LABEL maintainer = "Webank CTB Team"

ENV APP_HOME=/home/app/wecube-plugins-deploy
ENV LOG_PATH=$APP_HOME/logs

RUN mkdir -p $APP_HOME $LOG_PATH

ADD wecube-plugins-deploy $APP_HOME/
ADD build/start.sh /start.sh
ADD scripts/salt/minions /srv/salt/minions
ADD conf/s3conf /conf/s3conf
ADD template  /conf/template

ADD scripts/salt/user_manage.sh /srv/salt/base/user_manage.sh
ADD scripts/salt/formatAndMountDisk.py /srv/salt/base/formatAndMountDisk.py
ADD scripts/salt/getUnformatedDisk.py /srv/salt/base/getUnformatedDisk.py

ADD scripts/salt/install_minion.sh $APP_HOME/scripts/salt/install_minion.sh
ADD scripts/salt/remove_master_unused_key.sh $APP_HOME/scripts/salt/remove_master_unused_key.sh

RUN chmod +x  /start.sh
RUN chmod +x $APP_HOME/scripts/salt/install_minion.sh
RUN chmod +x $APP_HOME/scripts/salt/remove_master_unused_key.sh
ENTRYPOINT [ "/bin/bash","-c","/start.sh" ]

