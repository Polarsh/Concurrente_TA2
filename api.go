package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

type BabyGender struct {
	Baby_Gender     int     `json:"baby_gender"`
	Euclid_Distance float64 `json:"euclid_distance"`
}

type Baby struct {
	Gender  int `json:"gender"`
	Day     int `json:"day"`
	Month   int `json:"month"`
	Weight  int `json:"weight"`
	Edadmad int `json:"edadmad"`
	Totemba int `json:"totemba"`
}

//global
var babies_gender []BabyGender
var babies []Baby
var k int = 7
var numGoRoutines int = 5

//sort
type ByEuclid_Distance []BabyGender

func (a ByEuclid_Distance) Len() int           { return len(a) }
func (a ByEuclid_Distance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEuclid_Distance) Less(i, j int) bool { return a[i].Euclid_Distance < a[j].Euclid_Distance }

//

func knn(x1, x2, x3, x4, x5, first_posi, last_posi int, doneCh chan struct{}) {
	for i := first_posi; i <= last_posi; i++ {
		//calculamos la euclediana
		euclediana := euclid(x1, x2, x3, x4, x5, babies[i].Day, babies[i].Month, babies[i].Weight, babies[i].Edadmad, babies[i].Totemba)

		//obtenemos el genero del bebe del back
		genero := babies[i].Gender

		//metemos la eucleidana y el genero a un struct
		babies_gender = append(babies_gender, BabyGender{
			Baby_Gender:     genero,
			Euclid_Distance: euclediana})
	}
	//doneCh <- struct{}{}
}

func sorting_ascendent() {
	for i := 0; i < len(babies_gender); i++ {
		sort.Sort(ByEuclid_Distance(babies_gender))
	}
}

func knn_and_sorting(x1, x2, x3, x4, x5 int) {
	final := len(babies)

	doneCh := make(chan struct{})

	for i := 0; i <= final; i = i + (final / numGoRoutines) + 1 {
		salto := i + (final / numGoRoutines)
		if salto >= final {
			salto = final - 1
		}
		//concurrencia
		knn(x1, x2, x3, x4, x5, i, salto, doneCh)
	}
	//espera que todas las rutinas se terminen
	//doneChNum := 0
	//for doneChNum < numGoRoutines {
	//	<-doneCh
	//	doneChNum++
	//}
	//fmt.Println(doneChNum)

	//ordenamos ascendentemente la distancia
	sorting_ascendent()

	//solo los interesa los k primeros, asi que cortamos el array
	babies_gender = babies_gender[:k]
}

func euclid(x1, x2, x3, x4, x5, y1, y2, y3, y4, y5 int) float64 {
	distance := math.Pow(float64(x1-y1), 2) +
		math.Pow(float64(x2-y2), 2) +
		math.Pow(float64(x3-y3), 2) +
		math.Pow(float64(x4-y4), 2) +
		math.Pow(float64(x5-y5), 2)

	return math.Sqrt(distance)
}

func count_gender() string {
	Hom := 0
	var gender string
	for i := 0; i < k; i++ {
		//contar cuantos son niños
		//comparamos 1:hombre 2:mujer
		if babies_gender[i].Baby_Gender == 1 {
			Hom++
		}
	}
	if Hom >= k/2+1 {
		gender = "niño"
	} else {
		gender = "niña"
	}
	return gender
}

func load_data() {
	csvFile, _ := http.Get("https://raw.githubusercontent.com/Polarsh/Concurrente_TA2/main/Dataset/LM2000_v7_Muestra.csv")
	reader := csv.NewReader(csvFile.Body)
	defer csvFile.Body.Close()
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		gender, _ := strconv.Atoi(line[8])
		day, _ := strconv.Atoi(line[9])
		month, _ := strconv.Atoi(line[10])
		weight, _ := strconv.Atoi(line[15])
		edadmad, _ := strconv.Atoi(line[23])
		totemba, _ := strconv.Atoi(line[38])
		babies = append(babies, Baby{
			Gender:  gender,
			Day:     day,
			Month:   month,
			Weight:  weight,
			Edadmad: edadmad,
			Totemba: totemba,
		})
	}
	//ultimo valor reemplaza al primero
	babies[0] = babies[len(babies)-1]
	//elimianos ultimos
	babies = babies[:len(babies)-1]
}

func allHandler(w http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /all")
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(babies, "", " ")
	io.WriteString(w, string(jsonBytes))
}

func genderHandler(w http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /gender")
	//w.Header().Set("Content-Type", "text/html")
	//w.Header().Set("Content-Type", "application/json")

	//recuperar parametros del front y lo pasamos a int
	day_r, _ := strconv.Atoi(request.FormValue("day"))
	month_r, _ := strconv.Atoi(request.FormValue("month"))
	weight_r, _ := strconv.Atoi(request.FormValue("weight"))
	edadmad_r, _ := strconv.Atoi(request.FormValue("edadmad"))
	totemba_r, _ := strconv.Atoi(request.FormValue("totemba"))

	//knn sacará la euclediana y lo ordenará
	knn_and_sorting(day_r, month_r, weight_r, edadmad_r, totemba_r)
	//---------

	//inicializamos el json y guardamos el genero & euclediana
	var jsonBytes []byte
	jsonBytes, _ = json.MarshalIndent(babies_gender, "", " ")

	//imprime los k vecinos
	io.WriteString(w, string(jsonBytes))
	//------------------

	//contamos #niños genero y definimos
	gender := count_gender()
	io.WriteString(w, gender)

	//eliminar datos para proximas consultas
	jsonBytes = nil
	babies_gender = nil
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
}

//8 Implementar el algoritmo indicado de manera eficiente
//3 Implementar una API REST en GO que se comunique con la interfaz Web para recibir los parámetros
//	de configuración, ejecute el algoritmo implementado y devuelva los resultados obtenidos del algoritmo.
//1 Presentación de la documentación completa.
