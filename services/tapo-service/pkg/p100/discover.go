package p100

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	discoverConcurrency = 100
	tapoAppResponse     = "/app"
)

func DiscoverIPv4(subnet string) ([]net.IP, error) {
	subnetSize, err := strconv.Atoi(strings.Split(subnet, "/")[1])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrDiscoverInvalidSubnet)
	}

	ipStr := strings.Split(strings.Split(subnet, "/")[0], ".")
	subnetIP := make([]int, 0, 4)
	for _, str := range ipStr {
		n, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), ErrDiscoverInvalidSubnet)
		}
		subnetIP = append(subnetIP, n)
	}

	var ips []string

	for i := subnetIP[3]; i < int(math.Pow(2, float64(32-subnetSize))); i++ {
		if subnetSize >= 24 {
			ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", subnetIP[0], subnetIP[1], subnetIP[2], i))
		} else if subnetSize >= 16 {
			ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", subnetIP[0], subnetIP[1], subnetIP[2]+i/256%256, i%256))
		} else if subnetSize >= 8 {
			ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", subnetIP[0], subnetIP[1]+i/256/256%256, subnetIP[2]+i/256%256, i%256))
		} else if subnetSize >= 1 {
			ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", subnetIP[0]+i/256/256/256%256, subnetIP[1]+i/256/256%256, subnetIP[2]+i/256%256, i%256))
		}
	}

	type job struct {
		IP string
	}
	type result struct {
		IP string
		Ok bool
	}

	jobs := make(chan job, len(ips))
	results := make(chan result, len(ips))

	for i := 0; i < discoverConcurrency; i++ {
		go func(jobs chan job, results chan result) {
			client := &http.Client{Timeout: time.Second * 2}
			for job := range jobs {
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/app", job.IP), nil)
				if err != nil {
					continue
				}
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				err = resp.Body.Close()
				if err != nil {
					continue
				}
				if string(data) == tapoAppResponse {
					results <- result{
						IP: job.IP,
						Ok: true,
					}
					continue
				}
				results <- result{
					IP: job.IP,
					Ok: false,
				}
			}
		}(jobs, results)
	}
	for _, ip := range ips {
		jobs <- job{IP: ip}
	}
	close(jobs)

	var discovered []net.IP

	for i := 0; i < len(ips); i++ {
		result := <-results
		if result.Ok {
			discovered = append(discovered, net.ParseIP(result.IP))
		}
	}

	close(results)

	return discovered, nil
}

func DiscoverIPv4AndLogin(subnet, username, password string) ([]*Client, error) {
	tapoEndpoints, err := DiscoverIPv4(subnet)
	if err != nil {
		return nil, err
	}

	var clients []*Client

	for _, endpoint := range tapoEndpoints {
		client, err := New(endpoint.String(), username, password)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}
