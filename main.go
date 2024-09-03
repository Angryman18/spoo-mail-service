package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Server struct {
	Conn net.Conn
}

func main() {
	err := LoadEnv()
	if err != nil {
		panic("Error Loading Environment Variables")
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	connStr := fmt.Sprintf("%s:%s", host, port)

	listener, err := net.Listen("tcp", connStr)
	if err != nil {
		log.Fatal("Error Opening a TCP Conection ", err)
	}
	defer listener.Close()

	fmt.Printf("Server has started at %s\n", connStr)
	conn, err := listener.Accept()

	if err != nil {
		log.Fatal("Error Accepting Connection ", err)
	}

	server := Server{Conn: conn}

	server.loop()
	fmt.Println("Server has started")

}

func LoadEnv() error {
	return godotenv.Load(".env")
}

func Includes(d, cmd string) bool {
	return strings.Contains(d, cmd)
}

func (s *Server) loop() {
	data := make([]byte, 10)
	conn := s.Conn
	for {
		fmt.Println("Connected to ", conn.RemoteAddr().String()+"\n")
		_, err := conn.Read(data)
		conn.Write([]byte("WELCOME TO SMTP SERVER"))
		if err != nil {
			fmt.Println("Connection Closed")
			return
		}
		go s.handler(&data)
	}
}

func (s *Server) handler(data *[]byte) {
	str := string(*data)
	fmt.Println(str)
	switch {
	case Includes(str, "HELO"):
		s.Conn.Write([]byte("250 HELO\r\n"))
	case Includes(str, "MAIL FROM"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "RCPT TO"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "DATA"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "QUIT"):
		s.Conn.Write([]byte("221 Bye\r\n"))
	default:

	}
}
