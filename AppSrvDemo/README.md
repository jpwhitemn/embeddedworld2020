# Lab 3 - Using the App Functions SDK

This lab will walk you through leveraging the App Functions SDK to build your own custom functions. 

In this lab you will learn how to do the following:
    
- Leverage Built-In Functions (ex, Filtering)
- Create a custom SDK Function
- Leverage the context provided to a function
- Publish data to a new ZMQ Topic

## Context

For this lab we will pretend that the values coming in as Int64 from the virtual device is temperature data and will assume it is being published in celsius. Our goal is to filter out other values from the device service, convert celsius to fahrenheit, and publish to a new ZMQ topic. 

For converting Celsius to Fahrenheit:
```golang
result := (t * 9 / 5 + 32)
```
where `t` is your temperature value

## Steps


1) Open `main.go` in the lab3 directory. Go ahead and run `go build` and you should have a successful build just to make sure things are working. 

2) As mentioned above, the first thing we'll want to do is look for Int64. To do this we can leverage the built-in SDK function called `FilterByDeviceName` and use `FilterByValueDescriptor`.
> Hint: You'll need to create a []string{} with the device name you wish to filter for. Check out the readme located here: https://github.com/edgexfoundry/app-functions-sdk-go 

3) After you've defined the filter and added it to your pipeline in the SetFunctionPipeline call, let's define our custom function. We've already created the function signature for you, all you have to do is fill out the body add it to your pipeline.
> Hint: Use the conversion formula above and pay attention to the annotations in the main.go.

4) After you have filled out the body of the `Fahrenheit()` function and added it to the pipeline. You can optionally publish the result to a new ZMQ Topic by calling `edgexcontext.Complete([]byte)`. The settings for where this will be published is located in the `/res/configuration.toml` file under the `[MessageBus]` section. Checkout the readme on the SDK for more information. 
> As a stretch goal, try setting up another instance of the SDK and set the `[SubscribeHost]` to the settings that this instance publishes too. This is how you can chain various Application Services together via ZMQ.

5) If you have completed this lab successfully, you should now be able to see logs with your converted data to the console. 
>TIP: ensure the logging level in the configuration correctly matches the level you used in your function.


