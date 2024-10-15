/*
Package restasks provides a [Server] to handle asynchronous tasks with a RES service.

The [Server] is designed to be used as a standalone service or as part of a larger application.

You can create a server by providing a cache implementation and a RES service:

	server := restasks.NewServer(
		cache.NewRistretto(cache.MediumRistrettoCache),	// Create a cache or use an existing one.
		res.NewService("test"),							// Create a RES service or use an existing one.
	)

Then, you can manage tasks using the [Server].
Let's imagine that we have a RES handler that does some heavy work and we want to run it asynchronously:

	func (request res.CallRequest) {
		// Create a new task.
		taskRID, err := server.CreateTask()
		if err != nil {
			request.Error(err)
			return
		}

		// Run the heavy work in a goroutine.
		go func() {
			// Do the heavy work.
			result, err := doHeavyWork()
			if err != nil {
				server.FailTask(taskRID, err) // Fail the task if an error occurs.
				return
			}

			server.CompleteTask(taskRID, result) // Complete the task with the result.
		}

		// You can directly return the task RID to the client so it can subscribe to the task RES resource.
		return request.Resource(taskRID)
	}
*/
package restasks
