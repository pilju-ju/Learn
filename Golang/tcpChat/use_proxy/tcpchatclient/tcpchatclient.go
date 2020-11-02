// tcpcc
package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	var port string
	var nickName string
	proxy, err := net.Dial("tcp", "127.0.0.1:1200")
	if err != nil {
		return
	}

	if err := readList(proxy); err != nil { //readList()에서 err가 발생하면 종료
		fmt.Println(err.Error()) // 에러메시지 출력
		return
	}

	fmt.Println("입장할 채팅방 주소를 입력하세요")
	fmt.Scanln(&port)

	fmt.Println("사용할 닉네임을 입력하세요 ")
	fmt.Scanln(&nickName)
	myNickName := []byte(nickName)

	nickNameSize := make([]byte, binary.MaxVarintLen64)           //닉네임 길이를 저장할 []byte
	n := binary.PutUvarint(nickNameSize, uint64(len(myNickName))) // 닉네임 길이를 Uvarint인코딩
	chat(port, nickNameSize[:n], myNickName)

}
func readList(proxy net.Conn) error {
	defer proxy.Close()

	proxy.Write([]byte("client\n")) //프록시 서버에 접속한 것이 클라이언트임을 알 수 있도록 전달
	reader := bufio.NewReader(proxy)
	fmt.Println("현재 개설중인 채팅방 주소")
	for {
		list, err := reader.ReadString('\n')
		if err != nil { //비정상적인 종료시
			return errors.New("비정상적인 종료")

		}
		if list == "END\n" { // 서버 목록을 모두 받았다는 신호가 왔을 경우
			return nil
		}
		if list == "*NOROOM\n" { //개설된 채팅 서버가 없다는 오류 메시지가 왔을 경우 에러처리
			return errors.New("개설된 채팅서버가 없습니다. 프로그램을 종료합니다.")
		}
		fmt.Print(list)
	}
	return nil
}
func chat(port string, nickNameSize []byte, myNickName []byte) {
	client, err := net.Dial("tcp", port) // 입력받은 port 클라이언트 dial
	if err != nil {
		fmt.Println("Dial err", err)
		return
	}
	defer client.Close()

	client.Write(nickNameSize) // 리턴 받은 바이트 수를 이용하여 닉네임  길이를 보냄
	client.Write(myNickName)   // 닉네임 보내기

	done := make(chan struct{}) //빈 struct를 이용하여 고루틴이 종료되었음을 전달

	go func(client net.Conn) { //서버로부터 채팅 내용을 받아오는 고루틴, 닉네임 중복에러 발생 시 종료
		defer close(done) //닉네임 에러 발생시 done채널을 닫아 종료되었다는 신호를 보냄
		c := bufio.NewReader(client)
		for {
			nameSize, _ := binary.ReadUvarint(c) // 닉네임 길이 받음

			line, isPrefix, err := c.ReadLine()
			if isPrefix || err != nil {
				return
			}

			if int(nameSize) == len(myNickName) { //받은 닉네임 길이와 내 닉네임 길이가 같을 때,
				replaceMyName := strings.Replace(string(line), string(myNickName), "나", 1) //받은 메세지가 내가 보낸 메지시 일 때, 채팅창의 닉네임 부분을 "나" 로 수정
				fmt.Println(replaceMyName)                                                 // 닉네임은 다르지만 길이가 같을 경우, replace가 적용안됨
			} else if int(nameSize) != len(myNickName) {
				fmt.Println(string(line)) // 출력
			}

			if string(line) == "*overlapErr" { //닉네임 중복 에러를 서버로 부터 받을 시,
				fmt.Println(" 닉네임 중복 에러, 엔터를 눌러 다시 접속해주세요") //에러 메세지 출력, 엔터를 한 번 눌러서 기존의 입력 과정을 한번 수행하고 종료되도록 유도
				return                                      // for문 밖으로 나감
			}
		}
	}(client)

	for {
		select {
		case <-done: //닉네임 에러 발생 시 입력을 받지 않고 종료
			return
		default:
			reader := bufio.NewReader(os.Stdin)  // 콘솔 입력(표준입력)을 io.Reader를 따르는 reader로 만듬
			text, err := reader.ReadString('\n') // 개행(엔터)를 기준으로 입력을 문자열로 읽음
			if err != nil {
				return
			}
			client.Write([]byte(text)) // 읽은 문자열을 서버로 전송
			if text == "*ESC\n" {      //채팅방을 나가기 위해 *ESC를 입력할 경우
				fmt.Println("채팅방을 나갔습니다\n") //종료 메세지 출력
				return
			}
		}
	}
}
