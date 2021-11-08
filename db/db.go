package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	deg "github.com/ecomcrawler/degisken"
	yar "github.com/ecomcrawler/yardimci"

	_ "github.com/go-sql-driver/mysql"
)

func DbBaglan() (db *sql.DB) {
	//lokal db için
	/* var dbc deg.DbConn
	dbFile, errread := ioutil.ReadFile("db.json")
	yar.CheckErr(errread)
	errjson := json.Unmarshal(dbFile, &dbc)
	yar.CheckErr(errjson)
	db, err := sql.Open(dbc.DbDriver, dbc.DbUser+":"+dbc.DbPass+"@/"+dbc.DbName)
	yar.CheckErr(err) */
	var dbc deg.DbConn
	dbFile, errread := ioutil.ReadFile("db.json")
	yar.CheckErr(errread)
	errjson := json.Unmarshal(dbFile, &dbc)
	yar.CheckErr(errjson)
	var err error
	dbArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		dbc.DbUser, dbc.DbPass, dbc.DbServer, dbc.DbPort, dbc.DbName)
	db, err = sql.Open("mysql", dbArgs)
	if err != nil {
		log.Fatal("Database bağlantısında hata oluştu: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db

}

func Dbolustur() {
	db := DbBaglan()
	defer db.Close()
	createStatement := `
CREATE TABLE IF NOT EXISTS data (
	id INT AUTO_INCREMENT,
	Title TEXT NOT NULL,
	Price TEXT NOT NULL,
	Resim TEXT,
	Link TEXT NOT NULL,
	Kategori TEXT,
	Marka TEXT,
	UrunKodu varchar(255) NOT NULL UNIQUE,
	KayitTarih TEXT NOT NULL,
	Domain TEXT,
	SonucSay INT DEFAULT 1,
	PRIMARY KEY (id)
		
);`
	// tam arama için create stmt eklenmeli
	//FULLTEXT (Title)
	_, err5 := db.Exec(createStatement)
	yar.CheckErr(err5)
}

func DbKayit() {

	//mysql işlemleri
	db := DbBaglan()
	defer db.Close()
	tx, err2 := db.Begin()
	yar.CheckErr(err2)
	stmtk, err2 := tx.Prepare(`INSERT INTO data (Title, Price, Resim, Link, Kategori, Marka, UrunKodu, KayitTarih, Domain) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE Price=?, KayitTarih=?, SonucSay=SonucSay+1`)
	yar.CheckErr(err2)
	_, err2 = stmtk.Exec(deg.Dt.Title, deg.Dt.Price, deg.Dt.Resim, deg.Dt.Link, deg.Dt.Kategori, deg.Dt.Marka, deg.Dt.UrunKodu, deg.Dt.KayitTarih, deg.Dt.Domain, deg.Dt.Price, deg.Dt.KayitTarih)
	yar.CheckErr(err2)

	tx.Commit()

}

func DbSorgu(arama string) (satirlar []deg.Data, err error) {
	db := DbBaglan()
	defer db.Close()

	aramaDt, err := db.Query(`SELECT * FROM data WHERE MATCH(Title) AGAINST (? IN NATURAL LANGUAGE MODE)`, arama)

	if err != nil {
		return nil, err
	} else {
		satir := deg.Data{}
		satirlar := []deg.Data{}

		for aramaDt.Next() {

			err2 := aramaDt.Scan(&satir.Id, &satir.Title, &satir.Price, &satir.Resim, &satir.Link, &satir.Kategori, &satir.Marka, &satir.UrunKodu, &satir.KayitTarih, &satir.Domain, &satir.SonucSay)
			if err2 != nil {
				return nil, err2
			} else {
				satirlar = append(satirlar, satir)

			}
		}
		if len(satirlar) == 0 {

			return nil, errors.New("kayıt bulunamadı")
		}

		return satirlar, nil
	}

}
