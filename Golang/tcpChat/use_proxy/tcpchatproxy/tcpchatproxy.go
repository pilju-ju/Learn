// tcpps
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":1200")
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	defer ln.Close()

	allServer := make(map[net.Conn]int)   // 모든 서버를 구분하기 위한 맵
	serverAddress := make(map[int]string) //서버의 주소를 관리하기 위한 map
	inUser := make(chan net.Conn)         //유저(서버/클라이언트) 의 접속을 확인하는 채널
	inClient := make(chan net.Conn)       // 클라이언트 접속 신호 채널
	outUser := make(chan net.Conn)        // 유저의 퇴장을 확인하는 채널

	go func() { //고루틴을 이용하여 클라이언트 연결을 동시에 받음
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			defer conn.Close()
			fmt.Println("접속")
			inUser <- conn // 클라이언트가 접속되면 Inuser에 값을 전달
		}
	}()
	count := 0

	for {
		select {
		case conn := <-inUser: //유저(서버/클라이언트)가 접속하면
			count++ //카운트 변수
			go func(conn net.Conn, i int) {
				c := bufio.NewReader(conn) //conn을 io.Reader로 변환
				for {
					list, err := c.ReadString('\n')
					if err != nil {
						break
					}

					if list == "client\n" { //받은 데이터가 클라이언트일 경우
						fmt.Println("클라이언트")
						if len(serverAddress) == 0 { //서버주소를 관리하는 map이 빈 경우 == 개설된 채팅방이 없을 경우
							conn.Write([]byte("*NOROOM\n")) // 클라이언트로 오류 메시지 전송
							break
						}
						inClient <- conn //클라이언트가 접속했다는 신호 전달
						break
					}
					allServer[conn] = i     //서버 접속을 관리하는 allServer map에 추가
					serverAddress[i] = list // 서버 주소를 관리하는 serverAddress map에 추가
					fmt.Println("서버", serverAddress[i])
				}
				outUser <- conn
			}(conn, count)

		case client := <-inClient: //신호가  들어오면
			for _, address := range serverAddress { //serverAddress map을 순회하며
				client.Write([]byte(address)) // 서버 주소를 전달
			}
			client.Write([]byte("END\n")) //서버 주소를 다 보냈다는 신호 전달

		case out := <-outUser: //유저가 나갔다고 들어오면
			fmt.Println("서버종료", serverAddress[allServer[out]])
			delete(serverAddress, allServer[out])
			delete(allServer, out)

			if len(allServer) == 0 { // 개설중인 서버가 없을 경우
				fmt.Println("개설중인 서버 0개")
				return
			}
		}
	}
}
