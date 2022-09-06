#!/bin/bash

BASEDIR="`dirname $0`"
LIBDIR="$BASEDIR/lib"

function loadScript() {
  scriptPath="$LIBDIR/$1.sh"
  if [ ! -f $scriptPath ]; then
    echo "** Failed to load: $scriptPath"
  else
    source $scriptPath
  fi;
}

function print_usage() {
  printf "Usage: $0 [arguments]\n"
  printf "\n"
  printf "Arguments:\n"
  printf "\t configure\n"
  printf "\t setup-gradle\n"
  printf "\t setup-maven\n"
  printf "\t clean\n"
  printf "\n"
  printf "Environment variables:\n"
  printf "\t GATEWAY_URL: The base URL of the security gateway without trailing /\n"
  printf "\t GATEWAY_USERNAME: Username to authenticate with gateway\n"
  printf "\t GATEWAY_PASSWORD: Password to authenticate with gateway\n"
  printf "\n"
}

function validateRunningEnv() {
  if [ -z "$GATEWAY_URL" ]; then
    error_msg "Gateway URL is missing"
    exit -1
  fi;

  if [ -z "$GATEWAY_USERNAME" ]; then
    error_msg "Gateway Username is missing"
    exit -1
  fi;

  if [ -z "$GATEWAY_PASSWORD" ]; then
    error_msg "Gateway Password is missing"
    exit -1
  fi;
}

function main() {
  loadScript "utils"
  loadScript "conventions"
  loadScript "configuration"

  command=$1
  if [ -z "$command" ]; then
    print_usage
    exit -1
  fi;

  loadConfigurationIfPresent
  case $command in
    configure)
      interactiveConfiguration
      ;;
    setup-gradle)
      validateRunningEnv
      loadScript "gradle"
      ;;
    setup-maven)
      validateRunningEnv
      loadScript "maven"
      ;;
    clean)
      loadScript "clean"
      ;;
    *)
      error_msg "Unknown command: $command"
      ;;
  esac
}

main $@
