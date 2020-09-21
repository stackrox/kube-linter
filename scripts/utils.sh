#!/usr/bin/env bash

bold="$(tput bold)"
reset="$(tput sgr0)"
green="$(tput setaf 2)"
yellow="$(tput setaf 3)"
red="$(tput setaf 1)"
black="$(tput setaf 0; tput setab 7)"

function eecho() {
  echo >&2 "$@"
}

function einfo() {
  eecho -en "${bold}${green}[INFO]${black} "
  eecho -n "$@"
  eecho -e "$reset"
}

function ewarn() {
  eecho -en "${bold}${yellow}[WARN]${black} "
  eecho -n "$@"
  eecho -e "$reset"
}

function eerror() {
  eecho -en "${bold}${red}[ERROR]${black} "
  eecho -n "$@"
  eecho -e "$reset"
}

function efatal() {
  eecho -en "${bold}${red}[FATAL]${black} "
  eecho -n "$@"
  eecho -e "$reset"
}

function die() {
  efatal "$@"
  exit 1
}
