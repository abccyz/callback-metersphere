package net

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	accessKey = "CQ8wtOqcq2vpO7HG"
	secretKey = "IyvOqpu6eOJQiovi"
)

type HTTPClient struct {
	client  *http.Client
	BaseURL string
	Headers map[string]string
}

func NewCli(baseURL string) *HTTPClient {
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	comboxKey := fmt.Sprintf("%s|%s|%d", accessKey, uuid.New(), timeStamp)
	signature := aesEncrypt(comboxKey, secretKey, accessKey)
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"AccessKey":    accessKey,
		"Signature":    signature,
		"Connection":   "close",
	}
	return &HTTPClient{
		client:  &http.Client{},
		BaseURL: baseURL,
		Headers: headers,
	}
}

func (c *HTTPClient) setHeaders(req *http.Request) {
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
}

func (c *HTTPClient) DoRequest(method, path string, data []byte) (*http.Response, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) Get(path string) (string, error) {
	resp, err := c.DoRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *HTTPClient) Post(path string, data []byte) (string, error) {
	resp, err := c.DoRequest(http.MethodPost, path, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *HTTPClient) Put(path string, data []byte) (string, error) {
	resp, err := c.DoRequest(http.MethodPut, path, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *HTTPClient) Delete(path string) (string, error) {
	resp, err := c.DoRequest(http.MethodDelete, path, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 加密算法
func aesEncrypt(data string, key string, iv string) string {
	if len(key) != 16 {
		return errors.New("AES key must be 16 bytes").Error()
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err.Error()
	}
	// Use PKCS5Padding
	paddingSize := aes.BlockSize - (len(data) % aes.BlockSize)
	paddedData := append([]byte(data), bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)...)

	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, paddedData)
	// 进行Base64编码
	encodedString := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedString
}
