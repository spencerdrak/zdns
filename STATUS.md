## Current Status

The current status of ZDNS in this repo is that the main core of the logic (e.g., the library refactoring) has been completed. The main logic has also been manually tested, however, there are a variety of TODO items in the code that may need to be addressed. These are marked as `TODO(spencer)` to delineate from other issues in the codebase. 

A sample program has been created for this as well, to demonstrate how modules would use the core library. This sample can be found at this remote: `github.com:spencerdrak/sample-zdns-client.git`. Please let me know if there's any questions about how this library works or how the sample client uses the library. 

The main bulk of the work to be completed is to migrate the modules (A, AAAA, MX, MXLookup, ALookup, etc) to use the new core library. The logic in each of these modules can remain mostly the same - however, the semantics of how the library is called needs to change. In setting up the sample client, I found that there is a slight bit more machinery around looping/passing all the inputs to the core library, but the main logic was much simpler to follow using the new library. If you find that this is not the case, I am happy to discuss - it may simply make sense to me since I wrote it, but I'm happy adjust!

The documents in this repo named `PROPOSED_REFACTORING.md` and `PROPOSED_ZDNS_INTERFACE.md` will be useful as docs for this - the latter is the most useful. Again, let me know if there are any questions!