package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	//"net/url"
	"regexp"
)

func main() {
	dat, err := ioutil.ReadFile("Words.txt")
	if err != nil{
		fmt.Println("Read the keyword file failed.\n\r")
		return
	}
	words := strings.Split(string(dat), "\n")
	fmt.Println("Read the keyword file success.\n\r")
	for _, r := range words{
		fmt.Printf("Start to searching the keyword :%s\n\r", r)
		response, err := http.Get("http://tess2.uspto.gov/bin/showfield?f=toc&state=4807%3Ag9r1kt.1.1&p_search=searchss&p_L=50&BackReference=&p_plural=yes&p_s_PARA1=&p_tagrepl%7E%3A=PARA1%24LD&expr=PARA1+AND+PARA2&p_s_PARA2=lyra&p_tagrepl%7E%3A=PARA2%24COMB&p_op_ALL=AND&a_default=search&a_search=Submit+Query&a_search=Submit+Query")
		if err != nil{
			fmt.Printf("Load page failed:%s\n\r", r)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		rows := GetRows(string(body))
		for _, row := range rows{
			c := isCorrectData(row)
			if c == true{
				continue
			}
			links := GetLink(row)
			//link := links[0]
			fmt.Println(links)
		}
		//fmt.Println(rows[5])
		return
	}
}

func LoadTheFormData() (string){
	var UrlStrings string
	//read the from data
	return UrlStrings
}

func GetRows(respBody string) []string {
	r := regexp.MustCompile(`<TR>[.\s\S]*?</TR>`)
	return r.FindAllString(respBody,-1)
}

func GetLink(row string) [][]string{
	r := regexp.MustCompile(`href="([.\s\S]*?)">`)
	return r.FindAllStringSubmatch(row,-1)
}

func isCorrectData(row string) bool{
	r := regexp.MustCompile(`showfield`)
	c := r.FindString(row)
	if c == ""{
		return true
	}
	return false
}