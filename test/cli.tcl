#!/bin/sh
# -*- tcl -*-
# The next line is executed by /bin/sh, but not tcl \
exec tclsh "$0" ${1+"$@"}

package require tcltest
namespace import ::tcltest::*

source [file join [file dirname [info script]] "commands.test"]

runAllTests
