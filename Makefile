# remove default suffixes as we dont use them
.SUFFIXES:

# defaults to using -s, unless VERBOSE is set
ifeq ($(VERBOSE)_x, _x)
	MAKEFLAGS+=-s
endif

# set the shell to bash in case some environments use sh
SHELL := /bin/bash

# ==============================================================================

release:
	goreleaser release --snapshot
