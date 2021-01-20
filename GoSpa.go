package main

import (
	"fmt"
	"os"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io/ioutil"
	"strconv"
	"strings"
)

var path = getDirectorioActual();

var path_temp_input_file string = path+"/temp_input_file/"
var path_temp_out_files string = path+"/temp_out_files/"

var file_xlsx string = "DATAMAPPING.xlsx"
var file_yml string = "SERVICIO.yml"
var file_json string = "INPUT.json"
var file_json_2 string = "OUTPUT.json"

var inputTrama string = ""
var outputTrama string = ""
var cadenaYml = ""
var cadenaJson = ""
var cadenaJson_2 = ""

var name = ""
var tipo = ""

var filaInicial = 5

func main()  {
	// Lo primero es verificar el archivo excel exista en el directorio
	if verificaInputFileXlsx() {
		// Se lee archivo xlsx de entrada
		leerInputFileXlsx()
	}
}

func verificaInputFileXlsx() bool {
	var _, err = os.Stat(path_temp_input_file+file_xlsx)
	// Crea el archivo si no existe
	if os.IsNotExist (err) {
		fmt.Println("NO EXISTE ARCHIVO DATAMAPPING.xlsx EN LA CARPETA temp_input_file")
		return false
	}
	return true
}

func leerInputFileXlsx() {

	// se lee la cadena de entrada
	inputTrama = getInputTrama()
	inputMapping, err := excelize.OpenFile(path_temp_input_file+file_xlsx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Leer todas las celdas de la hoja inputMapping es super estricto
	fmt.Println("EL DATAMAPPING ES : ")
	rowsInputMapping, err := inputMapping.GetRows("inputMapping")
	fmt.Println(rowsInputMapping)
	fmt.Println(len(rowsInputMapping))

	// Secrea la cabeceradel .yml
	cadenaYml += "views:\n"
	cadenaYml += "  - input:\n"
	cadenaYml += "      fields:\n"
	
	// Es la fila del excel en donde comienza a mapear los name para input del .yml
	i := filaInicial

	cadenaJson += "{\n"

	for i < len(rowsInputMapping) {
		name, err = inputMapping.GetCellValue("inputMapping", "M"+strconv.Itoa(i+1))
		tipo, err = inputMapping.GetCellValue("inputMapping", "C"+strconv.Itoa(i+1))
		position, err := inputMapping.GetCellValue("inputMapping", "L"+strconv.Itoa(i+1))
		largo, err := inputMapping.GetCellValue("inputMapping", "K"+strconv.Itoa(i+1))
		required, err := inputMapping.GetCellValue("inputMapping", "N"+strconv.Itoa(i+1))
		if err != nil {
			fmt.Println(err)
			return
		}
	
		cadenaYml += "        - name: '"+name+"'\n"
		if(tipo == "CHAR") {
			tipo = "String"
			cadenaYml += "          type: '"+tipo+"'\n"
		// Fix string input	
		} else if tipo == "NUMERIC"{
			tipo = "String"
			cadenaYml += "          type: '"+tipo+"'\n"
		}		
		cadenaYml += "          position: "+position+"\n"
		cadenaYml += "          length: "+largo+"\n"

		// Valor por defecto tomamos la cadena de entrada
		pos, err := strconv.Atoi(position)
		if err != nil {
			fmt.Println(err)
			return
		}
		lar, err := strconv.Atoi(largo)
		substring := inputTrama[pos: pos+lar]

		if required == "true" {
			// no hace nada
		} else {
			cadenaYml += "          defaultValue: '"+substring+"'\n"
		}

		cadenaJson += "\t\""+name+"\": \""+substring+"\",\n"
		i++
	}

	str := strings.TrimRight(cadenaJson,",\n")
	cadenaJson = str+"\n}"

	// se lee la trama de salida y su mappignde salida
	outputTrama = getOutputTrama()
	outputMapping, err2 := excelize.OpenFile(path_temp_input_file+file_xlsx)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}
	// Leer todas las celdas de la hoja inputMapping es super estricto
	rowsOutputMapping, err3 := outputMapping.GetRows("outputMapping")
	if err3 != nil {
		fmt.Println(err.Error())
		return
	}

	// Secrea el output del .yml
	cadenaYml += "    output:\n"
	cadenaYml += "      matchers:\n"
	cadenaYml += "        - type: 'returnValue'\n"
	cadenaYml += "          condition: 'equal'\n"
	cadenaYml += "          value: '0'\n"
	cadenaYml += "          fields:\n"

	cadenaJson_2 += "{\n"
	
	j := filaInicial
	for j < len(rowsOutputMapping) {
		name, err = outputMapping.GetCellValue("outputMapping", "M"+strconv.Itoa(j+1))
		tipo, err = outputMapping.GetCellValue("outputMapping", "C"+strconv.Itoa(j+1))
		position, err := outputMapping.GetCellValue("outputMapping", "L"+strconv.Itoa(j+1))
		largo, err := outputMapping.GetCellValue("outputMapping", "K"+strconv.Itoa(j+1))
		decimal, err := outputMapping.GetCellValue("outputMapping", "E"+strconv.Itoa(j+1)) //DECIMAL
		equal, err := outputMapping.GetCellValue("outputMapping", "N"+strconv.Itoa(j+1)) //DECIMAL

		if err != nil {
			fmt.Println(err)
			return
		}

		if equal == "true" {
			cadenaYml += "            - name: '"+name+"'\n"
			if(tipo == "CHAR") {
				tipo = "String"
				cadenaYml += "              type: '"+tipo+"'\n"
			} else if tipo == "NUMERIC"{
				tipo = "Numeric"
				cadenaYml += "              type: '"+tipo+"'\n"
			}
			cadenaYml += "              position: "+position+"\n"
			cadenaYml += "              length: "+largo+"\n"

			// Valor por defecto tomamos la cadena de entrada
			pos2, err := strconv.Atoi(position)
			if err != nil {
				fmt.Println(err)
				return
			}
			lar2, err := strconv.Atoi(largo)

			// pregunto si es numerico para aplciar formatos
			if tipo == "Numeric" {
				// luego verifico si el numero tiene precision decimal
				if decimal != "[NULL]" {
					// aplico formato segun precision decimal
					str1 := "#"
					lar, err := strconv.Atoi(largo)
					dec, err := strconv.Atoi(decimal)
					if err != nil {
						fmt.Println(err)
						return
					}
					largo_formato := lar - dec
					str2 := strings.Repeat(str1, largo_formato)
					str3 := strings.Repeat(str1, dec)
					cadenaYml += "              format: '"+str2+"."+str3+"'\n"
				} else {
					cadenaYml += "              format: '#'\n"
				}
			} else if tipo == "String" {
				// agrega required true a los tipo string en la salida
				cadenaYml += "              trimRequired: true\n"
			} else {
				// no hace nada
			}

			substring_2 := outputTrama[pos2: pos2+lar2]
			cadenaJson_2 += "\t\""+name+"\": \""+substring_2+"\",\n"
		}
		
		j++;
	}

	str2 := strings.TrimRight(cadenaJson_2,",\n")
	cadenaJson_2 = str2+"\n}"

	cadenaYml += "        - type: 'returnValue'\n"
	cadenaYml += "          condition: 'notEqual'\n"
	cadenaYml += "          value: '0'\n"
	cadenaYml += "          fields:\n"

	k := filaInicial
	for k < len(rowsOutputMapping) {
		name, err = outputMapping.GetCellValue("outputMapping", "M"+strconv.Itoa(k+1))
		tipo, err = outputMapping.GetCellValue("outputMapping", "C"+strconv.Itoa(k+1))
		position, err := outputMapping.GetCellValue("outputMapping", "L"+strconv.Itoa(k+1))
		largo, err := outputMapping.GetCellValue("outputMapping", "K"+strconv.Itoa(k+1))
		decimal, err := outputMapping.GetCellValue("outputMapping", "E"+strconv.Itoa(k+1))
		notequal, err := outputMapping.GetCellValue("outputMapping", "O"+strconv.Itoa(k+1))

		if err != nil {
			fmt.Println(err)
			return
		}

		if notequal == "true" {
			cadenaYml += "            - name: '"+name+"'\n"
			if(tipo == "CHAR") {
				tipo = "String"
				cadenaYml += "              type: '"+tipo+"'\n"
			} else if tipo == "NUMERIC"{
				tipo = "Numeric"
				cadenaYml += "              type: '"+tipo+"'\n"
			}
			cadenaYml += "              position: "+position+"\n"
			cadenaYml += "              length: "+largo+"\n"

			// pregunto si es numerico para aplciar formatos
			if tipo == "Numeric" {
				// luego verifico si el numero tiene precision decimal
				if decimal != "[NULL]" {
					// aplico formato segun precision decimal
					str1 := "#"
					lar, err := strconv.Atoi(largo)
					dec, err := strconv.Atoi(decimal)
					if err != nil {
						fmt.Println(err)
						return
					}
					largo_formato := lar - dec
					str2 := strings.Repeat(str1, largo_formato)
					str3 := strings.Repeat(str1, dec)
					cadenaYml += "              format: '"+str2+"."+str3+"'\n"
				} else {
					cadenaYml += "              format: '#'\n"
				}
			} else if tipo == "String" {
				// agrega required true a los tipo string en la salida
				cadenaYml += "              trimRequired: true\n"
			} else {
				// no hace nada
			}
		}

		k++
	}

	cadenaYmlByte := []byte(cadenaYml)
	// se escribe el archivo 
	err = ioutil.WriteFile(path_temp_out_files+file_yml, cadenaYmlByte, 0644)
	if err != nil {
		panic(err)
	}

	cadenaJsonByte := []byte(cadenaJson)
	// se escribe el archivo 
	err = ioutil.WriteFile(path_temp_out_files+file_json, cadenaJsonByte, 0644)
	if err != nil {
		panic(err)
	}

	cadenaJson2Byte := []byte(cadenaJson_2)
	// se escribe el archivo 
	err = ioutil.WriteFile(path_temp_out_files+file_json_2, cadenaJson2Byte, 0644)
	if err != nil {
		panic(err)
	}

	// leer archivo
	data, err := ioutil.ReadFile(path_temp_out_files+file_yml)
	if err != nil {
		panic(err)
	}

	fmt.Println("Archivo SERVICIO.yml generado : " + string(data))
}

func getInputTrama() string {
	// Leer cadena de entrada de una hoja en especifico y de una celda en especifico
	f, err := excelize.OpenFile(path_temp_input_file+file_xlsx)
	if err != nil {
		fmt.Println(err.Error())
	}

	cell, err := f.GetCellValue("inputTrama","A1")
	if err != nil {
		println(err.Error())
	}
	fmt.Println("LA CADENA DE ENTRADA ES : ["+cell+"]")

	return cell
}

func getOutputTrama() string {
	// Leer cadena de entrada de una hoja en especifico y de una celda en especifico
	f, err := excelize.OpenFile(path_temp_input_file+file_xlsx)
	if err != nil {
		fmt.Println(err.Error())
	}

	cell, err := f.GetCellValue("outputTrama","A1")
	if err != nil {
		println(err.Error())
	}
	fmt.Println("LA CADENA DE salida ES : ["+cell+"]")

	return cell
}

func getDirectorioActual() string { 
	mydir, err := os.Getwd() 
	if err != nil { 
		fmt.Println(err) 
	} 
	fmt.Println(mydir) 
	return mydir
}