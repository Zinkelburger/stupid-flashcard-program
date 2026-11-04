# Stupid Flashcard Program

Idk I wanna do the RHCSA so I made this instead. Anki deck on the command line with a local LLM.

Questions taken from <https://github.com/chlebik/rhcsa-practice-questions> and given to Claude

## RHCSA deck (EX200, RHEL 10)

- `rhcsa-flashcards.tsv` — 505 cards written against the RHEL 10 exam objectives. Import it in Anki with File → Import; the header in the file sets the deck, the tags and the stable note IDs, so re-importing an updated file updates cards in place.
- `RHCSA-STUDY-PLAN.md` — what the cards won't teach you, resources, lab setup, and a training plan.

## Network+ deck (N10-009)

- `networkplus-flashcards.tsv` — 629 cards covering every N10-009 objective (1.1–5.5), tagged `netplus::<objective>`. Imports into its own Anki deck, "CompTIA Network+ N10-009".
- `NETWORKPLUS-STUDY-PLAN.md` — exam format, which resources are worth it, and a training plan.

(The CLI program still reads the older `flashcards.csv`.)

Ramalama is mean and makes me set up a server or something idk how to do that

```
curl -fsSL https://ollama.ai/install.sh | sh
ollama pull llama3.2
ollama serve
```
