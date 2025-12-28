# LaTeX Diagram Guide for MangaHub Report

This guide shows you how to add various types of diagrams to your LaTeX report.

## Table of Contents

1. [Using External Images](#1-using-external-images)
2. [TikZ Diagrams (In-LaTeX)](#2-tikz-diagrams-in-latex)
3. [Performance Charts](#3-performance-charts)
4. [Sequence Diagrams](#4-sequence-diagrams)
5. [Network Protocol Diagrams](#5-network-protocol-diagrams)

---

## 1. Using External Images

### Setup

Create a `figures/` directory in your `report/` folder:

```bash
mkdir -p /home/tnphuc/Workspace/Projects/MangaHub/report/figures
```

### Add Image to LaTeX

```latex
\begin{figure}[htbp]
    \centering
    \includegraphics[width=0.8\textwidth]{figures/architecture.png}
    \caption{System Architecture Diagram}
    \label{fig:architecture}
\end{figure}
```

**Options for includegraphics:**

- `width=0.8\textwidth` - 80% of page width
- `height=6cm` - Fixed height
- `scale=0.5` - 50% of original size
- `angle=90` - Rotate 90 degrees

**Supported formats:** PNG, JPG, PDF (recommended for vector graphics)

### Referencing in Text

```latex
As shown in Figure~\ref{fig:architecture}, the system consists of five servers.
```

---

## 2. TikZ Diagrams (In-LaTeX)

Already added to your `report.tex`. You can draw diagrams directly in LaTeX!

### Simple Flowchart Example

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}[node distance=2cm]
        \node (start) [rectangle, draw, fill=green!20] {Start};
        \node (process) [rectangle, draw, below of=start, fill=blue!20] {Process};
        \node (decision) [diamond, draw, below of=process, fill=yellow!20] {Decision?};
        \node (end) [rectangle, draw, below of=decision, fill=red!20] {End};

        \draw [arrow] (start) -- (process);
        \draw [arrow] (process) -- (decision);
        \draw [arrow] (decision) -- node[anchor=west] {Yes} (end);
    \end{tikzpicture}
    \caption{Process Flowchart}
    \label{fig:flowchart}
\end{figure}
```

### Network Topology Example

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}[node distance=3cm]
        % Nodes
        \node (client) [client] {Web Client};
        \node (lb) [server, right of=client] {Load Balancer};
        \node (api1) [server, above right of=lb] {API Server 1};
        \node (api2) [server, below right of=lb] {API Server 2};
        \node (db) [database, right of=lb, xshift=4cm] {Database};

        % Connections
        \draw [arrow] (client) -- (lb);
        \draw [arrow] (lb) -- (api1);
        \draw [arrow] (lb) -- (api2);
        \draw [arrow] (api1) -- (db);
        \draw [arrow] (api2) -- (db);
    \end{tikzpicture}
    \caption{Load Balanced Architecture}
    \label{fig:load-balance}
</figure>
```

---

## 3. Performance Charts

### Bar Chart with TikZ

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}
        \begin{axis}[
            ybar,
            width=12cm,
            height=6cm,
            ylabel={Response Time (ms)},
            xlabel={Number of Concurrent Users},
            symbolic x coords={100, 500, 1000},
            xtick=data,
            nodes near coords,
            nodes near coords align={vertical},
        ]
        \addplot coordinates {(100,15) (500,45) (1000,120)};
        \end{axis}
    \end{tikzpicture}
    \caption{HTTP API Response Time vs Load}
    \label{fig:performance}
\end{figure}
```

**Note:** Requires `\usepackage{pgfplots}` in preamble.

### Pie Chart Example

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}
        \pie[radius=2]{
            70/Unit Tests,
            20/Integration Tests,
            10/E2E Tests
        }
    \end{tikzpicture}
    \caption{Testing Distribution}
    \label{fig:test-distribution}
\end{figure}
```

**Note:** Uses `pgf-pie` package (already added).

---

## 4. Sequence Diagrams

### Enhanced Sequence Diagram

```latex
\begin{figure}[htbp]
    \centering
    \begin{sequencediagram}
        \newinst{client}{Client}
        \newinst[2]{server}{Server}
        \newinst[2]{db}{Database}

        \begin{call}{client}{Login Request}{server}{200 OK}
            \begin{call}{server}{Verify User}{db}{User Data}
            \end{call}
            \begin{call}{server}{Generate JWT}{server}{Token}
            \end{call}
        \end{call}
    \end{sequencediagram}
    \caption{Authentication Sequence}
    \label{fig:auth-sequence}
\end{figure}
```

**Note:** Requires `\usepackage{pgf-umlsd}` package.

---

## 5. Network Protocol Diagrams

### WebSocket Connection Diagram

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}[node distance=1.5cm]
        % Nodes
        \node (browser) [client] {Browser};
        \node (http) [server, right of=browser, xshift=2cm] {HTTP Server};
        \node (ws) [server, below of=http] {WebSocket Hub};
        \node (room1) [client, right of=ws, xshift=2cm, yshift=0.7cm] {Room 1};
        \node (room2) [client, right of=ws, xshift=2cm, yshift=-0.7cm] {Room 2};

        % Arrows
        \draw [arrow] (browser) -- node[above] {HTTP Upgrade} (http);
        \draw [arrow] (http) -- node[right] {Upgrade to WS} (ws);
        \draw [arrow, <->] (ws) -- (room1);
        \draw [arrow, <->] (ws) -- (room2);
    \end{tikzpicture}
    \caption{WebSocket Chat Architecture}
    \label{fig:websocket-arch}
\end{figure}
```

---

## 6. Creating Diagrams with External Tools

### Recommended Tools:

1. **Draw.io (diagrams.net)**
   - Free, web-based
   - Export as PNG or PDF
   - URL: https://app.diagrams.net/

2. **Lucidchart**
   - Professional diagrams
   - Export as PNG/PDF
   - URL: https://www.lucidchart.com/

3. **PlantUML**
   - Text-based diagrams
   - Export as PNG/PDF
   - Great for sequence diagrams

4. **Graphviz**
   - Graph visualization
   - Good for network topologies

### Workflow:

1. Create diagram in external tool
2. Export as PNG or PDF
3. Save to `report/figures/` directory
4. Include in LaTeX with `\includegraphics`

---

## 7. Quick Reference: Diagram Positions

### Figure Placement Options

```latex
\begin{figure}[placement]
```

- `h` - Here (approximately)
- `t` - Top of page
- `b` - Bottom of page
- `p` - Separate page
- `!` - Override LaTeX placement algorithm
- `H` - Exactly here (requires `\usepackage{float}`)

**Best practice:** Use `[htbp]` for most figures.

---

## 8. Tips for Your MangaHub Report

### Suggested Diagrams to Add:

1. **Chapter 3 - System Design:**
   - ✅ System architecture (already added)
   - ✅ Database ER diagram (already added)
   - Protocol interaction flow
   - Deployment architecture

2. **Chapter 4 - Implementation:**
   - ✅ TCP synchronization flow (already added)
   - WebSocket hub pattern
   - gRPC service structure
   - Authentication flow

3. **Chapter 5 - Testing:**
   - Test pyramid diagram
   - Performance graphs
   - Code coverage charts

4. **Appendices:**
   - Network topology
   - Data flow diagrams

### Performance Graph Example (Add to Chapter 5)

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}
        \begin{axis}[
            xlabel={Concurrent Users},
            ylabel={Throughput (req/s)},
            width=10cm,
            height=6cm,
            grid=major
        ]
        \addplot[color=blue, mark=square] coordinates {
            (100, 6500)
            (500, 11000)
            (1000, 8300)
        };
        \end{axis}
    \end{tikzpicture}
    \caption{API Server Throughput Under Load}
    \label{fig:throughput}
\end{figure}
```

---

## 9. Compilation Commands

After adding diagrams, compile with:

```bash
cd /home/tnphuc/Workspace/Projects/MangaHub/report

# Method 1: Standard pdflatex (run twice for references)
pdflatex report.tex
pdflatex report.tex

# Method 2: Using latexmk (recommended)
latexmk -pdf report.tex

# Method 3: Clean and rebuild
latexmk -pdf -pvc report.tex  # Auto-recompile on changes
```

### Troubleshooting Compilation Errors

**Common Issues:**

1. **Missing TikZ library:**

   ```latex
   \usetikzlibrary{shapes.geometric, arrows, positioning}
   ```

2. **PGFPlots for charts:**

   ```latex
   \usepackage{pgfplots}
   \pgfplotsset{compat=1.18}
   ```

3. **Float package for exact positioning:**
   ```latex
   \usepackage{float}
   ```

---

## 10. Example: Complete Testing Pyramid Diagram

Add this to Chapter 5 Testing section:

```latex
\begin{figure}[htbp]
    \centering
    \begin{tikzpicture}
        % Draw pyramid
        \fill[green!30] (0,0) -- (4,0) -- (3,2) -- (1,2) -- cycle;
        \fill[blue!30] (1,2) -- (3,2) -- (2.5,3.5) -- (1.5,3.5) -- cycle;
        \fill[red!30] (1.5,3.5) -- (2.5,3.5) -- (2,4.5) -- cycle;

        % Labels
        \node at (2,0.8) {\Large Unit Tests (70\%)};
        \node at (2,2.7) {\large Integration (20\%)};
        \node at (2,4) {E2E (10\%)};

        % Borders
        \draw[thick] (0,0) -- (4,0) -- (3,2) -- (1,2) -- cycle;
        \draw[thick] (1,2) -- (3,2) -- (2.5,3.5) -- (1.5,3.5) -- cycle;
        \draw[thick] (1.5,3.5) -- (2.5,3.5) -- (2,4.5) -- cycle;
    \end{tikzpicture}
    \caption{Testing Pyramid Strategy}
    \label{fig:testing-pyramid}
\end{figure}
```

---

## Additional Resources

- TikZ Documentation: http://mirrors.ctan.org/graphics/pgf/base/doc/pgfmanual.pdf
- TikZ Examples: https://texample.net/tikz/examples/
- LaTeX Wikibook: https://en.wikibooks.org/wiki/LaTeX/Floats,_Figures_and_Captions

---

**Pro Tip:** Keep your diagrams simple and clear. Academic reports value clarity over visual complexity!
