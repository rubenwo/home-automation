package rtsp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os/exec"
	"strings"
)

type Client struct {
	Id   int64
	Host url.URL
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return &Client{}, nil
}

func (c *Client) Connect() error {
	return nil
}

func Stream(conn net.Conn) error {
	cmdStr := `ffmpeg -i rtsp://user:pwd@somewhere/videoSub -c:v copy -c:a copy -bsf:v h264_mp4toannexb -maxrate 500k -f matroska -`
	args := strings.Split(cmdStr, " ")[1:]
	cmd := exec.Command("ffmpeg", args...)

	stdout, _ := cmd.StdoutPipe()

	stderr, _ := cmd.StderrPipe()

	go func(outPipe io.ReadCloser) {
		reader := bufio.NewReader(outPipe)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(b)
		}

	}(stdout)
	go func(errPipe io.ReadCloser) {
		reader := bufio.NewReader(errPipe)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(b)
		}
	}(stderr)

	err := cmd.Start()
	if err != nil {
		return err
	}
	//	fmt.Println("Started encoder")
	return nil
}
