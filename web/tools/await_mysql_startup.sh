#!/bin/sh

set -eu

host="$1"
shift
user="$1"
shift
password="$1"
shift
cmd="$@"

echo "awaiting mysql startup..."
until mysqladmin ping -h"$host" -u"$user" -p"$password"
do
        >&2 echo -n "."
        sleep 1
        echo "awaiting..."
done

>&2 echo "MySQL is up - executing command"
exec $cmd

