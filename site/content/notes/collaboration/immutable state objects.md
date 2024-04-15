---
title: immutable state objects
type: note
---


The state of data is *never* modified in place. To change it, we pass the current state of it through an `Automerge.change` function that returns a new object where that change is reflected along side a possible commit message that is with the change

```js
// pesudo code
state = Automerge.change(state, "Some kind of message") {
  (doc) => {
	  doc.todo.push({
		  "title": "Dry Laundry",
		  "done": false,
	  })
  }
}
```

the `doc` callback here is mutable only within the callback. 
