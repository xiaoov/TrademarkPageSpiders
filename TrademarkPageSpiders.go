package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"regexp"
	"strconv"
	"net/url"
	"bytes"
)

type Jar struct{
	cookies []*http.Cookie
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie{
	return jar.cookies
}

func main() {
	dat, err := ioutil.ReadFile("Words.txt")
	if err != nil{
		fmt.Println("Read the keyword file failed.\n")
		return
	}
	words := strings.Split(string(dat), "\n")
	r := words[0]
	keyword := words[1]
	fmt.Println("Read the keyword file success.\n")
	fmt.Println("Start to load the key and cookies.")
	key, _ := LoadTheFormData()
	cookie, err := ioutil.ReadFile("cookie")
	/*if err != nil{

	}*/
	fmt.Println("Load key and cookies success.")
	file := r
	//fmt.Println(r)
	client := http.Client{nil, nil, nil ,0}
	fmt.Printf("Start to search the keyword :%s\n\r", r)
	SEARCHURL := bytes.Buffer{}
	SEARCHURL.WriteString("http://tmsearch.uspto.gov/bin/showfield?f=toc&")
	SEARCHURL.WriteString(key)
	SEARCHURL.WriteString("&p_search=searchss&p_L=50&BackReference=&p_plural=yes&p_s_PARA1=&p_tagrepl%7E%3A=PARA1%24LD&expr=PARA1+AND+PARA2&p_s_PARA2=")
	SEARCHURL.WriteString(strings.Trim(r,"\r"))
	SEARCHURL.WriteString("&p_tagrepl%7E%3A=PARA2%24COMB&p_op_ALL=AND&a_default=search&a_search=Submit+Query&a_search=Submit+Query")
	request, _ := http.NewRequest("Get",SEARCHURL.String(),nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8");
	request.Header.Add("Accept-Encoding", "gzip, deflate, sdch");
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8");
	request.Header.Add("Connection", "keep-alive");
	request.Header.Add("Cache-Control", "max-age=0");
	request.Header.Add("Host", "tmsearch.uspto.gov");
	request.Header.Add("Cookie", string(cookie));
	//request.Header.Add("Referer", "http://tmsearch.uspto.gov/bin/gate.exe?f=searchss&"+key);
	request.Header.Add("Upgrade-Insecure-Requests", "1");
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36");
	response, _ :=  client.Do(request)
	//reqHeader := request.Header
	//fmt.Println(request.Header)
	//fmt.Println(SEARCHURL.String())
	//return
	if err != nil{
		fmt.Printf("Load page failed:%s\n\r", r)
	}
	body, _ := ioutil.ReadAll(response.Body)

	ro := regexp.MustCompile(`Record List Display`)
	c8 := ro.FindString(string(body))
	if c8 == ""{
		fmt.Println("Read faild")
		//fmt.Println(string(body))
		//fmt.Println(SEARCHURL.String())
		return
	}
	//fmt.Println(string(body))
	rows := GetRows(string(body))
	docs := GetDocs(string(body))
	//fmt.Println(docs)
	links := GetLink(rows)
	docs_int , _ := strconv.Atoi(docs)
	//fmt.Println(docs_int)
	var writeString string
	for i:=1;i<=docs_int;i++{
		URL := "http://tmsearch.uspto.gov" + links + "." + strconv.Itoa(i)
		fmt.Println(URL)
		//response_child_link, err := http.Get("http://tmsearch.uspto.gov//bin/showfield?f=doc&state=4808:qy8yvq.46.9")
		request_child_link, _ := http.NewRequest("Get",URL,nil)
		//request_child_link.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,**/*//*;q=0.8");
		//request_child_link.Header.Add("Accept-Encoding", "gzip, deflate, sdch");
		//request_child_link.Header.Add("Accept-Language", "zh-CN,zh;q=0.8");
		//request_child_link.Header.Add("Connection", "keep-alive");
		//request_child_link.Header.Add("Cache-Control", "max-age=0");
		//request_child_link.Header.Add("Host", "tmsearch.uspto.gov");
		//request_child_link.Header.Add("Cookie", string(cookie));
		////request_child_link.Header.Add("Referer", "http://tmsearch.uspto.gov/bin/gate.exe?f=searchss&"+key);
		//request_child_link.Header.Add("Upgrade-Insecure-Requests", "1");
		//request_child_link.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36");
		response_child_link, _ :=  client.Do(request_child_link)
		//request_child_link.Body.Close()
		if err != nil{
			fmt.Printf("Load page failed:%s\n\r", response_child_link)
		}
		body_child, _ := ioutil.ReadAll(response_child_link.Body)
		ro := regexp.MustCompile(`Trademark Electronic Search System`)
		c9 := ro.FindString(string(body_child))
		if c9 == ""{
			i--
			fmt.Println("Read faild")
			//return
		}else{
			fmt.Println("Read success")
		}

		cc := isCorrectLink(string(body_child), file, keyword)
		//fmt.Println(string(body_child))

		if cc == true{
			writeString += URL + "\r\n"
			fmt.Println("found one")
		}else{
			//fmt.Println("no such data")
		}
	}
	/*fmt.Println("file:")
	fmt.Println(file)*/
	filename := strings.Trim(file," \n\r")+".txt"
	fmt.Println("Saving results.")
	err1 := ioutil.WriteFile(strings.Trim(filename," \n\r"),  []byte(writeString), 0666)  //写入文件(字节数组)
	if err1 != nil{
		fmt.Printf("Write file failed:%s\n\r", err)
	}
}

func isCorrectLink(URL string, search_text string, keyword string)bool{
	search_reg := "(?i)" + search_text
	//fmt.Println(URL)
	r1 := regexp.MustCompile(strings.Trim(search_reg," \n\r"))
	c1 := r1.FindString(URL)
	r2 := regexp.MustCompile(`(?i)` + strings.Trim(keyword," \n\r"))
	c2 := r2.FindString(URL)
	if c1 == "" || c2 == ""{
		return false
	}
	return true
}

func LoadTheFormData() (string, string){
	var key_new string
	jar := new(Jar)
	client := http.Client{nil, nil, jar ,0}
	cookie, err4 := ioutil.ReadFile("cookie")
	if err4 != nil{
		fmt.Println("Read cookie failed")
	}
	key, err5 := ioutil.ReadFile("key")
	if err5 != nil{
		fmt.Println("Read key failed")
	}
	test_request, _ := http.NewRequest("Get","http://tmsearch.uspto.gov/bin/gate.exe?f=searchss&"+string(key),nil)
	test_request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8");
	test_request.Header.Add("Accept-Encoding", "gzip, deflate, sdch");
	test_request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8");
	test_request.Header.Add("Connection", "keep-alive");
	test_request.Header.Add("Cache-Control", "max-age=0");
	test_request.Header.Add("Host", "tmsearch.uspto.gov");
	//test_request.Header.Add("Cookie", string(cookie));
	//test_request.Header.Add("Referer", "http://tmsearch.uspto.gov/bin/gate.exe?f=searchss&"+string(key));
	test_request.Header.Add("Upgrade-Insecure-Requests", "1");
	test_request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36");
	test_resp, _ :=  client.Do(test_request)
	test_body, _ := ioutil.ReadAll(test_resp.Body)
	ro := regexp.MustCompile(`Trademark Electronic Search System`)
	c := ro.FindString(string(test_body))
	//fmt.Println("c:")
	//fmt.Println(string(test_body))
	if c != ""{
		fmt.Println("Old key is still useful.")
		return string(key), string(cookie)
	}
	fmt.Println("Old key and cookies has been expired.Request for a new one")

	Gate := "http://tmsearch.uspto.gov/bin/gate.exe?f=login&p_lang=english&p_d=trmk"
	request, _ := http.NewRequest("Get",Gate,nil)
	response, _ :=  client.Do(request)
	cookies := response.Cookies()
	cookie_string := bytes.Buffer{}
	for _, cookie := range cookies {
		cookie_string.WriteString(cookie.Name)
		cookie_string.WriteString(cookie.Value)
		cookie_string.WriteString("; ")
	}
	fmt.Println("Write cookie")
	//fmt.Println(cookie_string.String())
	err := ioutil.WriteFile("cookie",  []byte(cookie_string.Bytes()), 0666)  //写入文件(字节数组)
	if err != nil{
		fmt.Printf("Write cookie failed:%s\n\r", err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	r := regexp.MustCompile(`(state=[.\s\S]*?).1">`)
	links := r.FindStringSubmatch(string(body))
	if len(links) == 0{
		fmt.Println("Redirect failed")
	} else if links[0] != ""{
		//fmt.Println(links[1])
		key_new = links[1]
	}
	fmt.Println("Write key")
	err3 := ioutil.WriteFile("key",  []byte(key_new), 0666)  //写入文件(字节数组)
	if err3 != nil{
		fmt.Printf("Write cookie failed:%s\n\r", err)
	}
	return key_new, cookie_string.String()
}

func GetRows(respBody string) []string {
	r := regexp.MustCompile(`<TR>[.\s\S]*?</TR>`)
	return r.FindAllString(respBody,-1)
}

func GetLink(rows []string) string{
	for _, row := range rows{
		c := isCorrectData(row)
		if c == true{
			continue
		}
		r := regexp.MustCompile(`href="([.\s\S]*?).1">`)
		links := r.FindStringSubmatch(row)
		//link := links[0]
		if len(links) == 0{

		} else if links[0] != ""{
			//fmt.Println(links[1])
			return links[1]
		}
	}
	return ""
}

func isCorrectData(row string) bool{
	r := regexp.MustCompile(`showfield`)
	c := r.FindString(row)
	if c == ""{
		return true
	}
	return false
}

func GetDocs(respBody string) string{
	r := regexp.MustCompile(`docs: (.*?) `)
	docs := r.FindStringSubmatch(respBody)
	if len(docs) == 0{
		//fmt.Println(respBody)
		return ""
	}
	fmt.Println(docs[1]+" records")
	return docs[1]
}