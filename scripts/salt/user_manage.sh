#!/bin/bash

help="
Usage:
        "$0"
        --action [add |remove]
        --user  linux user name, add action need this param
        --password linux user password,add action need this param
        --userId linux userId,option param for add action
        --group linux user group
        --groupId linux user groupId,option param for add action
        --home linux user home,option param for add action
  
"

log(){
    echo $(date +"[%Y%m%d %H:%M:%S]: ") $1
}

# 检测并处理 base64 编码的密码
decode_password(){
    local password="$1"

    # 直接尝试 base64 解码
    # shellcheck disable=SC2155
    local decoded_password=$(echo "$password" | base64 -d 2>/dev/null)

    # 检查解码是否成功且结果不为空
    if [[ $? -eq 0 ]] && [[ -n "$decoded_password" ]]; then
        # 验证解码后的内容是否包含可打印字符
        if [[ $decoded_password =~ ^[[:print:]]*$ ]]; then
            # 双向验证：重新编码后应该与原密码一致
            # shellcheck disable=SC2155
            local original_encoded=$(echo -n "$decoded_password" | base64 2>/dev/null)
            if [[ "$original_encoded" == "$password" ]]; then
                log "Password detected as base64 encoded, decoding successfully" >&2
                echo "$decoded_password"
                return 0
            else
                log "WARNING: Password looks like base64 but re-encoding doesn't match original, using original password" >&2
            fi
        else
            log "WARNING: Password base64 decoded but contains non-printable characters, using original password" >&2
        fi
    else
        log "INFO: Password is not base64 encoded or decode failed, using original password" >&2
    fi

    # 解码失败或验证不通过，返回原密码
    echo "$password"
    return 0
}

printHelp(){
        echo "$help"
        exit 1
}

parse_args(){
    if [ $# -lt 4 ]; then
       printHelp
    fi

    while [[ $# -gt 0 ]]
    do
    key="$1"

    case $key in
        --action)
            export ACTION=$2
            shift
        ;;
        --user)
            export USER_NAME=$2
            shift
        ;;
        --password)
            export USER_PWD=$2
            shift
        ;;
        --group)
            export GROUP=$2
            shift
        ;;
        --userId)
            export USER_ID=$2
            shift
        ;;
        --groupId)
            export GROUP_ID=$2
            shift
        ;;
        --home)
            export USER_HOME=$2
            shift
        ;;
        --makeDir)
            export MAKE_DIR=$2
            shift
        ;;
        --rwFile)
            export RW_FILE=$2
            shift
        ;;
        *)
            # unknown option
            log "unkonw option [$key]"
            printHelp
        ;;
    esac
    shift
    done
}

check_args(){
    if [[ -z $ACTION ]]; then
       log "param error:empty action param"
       printHelp
    fi 

    if [[ -z $USER_NAME ]]; then
       log "param error:empty user name"
       printHelp
    fi 
}

addUser(){
    id -u $USER_NAME >& /dev/null
    
    if [ $? -ne 0 ]; then
        if [[ -z $USER_PWD ]]; then
            log "param error:empty user password"
            printHelp
        fi
        
        groupId=""
        if [[ -n $GROUP_ID ]]; then
            groupId="-g "$GROUP_ID 
        fi

	    group=""
        if [[ -n $GROUP ]]; then
            group="-g "$GROUP
	    grep -qw ^$GROUP /etc/group || groupadd $GROUP $groupId
        fi 

        uid=""
        if [[ -n $USER_ID ]]; then
            uid="-u "$USER_ID 
        fi 

        home=""
        if [[ -n $USER_HOME ]]; then
            home="-d "$USER_HOME
            if [ ! -d $USER_HOME ];then
            mkdir -p $USER_HOME
            fi
        fi
        useradd $USER_NAME  $uid $home -m -p $(echo $USER_PWD | openssl passwd -1 -stdin) $group
    fi   
}

addGroup(){
    if [[ -n $GROUP ]]; then
        grep -qw ^$GROUP /etc/group || groupadd $GROUP
        usermod -a -G $GROUP $USER_NAME
    fi
}

changePassword(){
    echo $USER_PWD | passwd --stdin $USER_NAME
}

removeUser(){
    userdel -rf  $USER_NAME
}

addDir(){
    if [[ ! -z $MAKE_DIR ]]; then
        arr=(${MAKE_DIR//,/ })
        for i in ${arr[@]}
        do
            mkdir -p $i && chmod 777 $i
        done
    fi
}

authorizeFile(){
    if [[ ! -z $RW_FILE ]]; then
        arr=(${RW_FILE//,/ })
        for i in ${arr[@]}
        do
            chmod 777 $i
        done
    fi
}


parse_args "$@"
check_args "$@"

if [[ $ACTION = "add" ]];then
    addUser
    addGroup
    addDir
    authorizeFile
elif [[ $ACTION = "remove" ]];then
    removeUser
elif [[ $ACTION = "change_password" ]];then
    changePassword
else 
    log "param error:unknown action($ACTION)"
    printHelp
fi
