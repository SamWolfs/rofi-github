##
# rofi-github
#
# @file
# @version 0.1
build:
	go build

.PHONY: format
format:
	go fmt github.com/SamWolfs/...

# end
