#  Wa-Tor Simulation

This project aims to recreate a **Wa-Tor ecological simulation**, originally described by  
*A.K. Dewdney in Scientific American*.

The world is simulated as a **toroidal grid** populated by fish and sharks.

##  Simulation Rules (Classical Wa-Tor Model)

- **Fish**  
  - Move randomly to adjacent empty squares  
  - Reproduce after surviving **FishBreed** chronons  
- **Sharks**  
  - Hunt nearby fish (move into fish cell & eat it)  
  - Lose 1 energy every chronon  
  - Starve when energy reaches zero  
  - Reproduce after surviving **SharkBreed** chronons  
- **World Behaviour**  
  - Full **edge-wrapping** (toroidal) movement  
  - Time progresses in discrete steps called **chronons**  
- Supports **multi-threaded execution** for performance speedups

---

##  Features

  working **Wa-Tor simulation**
  Correct **fish & shark behaviour**  
  Implements **energy**, **reproduction**, **starvation**
  Simple **ASCII grid display** for graphical output in the terminal
  **Multithreading with goroutines**
  **Benchmark mode** with CSV export
  Parallel **speedup analysis (1–8 threads)**
  fully documented with **Doxygen-compatible comments**

---

**   Viewing the Doxygen Documentation

The full Doxygen-generated documentation is included inside the html/ folder of this repository.
To view the documentation:
Open the html/ folder in the repository.
Locate the file named index.html.
Click "View Raw" → your browser will automatically open and display the full Doxygen site.

---

##  Requirements

- **Linux**
- **Go 1.20+**

---

##  How to Run

The Wa-Tor simulation is executed from the terminal.

It requires **7 positional parameters**, plus optional flags for drawing, chronons, and benchmarking.

### **Basic run example**
```sh
go run . 50 200 3 6 5 200 4
