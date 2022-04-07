### Task 1 - Matching Behaviour

Take a look at the program [matching.go](../src/matching.go). Explain what happens and why it happens if you make the following changes. Try first to reason about it, and then test your hypothesis by changing and running the program.

  * What happens if you remove the `go-command` from the `Seek` call in the `main` function?

    If you remove the go-command the `Seek` commands will not occur/execute concurrently. This means that the names will be printed in order, i.e. `Anna` will always send to `bob`, and `Cody` will always send to `Dave`. This is because Anna will send a message, in the next iteration of the loop it is `Bobs` turn to either send or receive a message. However since there is a message in the channel `Bob` will receive it and so on.

  * What happens if you switch the declaration `wg := new(sync.WaitGroup`) to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?
    
    This will create a deadlock since the `wg` value passed to the `seek` funcition will not change the value in the `main` method, which will make the program wait forever.
    
  * What happens if you remove the buffer on the channel match?
    
    This will result in a deadlock, since no "message" or data can be stored in the channel, which means that the `seek` channel will have nothing to read. and yeah deadlock :( (Since the data can not be received since it is not saved anywhere) 

  * What happens if you remove the default-case from the case-statement in the `main` function?

    Since the case of ``name := <-match`` will always be true in this program (if no changes are made), since there is an odd amount of people. This means that nothing will change if the `default` case is removed.

Hint: Think about the order of the instructions and what happens with arrays of different lengths.

## Mapreduce and singleworker

|        | Avg time/run (ms)     |
|---            |---                    |
| Single Worker |   ~9 ms             |
| Map Reduce    |   ~5.5 ms            |