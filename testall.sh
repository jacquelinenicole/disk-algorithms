go build diskScheduler.go

./diskScheduler fcfsPA.txt>fcfsPA.out
diff fcfsPA.base fcfsPA.out
./diskScheduler sstfPA.txt>sstfPA.out
diff sstfPA.base sstfPA.out
./diskScheduler scanPA.txt>scanPA.out
diff scanPA.base scanPA.out
./diskScheduler c-scanPA.txt>c-scanPA.out
diff c-scanPA.out c-scanPA.base
./diskScheduler lookPA.txt>lookPA.out
diff lookPA.base lookPA.out
./diskScheduler c-lookPA.txt>c-lookPA.out
diff c-lookPA.out c-lookPA.base

./diskScheduler fcfs01.txt>fcfs01.out
diff fcfs01.base fcfs01.out
./diskScheduler sstf01.txt>sstf01.out
diff sstf01.base sstf01.out
./diskScheduler scan01.txt>scan01.out
diff scan01.base scan01.out
./diskScheduler c-scan01.txt>c-scan01.out
diff c-scan01.base c-scan01.out
./diskScheduler look01.txt>look01.out
diff look01.base look01.out
./diskScheduler c-look01.txt>c-look01.out
diff c-look01.base c-look01.out

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
