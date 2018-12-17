package onelog

import (
	"os"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	var log = New(Stdout{os.Stdout}, InfoLevel, JsonPattern{})

	go log.Info().Bool("B", true).Int("INT", 10).Msg("abcdefg")
	log.Info().Bool("B", true).Int("INT", 10).Msg("abcdefg")
	go log.Info().Bool("B", true).Int("INT", 10).Msg("abcdefg")
	log.Info().Bool("B", true).Int("INT", 10).Msg("abcdefg")

	time.Sleep(50000000)
	//fmt.Printf("%p %d  %p  %d\n", in1, in1.test, in2, in2.test)
}

func TestWriteFile(t *testing.T) {
	//fw, _ := NewFileWriter("./log/a.log", 50000000)
	//var log = New(fw, TraceLevel)
	var log = New(Stdout{os.Stdout}, TraceLevel, JsonPattern{})

	//TimeFormat = ""

	log.Debug().AddRuntime(&CoroutineID{}).AddConstant("a1", "b2")

	go log.Error().Int("INT", 10).Bool("B", true).Msg("abc1231231defg")
	log.Panic().Int("INT", 10).Bool("B", true).Msg("abcdef23123g")
	go log.Debug().Int("INT", 10).Bool("B", true).Msg("abc31231defg")
	go log.Trace().Int("INT", 10).Bool("B", true).Msg("12312")
	log.Info().Bool("B", true).Int("INT", 10).Msg("54321")

	time.Sleep(50000000)
	log.Close()
	//fmt.Printf("%p %d  %p  %d\n", in1, in1.test, in2, in2.test)
}

func TestCaller(t *testing.T) {
	call()
}

func call() {
	fw, _ := NewFileWriter("./log/a.log", 50000000)
	var log = New(fw, TraceLevel, JsonPattern{})

	log.Debug().AddRuntime(&Caller{})

	go log.Error().Int("INT", 10).Bool("B", true).Msg("abc1231231defg")
	log.Panic().Int("INT", 10).Bool("B", true).Msg("abcdef23123g")
	go log.Debug().Int("INT", 10).Bool("B", true).Msg("abc31231defg")
	go log.Trace().Int("INT", 10).Bool("B", true).Msg("12312")
	log.Info().Bool("B", true).Int("INT", 10).Msg("54321")

	time.Sleep(50000000)
	log.Close()
}

func TestMulFile(t *testing.T) {
	fw, _ := NewFileWriter("./log/a.log", 50000000)
	//var log = New(fw, TraceLevel)
	mul := NewMultipleWriter(fw, Stdout{os.Stdout})

	var log = New(mul, TraceLevel, JsonPattern{})

	TimeFormat = ""

	log.Debug().AddRuntime(&CoroutineID{}).AddConstant("a1", "b2")
	log.Error().AddRuntime(&CoroutineID{}).AddRuntime(&Caller{}).AddConstant("a1", "b2")

	log.Error().Int("INT", 10).Bool("B", true).Msg("abc1231231defg")
	log.Panic().Int("INT", 10).Bool("B", true).Msg("abcdef23123g")
	go log.Debug().Int("INT", 10).Bool("B", true).Msg("abc31231defg")
	go log.Trace().Int("INT", 10).Bool("B", true).Msg("12312")
	go log.Debug().Int("INT", 10).Bool("谁知道是什么", true).Msg("中文")
	log.Info().Bool("B", true).Int("INT", 10).Msg("54321")

	time.Sleep(50000000)
	log.Close()
}

func TestPattern(t *testing.T) {
	fw, _ := NewFileWriter("./log/a.log", 50000000)
	//var log = New(fw, TraceLevel)
	mul := NewMultipleWriter(fw, Stdout{os.Stdout})

	var log = New(mul, TraceLevel, OldPattern{})

	//TimeFormat = ""

	log.Debug().AddRuntime(&CoroutineID{}).AddConstant("a1", "b2")
	log.Error().AddRuntime(&CoroutineID{}).AddRuntime(&Caller{}).AddConstant("a1", "b2")

	log.Error().Int("INT", 10).Bool("B", true).Msg("中文")
	log.Panic().Int("INT", 10).Bool("B", true).Msg("中文")
	go log.Debug().Int("INT", 10).Bool("谁知道是什么", true).Msg("中文")
	go log.Trace().Int("INT", 10).Bool("B", true).Msg("中文")
	log.Info().Bool("B", true).Int("INT", 10).Msg("中文")

	time.Sleep(50000000)
	log.Close()
}

func TestRuningTime(t *testing.T) {
	var log = New(Stdout{os.Stdout}, ErrorLeven, JsonPattern{})

	//TimeFormat = ""
	//log.Error().AddRuntime(&CoroutineID{}).AddRuntime(&Caller{}).AddConstant("a1", "b2")

	log.Error().Msg("中文")
	log.Error().Msg("中文")
	log.Error().Msg("中文")

	time.Sleep(50000000)
	log.Close()
}
