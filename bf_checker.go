package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var non_one_numbers int
var one_numbers int
var non_two_numbers int
var two_numbers int
var non_three_numbers int
var three_numbers int
var non_four_numbers int
var four_numbers int
var non_five_numbers int
var five_numbers int
var url_filename string
var s string

func bf_check(url string) {
	now := time.Now()

	if strings.Contains(url, "http") {
		fmt.Println("URL: " + url)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		str := strings.Replace(string(responseData), "'", "", -1)

		t := now.Unix()
		s = fmt.Sprintf("%v", t)
		url_filename = "tmp" + s + ".txt"
		err = ioutil.WriteFile(url_filename, []byte(str), 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	} else {
		url_filename = url
		fmt.Println("FILE: " + url)
	}

	benford_compliant := "BFC-0"
	numbers_on_page := exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}' | wc -l | tr -d '\n' | sed -e 's/^[[:space:]]*//'")
	fmt.Println("Total number of digits inside page ", numbers_on_page)

	fmt.Println("Distribution of first digits inside the page:")
	non_one_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep -v '^[1]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))
	one_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep  '^[1]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))

	non_two_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep -v '^[2]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))
	two_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep  '^[2]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))

	non_three_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep -v '^[3]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))
	three_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep  '^[3]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))

	non_four_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep -v '^[4]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))
	four_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep  '^[4]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))

	non_five_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep -v '^[5]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))
	five_numbers, _ = strconv.Atoi(exec_shell("cat " + url_filename + " | grep -Eo '[0-9]{1,4}'  | grep  '^[5]' -c | tr -d '\n' | sed -e 's/^[[:space:]]*//'"))

	first := (float64(one_numbers) / float64(non_one_numbers)) * 100
	second := (float64(two_numbers) / float64(non_two_numbers)) * 100
	third := (float64(three_numbers) / float64(non_three_numbers)) * 100
	fourth := (float64(four_numbers) / float64(non_four_numbers)) * 100
	fifth := (float64(five_numbers) / float64(non_five_numbers)) * 100

	fmt.Printf("1 -  %.2f%% \n", first)
	fmt.Printf("2 -  %.2f%% \n", second)
	fmt.Printf("3 -  %.2f%% \n", third)
	fmt.Printf("4 -  %.2f%% \n", fourth)
	fmt.Printf("5 -  %.2f%% \n\n", fifth)

	if first > second {
		benford_compliant = "BFC-1"
		if second > third {
			benford_compliant = "BFC-2"
			if third > fourth {
				benford_compliant = "BFC-3"
				if fourth > fifth {
					benford_compliant = "BFC-4"
				}
			}
		}
	}
	results_filename := "results-benford.csv"
	// write out in append results
	// https://example.org,BFC-1,timestamp
	// url,benford_compliant_level,timestamp,numbers_on_page
	f, err := os.OpenFile(results_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(url + "," + benford_compliant + "," + s + "," + numbers_on_page + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")

}

func exec_shell(command string) string {
	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var url string

	f, err := os.OpenFile("bf_checker.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	if len(os.Args) < 3 {

		raw_arg := os.Args[1]
		if strings.Contains(raw_arg, "list.txt") {
			// run multiple bf_check() per each line inside list.txt
			// open file .. for each line do ...
			file, err := os.Open("list.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				url = scanner.Text()
				bf_check(url)
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		} else {
			url = raw_arg
			bf_check(url)
			// ./bf_check URL/list.txt/file e.g
			//  Distribution of first digits on numbers inside the page:
			/// 1 - 30.2%
			/// 2 - 17%
			/// 3 - 9%
			/// The distribution is Benford Law compliant thus numbers are mostlikely not tampered

		}
		exec_shell("rm tmp*txt || true") //cleanup tmp files
	} // if contains

}
