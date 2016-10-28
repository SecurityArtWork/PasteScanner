//  This file is part of Pastescanner.
// 
//  Copyright (C) 2016, "Security Art Work" www.securityartwork.es
//  All rights reserved.
// 
//  Pastescanner is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
// 
//  Pastescanner is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Lesser General Public License for more details.
// 
//  You should have received a copy of the GNU Lesser General Public License
//  along with Pastescanner.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"net/http"
	"gopkg.in/xmlpath.v2"
	"os"
	"strconv"
)

var visited []string

func getPaste(link string) string{
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("error on httpRequest")
		return "-1"
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error sending request")
		return "-1"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error reading req Body")
		return "-1"
	}

	return string(body)
}

func notin(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return false
        }
    }
    return true
}

func find(link string, keywords []string){

	var raw string
	ttl := "Never"

	if strings.Contains(link,"pastebin.com") {
		ttl,raw = getDataPB(getPaste(link))
	
	} else if strings.Contains(link,"pastie") {
		raw = getDataPIE(getPaste(link+"/text"))	
	
	} else if strings.Contains(link,"pastebin") {
		raw = getPaste(link)
	
	} else {
		raw = getDataPLK(getPaste(link))
	
	}

	var contains []string
	for i:=0; i< len(keywords); i++ {
		if (strings.Contains(raw, keywords[i])){
			contains = append(contains,keywords[i])
		}
	}

	if len(contains) > 0 {
		
		path := "./pastes/"+link[len(link)-4:]+".txt"
		if ttl!="Never" {
			path = "./pastes/temp/"+link[len(link)-4:]+".txt"
		}

		f, err := os.Create(path)
		if err != nil {
	        fmt.Println(err)
	    }

		f.WriteString("-----------Meta----------\n")
		f.WriteString("LINK: "+link+"\n")
		f.WriteString("TTL: "+ttl+"\n")
		f.WriteString("contains: "+strings.Join(contains[:],",")+"\n")
		f.WriteString("-----------Meta----------\n\n")
		f.WriteString(raw)
		f.Sync()
		defer f.Close()

	} 
}

func main() {	
	//fmt.Println("[cargado configuracion]")
	dat, _ := ioutil.ReadFile("./paste.conf")
    
	file := strings.Split(string(dat),"\n")

	var keys []string
	var sites []string

	arekeys:=false
	aresites:=false
	for i:=0; i< len(file);i++ {
		if arekeys {
			if file[i]=="[sites]"{
				arekeys =false
				aresites = true
			}else {
				keys = append(keys, file[i])
			} 
		} else if file[i]=="[Keys]" {
			arekeys = true
		} else if aresites {
			sites = append(sites, file[i])
		}
	}

	
	for {
		// Latest pastes list from every site | Lista de ultimos pastes de cada sitio
 		var lista []string

		//fmt.Println("[Actualizando ultimos pastes]")
 		if !notin("pastebin.com", sites) {
	 		listaPB := getLastsPB()   //PASTEBIN
	 		for x:=0 ; x < len(listaPB); x++ {
	 			lista = append(lista, listaPB[x])
	 		}
 		}

 		if !notin("pastie.org", sites) {
	 		listaPIE := getLastsPIE() //PASTIE
	 		for x:=0 ; x < len(listaPIE); x++ {
	 			lista = append(lista, listaPIE[x])
	 		}
 		}
        
        if !notin("pastebin.ca", sites) {
	 		listaPCA := getLastsPCA() //PASTEBIN.CA
			for x:=0 ; x < len(listaPCA); x++ {
	 			lista = append(lista, listaPCA[x])
	 		}
 		}

 		if !notin("pastelink.net", sites) {
	 		listaPLK := getLastsPLK() //PASTELINK
			for x:=0 ; x < len(listaPLK); x++ {
	 			lista = append(lista, listaPLK[x])
	 		}
 		}

		//fmt.Println("[Revisando pastes]")
		for i := 0; i < len(lista); i++ {
			if notin(lista[i], visited) { //this has to be atomic so links are overlapped and lost D: | esto tiene que ser atomico para que solape links y pierda muchos D:
				visited = append(visited,lista[i])
				go find(lista[i],keys)
			}
		}

		time.Sleep(5000 * time.Millisecond);
	}


}


//****************************************************************
// 							PASTEBIN
//****************************************************************

func getLastsPB() [8]string{

	var links [8]string
	reader := strings.NewReader(getPaste("http://pastebin.com"))
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return links
	}

	for i := 1; i < 9; i++ { 
	  	path := xmlpath.MustCompile("//*[@id=\"menu_2\"]/ul/li["+strconv.Itoa(i)+"]/a/@href")
	  	if value, ok := path.String(xmlroot); ok {
	 	    links[i-1]="http://pastebin.com"+value
	 	}
	}

	return links;
}

func getDataPB(body string)(string,string){

	reader := strings.NewReader(body)
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return "-1","-1"
	}

	var ttl, raw string
	path := xmlpath.MustCompile("//*[@id=\"content_left\"]/div[3]/div[3]/div[2]")
	if value, ok := path.String(xmlroot); ok {
	    ttl = strings.TrimSpace(strings.Split(value,"\n")[4])
	}

	path = xmlpath.MustCompile("//*[@id=\"paste_code\"]")
	if value, ok := path.String(xmlroot); ok {
	    raw = value
	}

	return ttl,raw
}

//****************************************************************
// 							PASTIE.ORG
//****************************************************************


func getLastsPIE() [20]string{

	var links [20]string
	reader := strings.NewReader(getPaste("http://pastie.org/pastes"))
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return links
	}

	for i := 3; i < 23; i++ { 
	  	path := xmlpath.MustCompile("//*[@id=\"content\"]/div/div["+strconv.Itoa(i)+"]/p[2]/a/@href")
	  	if value, ok := path.String(xmlroot); ok {
	 	    links[i-3]=value
	 	}
	}

	return links;
}

func getDataPIE(body string) string {
	
	reader := strings.NewReader(body)
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return "-1"
	}

	path := xmlpath.MustCompile("/html/body/pre")
	if value, ok := path.String(xmlroot); ok {
	    return value
	}

	return "-1"
}


//****************************************************************
// 						    PASTEBIN.CA
//****************************************************************

func getLastsPCA() [15]string{

	var links [15]string
	reader := strings.NewReader(getPaste("http://pastebin.ca/"))
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return links
	}

	for i := 1; i < 16; i++ {        
	  	path := xmlpath.MustCompile("//*[@id=\"idmenurecent-collapse\"]/div["+strconv.Itoa(i)+"]/a/@href")
	  	if value, ok := path.String(xmlroot); ok {
	 	    links[i-1]="http://pastebin.ca/raw"+value
	 	}
	}

	return links;
}

//****************************************************************
// 							PASTELINK
//****************************************************************


func getLastsPLK() [20]string{

	var links [20]string
	reader := strings.NewReader(getPaste("https://pastelink.net/read"))
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return links
	}

	for i := 1; i < 21; i++ { 		
	  	path := xmlpath.MustCompile("//*[@id=\"listing\"]/tbody/tr["+strconv.Itoa(i)+"]/td[1]/a/@href")
	  	if value, ok := path.String(xmlroot); ok {
	 	    links[i-1]=value
	 	}
	}

	return links;
}

func getDataPLK(body string) string {

	reader := strings.NewReader(body)
	xmlroot, xmlerr := xmlpath.ParseHTML(reader)
	if xmlerr != nil {
	    return "-1"
	}

	path := xmlpath.MustCompile("//*[@id=\"body-display\"]")
	if value, ok := path.String(xmlroot); ok {
	    return value
	}

	return "-1"

}
