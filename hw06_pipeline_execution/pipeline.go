package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := pipelinePump(done, in)
	for _, stage := range stages {
		out = stage(pipelinePump(done, out))
	}
	return out
}

func pipelinePump(done In, in In) Out {
	valueStream := make(Bi)
	go func() {
		defer func() {
			close(valueStream)
			for range in {
				// drain in channel
			}
		}()
		for {
			select {
			case <-done:
				return
			default:
			}
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
	}()
	return valueStream
}
