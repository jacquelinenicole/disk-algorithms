go build diskScheduler.go

./diskScheduler fcfs20.txt>fcfs20.out
diff fcfs20.base fcfs20.out
./diskScheduler sstf20.txt>sstf20.out
diff sstf20.base sstf20.out
./diskScheduler scan20.txt>scan20.out
diff scan20.base scan20.out
./diskScheduler c-scan20.txt>c-scan20.out
diff c-scan20.base c-scan20.out
./diskScheduler look20.txt>look20.out
diff look20.base look20.out
./diskScheduler c-look20.txt>c-look20.out
diff c-look20.base c-look20.out
