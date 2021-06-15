package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Baby struct {
	Sex     string `json:"sex"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	Weight  string `json:"weight"`
	Edadmad string `json:"edadmad"`
	Totemba string `json:"totemba"`
}

//global
var babies []Baby

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
			Sex:     line[8],
			Day:     line[9],
			Month:   line[10],
			Weight:  line[15],
			Edadmad: line[23],
			Totemba: line[38],
		})
	}
}

func homeHandler(response http.ResponseWriter, request *http.Request) {
	htmlFile, error := template.ParseFiles("index.html")
	if error != nil {
		log.Fatal(error)
	}
	htmlFile.Execute(response, "")
}

func allHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /all")
	//tipo de contenido de rpta
	response.Header().Set("Content-Type", "application/json")
	//serializar, codificar el resultado a json
	jsonBytes, _ := json.MarshalIndent(babies, "", " ")
	io.WriteString(response, string(jsonBytes))

}

func sexHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /sex")
	//tipo de contenido de rpta
	response.Header().Set("Content-Type", "application/json")

	//recuperar parametros
	day := request.FormValue("day")
	month := request.FormValue("month")
	weight := request.FormValue("weight")
	edadmad := request.FormValue("edadmad")
	totemba := request.FormValue("totemba")

	day_float, _ := strconv.ParseFloat(day, 32)
	month_float, _ := strconv.ParseFloat(month, 32)
	weight_float, _ := strconv.ParseFloat(weight, 32)
	edadmad_float, _ := strconv.ParseFloat(edadmad, 32)
	totemba_float, _ := strconv.ParseFloat(totemba, 32)

	response.Write([]byte(day + month + weight + edadmad + totemba))
	fmt.Println(day_float, month_float, weight_float, edadmad_float, totemba_float)
	//logica del endpoint

	//for _, Baby := range babies {
	//if baby.DNI == dni {
	//	//codificarlo
	//	jsonBytes, _ := json.MarshalIndent(baby, "", " ")
	//	io.WriteString(response, string(jsonBytes))
	//}
	//}
}

func creditsHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /credits")
	response.Header().Set("Content-Type", "text/html")
	io.WriteString(response,
		`<doctype html>
	<html>
		<head>
			<title>API</title>
		</head>
		<body>
			<h2>
				Api desarrollada para TA2 de Concurrente
				by: u201723059 Deyvidyorch Sanchez
			</h2>
		</body>
	</html>
	`)
}

func handle_request() {
	//definir endpoints
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/all", allHandler)
	http.HandleFunc("/sex", sexHandler)
	http.HandleFunc("/credits", creditsHandler)

	//establecer puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	load_data()
	handle_request()

	//8 Implementar el algoritmo indicado de manera eficiente
	//4 Implementar una interfaz Web que permita configurar los parámetros y mostrar resultados.
	//3 Implementar una API REST en GO que se comunique con la interfaz Web para recibir los parámetros
	//	de configuración, ejecute el algoritmo implementado y devuelva los resultados obtenidos del algoritmo.
	//1 Presentación de la documentación completa.
}
