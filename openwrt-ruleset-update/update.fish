#!/bin/env fish

# https://github.com/MetaCubeX/meta-rules-dat/tree/sing

set -g rule_set_name $argv[1]
set -g site_name geosite_url
set -g ip_name geoip_url

set -g failed_urls

set -g router_user root
set -g router_ip   192.168.1.1
set -g router_path /etc/homeproxy/ruleset

set -g failed_uploads

function classify
	grep -r geosite $rule_set_name > $site_name
	grep -r geoip   $rule_set_name > $ip_name
end

function download -a target_dir target_urls
	rm -rf $target_dir
	mkdir $target_dir
	for srs in (cat $target_urls)
		wget $srs -P $target_dir || set failed_urls $failed_urls $srs
	end
	rm -f $target_urls
end

function upload -a target_dir
	for srs in (ls $target_dir)
		scp -O -o HostKeyAlgorithms=+ssh-rsa $target_dir/$srs $router_user@$router_ip:$router_path/$target_dir || set failed_uploads $failed_uploads $target_dir/$srs
	end
end

function download_conclude
	if test -z "$failed_urls"
		echo "Successfully download all rule sets!"
		return
	end

	echo "================================================================================"
	for url in $failed_urls
		echo "download failed: $url"
	end
	echo "================================================================================"
end

function upload_conclude
	if test -z "$failed_uploads"
		echo "Successfully upload all rule sets!"
		return
	end

	echo "================================================================================"
	for srs in $failed_uploads
		echo "upload failed: $srs"
	end
	echo "================================================================================"
end

function start
	classify

	download geoip   $ip_name

	download geosite $site_name

	download_conclude

	upload geoip

	upload geosite

	upload_conclude
end

if test -z "$argv"
	echo "Require a rule set file!"
	exit 1
end

if not test -e $rule_set_name
	echo "'$rule_set_name' does not exists in current directory"
	exit 1
end

start

