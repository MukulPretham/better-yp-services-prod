package main

import (
	"fmt"
	"net/smtp"
	// "net/smtp"
	"os"
)

func SendMain(toAddress []string, message string)error{
	auth := smtp.PlainAuth("",os.Getenv("FromEmail"),os.Getenv("AppPassword"),os.Getenv("SmptHost"))
	fmt.Print(os.Getenv("AppPassword"))
	err := smtp.SendMail(os.Getenv("SmptIP"),auth,os.Getenv("FromEmail"),toAddress,[]byte(message))
	if err != nil{
		return fmt.Errorf("%v",err)
	}
	return nil
}