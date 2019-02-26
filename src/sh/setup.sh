#!/bin/bash -eu

function setupLocalProject() {
    runForInitShCall initializeLocalProject
}

function initializeLocalProject() {
    prepareFreshView
    blueprintProject
    sweepImports
}

function runForInitShCall() {
    local initCmd=$*
    initializeInitShCall
    runForInitShCallSafeVar $initCmd
}

function runForInitShCallSafeVar() {
    local initCmd=$*
    if hasInitShBeenCalledHere ; then
	pushNewInitShProjectCall $initCmd
    else
	pushSubroutineInitShCall
    fi
}

function initializeInitShCall() {
    local emptyProjCallStr=""
    : ${INIT_SH_PROJECT_CALL:=$emptyProjCallStr}
    export INIT_SH_PROJECT_CALL
}

function hasInitShBeenCalledHere() {
    [[ "$INIT_SH_PROJECT_CALL" != $(discoverProjectDir) ]]
}

function pushNewInitShProjectCall() {
    local initCmd=$1
    estabilishStackForTopInitShCall
    $initCmd
}

function estabilishStackForTopInitShCall() {
    export INIT_SH_PROJECT_CALL=$(discoverProjectDir) 
    export INIT_SH_INSTANCE_DIR=$(randomizeInstanceDir)
    export INIT_SH_CALL_STACK_LEVEL=0
}

function pushSubroutineInitShCall() {
    export INIT_SH_CALL_STACK_LEVEL=$(($INIT_SH_CALL_STACK_LEVEL + 1))
}

function cleanupLocalProject() {
    : ${INIT_SH_DEBUG:=0}
    if [[ $INIT_SH_DEBUG -gt 0 ]] ; then
	notifyDebugTrace
    else
	rmLocalProject
    fi
}

function rmLocalProject() {
    if isTopLevelInitShCall ; then
	rm -r $INIT_SH_INSTANCE_DIR
    fi
}

function notifyDebugTrace() {
    if isTopLevelInitShCall; then
	echo "init.sh debugging: Retaining local inst at $INIT_SH_INSTANCE_DIR.\
(Unset INIT_SH_CLEANUP to disable behavior)" >&2
    fi
}

function isTopLevelInitShCall() {
    [[ $INIT_SH_CALL_STACK_LEVEL -eq 0 ]]
}

function nukeAllInstances() {
    rm -rf $(retrieveInstanceParent)
}