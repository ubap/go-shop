# backend
```bash
go run .
```
# frontend
Dev:
```bash
npm run dev
```
Test:
```bash
npm test
```



## Notes
* Front end communicated with backend over websocket.
* Item is identified by its name, case-insensitive.
* Goal: PWA, so has to support offline mode
* The work made offline should sync after getting back oline


### How to support item rename?
#### If the item is identified by its name 
This can be problematic when the item name is updated in offline mode.
It might be difficult to keep track of what item is what.
Unless we allow changing item name only when connected to the server,
and the change would be server-first (meaning that the local state is not updated until server doesn't confirm it)

Then, the clients that are offline, when getting online would receive the new name, update the name locally, 
and only then send updates made offline - already using the new name.

#### If the item is identified by id
It should be no problem, because we can still identify the item, no matter the name.

### How to add new item when offline?
The issue is with id, if two offline clients add new items, how to resolve it when they come online?
Ideally the client would already have an id.
#### Use unique ID for each item, generate it on the frontend.
What if two clients add the same item when offline? - The item will be duplicated.