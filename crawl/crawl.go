package crawl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	dbi "github.com/ecomcrawler/db"
	deg "github.com/ecomcrawler/degisken"
	yar "github.com/ecomcrawler/yardimci"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var sliceDt []deg.Data // json için ürün listesi

func TamTarama(site deg.SiteHtml) {
	dbi.Dbolustur()
	c := colly.NewCollector(
		colly.AllowedDomains(site.AllowedDomains),
		colly.UserAgent(site.UserAgent),
		//colly.Async(true),
		//colly.IgnoreRobotsTxt(),
	)

	c.Visit(site.AllowedDomains)

	kateg := make([]deg.Kategori, 0)

	c.OnHTML(site.KategoriAna, func(e *colly.HTMLElement) {

		var k deg.Kategori
		/* e.ForEach("li.feed-carousel-card a.a-link-normal", func(i int, h *colly.HTMLElement) {
			k.KategoriAdi = h.Text
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 0
			//c.Visit(k.KategoriLink)
			fmt.Println(k)

		}) */
		e.ForEach(site.Kategori5[0], func(i int, h *colly.HTMLElement) {
			k.KategoriAdi = h.Attr(site.KategoriAdi[0])
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 5

			kateg = append(kateg, k)
		})
		e.ForEach(site.Kategori4[0], func(i int, h *colly.HTMLElement) {
			k.KategoriAdi = h.Attr(site.KategoriAdi[0])
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 4

			kateg = append(kateg, k)

		})
		e.ForEach(site.Kategori3[0], func(i int, h *colly.HTMLElement) {

			k.KategoriAdi = h.Attr(site.KategoriAdi[0])
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 3
			kateg = append(kateg, k)

		})
		e.ForEach(site.Kategori2[0], func(i int, h *colly.HTMLElement) {

			k.KategoriAdi = h.Attr(site.KategoriAdi[0])
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 2
			kateg = append(kateg, k)

		})
		e.ForEach(site.Kategori1[0], func(i int, h *colly.HTMLElement) {

			k.KategoriAdi = h.Attr(site.KategoriAdi[0])
			k.KategoriLink = h.Request.AbsoluteURL(h.Attr(site.KategoriLink[0]))
			k.AltKategori = 1
			kateg = append(kateg, k)

		})

		a, err := json.Marshal(kateg)
		yar.CheckErr(err)
		err2 := ioutil.WriteFile(site.AllowedDomains+"_kategoriler.json", []byte(a), 0644)
		yar.CheckErr(err2)

	})
	tt := c.Clone()
	tt.Limit(&colly.LimitRule{
		//aynı andaki istek sayısı
		Parallelism: 2,
		//iki istek arası bekleme
		Delay: 10 * time.Second,
	})
	tt.OnHTML(site.UrunlerContainer, func(e *colly.HTMLElement) {
		e.ForEach(site.UrunContainer, func(i int, h *colly.HTMLElement) {
			// başında ve sonundaki boşlukları ve "TL" silme ve float dönüştürme
			rawPrice := h.ChildText(site.UrunFiyati[0])
			if rawPrice == "" {
				rawPrice = h.ChildText(site.UrunFiyati[1])

			}
			p, _ := strconv.ParseFloat(strings.ReplaceAll((strings.TrimSpace(strings.Trim(h.ChildText(rawPrice), "TL"))), ",", "."), 32)
			deg.Dt.Price = math.Round(p*100) / 100
			if deg.Dt.Price != 0 {
				currentTime := time.Now()
				deg.Dt.PageId = h.Index

				deg.Dt.Title = h.ChildText(site.UrunAdi)

				//tüm boşlukları silmek için
				/* dt.Price = strings.ReplaceAll(s.Find(".last-price").Text(), "\n", "") */
				deg.Dt.Resim = h.ChildAttr(site.UrunResmi[0], site.UrunResmi[1])

				if site.UrunLinki[0] != "" {
					deg.Dt.Link = h.Request.AbsoluteURL(h.ChildAttr(site.UrunLinki[0], site.UrunLinki[1]))
				} else {
					deg.Dt.Link = h.Request.AbsoluteURL(h.Attr(site.UrunLinki[1]))
				}

				deg.Dt.UrunKodu = h.Attr(site.UrunKodu[1])
				if deg.Dt.UrunKodu == "" {
					deg.Dt.UrunKodu = h.ChildAttr(site.UrunKodu[0], site.UrunKodu[1])
				}
				deg.Dt.KayitTarih = currentTime.Format("2006.01.02 15:04:05")
				deg.Dt.Domain = h.Request.URL.Host

				dbi.DbKayit()

				sliceDt = append(sliceDt, deg.Dt)
			}

		})
	})
	if !site.Pagination {

		c.OnHTML("div.dscrptn", func(h *colly.HTMLElement) {
			s := h.Text
			//s = s[len(s)-30 : len(s)-20]
			reg, err := regexp.Compile("[^0-9]")
			if err != nil {
				log.Fatal(err)
			}
			sonuc := reg.ReplaceAllString(s, "")

			sonucint, _ := strconv.Atoi(sonuc)
			pi := (sonucint / 24) + 2
			for i := 2; i < pi; i++ {
				arg := kateg[i].KategoriLink + "&pi=%d"
				sayfalar := fmt.Sprintf(arg, i)
				tt.Visit(sayfalar)
			}
		})

	}

	tt.OnHTML(site.SayfaSonuc[0], func(e *colly.HTMLElement) {
		pageLink := e.Attr(site.SayfaSonuc[1])

		tt.Visit(e.Request.AbsoluteURL(pageLink))

	})

	tt.OnRequest(func(r *colly.Request) {
		fmt.Println("Sayfa taranıyor...", r.URL.String())
		// karakter problemi olursa r.ResponseCharacterEncoding = "utf-8"

	})

	c.Visit("https://" + site.AllowedDomains)
	if !site.KategoriLinks {
		klist := HbKateg()

		for i := range klist {
			tt.Visit("https://" + site.AllowedDomains + klist[i].KategoriLink)
		}
	}
	for i := range kateg {

		tt.Visit(kateg[i].KategoriLink)
	}
}

func Tarama(arama string, site deg.SiteHtml) {
	//mysql işlemleri
	dbi.Dbolustur()
	// toplayıcı oluşturma
	c := colly.NewCollector(
		// Ziyaret edilecek domainler
		colly.AllowedDomains(site.AllowedDomains),
		//eşzamanlı sayfa tarama ile çoklu sayfaları hızlandırıyor
		colly.Async(true),
		//izin isteyen sayfalar için "Forbidden" hatası
		colly.UserAgent(site.UserAgent),
		//colly.IgnoreRobotsTxt(),
	)

	c.Limit(&colly.LimitRule{
		//aynı andaki istek sayısı
		Parallelism: 2,
		//iki istek arası bekleme
		Delay: 10 * time.Second,
	})
	// tarama esnasında banlanmamak için proxy ile ip değiştirme
	/* proxySwitcher, err := proxy.RoundRobinProxySwitcher("socks5://195.175.67.202:1080", "socks5://78.47.182.101:31907")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(proxySwitcher) */

	c.OnHTML(site.UrunlerContainer, func(e *colly.HTMLElement) {
		e.ForEach(site.UrunContainer, func(i int, h *colly.HTMLElement) {
			// başında ve sonundaki boşlukları ve "TL" silme ve float dönüştürme

			p, _ := strconv.ParseFloat(strings.ReplaceAll((strings.TrimSpace(strings.Trim(h.ChildText(site.UrunFiyati[0]), "TL"))), ",", "."), 32)
			if p == 0 {

				p, _ = strconv.ParseFloat(strings.ReplaceAll((strings.TrimSpace(strings.Trim(h.ChildText(site.UrunFiyati[1]), "TL"))), ",", "."), 32)
			}

			deg.Dt.Price = math.Round(p*100) / 100

			if deg.Dt.Price != 0 {
				currentTime := time.Now()
				deg.Dt.PageId = h.Index

				deg.Dt.Title = h.ChildText(site.UrunAdi)

				//tüm boşlukları silmek için
				/* dt.Price = strings.ReplaceAll(s.Find(".last-price").Text(), "\n", "") */
				deg.Dt.Resim = h.ChildAttr(site.UrunResmi[0], site.UrunResmi[1])

				if site.UrunLinki[0] != "" {
					deg.Dt.Link = h.Request.AbsoluteURL(h.ChildAttr(site.UrunLinki[0], site.UrunLinki[1]))
				} else {
					deg.Dt.Link = h.Request.AbsoluteURL(h.Attr(site.UrunLinki[1]))
				}

				deg.Dt.UrunKodu = h.Attr(site.UrunKodu[1])
				if deg.Dt.UrunKodu == "" {
					deg.Dt.UrunKodu = h.ChildAttr(site.UrunKodu[0], site.UrunKodu[1])
				}
				deg.Dt.KayitTarih = currentTime.Format("2006.01.02 15:04:05")
				deg.Dt.Domain = h.Request.URL.Host

				dbi.DbKayit()

				js, err := json.Marshal(deg.Dt)
				yar.CheckErr(err)
				fmt.Println(string(js))
				sliceDt = append(sliceDt, deg.Dt)
			}

		})

	})
	c.OnHTML(site.SayfaSonuc[0], func(e *colly.HTMLElement) {
		pageLink := e.Attr(site.SayfaSonuc[1])

		c.Visit(e.Request.AbsoluteURL(pageLink))

	})

	// Sayfa açılma öncesi bildirim "Sayfa taranıyor..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Sayfa taranıyor...", r.URL.String())
		// karakter problemi olursa r.ResponseCharacterEncoding = "utf-8"

	})

	// taramayı başlat

	baslat := fmt.Sprintf(site.BaslangicSayfa, url.PathEscape(arama))
	fmt.Println(yar.ReplaceTrChar(url.PathEscape(arama)))
	if !site.Pagination {

		c.OnHTML(site.PaginationTag, func(h *colly.HTMLElement) {
			s := h.Text
			//s = s[len(s)-30 : len(s)-20]
			reg, err := regexp.Compile("[^0-9]")
			if err != nil {
				log.Fatal(err)
			}
			sonuc := reg.ReplaceAllString(s, "")

			sonucint, _ := strconv.Atoi(sonuc)
			pi := (sonucint / 24) + 2
			for i := 2; i < pi; i++ {
				arg := baslat + site.PaginationArg
				sayfalar := fmt.Sprintf(arg, i)
				c.Visit(sayfalar)
			}
		})

	}

	c.Visit(baslat)

	c.Wait()
	//json dosyasına kaydetme
	file, _ := json.MarshalIndent(sliceDt, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)
	fmt.Printf("%d online kayıt bulundu", len(sliceDt))

}

func HbKateg() (klist []deg.Kategori) {
	file, _ := ioutil.ReadFile("hb.json")
	var hbkat deg.Hb
	var k deg.Kategori

	_ = json.Unmarshal(file, &hbkat)
	for i := 0; i < len(hbkat); i++ {
		k.KategoriLink = hbkat[i].URL
		k.KategoriAdi = hbkat[i].Title
		k.KategTip = hbkat[i].EntityType
		klist = append(klist, k)
		a := hbkat[i].Children

		for i := 0; i < len(a); i++ {
			c := a[i].Children
			k.KategoriAdi = a[i].Title
			k.KategoriLink = a[i].URL
			k.KategTip = a[i].EntityType
			klist = append(klist, k)

			for i := 0; i < len(c); i++ {
				k.KategoriLink = c[i].URL
				k.KategoriAdi = c[i].Title
				k.KategTip = c[i].EntityType
				klist = append(klist, k)
				e := c[i].Children

				for i := 0; i < len(e); i++ {
					k.KategoriLink = e[i].URL
					k.KategoriAdi = e[i].Title
					k.KategTip = e[i].EntityType
					klist = append(klist, k)
					f := e[i].Children
					for i := 0; i < len(f); i++ {
						k.KategoriLink = f[i].URL
						k.KategoriAdi = f[i].Title
						k.KategTip = f[i].EntityType
						klist = append(klist, k)

					}

				}
			}
		}

	}
	return klist
}
