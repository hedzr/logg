#!/bin/bash

# for d in *; do
# 	if [ -d $d ]; then
# 	fi
# done

# find . -type f -iname 'go.mod' -print0 | xargs -0I% echo "pushd \$(dirname %)>/dev/null && pwd && go mod tidy && popd >/dev/null; echo;echo;echo" | sh

set -e

MAIN_BUILD_PKG=(".")
MAIN_APPS=(cli)
APPS=(wzstore loop)

# for an in ${APPS[*]}; do
# 	echo "$an //"
# done

# LDFLAGS = -s -w -X 'github.com/hedzr/cmdr/v2/conf.Buildstamp=2024-10-25T18:09:06+08:00' -X 'github.com/hedzr/cmdr/v2/conf.GIT_HASH=580ca50' -X 'github.com/hedzr/cmdr/v2/conf.GitSummary=580ca50-dirty' -X 'github.com/hedzr/cmdr/v2/conf.GitDesc=580ca50 upgrade deps' -X 'github.com/hedzr/cmdr/v2/conf.BuilderComments=' -X 'github.com/hedzr/cmdr/v2/conf.GoVersion=go version go1.23.7 darwin/arm64' -X 'github.com/hedzr/cmdr/v2/conf.Version=0.5.1'

boot-cov() {
	go test ./... -v -race -cover -coverprofile=./logs/coverage-cl.txt -covermode=atomic -test.short -vet=off 2>&1 | tee ./logs/cover-cl.log
}

extract-app-version() {
	local DEFAULT_DOC_NAME="${DEFAULT_DOC_NAME:-${1:-slog/doc.go}}"
	APPNAME="$(grep -E "appName[ \t]+=[ \t]+" ${DEFAULT_DOC_NAME} | grep -Eo "\\\".+\\\"")"
	VERSION="$(grep -E "version[ \t]+=[ \t]+" ${DEFAULT_DOC_NAME} | grep -Eo "[0-9.]+")"
}

build-all-platforms() {
	local W_PKG=github.com/hedzr/cmdr/v2/conf
	local TIMESTAMP="$(date -u '+%Y-%m-%dT%H:%M:%S')"
	local GOVERSION="$(go version)"
	local GIT_VERSION="$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")"
	local GIT_REVISION="$(git rev-parse --short HEAD)"
	local GIT_SUMMARY="$(git describe --tags --dirty --always)"
	local GIT_DESC="$(git log --oneline -1)"
	local GIT_HASH="$(git rev-parse HEAD)"
	local GOBUILD_TAGS=-tags="hzstudio sec antonal"
	extract-app-version
	tip "APPNAME = $APPNAME, VERSION = $VERSION"
	local LDFLAGS="-s -w -X '${W_PKG}.Buildstamp=$TIMESTAMP' \
		-X '${W_PKG}.GIT_HASH=$GIT_REVISION' \
		-X '${W_PKG}.GitSummary=$GIT_SUMMARY' \
		-X '${W_PKG}.GitDesc=$GIT_DESC' \
		-X '${W_PKG}.BuilderComments=$BUILDER_COMMENT' \
		-X '${W_PKG}.GoVersion=$GOVERSION' \
		-X '${W_PKG}.Version=$VERSION' \
		-X '${W_PKG}.AppName=$APPNAME'"
	local CGO_ENABLED=0
	local LS_OPT="-G"
	local GOBIN="./bin" GOOS GOARCH
	local mbp ma an

	for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
		if [[ "$GOOS" != "aix" && "$GOOS" != "android" && "$GOOS" != "illumos" && "$GOOS" != "ios" ]]; then
			for mbp in ${MAIN_BUILD_PKG[*]}; do
				for ma in ${MAIN_APPS[*]}; do
					for an in ${APPS[*]}; do
						local ANAME="${mbp}/${ma}/${an}"
						echo -e "\n\n" && tip "BUILDING FOR ${ANAME} / $GOOS ...\n"
						if [ -d "$ANAME" ]; then
							for GOARCH in $(go tool dist list | grep -E "^$GOOS" | awk -F'/' '{print $2}' | sort -u); do
								SUFFIX="_${GOOS}-${GOARCH}"
								echo "  >> building ${ANAME} - ${GOARCH} ..."
								go build -trimpath -gcflags=all='-l -B' \
									"${GOBUILD_TAGS}" -ldflags "${LDFLAGS}" \
									-o ${GOBIN}/${an}${SUFFIX} \
									${ANAME}/ || exit
							done
							ls -la $LS_OPT $GOBIN/${an}_${GOOS}*
						fi
						return
					done
				done
			done
		fi
	done
}

test-all-platforms() {
	for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
		echo -e "\n\nTESTING FOR $GOOS ...\n"
		go test -v -race -coverprofile=coverage-$GOOS.txt -test.run=^TestDirTimestamps$ ./dir/ || exit
	done
}

cov1() {
	for GOOS in darwin linux windows; do
		go test -v -race -coverprofile=coverage-$GOOS.txt ./...
		go tool cover -html=coverage-$GOOS.txt -o cover-$GOOS.html
	done
	open cover-darwin.html
}

bf1() {
	# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	# Or: go tool dist list
	# Missed: posix
	for GOOS in $(go tool dist list | awk -F'/' '{print $1}' | sort -u); do
		echo -e "\n\nTESTING FOR $GOOS ...\n"
		go test -v -race -coverprofile=coverage-$GOOS.txt -test.run=^TestDirTimestamps$ ./dir/ || exit
	done
}

fmt() {
	echo fmt...
	gofmt -l -w -s .
}

lint() {
	echo lint...
	golint ./...
}

cyclo() {
	echo cyclo...
	gocyclo -top 10 .
}

all-ops() {
	fmt && lint && cyclo
}

all() { build-all-platforms "$@"; }

# if [[ $# -eq 0 ]]; then
# 	cmd=cov
# else
# 	cmd=${1:-cov} && shift
# fi
# $cmd "$@"

sleep() { tip "sleeping..."; }

######### SIMPLE BASH.SH FOOTER BEGIN #########

###
is_git_clean() { git diff-index --quiet "$@" HEAD -- 2>/dev/null; }
is_git_dirty() {
	if is_git_clean "$@"; then
		false
	else
		true
	fi
}

dbg() { ((DEBUG)) && printf ">>> \e[0;38;2;133;133;133m$@\e[0m\n" || :; }
tip() { printf "\e[0;38;2;133;133;133m>>> $@\e[0m\n"; }
wrn() { printf "\e[0;38;2;172;172;22m... [WARN] \e[0;38;2;11;11;11m$@\e[0m\n"; }
err() { printf "\e[0;33;1;133;133;133m>>> $@\e[0m\n" 1>&2; }
cmd_exists() { command -v $1 >/dev/null; } # it detects any builtin or external commands, aliases, and any functions
fn_exists() { LC_ALL=C type $1 2>/dev/null | grep -qE '(shell function)|(a function)'; }
fn_builtin_exists() { LC_ALL=C type $1 2>/dev/null | grep -q 'shell builtin'; }
fn_defined() { LC_ALL=C type $1 2>/dev/null | grep -qE '( shell function)|( a function)|( shell builtin)'; }

###
is_darwin() { [[ $OSTYPE == darwin* ]]; }
is_darwin_sillicon() { is_darwin && [[ $(uname_mach) == arm64 ]]; }
is_linux() { [[ $OSTYPE == linux* ]]; }
is_freebsd() { [[ $OSTYPE == freebsd* ]]; }
is_win() { in_wsl; }
in_wsl() { [[ "$(uname -r)" == *windows_standard* ]]; }

###
if_nix_typ() {
	case "$OSTYPE" in
	*linux* | *hurd* | *msys* | *cygwin* | *sua* | *interix*) sys="gnu" ;;
	*bsd* | *darwin*) sys="bsd" ;;
	*sunos* | *solaris* | *indiana* | *illumos* | *smartos*) sys="sun" ;;
	esac
	echo "${sys}"
}
if_nix() { [[ "$(if_nix_typ)" == "$1" ]]; }
if_mac() { [[ $OSTYPE == darwin* ]]; }
if_ubuntu() {
	if [[ $OSTYPE == linux* ]]; then
		[ -f /etc/os-release ] && grep -qi 'ubuntu' /etc/os-release
	fi
}
if_vagrant() { [ -d /vagrant ]; }
in_vagrant() { [ -d /vagrant ]; }
in_orb() { [[ -d /mnt/mac ]]; }
path_in_orb_host() { [[ "$1" = /mnt/mac/* ]]; }
if_centos() {
	if [[ $OSTYPE == linux* ]]; then
		if [ -f /etc/centos-release ]; then
			:
		else
			[ -f /etc/issue ] && grep -qEi '(centos|(Amazon Linux AMI))' /etc/issue
		fi
	fi
}
in_vmware() {
	if cmd_exists hostnamectl; then
		$SUDO hostnamectl | grep -E 'Virtualization: ' | grep -qEi 'vmware'
	else
		false
	fi
}
in_vm() {
	if cmd_exists hostnamectl; then
		# dbg "checking hostnamectl"
		if $SUDO hostnamectl | grep -iE 'chassis: ' | grep -q ' vm'; then
			true
		elif $SUDO hostnamectl | grep -qE 'Virtualization: '; then
			true
		fi
	else
		# dbg "without hostnamectl"
		false
	fi
}
if_upstart() { [[ $(/sbin/init --version) =~ upstart ]]; }
if_systemd() { [[ $(systemctl) =~ -\.mount ]]; }
if_sysv() { [[ -f /etc/init.d/cron && ! -L /etc/init.d/cron ]]; }

###
# The better consice way to get baseDir, ie. $CD, is:
#       CD=$(cd `dirname "$0"`;pwd)
# It will open a sub-shell to print the folder name of the running shell-script.

###
CD="$(cd $(dirname "$0") && pwd)" && BASH_SH_VERSION=v20251116 && DEBUG=${DEBUG:-0} && PROVISIONING=${PROVISIONING:-0}
SUDO=sudo && { [ "$(id -u)" = "0" ] && SUDO= || :; }
LS_OPT="--color" && { is_darwin && LS_OPT="-G" || :; }
(($#)) && {
	dbg "$# arg(s) | CD = $CD"
	check_entry() {
		local prefix="${1:-boot}" cmd="${2:-first}" && shift && shift
		if fn_exists "${prefix}_${cmd}_entry"; then
			eval "${prefix}_${cmd}_entry" "$@"
		elif fn_exists "${cmd}_entry"; then
			eval "${cmd}_entry" "$@"
		else
			prefix="${prefix}_${cmd}"
			if fn_exists $prefix; then
				eval $prefix "$@"
			elif fn_exists ${prefix//_/-}; then
				eval ${prefix//_/-} "$@"
			elif fn_exists $cmd; then
				eval $cmd "$@"
			elif fn_exists ${cmd//_/-}; then
				eval ${cmd//_/-} "$@"
			else
				err "command not found: $cmd $@"
				return 1
			fi
		fi
	}
	check_entry "${FN_PREFIX:-boot}" "$@"
} || { dbg "empty: $# | CD = $CD"; }
######### SIMPLE BASH.SH FOOTER END #########
