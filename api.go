package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

type BabyGender struct {
	Baby_Gender     string `json:"baby_gender"`
	Euclid_Distance string `json:"euclid_distance"`
}

type Baby struct {
	Gender  string `json:"gender"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	Weight  string `json:"weight"`
	Edadmad string `json:"edadmad"`
	Totemba string `json:"totemba"`
}

//global
var babies_gender []BabyGender
var babies []Baby
var k int = 7

//sort
type ByEuclid_Distance []BabyGender

func (a ByEuclid_Distance) Len() int      { return len(a) }
func (a ByEuclid_Distance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByEuclid_Distance) Less(i, j int) bool {
	aux1, _ := strconv.ParseFloat(a[i].Euclid_Distance, 32)
	aux2, _ := strconv.ParseFloat(a[j].Euclid_Distance, 32)
	return aux1 < aux2
}

//

func euclid(x1, x2, x3, x4, x5, y1, y2, y3, y4, y5 float64) string {
	distance := math.Pow(x1-y1, 2) +
		math.Pow(x2-y2, 2) +
		math.Pow(x3-y3, 2) +
		math.Pow(x4-y4, 2) +
		math.Pow(x5-y5, 2)
	//raiz
	distance = math.Sqrt(distance)
	return strconv.FormatFloat(distance, 'f', 5, 64)
}

func load_data() {
	csvFile, _ := http.Get("https://raw.githubusercontent.com/Polarsh/Concurrente_TA2/main/Dataset/LM2000_v7_Muestra.csv")
	reader := csv.NewReader(csvFile.Body)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		babies = append(babies, Baby{
			Gender:  line[8],
			Day:     line[9],
			Month:   line[10],
			Weight:  line[15],
			Edadmad: line[23],
			Totemba: line[38],
		})
	}
	//ultimo valor reemplaza al primero
	babies[0] = babies[len(babies)-1]
	//elimianos ultimos
	babies = babies[:len(babies)-1]
}

func allHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /all")
	response.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(babies, "", " ")
	io.WriteString(response, string(jsonBytes))
}

func genderHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /gender")
	response.Header().Set("Content-Type", "text/html")
	//response.Header().Set("Content-Type", "application/json")
	//recuperar parametros
	day_float, _ := strconv.ParseFloat(request.FormValue("day"), 32)
	month_float, _ := strconv.ParseFloat(request.FormValue("month"), 32)
	weight_float, _ := strconv.ParseFloat(request.FormValue("weight"), 32)
	edadmad_float, _ := strconv.ParseFloat(request.FormValue("edadmad"), 32)
	totemba_float, _ := strconv.ParseFloat(request.FormValue("totemba"), 32)
	//inicializamos el json
	var jsonBytes []byte
	//logica del endpoint
	for _, baby := range babies {
		day2_float, _ := strconv.ParseFloat(baby.Day, 32)
		month2_float, _ := strconv.ParseFloat(baby.Month, 32)
		weight2_float, _ := strconv.ParseFloat(baby.Weight, 32)
		edadmad2_float, _ := strconv.ParseFloat(baby.Edadmad, 32)
		totemba2_float, _ := strconv.ParseFloat(baby.Totemba, 32)

		//to string
		euclediana := euclid(day_float, month_float, weight_float, edadmad_float, totemba_float, day2_float, month2_float, weight2_float, edadmad2_float, totemba2_float)
		genero := baby.Gender

		//añadimos a un struct
		babies_gender = append(babies_gender, BabyGender{
			Baby_Gender:     genero,
			Euclid_Distance: euclediana,
		})
		//ordenamos
		sort.Sort(ByEuclid_Distance(babies_gender))

		//añadimos
		jsonBytes, _ = json.MarshalIndent(babies_gender, "", " ")
	}
	//printf arreglo con todos los niños
	//io.WriteString(response, string(jsonBytes))
	println(jsonBytes)
	//------------------
	Hom := 0
	var jsonBytesVecinos []byte
	var babies_gender_vecinos []BabyGender
	for i := 0; i < k; i++ {
		babies_gender_vecinos = append(babies_gender_vecinos, babies_gender[i])
		jsonBytesVecinos, _ = json.MarshalIndent(babies_gender_vecinos, "", "  ")
		//contar cuantos son niños
		tr, _ := strconv.ParseFloat(babies_gender[i].Baby_Gender, 32)
		if tr == 1 {
			Hom++
		}
	}
	var gender string
	if Hom >= k/2+1 {
		gender = "niño"
	} else {
		gender = "niña"
	}
	io.WriteString(response, gender)
	fmt.Fprintf(response, gender)
	//eliminar datos para proximas consultas
	fmt.Println(jsonBytesVecinos)
	jsonBytes = nil
	jsonBytesVecinos = nil
	babies_gender = nil
	babies_gender_vecinos = nil
}

func handle_request() {
	htmlFile := http.FileServer(http.Dir("./front"))

	//definir endpoints
	http.Handle("/", htmlFile)
	http.HandleFunc("/all", allHandler)
	http.HandleFunc("/gender", genderHandler)

	//establecer puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	load_data()
	handle_request()
	//8 Implementar el algoritmo indicado de manera eficiente
	//3 Implementar una API REST en GO que se comunique con la interfaz Web para recibir los parámetros
	//	de configuración, ejecute el algoritmo implementado y devuelva los resultados obtenidos del algoritmo.
	//1 Presentación de la documentación completa.
}
