# govulnreview
Golang script to perform vulnerability review


## Setup
```
git clone github.com/manasmbellani/goengine /opt/goengine
cd /opt/goengine
./setup.sh
```

## Usage
To execute commands on a folder such as `/tmp` e.g. code search/grep
```
echo folder:///tmp | goengine -m via_auto_grep -f vulnreview.yaml-test
```

To execute web checks on a URL/protocol path:
```
echo http://www.google.com | goengine -m via_auto_test_webrequest -f vulnreview.yaml-test
```