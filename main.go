package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Solat struct {
	PrayerTime []struct {
		Hijri   string `json:"hijri"`
		Date    string `json:"date"`
		Day     string `json:"day"`
		Imsak   string `json:"imsak"`
		Fajr    string `json:"fajr"`
		Syuruk  string `json:"syuruk"`
		Dhuhr   string `json:"dhuhr"`
		Asr     string `json:"asr"`
		Maghrib string `json:"maghrib"`
		Isha    string `json:"isha"`
	} `json:"prayerTime"`
	Status     string `json:"status"`
	ServerTime string `json:"serverTime"`
	PeriodType string `json:"periodType"`
	Lang       string `json:"lang"`
	Zone       string `json:"zone"`
	Bearing    string `json:"bearing"`
}

const API = "https://www.e-solat.gov.my/index.php?r=esolatApi/takwimsolat&period=month&zone="
const DEFAULT_ZONE = "PNG01"

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
  JHR01 Pulau Aur dan Pulau Pemanggil
  JHR02 Johor Bahru, Kota Tinggi, Mersing, Kulai
  JHR03 Kluang, Pontian
  JHR04 Batu Pahat, Muar, Segamat, Gemas Johor, Tangkak

Kedah:
  KDH01 Kota Setar, Kubang Pasu, Pokok Sena (Daerah Kecil)
  KDH02 Kuala Muda, Yan, Pendang
  KDH03 Padang Terap, Sik
  KDH04 Baling
  KDH05 Bandar Baharu, Kulim
  KDH06 Langkawi
  KDH07 Puncak Gunung Jerai

Kelantan:
  KTN01 Bachok, Kota Bharu, Machang, Pasir Mas, Pasir Puteh, Tanah Merah, Tumpat, Kuala Krai, Mukim Chiku
  KTN02 Gua Musang (Daerah Galas Dan Bertam), Jeli, Jajahan Kecil Lojing

Melaka:
  MLK01 SELURUH NEGERI MELAKA

Negeri Sembilan:
  NGS01 Tampin, Jempol
  NGS02 Jelebu, Kuala Pilah, Rembau
  NGS03 Port Dickson, Seremban

Pahang:
  PHG01 Pulau Tioman
  PHG02 Kuantan, Pekan, Rompin, Muadzam Shah
  PHG03 Jerantut, Temerloh, Maran, Bera, Chenor, Jengka
  PHG04 Bentong, Lipis, Raub
  PHG05 Genting Sempah, Janda Baik, Bukit Tinggi
  PHG06 Cameron Highlands, Genting Higlands, Bukit Fraser

Perlis:
  PLS01 Kangar, Padang Besar, Arau

Pulau Pinang:
  PNG01 Seluruh Negeri Pulau Pinang

Perak:
  PRK01 Tapah, Slim River, Tanjung Malim
  PRK02 Kuala Kangsar, Sg. Siput , Ipoh, Batu Gajah, Kampar
  PRK03 Lenggong, Pengkalan Hulu, Grik
  PRK04 Temengor, Belum
  PRK05 Kg Gajah, Teluk Intan, Bagan Datuk, Seri Iskandar, Beruas, Parit, Lumut, Sitiawan, Pulau Pangkor
  PRK06 Selama, Taiping, Bagan Serai, Parit Buntar
  PRK07 Bukit Larut

Sabah:
  SBH01 Bahagian Sandakan (Timur), Bukit Garam, Semawang, Temanggong, Tambisan, Bandar Sandakan, Sukau
  SBH02 Beluran, Telupid, Pinangah, Terusan, Kuamut, Bahagian Sandakan (Barat)
  SBH03 Lahad Datu, Silabukan, Kunak, Sahabat, Semporna, Tungku, Bahagian Tawau  (Timur)
  SBH04 Bandar Tawau, Balong, Merotai, Kalabakan, Bahagian Tawau (Barat)
  SBH05 Kudat, Kota Marudu, Pitas, Pulau Banggi, Bahagian Kudat
  SBH06 Gunung Kinabalu
  SBH07 Kota Kinabalu, Ranau, Kota Belud, Tuaran, Penampang, Papar, Putatan, Bahagian Pantai Barat
  SBH08 Pensiangan, Keningau, Tambunan, Nabawan, Bahagian Pendalaman (Atas)
  SBH09 Beaufort, Kuala Penyu, Sipitang, Tenom, Long Pasia, Membakut, Weston, Bahagian Pendalaman (Bawah)

Selangor:
  SGR01 Gombak, Petaling, Sepang, Hulu Langat, Hulu Selangor, S.Alam
  SGR02 Kuala Selangor, Sabak Bernam
  SGR03 Klang, Kuala Langat

Sarawak:
  SWK01 Limbang, Lawas, Sundar, Trusan
  SWK02 Miri, Niah, Bekenu, Sibuti, Marudi
  SWK03 Pandan, Belaga, Suai, Tatau, Sebauh, Bintulu
  SWK04 Sibu, Mukah, Dalat, Song, Igan, Oya, Balingian, Kanowit, Kapit
  SWK05 Sarikei, Matu, Julau, Rajang, Daro, Bintangor, Belawai
  SWK06 Lubok Antu, Sri Aman, Roban, Debak, Kabong, Lingga, Engkelili, Betong, Spaoh, Pusa, Saratok
  SWK07 Serian, Simunjan, Samarahan, Sebuyau, Meludam
  SWK08 Kuching, Bau, Lundu, Sematan
  SWK09 Zon Khas (Kampung Patarikan)

Terengganu:
  TRG01 Kuala Terengganu, Marang, Kuala Nerus
  TRG02 Besut, Setiu
  TRG03 Hulu Terengganu
  TRG04 Dungun, Kemaman

Wilayah Persekutuan:
  WLY01 Kuala Lumpur, Putrajaya
  WLY02 Labuan
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
	_, _, today := now.Date()

	for index, slt := range solat.PrayerTime {
		var current_day = index + 1

		fmt.Printf("Tarikh: %s, Imsak: %s, Subuh: %s, Syuruk: %s, Zohor: %s, Asar: %s, Maghrib: %s, Isyak: %s",
			slt.Date, slt.Imsak, slt.Fajr, slt.Syuruk, slt.Dhuhr, slt.Asr, slt.Maghrib, slt.Isha)

		if current_day == today {
			fmt.Println(" ****")
		} else {
			fmt.Println()
		}
	}
}
