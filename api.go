package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
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

var k int = 7	//cantidad de vecinos más cercanos a analizar
var numGoRoutines int = 5	//cantidad de gorutinas

//

func knn(x1, x2, x3, x4, x5, first_posi, last_posi int, doneCh chan struct{}) {
	for i := first_posi; i <= last_posi; i++ {
		//calculamos la euclediana enre los datos del front
		//	con cada dato registro del csvfile
		euclediana := euclid(x1, x2, x3, x4, x5, babies[i].Day,
			babies[i].Month, babies[i].Weight, babies[i].Edadmad, babies[i].Totemba)

		//obtenemos el genero del bebe del back
		genero := babies[i].Gender

		//metemos la eucleidana y el genero a un struct
		babies_gender = append(babies_gender, BabyGender{
			Baby_Gender:     genero,
			Euclid_Distance: euclediana})
	}
	//mandamos señal que la goroutine terminó
	doneCh <- struct{}{}
}

func sorting_ascendent() {
	var min float64
	var location int
	//como solo nos interesa saber los k vecinos más cercanos
	//	solo ordenaremos los k primeros más cercanos
	//		min guarda el minimo en cada interacion
	//			location guarda al locacion del minimo en esa
	//				iteración
	for i := 0; i < k; i++ {
		for j := i; j < len(babies_gender); j++ {
			if j == i {
				min = babies_gender[j].Euclid_Distance
			}
			if min > babies_gender[j].Euclid_Distance {
				min = babies_gender[j].Euclid_Distance
				location = j
			}
		}
		//para no perder un posible minimo intercambiaremos
		//	la data en un aux
		dist := babies_gender[i].Euclid_Distance
		gend := babies_gender[i].Baby_Gender

		babies_gender[i] = babies_gender[location]

		babies_gender[location].Euclid_Distance = dist
		babies_gender[location].Baby_Gender = gend
	}
	//recortaremos el total de registros a solo los k primeros
	babies_gender = babies_gender[:k]
}

func knn_and_sorting(x1, x2, x3, x4, x5 int) {
	//creamos el channel que nos dirá cada que acaba un knn
	doneCh := make(chan struct{})

	//este for divide en parte iguales el total de registros para que
	// cada goroutine tenga el mismo trabajo
	for i := 0; i <= len(babies); i = i + (len(babies) / numGoRoutines) + 1 {
		salto := i + (len(babies) / numGoRoutines)
		if salto >= len(babies) {
			salto = len(babies) - 1
		}
		//concurrencia
		//mandamos los datos del front y de que posicion a que posicion
		//	trabajará cada goroutine, tbm mandamos el canal para que nos
		//		mande un señal de finalizacion
		go knn(x1, x2, x3, x4, x5, i, salto, doneCh)
	}
	//espera que todas las rutinas se terminen
	doneChNum := 0
	for doneChNum < numGoRoutines {
		<-doneCh
		doneChNum++
	}
	//ordenamos ascendentemente segun la distancia
	sorting_ascendent()
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
	csvFile, _ := http.Get("https://raw.githubusercontent.com/Polarsh/Concurrente_TA2/main/Dataset/LM2000_v7.csv")
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

	//recuperar parametros del front y lo convertimos a int para hacer los cálculos
	day_r, _ := strconv.Atoi(request.FormValue("day"))
	month_r, _ := strconv.Atoi(request.FormValue("month"))
	weight_r, _ := strconv.Atoi(request.FormValue("weight"))
	edadmad_r, _ := strconv.Atoi(request.FormValue("edadmad"))
	totemba_r, _ := strconv.Atoi(request.FormValue("totemba"))

	//mandamos los datos para que se procesen y ordenen
	knn_and_sorting(day_r, month_r, weight_r, edadmad_r, totemba_r)
	//---------

	//contamos cuantos hay de cada sexo y definimos cual hay más
	gender := count_gender()

	//mandamos al front
	io.WriteString(w, "El bebe es "+gender)

	//eliminar datos para proximas consultas
	babies_gender = nil
}

func handle_request() {
	//abrimos el index.html
	htmlFile := http.FileServer(http.Dir("./front"))

	//definimos endpoints
	http.Handle("/", htmlFile)
	http.HandleFunc("/all", allHandler)
	http.HandleFunc("/gender", genderHandler)

	//establecemos puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	load_data()
	handle_request()
}
