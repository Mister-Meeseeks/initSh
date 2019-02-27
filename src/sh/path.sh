
# Use case: Projects of the form proj/lib/sh/init.sh
#                                           /[initSh scripts]
#                                        /py/[python libraries]
#                                        /R/[RLibs]
#                                        ...
function getCoLibDir() {
    local coLibName=$1
    echo $(discoverProjectDir)/../$coLibName/
}

function getPythonCoLibDir() {
    getCoLibDir py
}

function getHaskellCoLibDir() {
    getCoLibDir hs
}

function getRCoLibDir() {
    getCoLibDir R
}

function getParentLibDir() {
    echo $(discoverProjectDir)/../
}

# Use case: Projects of the form proj/src/sh/init.sh
#                                     etc/[data]
#                                     lib/[objs]
#                                     var/
function getProjEtcDir() {
    echo $(getParentProjectDir)/etc/
}

function getProjLibDir() {
    echo $(getParentProjectDir)/lib/
}

function getProjVarDir() {
    echo $(getParentProjectDir)/var/
}

function getParentProjectDir() {
    echo $(discoverProjectDir)/../../
}

function prepareFreshView() {
    resetViews
    setPathForProject
}

function setPathForProject() {
    if [[ $doesProjectHavePrecedence -gt 0 ]] ; then
	export PATH=$(getProjectPaths):$PATH
    else
	export PATH=$PATH:$(getProjectPaths)
    fi
}

function resetViews() {
    rm -r $(getProjDirs)
}

function getProjectPaths() {
    getProjDirs | sed 's+ +:+g'

}

function getProjDirs() {
    echo  $(retrieveProjBinView) \
	$(retrieveProjLibView)
}

function retrieveProjBinView() {
    retrieveInternalSubDir $viewBinSubDir
}

function retrieveProjLibView() {
    retrieveInternalSubDir $viewLibSubDir
}

function retrieveInternalSubDir() {
    local subDir=$1
    local outDir=$INIT_SH_INSTANCE_DIR/$subDir
    mkdir -p $outDir
    echo $outDir
}

