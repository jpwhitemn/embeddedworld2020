package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

func main() {

	// 1) First thing to do is to create an instance of the EdgeX SDK, giving it a service key
	edgexSdk := &appsdk.AppFunctionsSDK{
		ServiceKey: "TemperatureConverter", // Key used by Registry (Aka Consul)
	}

	// 2) Next, we need to initialize the SDK
	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	// Since we are using MQTT, we'll also need to set up the addressable model to
	// configure it to send to our broker. If you don't have a broker setup you can pull one from docker i.e:
	// docker run -it -p 1883:1883 -p 9001:9001  eclipse-mosquitto
	addressable := models.Addressable{
		Address:   "m12.cloudmqtt.com",
		Port:      14073,
		Protocol:  "tcp",
		Publisher: "TestPusher",
		User:      "kcmvjbjy",
		Password:  "WNkiuH80yqNX",
		Topic:     "TestTopic",
	}

	// Using default settings, so not changing any fields in MqttConfig
	mqttConfig := transforms.MqttConfig{}
	mqttSender := transforms.NewMQTTSender(edgexSdk.LoggingClient, addressable, nil, mqttConfig, false)


	// 3) Since our FilterByDeviceName Function requires the list of Device Names we would
	// like to search for, we'll go ahead and define that now.
	deviceNames := []string{"Random-Integer-Generator01"}
	vdNames := []string{"RandomValue_Int8"}

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(deviceNames).FilterByDeviceName,
		transforms.NewFilter(vdNames).FilterByValueDescriptor,
		Fahrenheit,
		mqttSender.MQTTSend,
	); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK SetPipeline failed: %v\n", err))
		os.Exit(-1)
	}

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events to trigger the pipeline.
	edgexSdk.MakeItRun()
}

func Fahrenheit(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	// Check that there is a result to work with
	if len(params) < 1 {
		// We didn't receive a result
		return false, errors.New("No Data Received")
	}
	event := params[0].(models.Event)
	temperatureInCelsius,_ := strconv.Atoi(event.Readings[0].Value)
	edgexcontext.LoggingClient.Info("Incoming temp in C: " + event.Readings[0].Value)
	// Convert the value to Fahrenheit
	result := (temperatureInCelsius* (9/5) + 32)
	// Print out the value for debugging purposes using the logging on the context
	edgexcontext.LoggingClient.Info("Outgoing temp in F: " + strconv.Itoa(result))
	// Publish the result to the ZMQ Topic (Step 4)
	//edgexcontext.Complete([]byte(strconv.Itoa(result)))

	return true, strconv.Itoa(result)
}
