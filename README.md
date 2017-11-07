# GoWeb
Website Templating Engine

Stuff you should never need to edit:

configurator/ -> automatically loads and parses configuration file

multiLogger/ -> logger on steroids, dispatches logs to multiple destinations, as specified in the configuration file

rpcDb/ -> helper function to connect to DB, for calls we use meddler anyways

Webserver.go -> main file


Helper stuff:

config_file.txt -> example configuration file

exampleDb.sqlite3 -> example database file, created using populateDb.sql schema file

populateDb.sql -> sql schema file

requirements.sh -> requirements file, showing all dependencies as go get commands


This is where you develop:

templates/ -> all http templates that get parsed

rpcListener/ -> this is where all the actual magic happens

 BUILTINDispatcher -> router dispatcher, add your pages here, as per example

 BUILTIN* -> builtins to serve http(s), etc. Should never need to edit, apart from the Dispatcher mentioned already

 structs.go -> structs, session, user, builtin, etc. Only add new ones and edit User one. You should not touch the rest.

 index/login/register -> the 3 example pages created, so you can see how index can redirect to login, register can work, etc


To develop pages:

1. add function name to BUILTINDispatcher.go

2. add file and function in the rpcListener (as per examples)

3. add template in templates/ if needed (pay attention to {{define}} names)

4. enjoy


To run:

If you want to see a LOT of debug lines, use --early-debug in parameters.

DO provide configuration file name in parameter, otherwise it won't start.

By default, config file will make use of exampleDb.sqlite3, so you can see it in action.


Clean:

It's best to put any logging lines in BUILTINLogMessage or create a LogMessages.go in rpcListener/ and put them there.

See the rest of code for example. This is keeping text out of code. CLEAN!


Helper variables:

I will not describe it here. See index/login/register go files. They will describe it in full.

----

Roadmap (what it does not do, that I wish it did):

* session handler needs to allow for choice of backends, use caching DBs like aerospike with expiry builtin (no SessionCleaner)

* modulize the configurator, multiLogger, rpcDb, rpcListener, Webserver basic (.go)

  * so that we could just go get all this, open new file, do a few imports, run a few Init and write our templates
