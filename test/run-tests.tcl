#!/bin/sh
# -*- tcl -*-
# The next line is executed by /bin/sh, but not tcl \
exec tclsh "$0" ${1+"$@"}

package require tcltest
namespace import ::tcltest::*

::tcltest::verbose {pass body error}
::tcltest::configure -testdir [file dirname [file normalize [info script]]]
# Expect is not compatible with multi-interp so run all the tests in the same process - unfortunately this dsiables
# test parallelism
::tcltest::configure -singleproc
eval ::tcltest::configure $argv

# Add the dist folder to the PATH so tests work both in and out of the launch
# container
set env(PATH) "$env(PATH):./dist/:";


# Workaround to make tcltest exit with a non-zero status code when a test fails
# or the tests crash - See https://groups.google.com/forum/#!topic/comp.lang.tcl/mAaGxQ1Die8
proc ::tcltest::cleanupTestsHook {} {
     variable numTests
     upvar 2 testFileFailures crashed
     set ::code [expr {$numTests(Failed) > 0}]
     if {[info exists crashed]} {
         set ::code [expr {$::code || [llength $crashed]}]
     }
}

runAllTests
exit $::code
