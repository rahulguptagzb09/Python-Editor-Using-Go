package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const temp_filename string = "script.py"

func create_file(text string) {
	err := os.WriteFile(temp_filename, []byte(text), 0755)
	if err != nil {
		log.Print("unable to write file: %w", err)
	}
}

func delete_file() {
	if _, err := os.Stat(temp_filename); err == nil {
		err := os.Remove(temp_filename) //remove the file
		if err != nil {
			log.Println("Error: ", err) //print the error if file is not removed
		} else {
			log.Println("Successfully deleted file: ", temp_filename) //print success if file is removed
		}
	} else {
		log.Println("python script does not exists.")
	}
}

func copy_output(r io.Reader, output *widget.Entry) {
	scanner := bufio.NewScanner(r)
	var full_out string = ""
	for scanner.Scan() {
		txt := scanner.Text()
		full_out += txt + "\n"
	}
	log.Println("Script Output : " + full_out)
	if full_out != "" {
		output.SetText(full_out)
	}
}

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("Python Editor")
	myWindow.CenterOnScreen()
	myWindow.SetOnClosed(delete_file)
	myWindow.Resize(fyne.NewSize(900, 600))

	heading := widget.NewLabel("Python Editor")
	input_label := widget.NewLabel("Input")
	output_label := widget.NewLabel("Output")

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter code...")
	input.SetMinRowsVisible(10)

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("No Output...")
	output.SetMinRowsVisible(10)

	run_button := widget.NewButton("Run", func() {
		log.Println("Input Text is :", input.Text)
		create_file(input.Text)
		cmd := exec.Command("python", temp_filename)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		copy_output(stdout, output)
		if stderr != nil {
			copy_output(stderr, output)
		}
		cmd.Wait()
	})

	clear_button := widget.NewButton("Clear", func() {
		input.SetText("")
		output.SetText("")
		delete_file()
	})

	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), heading, layout.NewSpacer())
	buttons := container.New(layout.NewHBoxLayout(), layout.NewSpacer(),
		run_button, layout.NewSpacer(), clear_button, layout.NewSpacer())
	content := container.New(layout.NewVBoxLayout(), centered, input_label, input, buttons,
		output_label, output)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

}
