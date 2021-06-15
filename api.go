package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Shed struct {
	Department  string `json:"department"`
	Province    string `json:"province"`
	District    string `json:"district"`
	Sector      string `json:"sector"`
	Beneficiary string `json:"beneficiary"`
	DNI         string `json:"dni"`
	Altitude    string `json:"altitude"`
	Ordinance   string `json:"ordinance"`
}

//global
var sheds []Shed

func load_data() {
	csvFile, _ := http.Get("https://raw.githubusercontent.com/Polarsh/Concurrente_TA2/main/TA2/Dataset/prueba.csv")
	reader := csv.NewReader(csvFile.Body)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		sheds = append(sheds, Shed{
			Department:  line[0],
			Province:    line[1],
			District:    line[2],
			Sector:      line[3],
			Beneficiary: line[4],
			DNI:         line[5],
			Altitude:    line[6],
			Ordinance:   line[7],
		})
	}
}

func solve_list(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /list")
	//tipo de contenido de rpta
	response.Header().Set("Content-Type", "application/json")
	//serializar, codificar el resultado a json
	jsonBytes, _ := json.MarshalIndent(sheds, "", " ")
	io.WriteString(response, string(jsonBytes))

}

func solve_search_dni(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /search_dni")
	//recuperar parametros
	dni := request.FormValue("dni")
	//tipo de contenido de rpta
	response.Header().Set("Content-Type", "application/json")
	//logica del endpoint
	for _, shed := range sheds {
		if shed.DNI == dni {
			//codificarlo
			jsonBytes, _ := json.MarshalIndent(shed, "", " ")
			io.WriteString(response, string(jsonBytes))
		}
	}
}

func solve_search_department(response http.ResponseWriter, request *http.Request) {
	log.Println("endpoint /search_department")
	//recuperar parametros
	depa := request.FormValue("department")
	//tipo de contenido de rpta
	response.Header().Set("Content-Type", "application/json")
	//logica del endpoint
	for _, shed := range sheds {
		if shed.Department == depa {
			//codificarlo
			jsonBytes, _ := json.MarshalIndent(shed, "", " ")
			io.WriteString(response, string(jsonBytes))
		}
	}
}

func solve_credits(response http.ResponseWriter, request *http.Request) {
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
	http.HandleFunc("/get/all", solve_list)
	http.HandleFunc("/get/dni", solve_search_dni)
	http.HandleFunc("/get/dep", solve_search_department)
	http.HandleFunc("/credits", solve_credits)

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
