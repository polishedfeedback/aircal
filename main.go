package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strconv"
	"strings"
)

type Pallet struct {
	Type   string
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Calculator interface {
	CalculateChargableWeight() float64
}

func (p Pallet) CalculateChargableWeight() float64 {
	volumetricWeight := (p.Length * p.Width * p.Height) / 6000
	if p.Weight > volumetricWeight {
		return p.Weight
	}
	return volumetricWeight
}

type Model struct {
	palletType    string
	textInput     textinput.Model
	step          int
	pallets       []Pallet
	selectedType  Pallet
	err           error
}

var (
	ukPallet = Pallet{
		Type:   "UK",
		Length: 120,
		Width:  100,
	}
	euPallet = Pallet{
		Type:   "EU",
		Length: 120,
		Width:  80,
	}
)

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your choice"
	ti.Focus()

	return Model{
		textInput: ti,
		step:      1,
		pallets:   make([]Pallet, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			switch m.step {
			case 1: 
				choice := strings.ToUpper(strings.TrimSpace(m.textInput.Value()))
				if choice == "UK" {
					m.selectedType = ukPallet
					m.step++
				} else if choice == "EU" {
					m.selectedType = euPallet
					m.step++
				} else {
					m.err = fmt.Errorf("invalid pallet type")
					return m, cmd
				}
				m.textInput.Reset()
				
			case 2:
				inputs := strings.FieldsFunc(m.textInput.Value(), func(r rune) bool {
					return r == ' ' || r == ';'
				})
				
				for _, input := range inputs {
					parts := strings.Split(strings.TrimSpace(input), ",")
					if len(parts) != 2 {
						m.err = fmt.Errorf("invalid input format. Use 'height,weight' for each pallet")
						return m, cmd
					}
					
					height, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
					if err != nil {
						m.err = fmt.Errorf("invalid height input")
						return m, cmd
					}
					
					weight, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					if err != nil {
						m.err = fmt.Errorf("invalid weight input")
						return m, cmd
					}
					
					pallet := m.selectedType
					pallet.Height = height
					pallet.Weight = weight
					m.pallets = append(m.pallets, pallet)
				}
				m.step++
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var s strings.Builder

	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %v\n", m.err))
		s.WriteString("Please try again\n")
		m.err = nil
	}

	switch m.step {
	case 1:
		s.WriteString("Select pallet type (UK/EU):\n")
		s.WriteString("UK: 120cm x 100cm\n")
		s.WriteString("EU: 120cm x 80cm\n")
		s.WriteString(m.textInput.View())
		
	case 2:
		s.WriteString(fmt.Sprintf("Selected %s pallet\n", m.selectedType.Type))
		s.WriteString("Enter pallet details (height,weight) separated by spaces or semicolons:\n")
		s.WriteString("Example: 150,350 160,400 or 150,350;160,400\n")
		s.WriteString(m.textInput.View())
		
	case 3:
		s.WriteString("Calculation Results:\n\n")
		var totalChargableWeight float64
		
		for i, pallet := range m.pallets {
			volumetricWeight := (pallet.Length * pallet.Width * pallet.Height) / 6000
			chargableWeight := pallet.CalculateChargableWeight()
			
			s.WriteString(fmt.Sprintf("Pallet %d:\n", i+1))
			s.WriteString(fmt.Sprintf("  Dimensions: %.0fcm x %.0fcm x %.0fcm\n", 
				pallet.Length, pallet.Width, pallet.Height))
			s.WriteString(fmt.Sprintf("  Actual Weight: %.2f kg\n", pallet.Weight))
			s.WriteString(fmt.Sprintf("  Volumetric Weight: %.2f kg\n", volumetricWeight))
			s.WriteString(fmt.Sprintf("  Chargable Weight: %.2f kg", chargableWeight))
			if chargableWeight == pallet.Weight {
				s.WriteString(" (Actual Weight)")
			} else {
				s.WriteString(" (Volumetric Weight)")
			}
			s.WriteString("\n\n")
			
			totalChargableWeight += chargableWeight
		}
		s.WriteString(fmt.Sprintf("Total Chargable Weight: %.2f kg\n", totalChargableWeight))
		s.WriteString("\nPress Ctrl+C to exit")
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
