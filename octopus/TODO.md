# Big todo list for exchange

 - Show job for storing in broker is hard to handle. in show time we have just a bunch of ids 
 - Must add each data into influx(or whatever time-series database we can handle) too. the data for that is very smaller than this.
 
 
# Things we need to consider

 - Never use start on an interface. interface is the generic type and using star on that normally is not correct
 - Never import and package for an interface (or multiple interface) go interface is some how is a way to prevent this.
   simply create your own interface with same methods (minimum methods that you need)
 - Global variable make your code hard to test
 
Fuck the code. 
 