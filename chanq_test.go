package chanq_test

import (
	"testing"

	"github.com/powerman/chanq"
)

type Small struct {
	S string
	I int
}

type Large struct {
	S1 Small
	S2 Small
	S3 Small
	S4 Small
	S5 Small
}

func makeSmall() Small {
	return Small{
		S: "something",
		I: 42,
	}
}

func newSmall() *Small {
	small := makeSmall()
	return &small
}

func makeLarge() Large {
	return Large{
		S1: makeSmall(),
		S2: makeSmall(),
		S3: makeSmall(),
		S4: makeSmall(),
		S5: makeSmall(),
	}
}

func newLarge() *Large {
	large := makeLarge()
	return &large
}

func queue[T any](in, out chan T, done chan struct{}) {
	q := chanq.NewQueue(out)
	for {
		select {
		case msg := <-in:
			q.Enqueue(msg)
		case q.C <- q.Elem:
			q.Dequeue()
		case <-done:
			return
		}
	}
}

func noqueue[T any](in, out chan T, done chan struct{}) {
	for {
		select {
		case msg := <-in:
			out <- msg
		case <-done:
			return
		}
	}
}

func benchmarkSlow[T any](b *testing.B, newT func() T, f func(in, out chan T, done chan struct{})) {
	b.Helper()

	done := make(chan struct{})
	in := make(chan T)
	out := make(chan T)

	go f(in, out, done)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in <- newT()
		<-out
	}
	close(done)
}

func benchmarkFast[T any](b *testing.B, newT func() T, f func(in, out chan T, done chan struct{})) {
	b.Helper()

	done := make(chan struct{})
	in := make(chan T)
	out := make(chan T)
	go func() {
		for i := 0; i < b.N; i++ {
			<-out
		}
		close(done)
	}()

	go f(in, out, done)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in <- newT()
	}
	<-done
}

func BenchmarkSlowNoSmall(b *testing.B) { benchmarkSlow(b, makeSmall, noqueue[Small]) }
func BenchmarkSlowSmall(b *testing.B)   { benchmarkSlow(b, makeSmall, queue[Small]) }
func BenchmarkFastNoSmall(b *testing.B) { benchmarkFast(b, makeSmall, noqueue[Small]) }
func BenchmarkFastSmall(b *testing.B)   { benchmarkFast(b, makeSmall, queue[Small]) }

func BenchmarkSlowNoRSmall(b *testing.B) { benchmarkSlow(b, newSmall, noqueue[*Small]) }
func BenchmarkSlowRSmall(b *testing.B)   { benchmarkSlow(b, newSmall, queue[*Small]) }
func BenchmarkFastNoRSmall(b *testing.B) { benchmarkFast(b, newSmall, noqueue[*Small]) }
func BenchmarkFastRSmall(b *testing.B)   { benchmarkFast(b, newSmall, queue[*Small]) }

func BenchmarkSlowNoLarge(b *testing.B) { benchmarkSlow(b, makeLarge, noqueue[Large]) }
func BenchmarkSlowLarge(b *testing.B)   { benchmarkSlow(b, makeLarge, queue[Large]) }
func BenchmarkFastNoLarge(b *testing.B) { benchmarkFast(b, makeLarge, noqueue[Large]) }
func BenchmarkFastLarge(b *testing.B)   { benchmarkFast(b, makeLarge, queue[Large]) }

func BenchmarkSlowNoRLarge(b *testing.B) { benchmarkSlow(b, newLarge, noqueue[*Large]) }
func BenchmarkSlowRLarge(b *testing.B)   { benchmarkSlow(b, newLarge, queue[*Large]) }
func BenchmarkFastNoRLarge(b *testing.B) { benchmarkFast(b, newLarge, noqueue[*Large]) }
func BenchmarkFastRLarge(b *testing.B)   { benchmarkFast(b, newLarge, queue[*Large]) }
