.DEFAULT_GOAL = help/short

EDITOR ?= vim
SHELL ?= /bin/bash

SELF = $(MAKE)

DEFAULT_HELP_TARGET ?= help/short
HELP_FILTER ?= .*

include $(CURDIR)/Makefile.*

## test the project
test: go/test

## lint the project
lint: go/lint go/fmt

## clean the project
clean: go/clean
