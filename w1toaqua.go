package main

import(
  "fmt"
  "io/ioutil"
  "log"
  "os/exec" 
  "strconv"
  "strings"
  "time"

  "github.com/dustin/go-coap"

  )


func check(e error){
  if e != nil {
    panic(e)
  }
}

func loadw1kos(){
  cmd1 := exec.Command("modprobe", "w1-gpio")
  err1 := cmd1.Run()
  check(err1)
  cmd2 := exec.Command("modprobe", "w1-therm")
  err2 := cmd2.Run()
  check(err2)

}


const devicesPath = "/sys/bus/w1/devices/"


func main(){

  for{

    var deviceIDs []string

    files, err := ioutil.ReadDir(devicesPath)
    check(err)
    for _, file := range files {
	if strings.HasPrefix(file.Name(), "28-"){
	  deviceIDs = append(deviceIDs, file.Name())
        }
    }

    for _, deviceID := range deviceIDs {
      data, err := ioutil.ReadFile(devicesPath + "/" + deviceID + "/w1_slave")
      check(err)

      sample := string(data)

      if strings.Contains(sample, "NO"){
        fmt.Print(string(data))
      } else {
        pos := strings.Index(sample, "t=")
        if pos >= 0 {
	  value := sample[pos+2:pos+5]
	  f, err := strconv.ParseFloat(value, 64)
	  check(err)
	  celsius := f * 0.1

	  temp := strconv.FormatFloat(celsius, 'f', 1, 64)
          fmt.Println(temp)
         sendMeasure(deviceID, "temperature", temp)
  } else {
          fmt.Print("cannot read data")
        }
      }
    }

    time.Sleep(time.Second)
  }
}

func sendMeasure(deviceID string, measure string, value string){

  req := coap.Message{
      Type:      coap.NonConfirmable,
      Code:      coap.POST,
      MessageID: 12345,
      Payload:   []byte(value),
  }

  path := "/aqua/" + deviceID + "/" + measure + "/"

  req.SetOption(coap.ETag, "weetag")
  req.SetOption(coap.MaxAge, 3)
  req.SetPathString(path)

  c, err := coap.Dial("udp", "192.168.0.45:5683")
  if err != nil {
	log.Fatal("Error dialing: %v", err)
  }

  rv, err := c.Send(req)
  if err != nil {
    log.Fatal("Error sending request: %v", err)
  }

  if rv != nil {
    log.Print("Response payload: %s", rv.Payload)
  }
}
