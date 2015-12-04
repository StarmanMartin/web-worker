package worker

import (
)

var ctasks = make(chan RequestTask, 20)

type Worker struct {
	MaxWorker        int
	MinWaitingWorker int
}

type RequestTask func()

func (w Worker) work(ctasks chan RequestTask, cwaiting chan taskWorker) {
	for {
		task := <-ctasks
		for i := len(cwaiting); i <= w.MinWaitingWorker; i++ {
			go func() {
				cwaiting <- taskWorker{}
				w.work(ctasks, cwaiting)
			}()
		}

		worker := <-cwaiting
		worker.doRun(task)
		select {
		case cwaiting <- worker:
			break
		default:
			return
		}

	}
}

func (w Worker) RunWorker() {
	cwaiting := make(chan taskWorker, w.MaxWorker)
	for i := 0; i < w.MinWaitingWorker; i++ {
		go func() {
			cwaiting <- taskWorker{}
			w.work(ctasks, cwaiting)
		}()
	}
}

func (w Worker) Pushtask(r RequestTask) {
	ctasks <- r
}

type taskWorker struct{}

func (w taskWorker) doRun(task RequestTask) {
	task()
}
