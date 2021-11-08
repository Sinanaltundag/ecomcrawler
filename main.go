package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	crw "projeler/collya/crawl"
	dbi "projeler/collya/db"
	deg "projeler/collya/degisken"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var tercih int

	for tercih < 1 || tercih > 2 {
		fmt.Println("Tam tarama için 1, arama yapmak için 2 girin")
		fmt.Scanln(&tercih)
	}

	if tercih == 1 {
		siteler := make([]deg.SiteHtml, 0, 50)
		file, _ := ioutil.ReadFile("siteler.json")

		_ = json.Unmarshal([]byte(file), &siteler)

		for i, v := range siteler {
			fmt.Printf("%d = %v\n", i+1, v.AllowedDomains)

		}
		fmt.Println("Tarama yapmak istediğiniz site:")
		fmt.Println("Tüm sitelerde tarama yapmak için Enter'a basınız. Sadece Bir sitede aramak için site numarasını tuşlayıp Enter'a basınız.")
		//Kullanıcı Girişi
		var siteNum int
		fmt.Scanln(&siteNum)
		if siteNum <= len(siteler) && siteNum > 0 {
			crw.TamTarama(siteler[siteNum-1])
		} else {

			for i := range siteler {
				crw.TamTarama(siteler[i])
			}
		}
	} else if tercih == 2 {
		fmt.Println("Arama yapmak istediğiniz ürün:")
		//Kullanıcı Girişi
		var arama string

		for len(arama) < 5 {
			fmt.Println("En az 5 karakter girin.")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				arama = scanner.Text()

			}
		}
		satirlar, err := dbi.DbSorgu(arama)

		if len(satirlar) != 0 {
			fmt.Printf("%d kayıt bulundu.\n", len(satirlar))
			for i := range satirlar {

				fmt.Println(satirlar[i])
			}

		}

		if err != nil {
			fmt.Println("Yerel kayıtlarımızda bulunamadı")
		}

		fmt.Println("Online aramak ister misiniz? Evet için ´e´ tuşuna basıp Enter'a basın.")
		//Kullanıcı Girişi
		var online string
		fmt.Scanln(&online)

		if online == "e" || online == "E" {
			fmt.Println("Arama yapmak istediğiniz site:")
			fmt.Println("Tüm sitelerde arama yapmak için Enter'a basınız. Sadece Bir sitede aramak için site numarasını tuşlayıp Enter'a basınız.")

			//json verilerini işlemek için struct slice oluşturuyoruz
			siteler := make([]deg.SiteHtml, 0, 50)
			file, _ := ioutil.ReadFile("siteler.json")

			_ = json.Unmarshal([]byte(file), &siteler)

			for i, v := range siteler {
				fmt.Printf("%d = %v\n", i+1, v.AllowedDomains)

			}
			//Kullanıcı Girişi
			var siteNum int
			fmt.Scanln(&siteNum)
			if siteNum <= len(siteler) && siteNum > 0 {
				crw.Tarama(arama, siteler[siteNum-1])
			} else {

				for i := range siteler {
					crw.Tarama(arama, siteler[i])
				}
			}
		}
	} else {
		fmt.Println("Lütfen 1 yada 2 girip Enter tuşuna basın.")
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		_ = sig
		done <- true
	}()

	<-done

	fmt.Println("Program Closed")
}
