# optimint-debug

This queries two different nodes for `/namespaced_data/[8]{0,1,2,3,4,5,6,8}` at `daHeight` between `os.Args[1]` to `os.Args[2]` then iterates over any Optimint blocks within the Celestia blocks and prints the heights.

Run with
```
go run main.go 55200 55200

daHeight: 55200, bridgeHeight: 1941, fullHeight: 1941
```
