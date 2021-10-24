package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		valueStream := make(chan interface{})
		go func(done In, in In) {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					valueStream <- v
				}
			}
		}(done, in)
		in = stage(valueStream)
	}
	return in
}
