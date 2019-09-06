salt_repo:
  file.recurse:
    - name: /etc/yum.repos.d
    - source: salt://minions/yum.repos.d
    - user: root
    - group: root
    - file_mode: 644
    - dir_mode: 755
    - include_empty: True
salt_minion_purge:
  pkg.purged:
    - pkgs:
      - salt-minion
salt_minion_install:
  pkg.installed:
    - pkgs:
      - salt-minion
    - require:
      - file: salt_repo
salt_minion_conf:
  file.managed:
    - name: /etc/salt/minion
    - source: salt://minions/conf/minion
    - user: root
    - group: root
    - mode: 640
    - template: jinja
    - defaults:
      minion_id: {{grains['id']}}
    - require:
      - pkg: salt_minion_install
salt_minion_service:
  service.running:
    - name: salt-minion
    - enable: True
    - require:
      - file: salt_minion_conf