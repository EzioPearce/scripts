package main

import (
    "bytes"
    "encoding/csv"
    "io"
    "os"
    "crypto/tls"
    "errors"
    "fmt"
    "net"
    "net/smtp"
    "time"
)

type loginAuth struct {
    username, password string
}

func LoginAuth(username, password string) smtp.Auth {
    return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
    return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
    if more {
        switch string(fromServer) {
        case "Username:":
            return []byte(a.username), nil
        case "Password:":
            return []byte(a.password), nil
        default:
            return nil, errors.New("Unknown from server")
        }
    }
    return nil, nil
}

func main() {

    // Sender data.
    username := "<Enter sender email>"
    password := "<Enter sender email password here>"

    //smtpFrom := "VocuniHelp <hello@vocuni.com>"

    // smtp server configuration.
    smtpHost := "smtp.office365.com"
    smtpPort := "587"

    conn, err := net.Dial("tcp", "smtp.office365.com:587")
    if err != nil {
        println(err)
    }

    c, err := smtp.NewClient(conn, smtpHost)
    if err != nil {
        println(err)
    }

    tlsconfig := &tls.Config{
        ServerName: smtpHost,
    }

    if err = c.StartTLS(tlsconfig); err != nil {
        println(err)
    }

    auth := LoginAuth(username, password)

    if err = c.Auth(auth); err != nil {
        println(err)
    }

    file , _ := os.Open( "<Enter location of the csv file with the reciepients name and email address>")
    	r := csv.NewReader(file)

    	/*records, err := r.ReadAll()
    	if err != nil {
    		log.Fatal(err)
    	}*/
    	r.Read()

      mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
      smtpFrom := []byte(fmt.Sprintf("From: Vocuni Help <%s>\n", username))
      counter := 0
      for {
    		record, err := r.Read()
    		if err == io.EOF {
    			break
    		}
    		name := record[0]
        email := record[1]
    		//ln := record[1]
    		//un := record[2]

    var smtpbody bytes.Buffer
    body:= bytes.Buffer{}
    fmt.Fprint(&body, "Hello " , name , "<br> Enter the email body here. The HTML break tag is used to create a newline charecter")

    smtpbody.Write(smtpFrom)
    smtpbody.Write([]byte(fmt.Sprintf("To: %s <%s>\n", name, email)))
    smtpbody.Write([]byte(fmt.Sprintf("Subject: Enter the Email Subject \n%s\n\n", mimeHeaders)))
    smtpbody.Write(body.Bytes())

    // Sending email.
    err = smtp.SendMail(smtpHost+":"+smtpPort, auth, username, []string {email}, smtpbody.Bytes())
    if err != nil {
        fmt.Println(err, email)
        continue
    }
    fmt.Println("Email Sent!" , email)
    counter++
    if(counter == 100) {
      time.Sleep(30*time.Second)
      counter = 0
    }
  }
}
