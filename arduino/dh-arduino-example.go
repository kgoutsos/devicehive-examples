package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
	"github.com/godbus/dbus"
)

type dbusWrapper struct {
	conn         *dbus.Conn
	path, iface  string
	handlers     map[string]signalHandlerFunc
	handlersSync sync.Mutex

	scanMap   map[string]bool

	readingsBuffer map[string]*[]float64

	deviceLock sync.Mutex
}

type signalHandlerFunc func(args ...interface{})
type cloudCommandHandler func(map[string]interface{}) (map[string]interface{}, error)

func NewdbusWrapper(path string, iface string) (*dbusWrapper, error) {
	d := new(dbusWrapper)

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Panic(err)
	}

	d.handlers = make(map[string]signalHandlerFunc)

	d.conn = conn
	d.path = path
	d.iface = iface
	d.readingsBuffer = make(map[string]*[]float64)
	d.scanMap = make(map[string]bool)

	filter := fmt.Sprintf("type='signal',path='%[1]s',interface='%[2]s',sender='%[2]s'", path, iface)
	log.Printf("Filter: %s", filter)

	conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus").Call("org.freedesktop.DBus.AddMatch", 0, filter)

	go func() {
		ch := make(chan *dbus.Signal, 1)
		conn.Signal(ch)
		for signal := range ch {
			if !((strings.Index(signal.Name, iface) == 0) && (string(signal.Path) == path)) {
				continue
			}
			if handler, ok := d.handlers[signal.Name]; ok {
				go handler(signal.Body...)
			} else {
				log.Printf("Unhandled signal: %s", signal.Name)
			}
		}
	}()

	return d, nil
}

func (d *dbusWrapper) call(name string, args ...interface{}) *dbus.Call {
	c := d.conn.Object(d.iface, dbus.ObjectPath(d.path)).Call(d.iface+"."+name, 0, args...)

	if c.Err != nil {
		log.Printf("Error calling %s: %s", name, c.Err)
	}

	return c
}

func (d *dbusWrapper) SendNotification(name string, parameters interface{}, priority uint64) {
	b, _ := json.Marshal(parameters)
	d.call("SendNotification", name, string(b), priority)
}

func (d *dbusWrapper) RegisterHandler(signal string, h signalHandlerFunc) {
	d.handlersSync.Lock()
	d.handlers[d.iface+"."+signal] = h
	d.handlersSync.Unlock()
}

func (d *dbusWrapper) CloudUpdateCommand(id uint32, status string, result map[string]interface{}) {
	b, _ := json.Marshal(result)
	d.call("UpdateCommand", id, status, string(b))
}

func main() {
	log.Printf("Establishing the serial connection...")
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second) // needed to give Arduino some time to initialise as it is reset when the port is opened

	cloud, err := NewdbusWrapper("/com/devicehive/cloud", "com.devicehive.cloud")
	if err != nil {
		log.Panic(err)
	}

	cloudHandlers := make(map[string]cloudCommandHandler)

	cloudHandlers["led/on"] = func(p map[string]interface{}) (map[string]interface{}, error) {
		log.Printf("Turning the LED on...")
		_, err := s.Write([]byte("1"))
		return map[string]interface{}{"led": "on"}, err
	}

	cloudHandlers["led/off"] = func(p map[string]interface{}) (map[string]interface{}, error) {
		log.Printf("Turning the LED off...")
		_, err := s.Write([]byte("0"))
		return map[string]interface{}{"led": "off"}, err
	}

	cloud.RegisterHandler("CommandReceived", func(args ...interface{}) {
		id := args[0].(uint32)
		command := args[1].(string)
		params := args[2].(string)

		log.Printf("Got command %d with name \"%s\" and parameters %s",id,command,params)
		var param_data map[string]interface{}
		b := []byte(params)
		json.Unmarshal(b, &param_data)

		//log.Printf("param_data %s", param_data)

		if h, ok := cloudHandlers[command]; ok {
			//At this point the client notifies the cloud about the success/failure of the command
			res, err := h(param_data)

			if err != nil {
				cloud.CloudUpdateCommand(id, fmt.Sprintf("ERROR: %s", err.Error()), nil)
			} else {
				cloud.CloudUpdateCommand(id, "success", res)
			}

		} else {
			log.Printf("Unhandled command: %s", command)
		}
	})

	select{}
}
