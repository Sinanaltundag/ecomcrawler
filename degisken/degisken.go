package degisken

type Data struct {
	Id         int     `json:"id"`
	PageId     int     `json:"pageid"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	Resim      string  `json:"resim"`
	Link       string  `json:"link"`
	Kategori   string  `json:"kategori"`
	Marka      string  `json:"marka"`
	UrunKodu   string  `json:"urunkodu"`
	KayitTarih string  `json:"kayittarih"`
	Domain     string  `json:"domain"`
	SonucSay   int     `json:"sonucsay"`
}
type DbConn struct {
	DbServer string `json:"dbServer"`
	DbUser   string `json:"dbUser"`
	DbPass   string `json:"dbPass"`
	DbName   string `json:"dbName"`
	DbPort   int    `json:"dbPort"`
}
type Siteler struct {
	Siteler []*SiteHtml `json:"Siteler"`
}

type SiteHtml struct {
	AllowedDomains   string    `json:"AllowedDomains"`
	BaslangicSayfa   string    `json:"BaslangicSayfa"`
	UrunlerContainer string    `json:"UrunlerContainer"`
	UrunContainer    string    `json:"UrunContainer"`
	UrunAdi          string    `json:"UrunAdi"`
	UrunFiyati       [2]string `json:"UrunFiyati"`
	UrunResmi        [2]string `json:"UrunResmi"`
	UrunLinki        [2]string `json:"UrunLinki"`
	UrunKodu         [2]string `json:"UrunKodu"`
	SayfaSonuc       [2]string `json:"SayfaSonuc"`
	Kategori1        [2]string `json:"Kategori1"`
	Kategori2        [2]string `json:"Kategori2"`
	Kategori3        [2]string `json:"Kategori3"`
	Kategori4        [2]string `json:"Kategori4"`
	Kategori5        [2]string `json:"Kategori5"`
	KategoriAna      string    `json:"KategoriAna"`
	KategoriSira     [2]string `json:"KategoriSira"`
	KategoriAdi      [2]string `json:"KategoriAdi"`
	KategoriLink     [2]string `json:"KategoriLink"`
	Pagination       bool      `json:"pagination"`
	KategoriLinks    bool      `json:"KategoriLinks"`
	PaginationTag    string    `json:"paginationTag"`
	PaginationArg    string    `json:"paginationArg"`
	UserAgent        string    `json:"userAgent"`
}

var Dt Data

type Kategori struct {
	Id           int
	KategoriAdi  string
	KategoriLink string
	AltKategori  int
	KategTip     string
}

type Hb []struct {
	ID         int    `json:"id"`
	Level      int    `json:"level"`
	ParentID   int    `json:"parentId"`
	CategoryID string `json:"categoryId"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Priority   int    `json:"priority"`
	EntityType string `json:"entityType"`
	Children   []struct {
		ID         int    `json:"id"`
		Level      int    `json:"level"`
		ParentID   int    `json:"parentId"`
		CategoryID string `json:"categoryId"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		Priority   int    `json:"priority"`
		EntityType string `json:"entityType"`
		Children   []struct {
			ID         int    `json:"id"`
			Level      int    `json:"level"`
			ParentID   int    `json:"parentId"`
			CategoryID string `json:"categoryId"`
			Title      string `json:"title"`
			URL        string `json:"url"`
			Priority   int    `json:"priority"`
			EntityType string `json:"entityType"`
			Children   []struct {
				ID         int    `json:"id"`
				Level      int    `json:"level"`
				ParentID   int    `json:"parentId"`
				CategoryID string `json:"categoryId"`
				Title      string `json:"title"`
				URL        string `json:"url"`
				Priority   int    `json:"priority"`
				EntityType string `json:"entityType"`
				Children   []struct {
					ID         int           `json:"id"`
					Level      int           `json:"level"`
					ParentID   int           `json:"parentId"`
					CategoryID string        `json:"categoryId,omitempty"`
					Title      string        `json:"title"`
					URL        string        `json:"url"`
					Priority   int           `json:"priority"`
					EntityType string        `json:"entityType"`
					Children   []interface{} `json:"children"`
				} `json:"children"`
			} `json:"children"`
		} `json:"children"`
	} `json:"children"`
	Banners []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Link      string `json:"link"`
		MediaPath string `json:"mediaPath"`
		SortID    int    `json:"sortId"`
	} `json:"banners"`
	ImageURL string `json:"imageUrl"`
}
