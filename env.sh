# default (empty file path) log to stdout
export LOG_FILE_PATH=
# default log to both file LOG_FILE_PATH and stdout
export LOG_NOT_STDOUT=
# default log level is debug
export LOG_LEVEL_INFO=

# the log file will be rotated if:
# * init
# * file size >100MB
# * mid night (location +07:00)
