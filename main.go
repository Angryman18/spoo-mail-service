package main

import (
	"bufio"
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection ", err)
			continue
		}

		server := Server{Conn: conn}

		go server.loop()
	}

}

func LoadEnv() error {
	return godotenv.Load(".env")
}

func Includes(d, cmd string) bool {
	return strings.Contains(d, cmd)
}

func (s *Server) loop() {
	// data := make([]byte, 10)
	conn := s.Conn
	for {
		conn.Write([]byte("220 WELCOME TO SMTP SERVER\r\n"))
		fmt.Println("Connected to ", conn.RemoteAddr().String()+"\n")
		reader := bufio.NewReader(conn)
		data, err := reader.ReadString('\n')
		// _, err := conn.Read(data)
		if err != nil {
			fmt.Println("Connection Closed", err)
			return
		}
		go s.handler(&data, reader)
	}
}

func (s *Server) handler(data *string, reader *bufio.Reader) {
	str := string(*data)
	fmt.Println("---> ", str)
	switch {
	case Includes(str, "HELO") || Includes(str, "helo"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "EHLO") || Includes(str, "ehlo"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "MAIL FROM"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "RCPT TO"):
		s.Conn.Write([]byte("250 OK\r\n"))
	case Includes(str, "DATA"):
		fmt.Println("DATA INVOKED ")
		s.Conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\r\n"))
		handleData(reader, s.Conn)
	case Includes(str, "QUIT"):
		s.Conn.Write([]byte("221 Bye\r\n"))
	default:
		s.Conn.Write([]byte("502 \r\n"))
		fmt.Println("502")
	}
}

func handleData(reader *bufio.Reader, conn net.Conn) {
	var data strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read email data: %v", err)
			return
		}

		if line == ".\r\n" {
			break
		}
		data.WriteString(line)
	}

	// Process the email data here
	log.Printf("Received email data:\n%s", data.String())
	conn.Write([]byte("250 OK\r\n"))
}
