# Maxsei's Gitlist  
  
This is a simple program is designed to list all of the repositories owned by a  
user.  
  
## Installation  
  
1. This requries you have [ Go ](https://github.com/golang/go) installed with `$` `GOPATH` [ set up ] (https://github.com/golang/go/wiki/SettingGOPATH)  
2. Once Go is installed just do `$` `go install https://github.com/maxsei/gitlist`  
3. You should be able to now call `$` `gitlist` from anywhere in your terminal  
  
## USAGE:  
Usage of gitlist:  
  -e    echo (default false)  
  -m int  
        maximum number of queries (default -1)  
  -o string  
        name of the outfile (default "repos.txt")  
  -p int  
        pagination; results per page (default 30)  
  -t float  
        maxlatency on failed request ( default 3.0s ) (default 3)  
  -u string  
        username (required)  
