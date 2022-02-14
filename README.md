# What is GoJo

GoJo is a Go library created during my bachelor thesis. 
The library implements Join Patterns as described by Cédric Fournet and Georges Gonthier [1].  

Join Patterns allow a declarative description of concurrent programs.

## Definitions

- A `Join Pattern` listens on messages on up to 3 `Ports` until the `Action` is executed.
- `Pattern Matching` is the process of matching messages to a specific `Join Pattern`.
- An `Action` is a function that will be executed once a pattern was matched.
- A `Port` is a Go channel that listens for signals used for `Pattern Matching`.
- A `Signal` is sending a message to a `Port`. A `Signal` can either be `Async` or `Sync`. 
  - `Async` signals are non-blocking and return immediately. 
  - `Sync` signals are blocking and will return a value. An `Action` sends it's return value to all registered `Sync` signals.
- A `Junction` is responsible to keep track of join patterns and registered `Signals`. A program can consist of multiple `Junctions`.

## Requirements
Go version 1.18 is required to use this library due to the usage of generics.

## Usage

To create new Join Patterns and Signals a Junction needs to be created first.

```
j := junction.NewJunction()
```

Afterwards, `AsyncSignals` and `SyncSignals` can be created through this junction. These functions return two values, a port and a signal. 
The port is used to register the signals with Join Patterns and signal is the function that can be called to send and receive values from the Join Pattern's action.

```
port, signal := junction.NewAsyncSignal[<TYPE>](j)
```

where `<TYPE>` is the data type of the value sent via the signal.


or

```
port, signal := junction.NewSyncSignal[<SEND_TYPE>,<RECV_TYPE>](j)
```

where `<SEND_TYPE>` is the data type of the Signal to the Join Pattern and `<RECV_TYPE>` is the data type that will be returned to the signal.

Once Signals are created, they can be used to create new Join Patterns. Join Patterns can listen on up to 3 signals and are not limited to only one synchronous port.

```
junction.NewBinaryAsyncJoinPattern[<FIRST_TYPE>,<SECOND_TYPE>](portA, portB).Action(func(a <FIRST_TYPE>, b <SECOND_TYPE>) {
  // Add function code here
})
```

In this example a binary async Join Pattern is created. This means that there exists now a function that will be executed in it's own goroutine once a message was received on both `portA` and `portB`.
Once messages are received on both ports, the function is executed and receives both messages as the function parameters.

To find more detailed examples, look into the [`examples`](https://github.com/GregorKonzett/GoJo/tree/master/examples) folder.

## Special Thanks
I want to thank my project's supervisor [Dr. Ian Stark](https://homepages.inf.ed.ac.uk/stark/). I would not have been able to complete this project without his guidance and help.

## References

[1] https://www.microsoft.com/en-us/research/wp-content/uploads/2017/01/join-tutorial.pdf<br />