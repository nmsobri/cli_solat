package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Solat struct {
	Day []struct {
		Date      string `json:"date"`
		Imsak     string `json:"imsak"`
		Subuh     string `json:"subuh"`
		Syuruk    string `json:"syuruk"`
		Zohor     string `json:"zohor"`
		Asar      string `json:"asar"`
		Maghrib   string `json:"maghrib"`
		Iswak     string `json:"iswak"`
		Direction string `json:"direction"`
	} `json:"solat"`
}

const API = "https://cms.waktusolat.digital/esolatjson.php?zon="
const DEFAULT_ZONE = "png01"

func getJson(url string) ([]byte, error) {
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	res, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func showHelp() {
	fmt.Println("Usage: solat [options]")
	fmt.Println("Options:")
	fmt.Println("  -help    Display this help message")
	fmt.Println("  -zones   Display all available zone to select")
	fmt.Println("  -zone    Select zone to display the prayer time")
}

func showZone() {
	fmt.Println("All available zone")

	zone := `
Johor:
  jhr01 Pemanggil
  jhr02 Kota Tinggi
  jhr02 Mersing
  jhr02 Johor Bahru
  jhr03 Kluang
  jhr03 Pontian
  jhr04 Batu Pahat
  jhr04 Muar
  jhr04 Segamat
  jhr04 Gemas

Kedah
  kdh01 Kota Setar
  kdh01 Kubang Pasu
  kdh01 Pokok Sena
  kdh02 Pendang
  kdh02 Kuala Muda
  kdh02 Yan
  kdh03 Padang Terap
  kdh03 Sik
  kdh04 Baling
  kdh05 Kulim
  kdh05 Bandar Baharu
  kdh06 Langkawi
  kdh07 Gunung Jerai

Kelantan
  ktn01 Kota Bahru
  ktn01 Bachok
  ktn01 Pasir Puteh
  ktn01 Tumpat
  ktn01 Pasir Mas
  ktn01 Tanah Merah
  ktn01 Machang
  ktn01 Kuala Krai
  ktn01 Mukim Chiku
  ktn03 Jeli
  ktn03 Gua Musang
  ktn03 Mukim Galas
  ktn03 Bertam

Melaka
  mlk01 Bandar Melaka
  mlk01 Alor Gajah
  mlk01 Jasin
  mlk01 Masjid Tanah
  mlk01 Merlimau
  mlk01 Nyalas

Negeri Sembilan
  ngs01 Jempol
  ngs01 Tampin
  ngs02 Port Dickson
  ngs02 Seremban
  ngs02 Kuala Pilah
  ngs02 Jelebu
  ngs02 Rembau

Pahang
  phg01 Pulau Tioman
  phg02 Kuantan
  phg02 Pekan
  phg02 Rompin
  phg02 Muadzam Shah
  phg03 Maran
  phg03 Chenor
  phg03 Temerloh
  phg03 Bera
  phg03 Jerantut
  phg04 Bentong
  phg04 Raub
  phg04 Kuala Lipis
  phg05 Genting Sempah
  phg05 Janda Baik
  phg05 Bukit Tinggi
  phg06 Bukit Fraser
  phg06 Genting Highlands
  phg06 Cameron Highlands

Perak
  prk01 Tapah
  prk01 Slim River
  prk01 Tanjung Malim
  prk02 Ipoh
  prk02 Batu Gajah
  prk02 Kampar
  prk02 Sungai Siput
  prk02 Kuala Kangsar
  prk03 Pengkalan Hulu
  prk03 Grik
  prk03 Lenggong
  prk04 Temenggung
  prk04 Belum
  prk05 Teluk Intan
  prk05 Bagan Datoh
  prk05 Kampung Gajah
  prk05 Sri Iskandar
  prk05 Beruas
  prk05 Parit
  prk05 Lumut
  prk05 Setiawan
  prk05 Pulau Pangkor
  prk06 Selama
  prk06 Taiping
  prk06 Bagan Serai
  prk06 Parit Buntar
  prk07 Bukit Larut

Perlis
  pls01 Kangar
  pls01 Padang Besar
  pls01 Arau

Pulau Pinang
  png01 Pulau Pinang

Sabah
  sbh01 Sandakan
  sbh01 Bandar Bukit Garam
  sbh01 Semawang
  sbh01 Temanggong
  sbh01 Tambisan
  sbh02 Pinangah
  sbh02 Terusan
  sbh02 Beluran
  sbh02 Kuamut
  sbh02 Telupit
  sbh03 Lahad Datu
  sbh03 Kunak
  sbh03 Silabukan
  sbh03 Tungku
  sbh03 Sahabat
  sbh03 Semporna
  sbh04 Tawau
  sbh04 Balong
  sbh04 Merotai
  sbh04 Kalabakan
  sbh05 Kudat
  sbh05 Kota Marudu
  sbh05 Pitas
  sbh05 Pulau Banggi
  sbh06 Gunung Kinabalu
  sbh07 Papar
  sbh07 Ranau
  sbh07 Kota Belud
  sbh07 Tuaran
  sbh07 Penampang
  sbh07 Kota Kinabalu
  sbh08 Pensiangan
  sbh08 Keningau
  sbh08 Tambunan
  sbh08 Nabawan
  sbh09 Sipitang
  sbh09 Membakut
  sbh09 Beaufort
  sbh09 Kuala Penyu
  sbh09 Weston
  sbh09 Tenom
  sbh09 Long Pa Sia

Sarawak
  swk01 Limbang
  swk01 Sundar
  swk01 Terusan
  swk01 Lawas
  swk02 Niah
  swk02 Belaga
  swk02 Sibuti
  swk02 Miri
  swk02 Bekenu
  swk02 Marudi
  swk03 Song
  swk03 Balingian
  swk03 Sebauh
  swk03 Bintulu
  swk03 Tatau
  swk03 Kapit
  swk04 Igan
  swk04 Kanowit
  swk04 Sibu
  swk04 Dalat
  swk04 Oya
  swk05 Belawai
  swk05 Matu
  swk05 Daro
  swk05 Sarikei
  swk05 Julau
  swk05 Bitangor
  swk05 Rajang
  swk06 Kabong
  swk06 Lingga
  swk06 Sri Aman
  swk06 Engkelili
  swk06 Betong
  swk06 Spaoh
  swk06 Pusa
  swk06 Saratok
  swk06 Roban
  swk06 Debak
  swk07 Samarahan
  swk07 Simunjan
  swk07 Serian
  swk07 Sebuyau
  swk07 Meludam
  swk08 Kuching
  swk08 Bau
  swk08 Lundu
  swk08 Sematan
  swk09 Zon Khas
      
Selangor
  sgr01 Gombak
  sgr01 Hulu Selangor
  sgr01 Rawang
  sgr01 Hulu Langat
  sgr01 Sepang
  sgr01 Petaling Jaya
  sgr01 Shah Alam
  sgr02 Sabak Bernam
  sgr02 Kuala Selangor
  sgr03 Klang
  sgr03 Kuala Langat

Terengganu
  trg01 Kuala Terengganu
  trg01 Marang
  trg02 Besut
  trg02 Setiu
  trg03 Hulu Terengganu
  trg04 Kemaman
  trg04 Dungun

Wilayah Persekutuan
  wly01 Putrajaya
  wly01 Kuala Lumpur
  wly02 Labuan
`

	fmt.Println(zone)
}

func main() {
	var help bool
	var zones bool
	var zone string

	flag.BoolVar(&help, "help", false, "display help message")
	flag.BoolVar(&zones, "zones", false, "display all available zone to select")
	flag.StringVar(&zone, "zone", "", "select zone to display the prayer time")
	flag.Parse()

	if help {
		showHelp()
		return
	}

	if zones {
		showZone()
		return
	}

	if zone == "" {
		zone = DEFAULT_ZONE
		fmt.Println("Warning no zone had been selected")
		fmt.Printf("Choosing default location (%s)\n", DEFAULT_ZONE)
	}

	requestURL := fmt.Sprintf("%s%s", API, zone)
	fmt.Println("Fetching...")
	body, err := getJson(requestURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	if len(body) == 0 {
		fmt.Fprintf(os.Stderr, "%s", "Cannot get response from the server. Probably invalid zone code.")
		return
	}

	var solat Solat

	if err := json.Unmarshal(body, &solat); err != nil {
		fmt.Fprintf(os.Stderr, "%s", "Can not unmarshal JSON")
		return
	}

	now := time.Now()
	year, month, day := now.Date()

	for _, slt := range solat.Day {
		thisDate := fmt.Sprintf("%s/%d/%d", slt.Date, month, year)
		currentDate, err := strconv.Atoi(slt.Date)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return
		}

		fmt.Printf("Tarikh: %s, Imsak: %s, Subuh: %s, Syuruk: %s, Zohor: %s, Asar: %s, Maghrib: %s, Isyak: %s",
			thisDate, slt.Imsak, slt.Subuh, slt.Syuruk, slt.Zohor, slt.Asar, slt.Maghrib, slt.Iswak)

		if day == currentDate {
			fmt.Println(" ****")
		} else {
			fmt.Println()
		}
	}
}
