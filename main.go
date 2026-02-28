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
				Nombre:          strings.ToUpper(strings.TrimSpace(r.FormValue("nombre"))),
				ApellidoPaterno: strings.ToUpper(strings.TrimSpace(r.FormValue("apellidoP"))),
				ApellidoMaterno: strings.ToUpper(strings.TrimSpace(r.FormValue("apellidoM"))),
				FechaNacimiento: r.FormValue("fecha"),
				Sexo:            r.FormValue("sexo"),
				Entidad:         r.FormValue("entidad"),
			}

			// 🔐 Validar nombres
			if !nombreValido(persona.Nombre) ||
				!nombreValido(persona.ApellidoPaterno) ||
				!nombreValido(persona.ApellidoMaterno) {

				resultado = "<h3 style='color:red'>Nombre y apellidos deben tener mínimo 3 letras y solo contener letras.</h3>"

			} else if !fechaValida(persona.FechaNacimiento) {

				resultado = "<h3 style='color:red'>Fecha inválida (máximo 120 años o fecha futura).</h3>"

			} else {

				curp := generarCURP(persona)
				resultado = "<h3 style='color:green'>CURP Generada: " + curp + "</h3>"
			}
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

func nombreValido(s string) bool {

	if len([]rune(s)) < 3 {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func generarCURP(p Persona) string {

	fecha, err := time.Parse("2006-01-02", p.FechaNacimiento)
	if err != nil {
		return ""
	}

	curp := string([]rune(p.ApellidoPaterno)[0])
	curp += primeraVocalInterna(p.ApellidoPaterno)
	curp += string([]rune(p.ApellidoMaterno)[0])
	curp += string([]rune(p.Nombre)[0])

	curp += fecha.Format("060102")

	curp += p.Sexo
	curp += p.Entidad

	curp += primeraConsonanteInterna(p.ApellidoPaterno)
	curp += primeraConsonanteInterna(p.ApellidoMaterno)
	curp += primeraConsonanteInterna(p.Nombre)

	curp += "00"

	return curp
}

func fechaValida(fechaStr string) bool {

	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		return false
	}

	hoy := time.Now().UTC()
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, time.UTC)

	if fecha.After(hoy) {
		return false
	}

	if fecha.Before(hoy.AddDate(-120, 0, 0)) {
		return false
	}

	return true
}

func primeraVocalInterna(s string) string {
	runes := []rune(s)
	for i := 1; i < len(runes); i++ {
		if strings.ContainsRune("AEIOU", runes[i]) {
			return string(runes[i])
		}
	}
	return "X"
}

func primeraConsonanteInterna(s string) string {
	runes := []rune(s)
	for i := 1; i < len(runes); i++ {
		if unicode.IsLetter(runes[i]) && !strings.ContainsRune("AEIOU", runes[i]) {
			return string(runes[i])
		}
	}
	return "X"
}
