package objstorage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"time"
)

func HttpPut(url string, body io.Reader) error {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http [%s] %s", resp.Status, data)
	}
	if len(data) > 0 {
		log.Println("[http body]", string(data))
	}
	return nil
}

func Md5(p []byte) string {
	t := md5.Sum(p)
	return hex.EncodeToString(t[:])
}

func Main() {
	ctx := context.Background()
	c, err := NewController(&minioImpl{}, NewKV())
	if err != nil {
		log.Fatalln(err)
	}

	name := "hello.txt"
	data := []byte("hello world")

	userID := "10000"

	name = path.Join("user_"+userID, name)

	addr, err := c.ApplyPut(ctx, &FragmentPutArgs{
		PutArgs: PutArgs{
			Name:          name,
			Size:          int64(len(data)),
			Hash:          Md5(data),
			EffectiveTime: time.Second * 60 * 60,
		},
		FragmentSize: 2,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println()
	fmt.Println()

	if addr.ResourceURL != "" {
		log.Println("服务器已经存在")
		return
	}
	var (
		start int
		end   = int(addr.FragmentSize)
	)

	for _, u := range addr.PutURLs {
		if start >= len(data) {
			break
		}
		if end > len(data) {
			end = len(data)
		}
		_ = u
		page := data[start:end]
		fmt.Print(string(page))
		start += int(addr.FragmentSize)
		end += int(addr.FragmentSize)
		err = HttpPut(u, bytes.NewReader(page))
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("[PUT_ID]", addr.PutID)

	info, err := c.ConfirmPut(ctx, addr.PutID)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v\n", info)

	log.Println("success")
}
