version: '2'
services:
  saltstack:
    image: wecube-plugins-saltstack:{version}
    container_name: wecube-plugins-saltstack-{version}
    restart: always
    volumes:
      - /etc/localtime:/etc/localtime
      - {path}/data/minions_pki:/etc/salt/pki/master/minions
      - {path}/saltstack/logs:/home/app/wecube-plugins-saltstack/logs
      - {path}/data:/home/app/data
    ports:
      - "9099:80"
      - "19098:8080"
      - "4505:4505"
      - "4506:4506"
      - "8082:8082"
    environment:
      - minion_master_ip={master_ip}
      - DEFAULT_S3_KEY={access_key}
      - DEFAULT_S3_PASSWORD={secret_key}
      - GATEWAY_URL={core_url}
      - SALTSTACK_DEFAULT_SPECIAL_REPLACE='@,#'
      - SALTSTACK_ENCRYPT_VARIBLE_PREFIX='!,%'
      - SALTSTACK_FILE_VARIBLE_PREFIX='^'
      - S3_SERVER_URL={s3_url}
      - SALTSTACK_LOG_LEVEL=info
      - SALTSTACK_RESET_ENV=N