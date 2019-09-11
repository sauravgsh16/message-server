Progress:

Failing when closing connection/channel. Both on server and client.
Assumption is that somewhere, I am closing something before it should. And then attempting to talk on that connection.

When connection close is called, without channel close, server breaks
    Assumption is that server tries to execute things on the closed channel
When channel close is called, client breaks.
    NEED TO FIND ROOT CAUSE.