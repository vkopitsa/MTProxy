package main

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Adress            string
	rconn             *net.TCPConn
	Client            *Client
	encryptKey_server []byte
	encryptIv_server  []byte
	decryptKey_server []byte
	decryptIv_server  []byte
	cipher_dec_server *Crypto
	cipher_enc_server *Crypto
	IdDc              int16
	Disabled          bool
}

func NewServer(adress string, idDc int16) *Server {
	return &Server{
		Adress:   adress,
		Disabled: false,
		IdDc:     idDc,
	}
}

func (s *Server) Run() error {
	raddr, err := net.ResolveTCPAddr("tcp", s.Adress)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.rconn, err = net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.rconn.SetKeepAlive(true)
	s.rconn.SetKeepAlivePeriod(time.Minute * 5)

	err = s.makeAuthKey()
	if err != nil {
		return err
	}

	go s.Read()

	return nil
}

func (s *Server) Read() {
	//64k buffer
	buff := make([]byte, 0xffff)
	for {
		n, err := s.rconn.Read(buff)
		if s.Client == nil {
			s.rconn.Close()
		}
		if err != nil {
			s.Client.Server = nil
			if s.rconn != nil {
				s.rconn.Close()
			}
			return
		}
		b := buff[:n]

		dec_packet := s.cipher_dec_server.Do(b)
		enc_packet := s.Client.cipher_enc_client.Do(dec_packet)

		s.Client.Conn.Write(enc_packet)
	}
}

func (s *Server) makeAuthKey() error {
	random_buf, _ := GenerateRandomBytes(64)

	for {
		var val int = int((random_buf[3] << 24) | (random_buf[2] << 16) | (random_buf[1] << 8) | (random_buf[0]))
		var val2 int = int((random_buf[7] << 24) | (random_buf[6] << 16) | (random_buf[5] << 8) | (random_buf[4]))
		if random_buf[0] != 0xef &&
			val != 0x44414548 &&
			val != 0x54534f50 &&
			val != 0x20544547 &&
			val != 0x4954504f &&
			val != 0xeeeeeeee &&
			val2 != 0x00000000 {

			random_buf[56] = 0xef
			random_buf[57] = 0xef
			random_buf[58] = 0xef
			random_buf[59] = 0xef
			break
		}
		random_buf, _ = GenerateRandomBytes(64)
	}

	keyIv := GenerateRandomBytes2(48)
	copy(keyIv[:], random_buf[8:])

	encryptKey_server := GenerateRandomBytes2(32)
	copy(encryptKey_server[:], keyIv[:])

	encryptIv_server := GenerateRandomBytes2(16)
	copy(encryptIv_server[:], keyIv[32:])

	reverseInplace2(&keyIv)

	decryptKey_server := GenerateRandomBytes2(32)
	copy(decryptKey_server[:], keyIv[:])

	decryptIv_server := GenerateRandomBytes2(16)
	copy(decryptIv_server[:], keyIv[32:])

	s.decryptKey_server = decryptKey_server
	s.decryptIv_server = decryptIv_server

	s.encryptKey_server = encryptKey_server
	s.encryptIv_server = encryptIv_server

	s.cipher_enc_server = NewCrypto(s.encryptKey_server, s.encryptIv_server)
	s.cipher_dec_server = NewCrypto(s.decryptKey_server, s.decryptIv_server)

	packet_enc := s.cipher_enc_server.Do(random_buf)

	copy(packet_enc[:], random_buf[:56])

	s.rconn.Write(packet_enc)

	return nil
}
