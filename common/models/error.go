package models

type ErrorMessageObj struct {
	Language                 string `json:"language"`
	Success                  string `json:"success"`
	ParamEmptyError          string `json:"param_empty_error"`
	ParamValidateError       string `json:"param_validate_error"`
	SysParamEmptyError       string `json:"sys_param_empty_error"`
	PasswordDecodeError      string `json:"password_decode_error"`
	PasswordEncodeError      string `json:"password_encode_error"`
	ExecRemoteCommandError   string `json:"exec_remote_command_error"`
	InstallMinionError       string `json:"install_minion_error"`
	UninstallMinionError     string `json:"uninstall_minion_error"`
	S3UrlValidateError       string `json:"s3_url_validate_error"`
	S3FileEmptyError         string `json:"s3_file_empty_error"`
	S3DownloadError          string `json:"s3_download_error"`
	S3UploadError            string `json:"s3_upload_error"`
	SaltApiTargetError       string `json:"salt_api_target_error"`
	SaltApiConnectError      string `json:"salt_api_connect_error"`
	DecompressSuffixError    string `json:"decompress_suffix_error"`
	UnpackFileError          string `json:"unpack_file_error"`
	MysqlConnectError        string `json:"mysql_connect_error"`
	AddMysqlDatabaseError    string `json:"add_mysql_database_error"`
	DeleteMysqlDatabaseError string `json:"delete_mysql_database_error"`
	RunMysqlCommandError     string `json:"run_mysql_command_error"`
	FileNotExistError        string `json:"file_not_exist_error"`
	RunMysqlScriptError      string `json:"run_mysql_script_error"`
	MysqlCreateUserError     string `json:"mysql_create_user_error"`
	RunRemoteScriptError     string `json:"run_remote_script_error"`
	RedisAddUserError        string `json:"redis_add_user_error"`
	RedisDeleteUserError     string `json:"redis_delete_user_error"`
	RedisGrantUserError      string `json:"redis_grant_user_error"`
	RedisRevokeUserError     string `json:"redis_revoke_user_error"`
	GenSSHKeyError           string `json:"gen_ssh_key_error"`
}
