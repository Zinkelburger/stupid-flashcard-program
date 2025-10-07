package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application states
type appState int

const (
	stateFlashcard appState = iota
	stateWaitingForAnswer
	stateShowingAnswer
	stateLoading
	stateQuit
)

// Main model
type model struct {
	state           appState
	flashcards      []flashcard
	currentIdx      int
	sessionProgress int // Track cards completed in this session
	totalCards      int // Total cards that were due at start
	textarea        textarea.Model
	spinner         spinner.Model
	question        string
	userAnswer      string
	evaluation      string
	correct         bool
	loadingMsg      string
	width           int
	height          int
}

// Flashcard represents a single flashcard with spaced repetition data
type flashcard struct {
	Front      string
	Back       string
	NextReview time.Time
	Difficulty string
}

type evaluationResponse struct {
	Result      string `json:"result"`
	Explanation string `json:"explanation"`
}

// Styles with dynamic terminal width wrapping
var (
	questionStyle = lipgloss.NewStyle().MarginBottom(1)
	answerStyle   = lipgloss.NewStyle().PaddingLeft(2)
)

// --- Ollama API and AI Logic ---

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type aiResponseMsg struct {
	content string
	msgType string // "question", "evaluation"
}

type errorMsg struct{ err error }

const (
	OLLAMA_URL = "http://localhost:11434/api/generate"
	MODEL_NAME = "llama3.2"
)

// --- Ollama Service Management ---

// checkOllamaRunning checks if Ollama is running by making a health check request
func checkOllamaRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// startOllama starts Ollama in the background if it's not already running
func startOllama() error {
	if checkOllamaRunning() {
		return nil // Already running
	}

	log.Println("Ollama not detected, starting Ollama service...")

	// Start ollama serve in the background
	cmd := exec.Command("ollama", "serve")
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start Ollama: %v", err)
	}

	// Wait a moment for Ollama to start
	time.Sleep(2 * time.Second)

	// Check if it's now running
	if !checkOllamaRunning() {
		return fmt.Errorf("Ollama failed to start properly")
	}

	log.Println("Ollama service started successfully")
	return nil
}

func callOllama(prompt, msgType string) tea.Cmd {
	return func() tea.Msg {
		reqBody := ollamaRequest{
			Model:  MODEL_NAME,
			Prompt: prompt,
			Stream: false,
		}
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return errorMsg{err}
		}
		resp, err := http.Post(OLLAMA_URL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return errorMsg{err}
		}
		defer resp.Body.Close()
		var ollamaResp ollamaResponse
		err = json.NewDecoder(resp.Body).Decode(&ollamaResp)
		if err != nil {
			return errorMsg{err}
		}
		return aiResponseMsg{
			content: strings.TrimSpace(ollamaResp.Response),
			msgType: msgType,
		}
	}
}

func getVariedQuestion(front string) tea.Cmd {
	prompt := fmt.Sprintf("You are helping with flashcard study. The flashcard front says: '%s'\n\nCreate a natural, conversational way to ask about this topic. Be engaging. Keep it under 2 sentences. Just provide the question, no extra formatting.", front)
	return callOllama(prompt, "question")
}

// --- UPDATED: Prompt now requests JSON output ---
func evaluateAnswer(back, userAnswer string) tea.Cmd {
	prompt := fmt.Sprintf(`
Analyze the user's answer compared to the correct answer.
Correct answer: "%s"
User answer: "%s"

Respond with ONLY a valid JSON object in the following format:
{"result": "correct" or "incorrect", "explanation": "a brief 1-sentence explanation"}
`, back, userAnswer)
	return callOllama(prompt, "evaluation")
}

func initialModel() model {
	cards, err := loadFlashcards("flashcards.csv")
	if err != nil {
		log.Fatalf("Could not load flashcards: %v", err)
	}

	ta := textarea.New()
	ta.Placeholder = "Type your answer..."
	ta.Focus()
	ta.SetHeight(3)
	ta.SetWidth(60)
	ta.ShowLineNumbers = false

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle()

	return model{
		state:           stateLoading,
		flashcards:      cards,
		sessionProgress: 0,
		totalCards:      len(cards),
		textarea:        ta,
		spinner:         s,
		loadingMsg:      "Getting your flashcard ready",
	}
}

func (m model) Init() tea.Cmd {
	if len(m.flashcards) > 0 {
		return tea.Batch(
			m.spinner.Tick,
			getVariedQuestion(m.flashcards[m.currentIdx].Front),
		)
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4) // Leave some padding
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateWaitingForAnswer:
			switch msg.String() {
			case "ctrl+c":
				m.state = stateQuit
				return m, tea.Quit
			case "enter":
				if m.textarea.Value() != "" {
					m.userAnswer = m.textarea.Value()
					m.state = stateLoading
					m.loadingMsg = "Evaluating your answer"
					m.textarea.Reset()
					return m, tea.Batch(
						m.spinner.Tick,
						evaluateAnswer(m.flashcards[m.currentIdx].Back, m.userAnswer),
					)
				}
			case "ctrl+s":
				m.userAnswer = "(showed answer)"
				m.state = stateShowingAnswer
				m.evaluation = "You chose to see the answer."
				m.correct = false
				return m, nil
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}

		case stateShowingAnswer:
			switch msg.String() {
			case "ctrl+c":
				m.state = stateQuit
				return m, tea.Quit
			case "1", "2", "3", "4":
				difficulty := map[string]string{
					"1": "Again", "2": "Hard", "3": "Good", "4": "Easy",
				}[msg.String()]
				m.updateCardDifficulty(difficulty)

				m.sessionProgress++

				cards, err := loadFlashcards("flashcards.csv")
				if err != nil {
					m.evaluation = "Error loading flashcards: " + err.Error()
					return m, nil
				}
				m.flashcards = cards

				if len(m.flashcards) == 0 {
					m.evaluation = "No more cards due for review! Great job! 🎉"
					return m, nil
				}

				m.currentIdx = 0
				m.state = stateLoading
				m.loadingMsg = "Getting next flashcard"
				return m, tea.Batch(
					m.spinner.Tick,
					getVariedQuestion(m.flashcards[m.currentIdx].Front),
				)
			}
		}

	case aiResponseMsg:
		switch msg.msgType {
		case "question":
			m.question = msg.content
			m.state = stateWaitingForAnswer
			m.textarea.Focus()
			return m, textarea.Blink
		case "evaluation":
			var evalResp evaluationResponse
			// Attempt to unmarshal the AI's response into our struct
			err := json.Unmarshal([]byte(msg.content), &evalResp)
			if err != nil {
				// If JSON parsing fails, fall back to a safe default
				m.evaluation = "Could not parse AI response. " + msg.content
				m.correct = false
			} else {
				// If successful, populate the model from the parsed JSON
				m.evaluation = evalResp.Explanation
				m.correct = (evalResp.Result == "correct")
			}
			m.state = stateShowingAnswer
		}
		return m, nil

	case errorMsg:
		m.evaluation = "Error: Could not connect to Ollama. Please ensure it's running."
		m.state = stateShowingAnswer
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.state == stateQuit {
		return "Saving progress... Goodbye! 👋\n"
	}

	var s strings.Builder

	width := m.width
	if width <= 0 {
		width = 80
	}
	// Helper style for simple width-constrained rendering
	renderWidth := lipgloss.NewStyle().Width(width - 4)

	switch m.state {
	case stateLoading:
		s.WriteString(fmt.Sprintf("\n\n  %s %s...\n", m.spinner.View(), m.loadingMsg))

	case stateWaitingForAnswer:
		s.WriteString(questionStyle.Width(width-4).Render(m.question) + "\n\n")
		s.WriteString(m.textarea.View() + "\n\n")
		s.WriteString(renderWidth.Render("Press Enter to submit • Ctrl+S to show answer • Ctrl+C to quit"))

	case stateShowingAnswer:
		if m.correct {
			s.WriteString(renderWidth.Render("✅ Correct! "+m.evaluation) + "\n\n")
		} else {
			s.WriteString(renderWidth.Render("❌ Not quite. "+m.evaluation) + "\n\n")
		}
		s.WriteString(answerStyle.Width(width-4).Render("Correct answer: "+m.flashcards[m.currentIdx].Back) + "\n\n")
		if m.userAnswer != "(showed answer)" {
			s.WriteString(answerStyle.Width(width-4).Render("Your answer: "+m.userAnswer) + "\n\n")
		}
		s.WriteString("How well did you know this?\n")
		s.WriteString(renderWidth.Render("  [1] Again (didn't know)   [2] Hard   [3] Good   [4] Easy\n"))
		s.WriteString(renderWidth.Render("\nPress 1-4 to continue • Ctrl+C to quit"))
	}

	s.WriteString(fmt.Sprintf("\n\n%s Card %d of %d", "Progress:", m.sessionProgress+1, m.totalCards))
	return s.String()
}

func (m *model) updateCardDifficulty(difficulty string) {
	allCards, err := loadAllFlashcards("flashcards.csv")
	if err != nil {
		return
	}

	currentCard := m.flashcards[m.currentIdx]
	for i, card := range allCards {
		if card.Front == currentCard.Front && card.Back == currentCard.Back {
			allCards[i].Difficulty = difficulty
			now := time.Now()
			switch difficulty {
			case "Again":
				allCards[i].NextReview = now.Add(1 * time.Minute)
			case "Hard":
				allCards[i].NextReview = now.Add(6 * time.Hour)
			case "Good":
				allCards[i].NextReview = now.Add(24 * time.Hour)
			case "Easy":
				allCards[i].NextReview = now.Add(4 * 24 * time.Hour)
			}
			break
		}
	}
	saveFlashcards("flashcards.csv", allCards)
}

func loadAllFlashcards(filename string) ([]flashcard, error) {
	f, err := os.Open(filename)
	if err != nil {
		// If the file doesn't exist, create it with a header
		if os.IsNotExist(err) {
			return []flashcard{}, saveFlashcards(filename, []flashcard{})
		}
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var cards []flashcard
	for i, record := range records {
		if i == 0 {
			continue
		} // Skip header

		card := flashcard{
			Front:      record[0],
			Back:       record[1],
			NextReview: time.Now(), // Default to now for new cards
			Difficulty: "New",
		}

		if len(record) > 2 && record[2] != "" {
			if t, err := time.Parse(time.RFC3339, record[2]); err == nil {
				card.NextReview = t
			}
		}
		if len(record) > 3 && record[3] != "" {
			card.Difficulty = record[3]
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func loadFlashcards(filename string) ([]flashcard, error) {
	cards, err := loadAllFlashcards(filename)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var dueCards []flashcard
	for _, card := range cards {
		if card.NextReview.Before(now) || card.NextReview.Equal(now) {
			dueCards = append(dueCards, card)
		}
	}

	sort.Slice(dueCards, func(i, j int) bool {
		return dueCards[i].NextReview.Before(dueCards[j].NextReview)
	})

	return dueCards, nil
}

func saveFlashcards(filename string, cards []flashcard) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"Front", "Back", "NextReview", "Difficulty"})
	for _, card := range cards {
		w.Write([]string{
			card.Front,
			card.Back,
			card.NextReview.Format(time.RFC3339),
			card.Difficulty,
		})
	}
	return nil
}

func main() {
	// Ensure Ollama is running before starting the application
	if err := startOllama(); err != nil {
		log.Fatalf("Failed to start Ollama: %v\nPlease ensure Ollama is installed and accessible in PATH", err)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
