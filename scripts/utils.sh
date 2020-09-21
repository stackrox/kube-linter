#!/usr/bin/env bash

function safe_tput() {
  tput "$@" 2>/dev/null || true
}

bold="$(safe_tput bold)"
reset="$(safe_tput sgr0)"
green="$(safe_tput setaf 2)"
yellow="$(safe_tput setaf 3)"
red="$(safe_tput setaf 1)"
black="$(safe_tput setaf 0; safe_tput setab 7)"

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
