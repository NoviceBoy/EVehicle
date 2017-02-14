package email

import (
    "gopkg.in/gomail.v2"
)

const EmailSubject = "Vehicle Charge Point Notification"
var EmailBody = "Vehicle Charging Completed, Please move the Vehicle"

func SendEmail( user string, licensePlate string ){
    m := gomail.NewMessage()
    m.SetAddressHeader("From", "nobody@arista.com", "EV Charge Admin")
    m.SetAddressHeader("To", user, "" )
    m.SetHeader("Subject", EmailSubject )
    m.SetBody("text/plain", licensePlate+ " " + EmailBody )

    //d := gomail.NewPlainDialer("smtp.gmail.com", 465, "evchargepoint@gmail.com" , "")

    d := gomail.NewPlainDialer("127.0.0.1", 25, "nobody@aristanetworks.com" , "")
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}