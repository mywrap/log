# default (empty file path) log to stdout
export LOG_FILE_PATH=

# maximum number of days to retain old log files, default is 7
export LOG_MAX_DAY=
# maximum size in megabytes of a log file, default 100
export LOG_MAX_MB=
# default log to both file LOG_FILE_PATH and stdout
export LOG_NOT_STDOUT=
# default log level is debug
export LOG_LEVEL_INFO=

# the log file will be rotated if:
# * initially run this package
# * file size > 100MB
# * mid night
