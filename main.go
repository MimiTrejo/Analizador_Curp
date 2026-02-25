package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type Persona struct {
	Nombre          string
	ApellidoPaterno string
	ApellidoMaterno string
	FechaNacimiento string
	Sexo            string
	Entidad         string
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var resultado string

		if r.Method == "POST" {

			persona := Persona{
				Nombre:          strings.ToUpper(r.FormValue("nombre")),
				ApellidoPaterno: strings.ToUpper(r.FormValue("apellidoP")),
				ApellidoMaterno: strings.ToUpper(r.FormValue("apellidoM")),
				FechaNacimiento: r.FormValue("fecha"),
				Sexo:            r.FormValue("sexo"),
				Entidad:         r.FormValue("entidad"),
			}

			curp := generarCURP(persona)
			resultado = "<h3 style='color:green'>CURP Generada: " + curp + "</h3>"
		}

		fmt.Fprintf(w, `
		<html>
		<body style="font-family:Arial; margin:40px; max-width:600px;">
			<h2>Generador de CURP</h2>
			<form method="POST">
				<input type="text" name="nombre" placeholder="Nombre" required><br><br>
				<input type="text" name="apellidoP" placeholder="Apellido Paterno" required><br><br>
				<input type="text" name="apellidoM" placeholder="Apellido Materno" required><br><br>
				<input type="date" name="fecha" required><br><br>

				<select name="sexo">
					<option value="H">Hombre</option>
					<option value="M">Mujer</option>
				</select><br><br>

				<select name="entidad">
					<option value="CL">Coahuila</option>
					<option value="QO">Querétaro</option>
				</select><br><br>

				<button type="submit">Generar CURP</button>
			</form>
			<hr>
			%s
		</body>
		</html>
		`, resultado)
	})

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func generarCURP(p Persona) string {

	fecha, _ := time.Parse("2006-01-02", p.FechaNacimiento)

	curp := string(p.ApellidoPaterno[0])
	curp += primeraVocalInterna(p.ApellidoPaterno)
	curp += string(p.ApellidoMaterno[0])
	curp += string(p.Nombre[0])

	curp += fecha.Format("060102")

	curp += p.Sexo
	curp += p.Entidad

	curp += primeraConsonanteInterna(p.ApellidoPaterno)
	curp += primeraConsonanteInterna(p.ApellidoMaterno)
	curp += primeraConsonanteInterna(p.Nombre)

	curp += "00"

	return curp
}

func primeraVocalInterna(s string) string {
	for i := 1; i < len(s); i++ {
		if strings.ContainsRune("AEIOU", rune(s[i])) {
			return string(s[i])
		}
	}
	return "X"
}

func primeraConsonanteInterna(s string) string {
	for i := 1; i < len(s); i++ {
		if unicode.IsLetter(rune(s[i])) && !strings.ContainsRune("AEIOU", rune(s[i])) {
			return string(s[i])
		}
	}
	return "X"
}
