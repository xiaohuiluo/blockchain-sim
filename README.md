
## blockchain consensus simulate

A simple demo to simulate blockchain consensus algorithms(pos and dpos) with p2p network.

## build
```bash
go build -o build/simulate main/simulate.go
```
## run
```bash
cd build

# simulate cmd help
./simulate --help

# run simulate cmd
# -l option is to set log_level, now support error,warn,info,debug
./simulate.go

# sub cmd help
simulate » help

# sim sub cmd help
simulate » sim --help

# simulate blockchain with pos consensus, 3 nodes and run 2 rounds
simulate » sim pos 3 2

# simulate blockchain with dpos consensus, 3 nodes and run 2 rounds
simulate » sim dpos 3 2
```
