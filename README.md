# Goengine
Golang script to perform vulnerability review and execute other checks via `goengine`

## Setup
### Via go-get
```
go get -u github.com/manasmbellani/goengine/src/goengine
```

### Via git clone
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

To execute `awscli` commands:
```
echo "gcp://athenaenterprises2021@gmail.com:clean-road-305712:us-central-1:us-central-1f" | goengine -f ../vulnreview.yaml-test -c test_echo -m via_auto_test_gcp
```

To execute `gcloud` commands:
```
echo "gcp://athenaenterprises2021@gmail.com:clean-road-305712:us-central-1:us-central-1f" | goengine -f ../vulnreview.yaml-test -c test_echo -m via_auto_test_gcp
```