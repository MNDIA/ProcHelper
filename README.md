一个小工具用于启动进程以其他uid,gid在安卓arm64平台上。
使用示例
```sh
cd /data/u/naiveDir/
/data/u/tools/prochelper -name "./naive"
ps -A -o pid,uid,gid,comm | grep naive   #result：114514 0 3005 naive
```
等效于
```sh
/data/u/tools/prochelper -name "/data/u/naiveDir/naive" -uid 0 -gid 3005 -outFile "null" -errFile "null" "-c /data/u/naiveDir/config.json"
ps -A -o pid,uid,gid,comm | grep naive   #result：114514 0 3005 naive
```
