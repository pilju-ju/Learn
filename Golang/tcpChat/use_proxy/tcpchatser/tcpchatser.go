// tcpss
package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func main() {
	proxy, err := net.Dial("tcp", "127.0.0.1:1200")
	if err != nil {
		return
	}

	var port string
	fmt.Println("만들고 싶은 채팅방의 주소를 입력하시오")
	fmt.Scanln(&port)

	proxy.Write([]byte(port + "\n"))

	chatRoom(port)
	proxy.Close()
}
func overlapCheck(nameList map[string]bool, name string) error {
	if _, ok := nameList[name]; ok {
		return errors.New("닉네임이 중복되었습니다.")
	}
	nameList[name] = true
	return nil
}

func chatRoom(port string) {
	fmt.Println(" 채팅방이 만들어졌습니다.")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Listen err", err)
		return
	}
	defer ln.Close()

	allUser := make(map[net.Conn]int) // 모든 유저를 구분하기 위한 맵
	userName := make(map[int]string)  // 유저 이름 저장하는 맵
	inUser := make(chan net.Conn)     //유저의 접속을 확인하는 채널
	message := make(chan string)      // 메세지 채널
	outUser := make(chan net.Conn)    // 유저의 퇴장을 확인하는 채널
	nameList := make(map[string]bool) //닉네임 중복을 검사할 때 사용하는 map
	nameLength := make(chan int)      //  닉네임 길이를 전달하는 채널

	go func() { //고루틴을 이용하여 클라이언트 연결을 동시에 받음
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			defer ln.Close()
			fmt.Println("접속")
			inUser <- conn // 클라이언트가 접속되면 Inuser에 값을 전달
		}
	}()
	countUser := 0 // 유저구분을 위한 변수

	for { // 채널 같은 경우 switch가 아닌 select문 이용하여 case를 사용
		select {
		case conn := <-inUser: // Inuser의 값이 들어오면
			countUser++                     //유저수 판단
			go func(conn net.Conn, i int) { // 들어온 유저로 고루틴 실행
				c := bufio.NewReader(conn)           // 받은 데이터를 io.Reader를 만족하는 bufio.Reader로 선언
				nameSize, _ := binary.ReadUvarint(c) // 닉네임 길이 받음
				name := make([]byte, nameSize)       // 닉네임을 담을 버퍼
				c.Read(name)                         // 닉네임을 읽음
				userName[i] = string(name)           // 유저 이름을 저장
				allUser[conn] = i                    // 유저를 구분

				if err := overlapCheck(nameList, string(name)); err != nil { // 닉네임 중복체크에서 에러가 생기면
					fmt.Println(err.Error())            //에러 메시지 출력
					delete(userName, i)                 // 해당 유저를 유저 이름 map에서 삭제
					delete(allUser, conn)               //해당 유저를 모든 유저를 저장한 맵에서 삭제
					conn.Write([]byte("*overlapErr\n")) //유저들에게 닉네임 중복오류 전송
				} else {
					fmt.Println("닉네임 ", string(name), " 님이 입장하셨습니다") // 입장 확인 메세지
					nameLength <- len(name)                          //메시지를 보내기 전, 보내는 사람의 닉네임 길이를 전달
					message <- userName[i] + " 님이 입장하셨습니다" + "---\n" // 입장 메세지를 클라이언트에게 전달
					for {
						m, err := c.ReadString('\n') //ReadString을 이용하여 받은 데이터에서 메세지를 읽음
						if m == "*ESC\n" {           // *ESC를 입력 받은 경우, 유저가 나간 것을 전달하기 위해 for문을 종료
							break
						}
						if err != nil {
							fmt.Println("비정상적인 종료가 감지 ")
							break
						}
						fmt.Print(userName[i], " : ", m) // 닉네임과 함께 출력
						nameLength <- len(name)
						message <- userName[i] + " : " + m // 메세지를 전달
					}
					outUser <- conn // 유저가 나간 것을 판단하기 위하여 outUser 채널에 값을 전달
					nameLength <- len(name)
					message <- userName[allUser[conn]] + " 님이 퇴장하셨습니다--- \n" // 퇴장메세지 출력
				}
			}(conn, countUser)
		case nameLen := <-nameLength: //닉네임 길이를 전달 받으면
			nickNameSize := make([]byte, binary.MaxVarintLen64)   //닉네임 길이를 저장할 []byte 생
			n := binary.PutUvarint(nickNameSize, uint64(nameLen)) // 닉네임 길이를 Uvarint인코딩
			for conn, _ := range allUser {                        // 메시지를 보내기 전 모든 유저에게 보내는 유저의 닉네임 길이를 보냄
				conn.Write(nickNameSize[:n])
			}

		case msg := <-message: //메세지가 들어오면
			for conn, _ := range allUser { // 모든 유저에게 메세지를 전달 하기위한 루프
				conn.Write([]byte(msg)) //메세지 전달
			}

		case out := <-outUser: // 유저가 *ESC를 입력하고 나가면
			fmt.Println(userName[allUser[out]], "님이 퇴장하셨습니다.")
			delete(userName, allUser[out])
			delete(allUser, out) // 모든 유저를 저장한 맵에서 나간 유저 삭제
		}
	}

}
