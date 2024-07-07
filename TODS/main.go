package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	log.Println("Starting go application ...")
	for {
		runLoadTest()
	}
}

func runLoadTest() {
	for j := 0; j < 30; j++ {
		log.Println("Printing ... ", j)
		time.Sleep(time.Second)
	}
	namespaces := []string{
		"default",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"monitoring",
	}
	endpoint := getEnv("ENDPOINT", "prometheus-kube-prometheus-prometheus.monitoring.svc.cluster.local")
	port := getEnv("PORT", "9090")
	noOfRoutines := getEnv("GOROUTINES", "5")
	noOfRoutinesInt, _ := strconv.Atoi(noOfRoutines)

	wg := new(sync.WaitGroup)
	wg.Add(noOfRoutinesInt)

	for i := 0; i < noOfRoutinesInt; i++ {
		go makeHttpCall(wg, endpoint, port, namespaces[i])
		time.Sleep(time.Second)
	}
	wg.Wait()
}

func makeHttpCall(wg *sync.WaitGroup, endpoint, port, namespace string) {
	defer wg.Done()
	metrics := []string{
		"kube_pod_container_resource_requests",
		"kube_pod_container_resource_limits",
		"kube_pod_container_info",
		"kube_pod_container_status_last_terminated_reason",
		"container_cpu_cfs_periods_total",
		"container_cpu_usage_seconds_total",
		"container_cpu_cfs_throttled_periods_total",
		"container_memory_usage_bytes",
		"container_memory_rss",
		"container_threads",
	}
	start := getEnv("START_TIME", "1712736893")
	end := getEnv("END_TIME", "1712866493")
	for _, metric := range metrics {
		url := fmt.Sprintf("http://%s:%s/api/v1/query_range?query=%s&namespace=%s&start=%s&end=%s&step=518", endpoint, port, metric, namespace, start, end)
		callEndpoint(url)
	}
}

func callEndpoint(url string) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("Connection", "close")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}
