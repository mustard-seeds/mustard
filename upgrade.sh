#!/bin/bash

package="golang.org/x/text"

# govendor remove
govendor list|grep "$package"|awk '{print $2}'|while read line;
do
  echo "govendor remove $line"
  govendor remove $line
done

# rm dir
echo "rm -rf ./internal/$package"
rm -rf "./internal/$package"

# update
govendor add +external
