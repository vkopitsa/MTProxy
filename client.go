package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"net"
)

type Client struct {
	Conn              *net.TCPConn
	init              bool
	encryptKey_client []byte
	encryptIv_client  []byte
	decryptKey_client []byte
	decryptIv_client  []byte
	Server            *Server
	Network           *Network

	IdDc int16

	cipher_dec_client *Crypto
	cipher_enc_client *Crypto

	secret []byte
}

func NewClient(conn *net.TCPConn, network *Network, secret []byte) *Client {
	return &Client{
		Conn:    conn,
		Network: network,
		secret:  secret,
	}
}

func (c *Client) Do() {

	//64k buffer
	buff := make([]byte, 0xffff)
	for {
		n, err := c.Conn.Read(buff)
		if err != nil {
			if c.Conn != nil {
				c.Conn.Close()
			}

			c.Server = nil

			return
		}
		data := buff[:n]

		if !c.init {
			var err error
			data, err = c.GenerateAuthPacket(data)
			if err != nil {
				if c.Conn != nil {
					c.Conn.Close()
				}
				return
			}

			c.init = true
		}

		payload := c.cipher_dec_client.Do(data)

		if c.Server == nil {
			c.Server = c.Network.GetServer(c.IdDc)
			c.Server.Client = c

		}

		enc_payload := c.Server.cipher_enc_server.Do(payload)

		c.Server.rconn.Write(enc_payload)
	}
}

func (c *Client) GenerateAuthPacket(data []byte) ([]byte, error) {
	data_len := len(data)
	if data_len == 41 || data_len == 56 || data_len < 64 {
		return nil, errors.New("")
	}

	buf64 := GenerateRandomBytes2(64)
	copy(buf64[:], data[:])

	keyIv := GenerateRandomBytes2(48)
	copy(keyIv[:], buf64[8:])

	encryptKey_server := GenerateRandomBytes2(32)
	copy(encryptKey_server[:], keyIv[:])

	encryptIv_server := GenerateRandomBytes2(16)
	copy(encryptIv_server[:], keyIv[32:])

	reverseInplace2(&keyIv)

	decryptKey_server := GenerateRandomBytes2(32)
	copy(decryptKey_server[:], keyIv[:])

	decryptIv_server := GenerateRandomBytes2(16)
	copy(decryptIv_server[:], keyIv[32:])

	c.decryptKey_client = decryptKey_server
	c.decryptIv_client = decryptIv_server

	c.encryptKey_client = encryptKey_server
	c.encryptIv_client = encryptIv_server

	var decryptKey bytes.Buffer
	decryptKey.Write(decryptKey_server)
	decryptKey.Write(c.secret)

	sdfs := sha256.Sum256(decryptKey.Bytes())
	c.decryptKey_client = sdfs[:]

	var encryptKey bytes.Buffer
	encryptKey.Write(encryptKey_server)
	encryptKey.Write(c.secret)

	dfsdf := sha256.Sum256(encryptKey.Bytes())
	c.encryptKey_client = dfsdf[:]

	c.cipher_enc_client = NewCrypto(c.decryptKey_client, c.decryptIv_client)
	c.cipher_dec_client = NewCrypto(c.encryptKey_client, c.encryptIv_client)

	dec_auth_packet := c.cipher_dec_client.Do(buf64)

	dcId := abs(int16(binary.LittleEndian.Uint16(dec_auth_packet[60:]))) - 1

	for i := 0; i < 4; i++ {
		if dec_auth_packet[56+i] != 0xef {
			return nil, errors.New("")
		}
	}

	if dcId > 4 || dcId < 0 {
		return nil, errors.New("")
	}

	c.IdDc = dcId

	return data[64:len(data)], nil
}
