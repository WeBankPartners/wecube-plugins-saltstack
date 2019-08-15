[salt]

name=salt repo
baseurl=http://{{minion_master_ip}}:9099/salt-repo
enabled=1
gpgcheck=1
gpgkey=http://{{minion_master_ip}}:9099/salt-repo/saltstack-signing-key
