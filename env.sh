# default config (empty file path) only log to stdout,
# make sure log directory was created and the user has permission to write
export LOG_FILE_PATH=

# maximum number of days to retain old log files, default is 32 days
export LOG_MAX_DAY=

# maximum size in megabytes of a log file, default 100 MBs
export LOG_MAX_MB=

# config LOG_NOT_STDOUT only works if LOG_FILE_PATH is set,
# if this config is set to "true", the logger will only log to file (not to stdout)
export LOG_NOT_STDOUT=

# set LOG_LEVEL_INFO to "true" will skip debug level log
export LOG_LEVEL_INFO=

# the log file will be rotated if the logger writes and one of the following
# condition is met:
# * initially run this package
# * file size > 100MB
# * midnight in UTC (7AM in Vietnam)
