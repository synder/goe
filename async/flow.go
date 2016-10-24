package async

import "sync"

func Parallel(tasks map[string]func()(interface{}, error)) map[string]interface{}{

	var wg sync.WaitGroup
	var rm sync.RWMutex

	var results = make(map[string]interface{})
	var errs error

	var handle = func(key string, task func()(interface{}, error)) {
		if temp, err := task(); err != nil {
			errs = err
		}else{
			rm.Lock()
			results[key] = temp
			rm.RUnlock()
		}
		wg.Done()
	}

	for key, task := range tasks {
		wg.Add(1)
		go handle(key, task)
	}

	wg.Wait()

	return results
}