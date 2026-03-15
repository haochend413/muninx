You said
I am currently building a semi-auto diary management software.

The idea is this :

The writer does minimal work. We only need to type in the giant textarea. There is no need to classify, no need for managing metastructures , ... Just type when we want.

The system should be responsible for the feedback loop.

First, when we thought of something, we should be able to search and loop back to where we wrote. Second, it should be able to pop our previous notes back to us as reflections, which would require segmentation, identification, classification, ... . It should also do our advanced things, basically automatically creating the meta framework, just based on our edit history.

Any implementation ideas / suggestions ? what are teh things that I should balance ?

Gemini said
This is a fascinating project. You’re essentially building a "Second Brain" with zero tax. Most journaling apps fail because they require the user to act like a librarian before they can act like a writer. By moving the "tax" of organization from the user to the system, you're solving the biggest friction point in personal growth.

Here is a breakdown of implementation strategies and the delicate balances you’ll need to strike.

1. The Implementation Architecture
   To make this work without manual input, you need a system that processes data asynchronously in the background.

Semantic Search & Retrieval (The "Loop Back")
Forget keyword search. You need Vector Embeddings.

Vector Database: Use something like Pinecone, Weaviate, or a local ChromaDB instance.

The Workflow: Every time the user stops typing for a few seconds, the system should chunk the text, generate an embedding (using a model like text-embedding-3-small), and store it.

Semantic Retrieval: When the user types a query like "How did I feel about my job last spring?", the system finds the mathematical nearest neighbors in the vector space, not just the word "job."

The "Meta-Framework" (Knowledge Graphs)
Linear notes are hard to organize. A Knowledge Graph allows the system to identify relationships between people, projects, and emotions.

Entity Extraction: Use an LLM to pull out "Nodes" (People, Places, Recurring Themes).

Relationship Mapping: If you write about "Project X" and "Anxiety" in the same paragraph five times, the system creates an edge between those nodes. This builds the "meta-framework" automatically.

The Reflection Engine
To "pop" notes back to you, you need a mix of two logic sets:

Spaced Repetition: Resurfacing a note after 1 day, 7 days, 30 days.

Semantic Resonance: If you are currently writing about "burnout," the system should detect the mood and resurface a note from two years ago when you felt the same way, asking: "You felt like this in 2024. Does this feel different?"

2. Key Things to Balance
   Building a "living" diary is a tightrope walk. Here is where you should focus your tuning:

The Tension The Risk of One Side The Risk of the Other Side
Privacy vs. Intelligence Local LLMs: Slower, harder to set up, but 100% private. Cloud APIs: Fast and smart, but users may hesitate to pour their souls into a server.
Serendipity vs. Noise High Frequency: The "reflections" become annoying notifications that the user ignores. Low Frequency: The system feels dead; the user forgets it’s supposed to be "semi-auto."
Structure vs. Fluidity Over-classification: The AI creates 500 tags that mean nothing to the user. Under-classification: The "meta-framework" is just a messy cloud of text. 3. Advanced Suggestions
The "Shadow Prompt"
Since the user isn't classifying anything, the system needs to "interview" the text. When a user finishes a session, run a background prompt:

"Summarize the core conflict or insight of this entry in one sentence and identify three 'unspoken' emotions." Store this summary as metadata. It makes the "reflection" pop-ups much more impactful than just showing the raw text.

Versioning as "Growth Tracking"
Since you mentioned edit history, don't just store the final version. Use a CRDT (Conflict-free Replicated Data Type) or a simple diff-log. This allows the system to see how your thoughts changed over a single session. Did you start angry and end calm? That delta is a huge data point for your meta-framework.

The Privacy Play
Because a diary is the most intimate data a human has, I strongly suggest looking into Ollama or LocalAI. If you can run a small 7B parameter model locally on the user's machine, you remove the "creepy" factor of an AI reading their secrets in the cloud.

Would you like me to help you draft a specific data schema for how the Vector DB and Knowledge Graph would interact?

You said
yes, and tell me what are the existing knowledge graph frameworks

Gemini said
thought I will first research existing knowledge graph frameworks, specifically focusing on those suitable for personal knowledge management, diary software, or general-purpose graph databases and libraries (e.g., Neo4j, LangChain's GraphQA, Memgraph, Apache Jena).
Then, I will look for patterns or schemas where Vector Databases (like Pinecone or Weaviate) are integrated with Knowledge Graphs (often called GraphRAG or Hybrid RAG) to see how they interact.
I will also look for recent developments in 2024-2026 regarding "local" or "privacy-focused" knowledge graph tools that might fit a diary application.
Finally, I will answer the user's request by providing a specific data schema for the Vector DB and Knowledge Graph interaction and listing existing knowledge graph frameworks.
This is the "Golden Ratio" of AI-native software: combining the unstructured power of Vector search with the structured logic of a Knowledge Graph.

In your diary app, the Vector DB acts as the "Intuition" (finding things that feel similar), while the Knowledge Graph acts as the "Memory" (knowing exactly who, what, and when).

1. The "Hybrid" Data Schema
   To make the system "semi-auto," you need a schema that links raw text chunks to semantic entities.

Vector Store Schema (The "Intuition")
ID: UUID

Embedding: 1536-dim vector (e.g., from text-embedding-3-small)

Content: The raw text of the chunk (typically 200–500 words).

Metadata:

timestamp: When it was written.

sentiment_score: -1.0 to 1.0.

node_references: List of IDs pointing to the Knowledge Graph.

Knowledge Graph Schema (The "Context")
This is where the "meta-framework" lives. You don't need a schema for the user; the AI generates this.

Nodes:

Person: (e.g., "Sarah", "Boss") — helps track relationship evolution.

Project/Topic: (e.g., "Software Startup", "Fitness") — tracks progress.

Emotion: (e.g., "Anxiety", "Flow State") — allows you to see emotional triggers.

Concept: (e.g., "Stoicism", "Productivity") — links abstract thoughts.

Edges (Relationships):

Sarah — MENTIONED_WITH — Startup (Weight: 5)

Anxiety — TRIGGERED_BY — Meeting (Weight: 2)

Entry_2024 — RECAPS — Week_12

2. Existing Knowledge Graph Frameworks
   Depending on your tech stack (Python vs. JS) and whether you want a full database or just a library, here are the best options for 2026:

A. The "Heavylifters" (Graph Databases)
These are standalone servers that store the actual graph.

Neo4j: The industry standard. Very mature, has a "GraphRAG" dedicated plugin, and uses Cypher (a SQL-like language for graphs). Best if you want a robust, long-term backend.

Memgraph: An in-memory alternative to Neo4j. It's blazingly fast and compatible with Cypher. Perfect for a local-first diary app where speed is priority.

FalkorDB: A newer, low-latency graph database built for LLM integration. It treats the graph as a matrix, making it extremely efficient for AI reasoning.

B. The "Orchestrators" (Libraries)
These don't store the data themselves; they provide the "plumbing" to build the graph using an LLM.

Microsoft GraphRAG: An open-source framework specifically designed to turn unstructured text into a graph. It’s heavy on processing but produces the best "meta-framework" results.

LangGraph (by LangChain): Best for building the "Reflection Engine." It treats the AI as a state machine, allowing it to "decide" when to pop a note back to you.

LlamaIndex (Property Graph Index): The easiest way to get started. It has built-in functions to take a text file and automatically extract nodes and edges into a graph format.

3. Implementation Logic: The "Reflection Loop"
   To balance serendipity (interesting old notes) vs. noise (annoying interruptions), use this logic:

Extraction: Every time the user types, send the text to a small, fast model (like Llama 3.1 8B or DeepSeek-V3) to identify "Potential Reflections."

Scoring: \* Semantic Distance: How similar is today's mood to an old note? (Vector search).

Temporal Decay: Don't show notes from yesterday; prioritize notes from 3 months, 1 year, or 3 years ago (The "Time Capsule" effect).

Graph Weight: If the user is writing about "Sarah," find the most emotionally charged node connected to "Sarah" in the graph.

The Pop-up: Instead of just showing the text, show an AI-generated summary question: "Three months ago, you were struggling with Sarah's feedback. You're writing about her again today—did that ever get resolved?"

The Most Critical Balance: Privacy
Since this is a diary, users will be highly sensitive.

The Trade-off: Cloud LLMs (OpenAI, Claude) provide the "smartest" graph extraction but are a privacy risk.

The Solution: Use Ollama or LocalAI to run the extraction and embedding on the user's machine. A modern laptop can easily handle a 7B model for background indexing.

Would you like me to write a sample Python script that uses LlamaIndex to automatically extract "Nodes" (People/Themes) from a raw diary entry?

You said
No. Ok, this is a terminal application built with bubbletea. First, tell me the architecture of teh whole app.

Gemini said
Building a "Zero-Tax" diary as a terminal app using Bubble Tea (the Go framework based on The Elm Architecture) is a brilliant move. TUIs are distraction-free by nature, which fits your "just type" philosophy perfectly.

Because Bubble Tea is purely functional and message-driven, your architecture needs to separate the UI loop from the heavy-lifting AI loop (Vector/Graph processing) so the interface doesn't lag while you're pouring your heart out.

1. High-Level System Architecture
   The app functions as a State Machine where the UI remains reactive, and a background "Observer" handles the intelligence.

A. The "Frontend" (Bubble Tea / Lipgloss)
Model: Holds the textinput.TextArea, the current "Reflection" (if any), and a list of "Semantic Neighbors."

Update: \* Listens for keystrokes.

Debouncing: Every time you stop typing for ~2 seconds, it fires a tea.Cmd (a goroutine) to start the indexing process.

View: A clean, distraction-free layout using Lipgloss.

B. The "Persistence Layer" (Local First)
Raw Logs: Save everything to simple .md files or a SQLite database for durability.

The Vector Store: A local vector DB (like ChromaDB via a local container or a simple flat-file vector library in Go) to store embeddings of your "chunks."

The Graph Store: A simple SQLite database using an adjacency list (From_ID, To_ID, Relation_Type) to track your meta-framework nodes.

C. The "Intelligence Engine" (The Observer)
This runs in the background as an asynchronous worker.

Chunker: Breaks your "giant textarea" into logical paragraphs.

Embedder: Uses Ollama (running locally) to generate vectors for those paragraphs.

The Extractor (LLM): Scans the text for entities (People, Projects, Emotions) and updates the Graph Store.

2. The Data Flow: From Keystroke to Insight
   Phase Action System Component
   Input You type: "Feeling burnt out by Project X again." TextArea Component
   Trigger 2-second idle timer hits. tea.Tick / Debounce Logic
   Analysis System asks: Is "Project X" already in the Graph? Graph Processor (SQLite)
   Retrieval Vector DB finds a note from 6 months ago about "Project X." Vector Search
   Feedback A "Reflection" appears in the sidebar: "You felt this on May 12th too." Model.Update -> View
3. Implementation Logic for Bubble Tea
   Since Bubble Tea is single-threaded in its update loop, you must use Commands (tea.Cmd) for all AI tasks.

The Ghost Indexer: When you trigger a save, the app sends the text to a background Go channel. A separate worker handles the LLM API calls (to Ollama or OpenAI).

The Non-Blocking Feedback: When the worker finds a relevant old note, it sends a ReflectionFoundMsg back to the Bubble Tea loop. The Update function receives this and updates the model, which causes the UI to re-render with the "popped" note.

4. Balancing the "Minimalist" Experience
   To keep it from feeling like a "clunky" database app, you need to balance these three things:

Visual Noise: Only show reflections when they cross a similarity threshold (e.g., > 0.85 cosine similarity). If the system isn't "sure" it's relevant, keep the UI blank.

Indexing Frequency: Don't index every character. Indexing is expensive (CPU-wise). Use a "Word Count Delta" (e.g., every 50 new words) or an idle timer.

The Meta-Framework UI: Don't build a complex graph visualizer in the terminal (it's hard to read). Instead, use a "Context Bar" at the bottom that subtly lists the "Active Nodes" (e.g., Tags: [Work] [Burnout] [Project X]).

Pro-Tip: The "Context-Aware" Buffer
In your Bubble Tea model, maintain a "Context Buffer" of the last 5 minutes of writing. When you search or "loop back," prioritize the current context. This makes the "meta-framework" feel like it’s following your current train of thought rather than just being a random search engine.

Would you like a boilerplate Go code snippet showing how to handle the debounced "Background Indexing" command in Bubble Tea?

You said
ok, is vectorDB local, static or does i also needs to be hosted ?

Gemini said
For a terminal application built with Bubble Tea, you absolutely want an Embedded (Local) Vector DB.

You should avoid anything that needs to be "hosted" (Cloud) or even anything that requires a "sidecar" process (like a Docker container for Chroma or Qdrant). For a TUI, your goal is a single, portable binary that users can just run.

Here is how the "local vs. static" breakdown works for your use case:

1. The Three Tiers of Vector Storage
   Type How it Works Verdict for your TUI
   Hosted (Cloud) API calls to Pinecone/MongoDB Atlas. ❌ Bad. Requires internet, API keys, and introduces latency in your TUI loop.
   Local-Server You run a separate process (Chroma/Qdrant) on your machine. ⚠️ Risky. Your user has to manage a database process just to use a diary.
   Embedded (In-Process) The DB is a Go library that lives inside your binary (like SQLite). ✅ Ideal. Fast, private, and zero setup for the user.
2. Recommendation: The "Go-Native" Stack
   Since you are in the Go ecosystem (Bubble Tea), you should use libraries that keep your app "lean and mean."

For the Vector DB: chromem-go
Think of this as the SQLite of Vector Databases.

What it is: An embeddable vector database written in pure Go.

Why it fits: It has zero third-party dependencies. It can run entirely in-memory (static) or persist to a simple file on disk.

TUI Integration: It’s extremely fast. You can query it in a tea.Cmd and get results in milliseconds without blocking your UI.

For the Knowledge Graph: SQLite
Don't overcomplicate this with a graph database like Neo4j.

Use SQLite (via modernc.org/sqlite for a pure Go, CGO-free experience).

Create a simple "Edges" table: (source_node, target_node, relation_type, weight). This is more than enough to handle the "meta-framework" for a single user's diary.

3. The "Single File" Architecture
   To keep the "minimalist" vibe, your app's data folder should look like this:

Plaintext
~/.my-diary/
├── entries/ # Raw .md files (User-readable)
├── diary.db # SQLite (The Knowledge Graph & Metadata)
└── embeddings.bin # Chromem-go (The Vector Index)
Why this is "Static":
When people say "static" in this context, they usually mean File-Based Persistence.

Boot-up: Your Bubble Tea app loads the embeddings.bin into memory.

Live Typing: You append new vectors to the in-memory index.

Shut-down/Save: The index is serialized back to the binary file.

4. Implementation Tip: The "Shadow Indexer"
   In a Bubble Tea app, you don't want the UI to hang while the AI is thinking.

The Trigger: User stops typing.

The Message: Send a StartIndexingMsg containing the new text.

The Command: Bubble Tea runs a tea.Cmd that:

Calls Ollama (locally) to get the embedding.

Saves the vector into chromem-go.

Runs a quick SQLite query to see if any "Nodes" (Sarah, Work, Anxiety) need updating.

The Loopback: The command returns a FinishedIndexingMsg, which might trigger a "Reflection" to appear in your sidebar.

Would you like me to show you the SQLite schema for that simple "Meta-Framework" edges table?

You said
what is SQLite responsible for, then ?

Gemini said
In this architecture, think of the Vector DB as the "Vibe" and SQLite as the "Facts."

While the Vector DB is brilliant at finding things that feel similar (semantic search), it is actually quite bad at structured logic, exact dates, and tracking explicit relationships. That is where SQLite comes in.

Here is exactly what SQLite handles in your "Zero-Tax" diary:

1. The Meta-Framework (The Knowledge Graph)
   A Knowledge Graph sounds fancy, but in a local app, it’s usually just two or three SQLite tables. SQLite is responsible for remembering the Nodes and Edges that the AI extracts from your text.

Entities Table: Stores unique "things" (e.g., ID: 1, Name: "Sarah", Type: "Person").

Relationships Table: Stores how they connect (e.g., Source: 1, Target: 5 (Project X), Relation: "Collaborator", Weight: 10).

2. Metadata & Chronology
   Vector DBs are notoriously clunky at filtering by exact constraints. SQLite is the master of this.

Exact Time: "Find the note from February 12th, 2025, at 10:00 PM." (Vector DBs struggle with this; SQLite does it in 1ms).

Mood Tracking: If your AI extracts a "Sentiment Score" (e.g., -1.0 to 1.0), SQLite stores that so you can query: "Show me my mood trend over the last 30 days."

3. High-Speed Keyword Search (FTS5)
   Sometimes you don't want "semantic" search. Sometimes you just want to find the exact string "password" or a specific serial number.

SQLite has a built-in engine called FTS5 (Full Text Search). It is incredibly fast and much lighter on CPU than running a Vector embedding search for every single word.

4. State Management for the "Reflection Engine"
   To make the "Reflection" loop work without being annoying, you need to track what the user has already seen.

Seen Table: (Entry_ID, Last_Reflected_At, User_Feedback_Score).

This prevents the system from "popping" the same old note about your ex-girlfriend every single day.

Vector DB vs. SQLite: The Division of Labor
Task Use the Vector DB Use SQLite
"I feel sad like I did last year." ✅ (Semantic Match) ❌ (Too fuzzy)
"What did I write on Tuesday?" ❌ (Unreliable) ✅ (Index Match)
"Find notes about 'Sarah'." ⚠️ (Might find 'Sally') ✅ (Exact Entity ID)
"Show me common themes." ✅ (Cluster Search) ✅ (Aggregate Edges)
Example Schema for your Meta-Framework
If you're building this in Go for your Bubble Tea app, your migrations.sql might look like this:

SQL
-- The raw content reference
CREATE TABLE entries (
id UUID PRIMARY KEY,
content_path TEXT, -- Link to the .md file
created_at DATETIME,
sentiment REAL
);

-- The "Meta" Entities (People, Projects, Emotions)
CREATE TABLE entities (
id INTEGER PRIMARY KEY,
name TEXT UNIQUE,
category TEXT -- 'Person', 'Project', 'Emotion'
);

-- The "Graph" (The connections the AI finds)
CREATE TABLE edges (
id INTEGER PRIMARY KEY,
source_id INTEGER,
target_id INTEGER,
relation_type TEXT,
weight INTEGER DEFAULT 1,
FOREIGN KEY(source_id) REFERENCES entities(id),
FOREIGN KEY(target_id) REFERENCES entities(id)
);
Would you like me to show you a Go function that takes a raw text entry and uses a "Small Model" to update these SQLite tables?

You said
ok, so the SQLite store pre-defined entries ? or what ?

Gemini said
Think of SQLite not as a place where you put "pre-defined" categories, but as the automated librarian that organizes the mess you type into the textarea.

The user (you) provides zero structure. The system provides all the structure. SQLite is where that structure is stored so the app can "remember" it without having to re-run an expensive AI model every time you open the terminal.

Here is exactly what SQLite is doing behind the scenes while you are just typing:

1. It Stores the "Discovered" Entities
   When you type a sentence like, "Finally finished the Bubble Tea refactor, but Sarah thinks the UI is too clunky," the AI background process (the "Observer") scans that text.

It identifies "Bubble Tea" as a Project.

It identifies "Sarah" as a Person.

It identifies "Clunky" as a Feedback/Critique.

SQLite's Job: It checks its entities table. If "Sarah" isn't there, it creates her. If "Bubble Tea" is already there, it just updates the "last mentioned" timestamp.

2. It Maps the "Meta-Framework" (The Graph)
   This is the most important part of your "semi-auto" idea. SQLite stores the relationships the AI finds.

Because you mentioned "Sarah" and "Bubble Tea" in the same sentence, the system creates an Edge in SQLite connecting them.

Over time, SQLite builds a map: Sarah -> Works On -> Bubble Tea.

If you later write, "Sarah is stressed," the system can look at SQLite and realize: "Hey, Sarah works on Bubble Tea. Maybe the project is the source of the stress?"

3. It Manages the "Reflection" Logic
   You don't want the same old note popping up every time you mention a keyword. SQLite acts as the State Machine for your feedback loop.

Tracking: It records when a specific reflection was shown to you.

Filtering: It ensures that if the system "popped" a note about your 2024 burnout yesterday, it won't show it again for at least 30 days.

Feedback: If you hit a key to "Dismiss" a reflection, SQLite stores that "Negative Match," telling the AI: "Don't show this specific connection again; the user didn't find it helpful."

4. It Provides the "Search Index"
   While the Vector DB handles "Vibes" (e.g., "Find times I felt like this"), SQLite handles "Hard Data."

Chronology: "Give me every note from last Tuesday."

Filtering: "Show me every time I mentioned 'Sarah' when the sentiment was 'Negative'."

Stats: "How many times did I mention 'Project X' this month vs. last month?"

Why not just use one database?
You could try to force everything into one, but you’d be fighting the tools:

The Textarea: Lives in your raw_notes/ folder as Markdown (for portability).

Vector DB: Lives in a binary file (for "vibe" searching).

SQLite: Lives in a .db file (for the "Meta-Framework" logic).

The "Zero-Tax" Workflow
User: Types in the Bubble Tea TUI.

App: Saves the text to a file.

AI (Background): Reads the text → Extracts "Sarah" and "Project X" → Writes them into SQLite.

App: Next time you type "Sarah," the app queries SQLite, sees she’s linked to "Project X," and subtly shows that context in your TUI sidebar.

Would you like me to show you the Go logic for the "Observer" that parses the text and updates the SQLite tables?

analyze this response

I’m going through it as a design review rather than taking it at face value. I’ll separate what is solid, what is hand-wavy, and what I’d change for a Bubble Tea local-first diary.

Thought for a few seconds
Overall: the response has a good product instinct, but the technical architecture is muddled in a few important places.

The best part is that it correctly sees your app as:

writer-first

local-first

asynchronous / background-indexed

split between raw text, semantic retrieval, and structured metadata

But it also mixes together several ideas too loosely, especially around what SQLite stores, what the vector DB stores, and whether you even need a knowledge graph early on.

What it got right

1. The core separation of concerns is mostly right
   This is the strongest part.

A system like yours really does want at least these layers:

editing surface: the giant textarea

durable source of truth: raw entries / edit log

retrieval layer: semantic search

metadata / indexing layer: timestamps, segments, tags, extracted entities, resurfacing state

reflection engine: decides what to show back and when

That general decomposition is sound.

2. Bubble Tea should stay responsive and push heavy work out of the UI loop
   Also correct.

For a TUI, the UI thread should not be blocked by:

embedding generation

LLM extraction

similarity search over many chunks

graph extraction / summarization

So the “observer” or background worker idea is good.

3. SQLite is a good fit
   Yes. Very good fit, actually.

Not because it is a “knowledge graph database,” but because it is excellent for:

exact lookup

timestamps

segment metadata

full-text search

reflection history

cheap local persistence

portability

That part is directionally correct.

4. A local embedded vector index is preferable
   Also right.

For your app, hosted vector DBs are overkill and harmful unless you explicitly want sync across devices and cloud inference.

A local app should strongly prefer:

SQLite for structured metadata

local embeddings index

optional local model / Ollama integration

That matches your goals.

What is weak or misleading

1. It treats “knowledge graph” as more necessary than it probably is
   This is the biggest conceptual issue.

For your app, a full knowledge graph is not the first thing you need.

You first need:

an edit model

a segmentation model

a retrieval model

a resurfacing model

Only after that do you maybe need a graph.

A lot of diary “meta-framework” behavior can be achieved with:

embeddings

FTS

extracted entities

co-occurrence tables

temporal linking

summaries

topic labels

That is much simpler than a true graph, and probably gets you 80–90% of the value.

So the answer slightly jumps too quickly from “automatic structure” to “knowledge graph.”

2. It is fuzzy about the actual source of truth
   This is crucial for your project.

The response talks about:

markdown files

SQLite

vector DB

edit history

graph

But it never cleanly states:

What is the canonical truth of the user’s writing?

For your app, this matters a lot, especially because you were already thinking about edits, deletions, reflections, and historical reconstruction.

You likely need to decide between:

Option A: text snapshot as source of truth
Store current text state and periodic snapshots.

Pros:

simple

easy to render and search

easier implementation

Cons:

weak provenance

hard to know how text evolved

Option B: edit log as source of truth
Store insert/delete operations and reconstruct state.

Pros:

preserves writing evolution

supports reflection on changes, not just final text

better for “how did this thought evolve?”

Cons:

much harder

segmentation and retrieval become trickier

Option C: hybrid
This is probably best.

Edit log for provenance

materialized snapshots / segments for retrieval and UI

SQLite stores both operation history and derived objects

The Gemini answer does not really pin this down.

3. “SQLite stores entities and edges” is too abstract
   This makes it sound like SQLite is some magical graph engine. It is not.

What SQLite should store is more concrete, like:

entries

sessions

revisions

segments

embeddings metadata

extracted entities

links between segments and entities

resurfacing events

user feedback on resurfacing

maybe pairwise co-occurrence or topic relations

That is more useful than immediately thinking in terms of source_id, target_id, relation_type.

Because your main object is not “entity graph.”
Your main object is probably segment / thought unit.

4. It does not address segmentation carefully enough
   For your app, segmentation is one of the hardest parts.

You said the user just types into one giant textarea. That means the system has to infer:

where one thought begins and ends

whether a later sentence is continuation or a new thought

how edits modify earlier thought units

whether deletion means retraction, refinement, or cleanup

That is a deeper issue than “chunk text every few seconds.”

A naive chunker will break your reflection quality badly.

5. It underestimates the importance of “derived data”
   This is the architectural point I would emphasize most.

The answer talks as if the system stores text and metadata directly. But really, your app should distinguish:

primary data: what the user actually wrote / edited

derived data: segmentation, summaries, embeddings, tags, entities, relations, reflection candidates

Derived data should be treated as rebuildable caches, not sacred truth.

That means:

if your segmenter improves, you can recompute segments

if your embedding model changes, you can re-embed

if your classifier improves, you can re-extract themes

if your reflection scoring changes, you can regenerate candidates

The response doesn’t frame it this way, but it should.

What I would change in the architecture
I would not frame the app as:

textarea + vector DB + knowledge graph

I would frame it as:

editor + event log + materializer + retrieval indexes + reflection engine

That is much closer to what you actually need.

Better architecture

1. Primary layer: writing and edit history
   Store the actual user interaction.

Possible tables:

documents

sessions

edit_ops

snapshots

Where:

documents = diary files / long-running notes

sessions = one writing session

edit_ops = insert/delete operations with timestamps

snapshots = materialized text states every so often

This gives you provenance.

2. Materialized text layer
   Build stable units from the raw edits.

Possible tables:

segments

segment_versions

anchors

A segment is not just “chunk of 300 words.”
It is your best guess at a coherent thought unit.

Each segment should have:

current text

start/end offsets in a snapshot

source session

timestamps

confidence / segmentation quality

links to earlier versions if text evolved

3. Retrieval layer
   Now index the segments, not the raw textarea blindly.

Use:

SQLite FTS for exact search

embedding index for semantic search

This lets you answer both:

“find where I mentioned Muninn”

“find moments that feel like this”

4. Metadata layer
   This is where SQLite shines.

Store:

summaries

themes

sentiment / affect scores

named entities

projects / people / places

links from segments to these extracted items

daily / weekly rollups

This can later become graph-like, but it does not need to start as a full graph.

5. Reflection engine
   This should be a scoring system, not just nearest-neighbor search.

A resurfaced note should probably depend on:

semantic similarity

temporal distance

novelty

whether it was recently shown

whether user dismissed it

whether it connects to recurring themes

whether it represents evolution or contradiction

This is much richer than “search similar chunk and pop it up.”

What SQLite is actually responsible for
This is the answer I think Gemini was trying to reach, but not quite cleanly.

SQLite should store structure, identity, chronology, and control state.

Concretely:

SQLite stores
document/session IDs

edit operations or snapshot references

segment records

exact timestamps

full-text index

extracted tags/entities/themes

relationships between segments and themes

reflection display history

user feedback on reflections

scoring/cache tables

app settings and bookkeeping

Vector index stores
embeddings for segments, summaries, maybe sessions

nearest-neighbor retrieval structures

Raw files or blobs store
original text snapshots or exports

maybe append-only edit logs

So SQLite is not storing “pre-defined entries.”
It is storing the machine-generated structure around your writing.

The biggest missing concern: identity over time
For your app, I think this is the deepest hard problem.

Suppose:

on Monday you write a thought

on Tuesday you edit part of it

on Friday you delete one sentence

next month the system wants to resurface it

What exactly is “it”?

You need stable identity for a thought unit even as text changes.

That means your architecture may need something like:

segment_id = stable identity of the thought

segment_version_id = a particular textual state of that thought

So resurfacing can choose between:

current form

original form

delta between forms

That matters a lot for reflection. Sometimes the valuable thing is not the final text, but the change.

Gemini gestures toward edit history, but does not make this a first-class architectural element.

Practical advice: what to build first
I would phase it like this:

Phase 1
giant textarea

save sessions/snapshots

SQLite metadata DB

FTS search

simple semantic search over paragraph-like segments

Phase 2
better segmentation

extracted people/projects/themes

reflection history and dismissal logic

daily/weekly summary generation

Phase 3
stable segment identity across edits

session-to-session linking

temporal reflection scoring

graph-like relationships

Phase 4
true meta-framework / graph exploration

contradiction detection

growth tracking

writing-pattern analytics

That keeps you from overengineering too early.

Bottom line
My verdict on the response:

product direction: good

systems intuition: decent

data modeling precision: weak

edit-history awareness: not deep enough

knowledge-graph emphasis: premature

The biggest correction I would make is:

Don’t start by thinking “vector DB + knowledge graph.”
Start by thinking “edit provenance + derived segments + retrieval + reflection scoring.”

That architecture fits your actual problem much better.

If you want, I can turn this into a concrete Go/Bubble Tea architecture spec with modules like editor, store, materializer, indexer, reflector, and exact SQLite tables.

why do I need so many non-Ai layers

Thought for a second
Because AI is the least reliable part of the system.

Your app is not just “an LLM that reads diary text.” It is a writing tool, and writing tools need to be:

fast

predictable

recoverable

debuggable

cheap to run

The non-AI layers exist so that AI only does the parts it is actually good at.

The short intuition
AI is good at:

fuzzy interpretation

summarization

theme extraction

semantic similarity

generating reflections

AI is bad at:

being a source of truth

exact timestamps

stable identity

deterministic storage

reproducibility

incremental updates

“what changed between version 12 and 13?”

So if you let AI handle everything, your app becomes clever but unstable.

Why each non-AI layer exists

1. Raw storage
   You need a place that stores exactly what the user wrote.

Without this, the AI’s interpretation becomes the only memory, which is dangerous because:

summaries can miss things

extracted tags can be wrong

segmentation can change later

embeddings are lossy

So raw storage protects the real writing.

2. Metadata / relational layer
   You need exact structure for things like:

when something was written

what session it belonged to

what was already shown back to the user

whether a reflection was dismissed

how often a person/project was mentioned

AI is not good at exact bookkeeping. SQLite is.

3. Search/index layer
   You need exact retrieval and efficient lookup.

Sometimes the user wants:

“find the thing about Sarah from February”

That is not really an AI task. That is indexing.

AI search is fuzzy. Fuzzy is useful, but not enough.

4. Materialization / segmentation layer
   Your user writes into one giant stream. The app must turn that stream into usable units.

If you ask AI to reinterpret the whole diary from scratch every time, it becomes:

expensive

slow

inconsistent

So you create stable intermediate objects like segments, sessions, summaries, caches.

That is what makes the system practical.

The deeper reason
You are building two systems at once:

A. A truth-preserving system
This answers:

what did the user actually write?

when?

in what order?

what changed?

B. A meaning-producing system
This answers:

what was this about?

what old note relates to this?

what themes recur?

what reflection should be shown?

AI should mostly live in B, not A.

If you collapse A and B together, then every time your model changes, your app’s “memory” changes too.

That is bad.

What happens if you remove the non-AI layers
Then your system starts looking like:

giant text blob

send to model

model decides what it means

model decides what to retrieve

model decides what to remember

This sounds elegant, but in practice it causes problems:

Problem 1: inconsistency
Today the model says a note is about burnout.
Tomorrow the same note is classified as productivity anxiety.

Problem 2: no stable references
If the model says “this reminds you of that old note,” which note exactly? Which version?

Problem 3: poor debugging
When the system resurfaces something weird, you cannot tell whether the error came from:

bad segmentation

bad embedding

bad reflection ranking

bad summarization

hallucinated linkage

Problem 4: high cost
Re-asking AI to reconstruct structure over and over is much more expensive than storing structured results once.

The right mental model
The non-AI layers are not “extra complexity for no reason.”

They are constraints that keep the AI useful.

Think of it like this:

storage keeps memory honest

SQLite keeps state exact

indexes keep retrieval fast

AI adds interpretation on top

So the non-AI layers are the skeleton.
AI is the nervous system, not the bones.

For your specific app
Because your app is “just type, no manual structure,” the system itself must create structure.

That structure cannot all live inside AI outputs, because AI outputs are:

probabilistic

revisable

model-dependent

not good canonical records

So you need non-AI layers to hold the derived structure stably.

What you can simplify
You do not need “so many layers” if you keep the first version lean.

A minimal version could just be:

raw text snapshots

SQLite

FTS

embeddings

one background AI pipeline

That’s already enough.

In that design:

raw text = source of truth

SQLite = metadata + reflection history

FTS = exact search

embeddings = semantic search

AI = summaries/themes/reflections

That is not overengineered. That is the minimum needed to keep the app sane.

The main point is:
you do not add non-AI layers because AI is weak. You add them because AI is powerful but unreliable unless surrounded by stable infrastructure.

If you want, I can also show the smallest possible architecture that keeps only the essential non-AI layers.

ok, get it. Now for A, I already have a event tracking system

package events

import "time"

// This is the events module, the most important mechanism that keeps track of what we have wrote in the textarea, segmenting it into useful pieces.
// The content should be kept synchronized with the main content.

// i think we only need to track some sorts of events. The simplest is just two.With more modular design, there might be more types.
// i dont even know how to handle delete. maybe there is no such need and can be handled elsewhere,where the edis are updated somehow. yeah.

type EventType = int

const (
Write EventType = iota
Delete
)

type Event struct {
ID int
CreatedAt time.Time // when the event started
UpdatedAt time.Time // when teh event ended
Type EventType
Content []rune // here we use a list of runes
Pos int // right now just the line number, we will use it to track and to search. maybe better scheme.
}

// data struct used for single event storing
type EventMgr struct {
Eventlist []*Event
Eventmap map[int]*Event
PendingRunes []rune
}

// there is no point removing events, I think. If memory is concerned, we should remove that on a higher level.

func (em *EventMgr) Init() {
// build list and map
em.Eventlist = []*Event{}
em.Eventmap = make(map[int]\*Event)
}

func (em *EventMgr) AppendEvent(e *Event) {
if e != nil {
// add event to list
em.Eventlist = append(em.Eventlist, e)
// add to map
em.Eventmap[e.ID] = e
}
}

func (em *EventMgr) GetEvent(id int) *Event {
if value, ok := em.Eventmap[id]; ok {
return value
}
return nil
}

what else do i probably need

You already have the seed of the right idea: an append-only event manager. The main thing missing is not more AI, but enough structure so edits can be replayed and grouped consistently.

Thought for a couple of seconds
You probably need less “more stuff” than it looks like. But you do need a few critical invariants that your current struct does not yet capture.

Right now your model is close to:

“something was written or deleted”

some content

some position

a start/end time

That is a good start, but it is still too weak to reliably answer:

what exactly changed?

where in the text did it change?

what text state did this event apply to?

can I replay all events and reconstruct the textarea?

can I merge adjacent writes/deletes safely?

can I later segment reflections from this history?

The first big thing: decide what an event means
Your current Event is ambiguous.

For example, for delete:

type Event struct {
Type EventType
Content []rune
Pos int
}
What does Pos mean?

line number?

rune offset in the full buffer?

cursor position before delete?

cursor position after delete?

And for Delete, what does Content mean?

the text that was deleted?

empty, because deletion is described only by span?

the text currently at Pos?

You want events to be fully replayable. That means each event should carry enough information that, given an older buffer, you can deterministically apply it.

Minimal replayable model
I would strongly recommend thinking in terms of spans in a flat rune buffer, not line numbers.

Instead of:

Pos int // line number
prefer something like:

Start int // rune offset in document before applying event
End int // rune offset in document before applying event
Text []rune
Then the meaning becomes:

Write/Insert: replace [Start:End] with Text

Delete: replace [Start:End] with empty text

or even better: unify both as Replace

That gives you one core event type.

Strong recommendation: use Replace, not Write/Delete
Instead of:

const (
Write EventType = iota
Delete
)
I would seriously consider:

const (
Replace EventType = iota
)
with:

type Event struct {
ID int64
CreatedAt time.Time
UpdatedAt time.Time

    Start int       // inclusive rune offset in old buffer
    End   int       // exclusive rune offset in old buffer

    Before []rune   // text that used to be there
    After  []rune   // text that is there after the event

}
Why this is much better:

insert = Start == End, Before == "", After != ""

delete = Start < End, Before != "", After == ""

replace/paste-over = Start < End, Before != "", After != ""

This makes your system much more future-proof.

It also handles:

backspace

selection replace

paste over selection

auto-indent

newline insertion

undo/redo recording

IME / multi-rune input

all with one consistent rule.

You need document state versioning
Your current events do not know which text state they apply to.

That is dangerous, because if you later merge or reorder things, Pos may no longer be meaningful.

Add a revision number:

type Event struct {
ID int64
Revision int64 // document revision before applying this event
...
}
or store both:

BaseRevision int64
NextRevision int64
This helps a lot when debugging replay issues.

You need a canonical document buffer
PendingRunes is not enough by itself.

You probably want the event manager, or some sibling struct, to own the current materialized text:

type EventMgr struct {
Events []*Event
EventMap map[int64]*Event
Buffer []rune
Revision int64
PendingEvent \*Event
}
Why:

every event should be applied against a known buffer

delete events need to know what text was removed

merging events depends on adjacency in the current buffer

reflection/segmentation later will likely work off snapshots

Without a canonical buffer, your event history becomes descriptive instead of executable.

You need ApplyEvent
This is the most important missing function.

Something like:

func (em *EventMgr) ApplyEvent(e *Event) error
This should:

validate span bounds

validate Before matches current buffer span

mutate the buffer

append to event log

increment revision

That validation is huge. It catches bugs early.

For example, if the event says:

replace [10:13]

Before = "abc"

but the current buffer at [10:13] is actually "abd"

then something is wrong.

You need coalescing / merging logic
You were already thinking about automerging. Yes, you likely need it.

Not because correctness depends on it, but because raw keystroke-level events are too noisy.

Examples:

typing hello should probably become one logical insert event, not 5

holding backspace should become one delete burst, not 20 tiny deletes

typing within a short time window at adjacent cursor positions should merge

So add some notion of a pending logical event.

For example:

type EventMgr struct {
Events []*Event
EventMap map[int64]*Event
Buffer []rune
Revision int64

    Pending      *Event
    MergeWindow  time.Duration

}
Then rules like:

merge adjacent inserts if:

same event type / replace shape

close in time

contiguous positions

merge repeated backspaces if they extend the same delete region

flush pending event on cursor jump, pause, or mode change

This gives you cleaner higher-level structure.

You need cursor/context metadata
For a diary/editor, event meaning depends on cursor behavior.

At minimum, store:

CursorBefore int
CursorAfter int
Possibly also selection:

SelStartBefore int
SelEndBefore int
SelStartAfter int
SelEndAfter int
This helps later if you want:

robust replay

command grouping

selection replacements

intent inference

You should not use line number as primary position
You already suspected this.

Yes, line number is too weak.

Use:

rune offsets in the whole document as canonical positions

derive line/column only for UI display

Why:

line numbers change too easily

deletes/inserts shift downstream lines

segmentation/search will want stable offsets into a snapshot

You may want event source / cause
This becomes useful surprisingly quickly.

type EventSource int

const (
SourceUserTyping EventSource = iota
SourcePaste
SourceBackspace
SourceDeleteKey
SourceUndo
SourceRedo
SourceProgrammatic
)
This lets you treat things differently later.

For example:

pasted text may deserve different segmentation

undo should maybe not count as “new writing”

programmatic formatting should maybe be ignored in diary analytics

You probably need snapshots too
Even with a perfect event log, replaying from the beginning every time is annoying.

So every so often, store a snapshot:

type Snapshot struct {
Revision int64
CreatedAt time.Time
Buffer []rune
}
Then you can reconstruct by:

loading nearest snapshot

replaying later events

This matters once entries get long.

You need a clear split between raw events and derived segments
Do not make your Event struct carry the burden of segmentation logic.

Events are raw edit facts.

Segments are derived objects, later.

For example:

type Segment struct {
ID int64
StartRev int64
EndRev int64
StartOffset int
EndOffset int
Text []rune
Kind string
}
You do not want to bake “useful pieces” directly into the event struct. That should be a later layer.

Your current code also needs some basic mechanical cleanup
A few practical notes.

1. Init should reset everything
   Right now PendingRunes is not reset.

2. IDs should probably be generated by manager
   Do not rely on caller to always provide a unique ID.

3. Eventmap and Eventlist names
   In Go style, these would usually be:

Events []*Event
EventMap map[int64]*Event 4. UpdatedAt meaning
For a merged event, CreatedAt and UpdatedAt are useful.
For a single atomic event, UpdatedAt is maybe just CreatedAt.
That is fine, but define it clearly.

The minimal set of things you probably still need
If I strip it down to the essentials, you likely need:

In the event itself
stable ID

timestamps

start/end rune offsets

before text

after text

revision

maybe cursor-before/cursor-after

maybe source/cause

In the manager
current document buffer

revision counter

append/apply function

merge/coalescing logic

snapshot support

flush-pending function

Invariants
events are append-only

events are replayable

applying an event is deterministic

positions are rune offsets, not line numbers

segmentation is derived later, not encoded into raw events

A better skeleton
Here is a much stronger base shape:

package events

import "time"

type EventSource int

const (
SourceTyping EventSource = iota
SourceBackspace
SourceDeleteKey
SourcePaste
SourceUndo
SourceRedo
SourceProgrammatic
)

type Event struct {
ID int64
CreatedAt time.Time
UpdatedAt time.Time

    Revision int64 // document revision before applying this event

    Start int // inclusive rune offset in old buffer
    End   int // exclusive rune offset in old buffer

    Before []rune // text replaced/removed
    After  []rune // text inserted/replaced

    CursorBefore int
    CursorAfter  int

    Source EventSource

}

type Snapshot struct {
Revision int64
CreatedAt time.Time
Buffer []rune
}

type EventMgr struct {
Events []*Event
EventMap map[int64]*Event
Buffer []rune
Revision int64
NextID int64
Pending *Event
Snapshots []*Snapshot
MergeWindow time.Duration
}

func (em *EventMgr) Init() {
em.Events = []*Event{}
em.EventMap = make(map[int64]_Event)
em.Buffer = nil
em.Revision = 0
em.NextID = 1
em.Pending = nil
em.Snapshots = nil
em.MergeWindow = 2 _ time.Second
}
Then later add:

ApplyEvent

TryMergePending

FlushPending

MakeSnapshot

ReplayFromSnapshot

The key design advice
Your event system should answer only one question well:

What exact edits happened to the document over time?

It should not yet try to answer:

what idea this belongs to

what theme it is

whether it is important

whether it should be reflected back

That comes later.

Build the event layer as a boring, correct edit log.
That will make everything above it much easier.

If you want, I can next help you turn this into an actual Go implementation of ApplyEvent and TryMergePending for Bubble Tea textarea input.

// Package textarea provides a multi-line text input component for Bubble Tea
// applications.
package textarea

import (
"crypto/sha256"
"fmt"
"image/color"
"slices"
"strconv"
"strings"
"time"
"unicode"

    tea "charm.land/bubbletea/v2"
    "github.com/atotto/clipboard"
    "github.com/charmbracelet/x/ansi"
    "github.com/haochend413/bubbles/v2/cursor"
    "github.com/haochend413/bubbles/v2/internal/memoization"
    "github.com/haochend413/bubbles/v2/internal/runeutil"
    "github.com/haochend413/bubbles/v2/key"
    "github.com/haochend413/bubbles/v2/statusbar"
    "github.com/haochend413/bubbles/v2/viewport"
    "github.com/haochend413/lipgloss/v2"
    rw "github.com/mattn/go-runewidth"
    "github.com/rivo/uniseg"

)

const (
minHeight = 1
defaultHeight = 6
defaultWidth = 40
defaultCharLimit = 0 // no limit
defaultMaxHeight = 99
defaultMaxWidth = 500

    // XXX: in v2, make max lines dynamic and default max lines configurable.
    maxLines = 10000

)

// Internal messages for clipboard operations.
type (
pasteMsg string
pasteErrMsg struct{ error }
)

// KeyMap is the key bindings for different actions within the textarea.
type KeyMap struct {
CharacterBackward key.Binding
CharacterForward key.Binding
DeleteAfterCursor key.Binding
DeleteBeforeCursor key.Binding
DeleteCharacterBackward key.Binding
DeleteCharacterForward key.Binding
DeleteWordBackward key.Binding
DeleteWordForward key.Binding
InsertNewline key.Binding
LineEnd key.Binding
LineNext key.Binding
LinePrevious key.Binding
LineStart key.Binding
PageUp key.Binding
PageDown key.Binding
Paste key.Binding
WordBackward key.Binding
WordForward key.Binding
InputBegin key.Binding
InputEnd key.Binding

    UppercaseWordForward  key.Binding
    LowercaseWordForward  key.Binding
    CapitalizeWordForward key.Binding

    TransposeCharacterBackward key.Binding

    // Keys for toggling between VIEW and INSERT modes
    EnterInsertMode key.Binding
    EnterViewMode   key.Binding

}

// DefaultKeyMap returns the default set of key bindings for navigating and acting
// upon the textarea.
func DefaultKeyMap() KeyMap {
return KeyMap{
CharacterForward: key.NewBinding(key.WithKeys("right", "ctrl+f"), key.WithHelp("right", "character forward")),
CharacterBackward: key.NewBinding(key.WithKeys("left", "ctrl+b"), key.WithHelp("left", "character backward")),
WordForward: key.NewBinding(key.WithKeys("alt+right", "alt+f"), key.WithHelp("alt+right", "word forward")),
WordBackward: key.NewBinding(key.WithKeys("alt+left", "alt+b"), key.WithHelp("alt+left", "word backward")),
LineNext: key.NewBinding(key.WithKeys("down", "ctrl+n"), key.WithHelp("down", "next line")),
LinePrevious: key.NewBinding(key.WithKeys("up", "ctrl+p"), key.WithHelp("up", "previous line")),
DeleteWordBackward: key.NewBinding(key.WithKeys("alt+backspace", "ctrl+w"), key.WithHelp("alt+backspace", "delete word backward")),
DeleteWordForward: key.NewBinding(key.WithKeys("alt+delete", "alt+d"), key.WithHelp("alt+delete", "delete word forward")),
DeleteAfterCursor: key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "delete after cursor")),
DeleteBeforeCursor: key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "delete before cursor")),
InsertNewline: key.NewBinding(key.WithKeys("enter", "ctrl+m"), key.WithHelp("enter", "insert newline")),
DeleteCharacterBackward: key.NewBinding(key.WithKeys("backspace", "ctrl+h"), key.WithHelp("backspace", "delete character backward")),
DeleteCharacterForward: key.NewBinding(key.WithKeys("delete", "ctrl+d"), key.WithHelp("delete", "delete character forward")),
LineStart: key.NewBinding(key.WithKeys("home", "ctrl+a"), key.WithHelp("home", "line start")),
LineEnd: key.NewBinding(key.WithKeys("end", "ctrl+e"), key.WithHelp("end", "line end")),
PageUp: key.NewBinding(key.WithKeys("pgup"), key.WithHelp("pgup", "page up")),
PageDown: key.NewBinding(key.WithKeys("pgdown"), key.WithHelp("pgdown", "page down")),
Paste: key.NewBinding(key.WithKeys("ctrl+v"), key.WithHelp("ctrl+v", "paste")),
InputBegin: key.NewBinding(key.WithKeys("alt+<", "ctrl+home"), key.WithHelp("alt+<", "input begin")),
InputEnd: key.NewBinding(key.WithKeys("alt+>", "ctrl+end"), key.WithHelp("alt+>", "input end")),

    	CapitalizeWordForward: key.NewBinding(key.WithKeys("alt+c"), key.WithHelp("alt+c", "capitalize word forward")),
    	LowercaseWordForward:  key.NewBinding(key.WithKeys("alt+l"), key.WithHelp("alt+l", "lowercase word forward")),
    	UppercaseWordForward:  key.NewBinding(key.WithKeys("alt+u"), key.WithHelp("alt+u", "uppercase word forward")),

    	TransposeCharacterBackward: key.NewBinding(key.WithKeys("ctrl+t"), key.WithHelp("ctrl+t", "transpose character backward")),

    	EnterInsertMode: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "enter insert mode")),
    	EnterViewMode:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "enter view mode")),
    }

}

// LineInfo is a helper for keeping track of line information regarding
// soft-wrapped lines.
type LineInfo struct {
// Width is the number of columns in the line.
Width int

    // CharWidth is the number of characters in the line to account for
    // double-width runes.
    CharWidth int

    // Height is the number of rows in the line.
    Height int

    // StartColumn is the index of the first column of the line.
    StartColumn int

    // ColumnOffset is the number of columns that the cursor is offset from the
    // start of the line.
    ColumnOffset int

    // RowOffset is the number of rows that the cursor is offset from the start
    // of the line.
    RowOffset int

    // CharOffset is the number of characters that the cursor is offset
    // from the start of the line. This will generally be equivalent to
    // ColumnOffset, but will be different there are double-width runes before
    // the cursor.
    CharOffset int

}

// PromptInfo is a struct that can be used to store information about the
// prompt.
type PromptInfo struct {
LineNumber int
Focused bool
}

// CursorStyle is the style for real and virtual cursors.
type CursorStyle struct {
// Style styles the cursor block.
//
// For real cursors, the foreground color set here will be used as the
// cursor color.
Color color.Color

    // Shape is the cursor shape. The following shapes are available:
    //
    // - tea.CursorBlock
    // - tea.CursorUnderline
    // - tea.CursorBar
    //
    // This is only used for real cursors.
    Shape tea.CursorShape

    // CursorBlink determines whether or not the cursor should blink.
    Blink bool

    // BlinkSpeed is the speed at which the virtual cursor blinks. This has no
    // effect on real cursors as well as no effect if the cursor is set not to
    // [CursorBlink].
    //
    // By default, the blink speed is set to about 500ms.
    BlinkSpeed time.Duration

}

// Styles are the styles for the textarea, separated into focused and blurred
// states. The appropriate styles will be chosen based on the focus state of
// the textarea.
type Styles struct {
Focused StyleState
Blurred StyleState
Cursor CursorStyle
}

// StyleState that will be applied to the text area.
//
// StyleState can be applied to focused and unfocused states to change the styles
// depending on the focus state.
//
// For an introduction to styling with Lip Gloss see:
// https://github.com/charmbracelet/lipgloss
type StyleState struct {
Base lipgloss.Style
Text lipgloss.Style
LineNumber lipgloss.Style
CursorLineNumber lipgloss.Style
CursorLine lipgloss.Style
EndOfBuffer lipgloss.Style
Placeholder lipgloss.Style
Prompt lipgloss.Style
}

func (s StyleState) computedCursorLine() lipgloss.Style {
return s.CursorLine.Inherit(s.Base).Inline(true)
}

func (s StyleState) computedCursorLineNumber() lipgloss.Style {
return s.CursorLineNumber.
Inherit(s.CursorLine).
Inherit(s.Base).
Inline(true)
}

func (s StyleState) computedEndOfBuffer() lipgloss.Style {
return s.EndOfBuffer.Inherit(s.Base).Inline(true)
}

func (s StyleState) computedLineNumber() lipgloss.Style {
return s.LineNumber.Inherit(s.Base).Inline(true)
}

func (s StyleState) computedPlaceholder() lipgloss.Style {
return s.Placeholder.Inherit(s.Base).Inline(true)
}

func (s StyleState) computedPrompt() lipgloss.Style {
return s.Prompt.Inherit(s.Base).Inline(true)
}

func (s StyleState) computedText() lipgloss.Style {
return s.Text.Inherit(s.Base).Inline(true)
}

// line is the input to the text wrapping function. This is stored in a struct
// so that it can be hashed and memoized.
type line struct {
runes []rune
width int
}

// Hash returns a hash of the line.
func (w line) Hash() string {
v := fmt.Sprintf("%s:%d", string(w.runes), w.width)
return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))
}

// Model is the Bubble Tea model for this text area element.
type Model struct {
Err error

    // General settings.
    cache *memoization.MemoCache[line, [][]rune]

    // Prompt is printed at the beginning of each line.
    //
    // When changing the value of Prompt after the model has been
    // initialized, ensure that SetWidth() gets called afterwards.
    //
    // See also [SetPromptFunc] for a dynamic prompt.
    Prompt string

    // Placeholder is the text displayed when the user
    // hasn't entered anything yet.
    Placeholder string

    // ShowLineNumbers, if enabled, causes line numbers to be printed
    // after the prompt.
    ShowLineNumbers bool

    // EndOfBufferCharacter is displayed at the end of the input.
    EndOfBufferCharacter rune

    // KeyMap encodes the keybindings recognized by the widget.
    KeyMap KeyMap

    // virtualCursor manages the virtual cursor.
    virtualCursor cursor.Model

    // CharLimit is the maximum number of characters this input element will
    // accept. If 0 or less, there's no limit.
    CharLimit int

    // MaxHeight is the maximum height of the text area in rows. If 0 or less,
    // there's no limit.
    MaxHeight int

    // MaxWidth is the maximum width of the text area in columns. If 0 or less,
    // there's no limit.
    MaxWidth int

    // Styling. Styles are defined in [Styles]. Use [SetStyles] and [GetStyles]
    // to work with this value publicly.
    styles Styles

    // useVirtualCursor determines whether or not to use the virtual cursor.
    // Use [SetVirtualCursor] and [VirtualCursor] to work with this this
    // value publicly.
    useVirtualCursor bool

    // If promptFunc is set, it replaces Prompt as a generator for
    // prompt strings at the beginning of each line.
    promptFunc func(PromptInfo) string

    // promptWidth is the width of the prompt.
    promptWidth int

    // width is the maximum number of characters that can be displayed at once.
    // If 0 or less this setting is ignored.
    width int

    // height is the maximum number of lines that can be displayed at once. It
    // essentially treats the text field like a vertically scrolling viewport
    // if there are more lines than the permitted height.
    height int

    // Underlying text value.
    value [][]rune

    // focus indicates whether user input focus should be on this input
    // component. When false, ignore keyboard input and hide the cursor.
    focus bool

    // Cursor column.
    col int

    // Cursor row.
    row int

    // Last character offset, used to maintain state when the cursor is moved
    // vertically such that we can maintain the same navigating position.
    lastCharOffset int

    // viewport is the vertically-scrollable viewport of the multi-line text
    // input.
    viewport *viewport.Model

    // rune sanitizer for input.
    rsan runeutil.Sanitizer

    // InsertMode controls whether the textarea is in INSERT mode (true) or VIEW mode (false)
    InsertMode bool

    // Statusbar displays mode and word/character count
    Statusbar *statusbar.Model

}

// New creates a new model with default settings.
func New() Model {
vp := viewport.New()
vp.KeyMap = viewport.KeyMap{}
cur := cursor.New()

    styles := DefaultDarkStyles()

    m := Model{
    	CharLimit:            defaultCharLimit,
    	MaxHeight:            defaultMaxHeight,
    	MaxWidth:             defaultMaxWidth,
    	Prompt:               lipgloss.ThickBorder().Left + " ",
    	styles:               styles,
    	cache:                memoization.NewMemoCache[line, [][]rune](maxLines),
    	EndOfBufferCharacter: ' ',
    	ShowLineNumbers:      true,
    	useVirtualCursor:     true,
    	virtualCursor:        cur,
    	KeyMap:               DefaultKeyMap(),

    	value: make([][]rune, minHeight, maxLines),
    	focus: false,
    	col:   0,
    	row:   0,

    	viewport: &vp,

    	InsertMode: true, // Start in INSERT mode by default
    }

    // Initialize statusbar
    sb := statusbar.New(
    	statusbar.WithHeight(1),
    	statusbar.WithWidth(defaultWidth),
    )
    modeElem := sb.AddLeft(8, "INSERT")
    modeElem.SetColors("0", "34")
    sb.SetTag(modeElem, "mode")

    countElem := sb.AddLeft(25, "0 chars | 0 words")
    countElem.SetColors("252", "236")
    sb.SetTag(countElem, "count")

    m.Statusbar = &sb

    m.SetHeight(defaultHeight)
    m.SetWidth(defaultWidth)

    return m

}

// DefaultStyles returns the default styles for focused and blurred states for
// the textarea.
func DefaultStyles(isDark bool) Styles {
lightDark := lipgloss.LightDark(isDark)

    var s Styles
    s.Focused = StyleState{
    	Base:             lipgloss.NewStyle(),
    	CursorLine:       lipgloss.NewStyle().Background(lightDark(lipgloss.Color("255"), lipgloss.Color("0"))),
    	CursorLineNumber: lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("240"), lipgloss.Color("240"))),
    	EndOfBuffer:      lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("254"), lipgloss.Color("0"))),
    	LineNumber:       lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("249"), lipgloss.Color("7"))),
    	Placeholder:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
    	Prompt:           lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
    	Text:             lipgloss.NewStyle(),
    }
    s.Blurred = StyleState{
    	Base:             lipgloss.NewStyle(),
    	CursorLine:       lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("245"), lipgloss.Color("7"))),
    	CursorLineNumber: lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("249"), lipgloss.Color("7"))),
    	EndOfBuffer:      lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("254"), lipgloss.Color("0"))),
    	LineNumber:       lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("249"), lipgloss.Color("7"))),
    	Placeholder:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
    	Prompt:           lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
    	Text:             lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("245"), lipgloss.Color("7"))),
    }
    s.Cursor = CursorStyle{
    	Color: lipgloss.Color("7"),
    	Shape: tea.CursorBlock,
    	Blink: true,
    }
    return s

}

// DefaultLightStyles returns the default styles for a light background.
func DefaultLightStyles() Styles {
return DefaultStyles(false)
}

// DefaultDarkStyles returns the default styles for a dark background.
func DefaultDarkStyles() Styles {
return DefaultStyles(true)
}

// Styles returns the current styles for the textarea.
func (m Model) Styles() Styles {
return m.styles
}

// SetStyles updates styling for the textarea.
func (m \*Model) SetStyles(s Styles) {
m.styles = s
m.updateVirtualCursorStyle()
}

// VirtualCursor returns whether or not the virtual cursor is enabled.
func (m Model) VirtualCursor() bool {
return m.useVirtualCursor
}

// SetVirtualCursor sets whether or not to use the virtual cursor.
func (m \*Model) SetVirtualCursor(v bool) {
m.useVirtualCursor = v
m.updateVirtualCursorStyle()
}

// updateVirtualCursorStyle sets styling on the virtual cursor based on the
// textarea's style settings.
func (m \*Model) updateVirtualCursorStyle() {
if !m.useVirtualCursor {
m.virtualCursor.SetMode(cursor.CursorHide)
return
}

    m.virtualCursor.Style = lipgloss.NewStyle().Foreground(m.styles.Cursor.Color)

    // By default, the blink speed of the cursor is set to a default
    // internally.
    if m.styles.Cursor.Blink {
    	if m.styles.Cursor.BlinkSpeed > 0 {
    		m.virtualCursor.BlinkSpeed = m.styles.Cursor.BlinkSpeed
    	}
    	m.virtualCursor.SetMode(cursor.CursorBlink)
    	return
    }
    m.virtualCursor.SetMode(cursor.CursorStatic)

}

// SetValue sets the value of the text input.
func (m \*Model) SetValue(s string) {
m.Reset()
m.InsertString(s)
}

// InsertString inserts a string at the cursor position.
func (m \*Model) InsertString(s string) {
m.insertRunesFromUserInput([]rune(s))
}

// InsertRune inserts a rune at the cursor position.
func (m \*Model) InsertRune(r rune) {
m.insertRunesFromUserInput([]rune{r})
}

// SetInsertMode sets the textarea to INSERT mode
func (m \*Model) SetInsertMode() {
m.InsertMode = true
if m.Statusbar != nil {
if elem := m.Statusbar.GetTag("mode"); elem != nil {
elem.SetValue("INSERT").SetColors("0", "34")
}
}
}

// SetViewMode sets the textarea to VIEW mode
func (m \*Model) SetViewMode() {
m.InsertMode = false
if m.Statusbar != nil {
if elem := m.Statusbar.GetTag("mode"); elem != nil {
elem.SetValue("VIEW").SetColors("0", "39")
}
}
}

// updateWordCount updates the character and word count in the statusbar
func (m \*Model) updateWordCount() {
if m.Statusbar == nil {
return
}
elem := m.Statusbar.GetTag("count")
if elem == nil {
return
}
text := m.Value()
chars := len([]rune(text))
words := 0
if len(strings.TrimSpace(text)) > 0 {
words = len(strings.Fields(text))
}
elem.SetValue(fmt.Sprintf("%d chars | %d words", chars, words))
}

// insertRunesFromUserInput inserts runes at the current cursor position.
func (m \*Model) insertRunesFromUserInput(runes []rune) {
// Clean up any special characters in the input provided by the
// clipboard. This avoids bugs due to e.g. tab characters and
// whatnot.
runes = m.san().Sanitize(runes)

    if m.CharLimit > 0 {
    	availSpace := m.CharLimit - m.Length()
    	// If the char limit's been reached, cancel.
    	if availSpace <= 0 {
    		return
    	}
    	// If there's not enough space to paste the whole thing cut the pasted
    	// runes down so they'll fit.
    	if availSpace < len(runes) {
    		runes = runes[:availSpace]
    	}
    }

    // Split the input into lines.
    var lines [][]rune
    lstart := 0
    for i := range runes {
    	if runes[i] == '\n' {
    		// Queue a line to become a new row in the text area below.
    		// Beware to clamp the max capacity of the slice, to ensure no
    		// data from different rows get overwritten when later edits
    		// will modify this line.
    		lines = append(lines, runes[lstart:i:i])
    		lstart = i + 1
    	}
    }
    if lstart <= len(runes) {
    	// The last line did not end with a newline character.
    	// Take it now.
    	lines = append(lines, runes[lstart:])
    }

    // Obey the maximum line limit.
    if maxLines > 0 && len(m.value)+len(lines)-1 > maxLines {
    	allowedHeight := max(0, maxLines-len(m.value)+1)
    	lines = lines[:allowedHeight]
    }

    if len(lines) == 0 {
    	// Nothing left to insert.
    	return
    }

    // Save the remainder of the original line at the current
    // cursor position.
    tail := make([]rune, len(m.value[m.row][m.col:]))
    copy(tail, m.value[m.row][m.col:])

    // Paste the first line at the current cursor position.
    m.value[m.row] = append(m.value[m.row][:m.col], lines[0]...)
    m.col += len(lines[0])

    if numExtraLines := len(lines) - 1; numExtraLines > 0 {
    	// Add the new lines.
    	// We try to reuse the slice if there's already space.
    	var newGrid [][]rune
    	if cap(m.value) >= len(m.value)+numExtraLines {
    		// Can reuse the extra space.
    		newGrid = m.value[:len(m.value)+numExtraLines]
    	} else {
    		// No space left; need a new slice.
    		newGrid = make([][]rune, len(m.value)+numExtraLines)
    		copy(newGrid, m.value[:m.row+1])
    	}
    	// Add all the rows that were after the cursor in the original
    	// grid at the end of the new grid.
    	copy(newGrid[m.row+1+numExtraLines:], m.value[m.row+1:])
    	m.value = newGrid
    	// Insert all the new lines in the middle.
    	for _, l := range lines[1:] {
    		m.row++
    		m.value[m.row] = l
    		m.col = len(l)
    	}
    }

    // Finally add the tail at the end of the last line inserted.
    m.value[m.row] = append(m.value[m.row], tail...)

    m.SetCursorColumn(m.col)

}

// Value returns the value of the text input.
func (m Model) Value() string {
if m.value == nil {
return ""
}

    var v strings.Builder
    for _, l := range m.value {
    	v.WriteString(string(l))
    	v.WriteByte('\n')
    }

    return strings.TrimSuffix(v.String(), "\n")

}

// Length returns the number of characters currently in the text input.
func (m \*Model) Length() int {
var l int
for \_, row := range m.value {
l += uniseg.StringWidth(string(row))
}
// We add len(m.value) to include the newline characters.
return l + len(m.value) - 1
}

// LineCount returns the number of lines that are currently in the text input.
func (m \*Model) LineCount() int {
return len(m.value)
}

// Line returns the 0-indexed row position of the cursor.
func (m Model) Line() int {
return m.row
}

// Column returns the 0-indexed column position of the cursor.
func (m Model) Column() int {
return m.col
}

// ScrollYOffset returns the Y offset (top row) index of the current view, which
// can be used to calculate the current scroll position.
func (m Model) ScrollYOffset() int {
return m.viewport.YOffset()
}

// ScrollPercent returns the amount of the textarea that is currently scrolled
// through, clamped between 0 and 1.
func (m Model) ScrollPercent() float64 {
return m.viewport.ScrollPercent()
}

// setCursorLineRelative moves the cursor by the given number of lines. Negative
// values move the cursor up, positive values move the cursor down.
func (m \*Model) setCursorLineRelative(delta int) {
if delta == 0 {
return
}

    li := m.LineInfo()
    charOffset := max(m.lastCharOffset, li.CharOffset)
    m.lastCharOffset = charOffset

    // 2 columns to account for the trailing space wrapping.
    const trailingSpace = 2

    if delta > 0 { //nolint:nestif
    	// Moving down.
    	for range delta {
    		if li.RowOffset+1 >= li.Height && m.row < len(m.value)-1 {
    			m.row++
    			m.col = 0
    		} else {
    			// Move the cursor to the start of the next virtual line.
    			m.col = min(li.StartColumn+li.Width+trailingSpace, len(m.value[m.row])-1)
    		}
    		li = m.LineInfo()
    	}
    } else {
    	// Moving up.
    	for range -delta {
    		if li.RowOffset <= 0 && m.row > 0 {
    			m.row--
    			m.col = len(m.value[m.row])
    		} else {
    			// Move the cursor to the end of the previous line.
    			m.col = li.StartColumn - trailingSpace
    		}
    		li = m.LineInfo()
    	}
    }

    nli := m.LineInfo()
    m.col = nli.StartColumn

    if nli.Width <= 0 {
    	m.repositionView()
    	return
    }

    offset := 0
    for offset < charOffset {
    	if m.row >= len(m.value) || m.col >= len(m.value[m.row]) || offset >= nli.CharWidth-1 {
    		break
    	}
    	offset += rw.RuneWidth(m.value[m.row][m.col])
    	m.col++
    }
    m.repositionView()

}

// CursorDown moves the cursor down by one line.
func (m \*Model) CursorDown() {
m.setCursorLineRelative(1)
}

// CursorUp moves the cursor up by one line.
func (m \*Model) CursorUp() {
m.setCursorLineRelative(-1)
}

// SetCursorColumn moves the cursor to the given position. If the position is
// out of bounds the cursor will be moved to the start or end accordingly.
func (m \*Model) SetCursorColumn(col int) {
m.col = clamp(col, 0, len(m.value[m.row]))
// Any time that we move the cursor horizontally we need to reset the last
// offset so that the horizontal position when navigating is adjusted.
m.lastCharOffset = 0
}

// CursorStart moves the cursor to the start of the input field.
func (m \*Model) CursorStart() {
m.SetCursorColumn(0)
}

// CursorEnd moves the cursor to the end of the input field.
func (m \*Model) CursorEnd() {
m.SetCursorColumn(len(m.value[m.row]))
}

// Focused returns the focus state on the model.
func (m Model) Focused() bool {
return m.focus
}

// activeStyle returns the appropriate set of styles to use depending on
// whether the textarea is focused or blurred.
func (m Model) activeStyle() \*StyleState {
if m.focus {
return &m.styles.Focused
}
return &m.styles.Blurred
}

// Focus sets the focus state on the model. When the model is in focus it can
// receive keyboard input and the cursor will be hidden.
func (m \*Model) Focus() tea.Cmd {
m.focus = true
return m.virtualCursor.Focus()
}

// Blur removes the focus state on the model. When the model is blurred it can
// not receive keyboard input and the cursor will be hidden.
func (m \*Model) Blur() {
m.focus = false
m.virtualCursor.Blur()
}

// Reset sets the input to its default state with no input.
func (m \*Model) Reset() {
m.value = make([][]rune, minHeight, maxLines)
m.col = 0
m.row = 0
m.viewport.GotoTop()
m.SetCursorColumn(0)
}

// Word returns the word at the cursor position.
// A word is delimited by spaces or line-breaks.
func (m \*Model) Word() string {
line := m.value[m.row]
col := m.col - 1

    if col < 0 {
    	return ""
    }

    // If cursor is beyond the line, return empty string
    if col >= len(line) {
    	return ""
    }

    // If cursor is on a space, return empty string
    if unicode.IsSpace(line[col]) {
    	return ""
    }

    // Find the start of the word by moving left
    start := col
    for start > 0 && !unicode.IsSpace(line[start-1]) {
    	start--
    }

    // Find the end of the word by moving right
    end := col
    for end < len(line) && !unicode.IsSpace(line[end]) {
    	end++
    }

    return string(line[start:end])

}

// san initializes or retrieves the rune sanitizer.
func (m \*Model) san() runeutil.Sanitizer {
if m.rsan == nil {
// Textinput has all its input on a single line so collapse
// newlines/tabs to single spaces.
m.rsan = runeutil.NewSanitizer()
}
return m.rsan
}

// deleteBeforeCursor deletes all text before the cursor. Returns whether or
// not the cursor blink should be reset.
func (m \*Model) deleteBeforeCursor() {
m.value[m.row] = m.value[m.row][m.col:]
m.SetCursorColumn(0)
}

// deleteAfterCursor deletes all text after the cursor. Returns whether or not
// the cursor blink should be reset. If input is masked delete everything after
// the cursor so as not to reveal word breaks in the masked input.
func (m \*Model) deleteAfterCursor() {
m.value[m.row] = m.value[m.row][:m.col]
m.SetCursorColumn(len(m.value[m.row]))
}

// transposeLeft exchanges the runes at the cursor and immediately
// before. No-op if the cursor is at the beginning of the line. If
// the cursor is not at the end of the line yet, moves the cursor to
// the right.
func (m \*Model) transposeLeft() {
if m.col == 0 || len(m.value[m.row]) < 2 {
return
}
if m.col >= len(m.value[m.row]) {
m.SetCursorColumn(m.col - 1)
}
m.value[m.row][m.col-1], m.value[m.row][m.col] = m.value[m.row][m.col], m.value[m.row][m.col-1]
if m.col < len(m.value[m.row]) {
m.SetCursorColumn(m.col + 1)
}
}

// deleteWordLeft deletes the word left to the cursor. Returns whether or not
// the cursor blink should be reset.
func (m \*Model) deleteWordLeft() {
if m.col == 0 || len(m.value[m.row]) == 0 {
return
}

    // Linter note: it's critical that we acquire the initial cursor position
    // here prior to altering it via SetCursor() below. As such, moving this
    // call into the corresponding if clause does not apply here.
    oldCol := m.col

    m.SetCursorColumn(m.col - 1)
    for unicode.IsSpace(m.value[m.row][m.col]) {
    	if m.col <= 0 {
    		break
    	}
    	// ignore series of whitespace before cursor
    	m.SetCursorColumn(m.col - 1)
    }

    for m.col > 0 {
    	if !unicode.IsSpace(m.value[m.row][m.col]) {
    		m.SetCursorColumn(m.col - 1)
    	} else {
    		if m.col > 0 {
    			// keep the previous space
    			m.SetCursorColumn(m.col + 1)
    		}
    		break
    	}
    }

    if oldCol > len(m.value[m.row]) {
    	m.value[m.row] = m.value[m.row][:m.col]
    } else {
    	m.value[m.row] = append(m.value[m.row][:m.col], m.value[m.row][oldCol:]...)
    }

}

// deleteWordRight deletes the word right to the cursor.
func (m \*Model) deleteWordRight() {
if m.col >= len(m.value[m.row]) || len(m.value[m.row]) == 0 {
return
}

    oldCol := m.col

    for m.col < len(m.value[m.row]) && unicode.IsSpace(m.value[m.row][m.col]) {
    	// ignore series of whitespace after cursor
    	m.SetCursorColumn(m.col + 1)
    }

    for m.col < len(m.value[m.row]) {
    	if !unicode.IsSpace(m.value[m.row][m.col]) {
    		m.SetCursorColumn(m.col + 1)
    	} else {
    		break
    	}
    }

    if m.col > len(m.value[m.row]) {
    	m.value[m.row] = m.value[m.row][:oldCol]
    } else {
    	m.value[m.row] = append(m.value[m.row][:oldCol], m.value[m.row][m.col:]...)
    }

    m.SetCursorColumn(oldCol)

}

// characterRight moves the cursor one character to the right.
func (m \*Model) characterRight() {
if m.col < len(m.value[m.row]) {
m.SetCursorColumn(m.col + 1)
} else {
if m.row < len(m.value)-1 {
m.row++
m.CursorStart()
}
}
}

// characterLeft moves the cursor one character to the left.
// If insideLine is set, the cursor is moved to the last
// character in the previous line, instead of one past that.
func (m \*Model) characterLeft(insideLine bool) {
if m.col == 0 && m.row != 0 {
m.row--
m.CursorEnd()
if !insideLine {
return
}
}
if m.col > 0 {
m.SetCursorColumn(m.col - 1)
}
}

// wordLeft moves the cursor one word to the left. Returns whether or not the
// cursor blink should be reset. If input is masked, move input to the start
// so as not to reveal word breaks in the masked input.
func (m _Model) wordLeft() {
for {
m.characterLeft(true /_ insideLine \*/)
if m.col < len(m.value[m.row]) && !unicode.IsSpace(m.value[m.row][m.col]) {
break
}
}

    for m.col > 0 {
    	if unicode.IsSpace(m.value[m.row][m.col-1]) {
    		break
    	}
    	m.SetCursorColumn(m.col - 1)
    }

}

// wordRight moves the cursor one word to the right. Returns whether or not the
// cursor blink should be reset. If the input is masked, move input to the end
// so as not to reveal word breaks in the masked input.
func (m _Model) wordRight() {
m.doWordRight(func(int, int) { /_ nothing \*/ })
}

func (m \*Model) doWordRight(fn func(charIdx int, pos int)) {
// Skip spaces forward.
for m.col >= len(m.value[m.row]) || unicode.IsSpace(m.value[m.row][m.col]) {
if m.row == len(m.value)-1 && m.col == len(m.value[m.row]) {
// End of text.
break
}
m.characterRight()
}

    charIdx := 0
    for m.col < len(m.value[m.row]) {
    	if unicode.IsSpace(m.value[m.row][m.col]) {
    		break
    	}
    	fn(charIdx, m.col)
    	m.SetCursorColumn(m.col + 1)
    	charIdx++
    }

}

// uppercaseRight changes the word to the right to uppercase.
func (m \*Model) uppercaseRight() {
m.doWordRight(func(\_ int, i int) {
m.value[m.row][i] = unicode.ToUpper(m.value[m.row][i])
})
}

// lowercaseRight changes the word to the right to lowercase.
func (m \*Model) lowercaseRight() {
m.doWordRight(func(\_ int, i int) {
m.value[m.row][i] = unicode.ToLower(m.value[m.row][i])
})
}

// capitalizeRight changes the word to the right to title case.
func (m \*Model) capitalizeRight() {
m.doWordRight(func(charIdx int, i int) {
if charIdx == 0 {
m.value[m.row][i] = unicode.ToTitle(m.value[m.row][i])
}
})
}

// LineInfo returns the number of characters from the start of the
// (soft-wrapped) line and the (soft-wrapped) line width.
func (m Model) LineInfo() LineInfo {
grid := m.memoizedWrap(m.value[m.row], m.width)

    // Find out which line we are currently on. This can be determined by the
    // m.col and counting the number of runes that we need to skip.
    var counter int
    for i, line := range grid {
    	// We've found the line that we are on
    	if counter+len(line) == m.col && i+1 < len(grid) {
    		// We wrap around to the next line if we are at the end of the
    		// previous line so that we can be at the very beginning of the row
    		return LineInfo{
    			CharOffset:   0,
    			ColumnOffset: 0,
    			Height:       len(grid),
    			RowOffset:    i + 1,
    			StartColumn:  m.col,
    			Width:        len(grid[i+1]),
    			CharWidth:    uniseg.StringWidth(string(line)),
    		}
    	}

    	if counter+len(line) >= m.col {
    		return LineInfo{
    			CharOffset:   uniseg.StringWidth(string(line[:max(0, m.col-counter)])),
    			ColumnOffset: m.col - counter,
    			Height:       len(grid),
    			RowOffset:    i,
    			StartColumn:  counter,
    			Width:        len(line),
    			CharWidth:    uniseg.StringWidth(string(line)),
    		}
    	}

    	counter += len(line)
    }
    return LineInfo{}

}

// repositionView repositions the view of the viewport based on the defined
// scrolling behavior.
func (m \*Model) repositionView() {
minimum := m.viewport.YOffset()
maximum := minimum + m.viewport.Height() - 1
if row := m.cursorLineNumber(); row < minimum {
m.viewport.ScrollUp(minimum - row)
} else if row > maximum {
m.viewport.ScrollDown(row - maximum)
}
}

// Width returns the width of the textarea.
func (m Model) Width() int {
return m.width
}

// MoveToBegin moves the cursor to the beginning of the input.
func (m \*Model) MoveToBegin() {
m.row = 0
m.SetCursorColumn(0)
m.repositionView()
}

// MoveToEnd moves the cursor to the end of the input.
func (m \*Model) MoveToEnd() {
m.row = len(m.value) - 1
m.SetCursorColumn(len(m.value[m.row]))
m.repositionView()
}

// PageUp moves the cursor up by one page. First call snaps to the first visible
// line, subsequent calls move up by a full page.
func (m \*Model) PageUp() {
// If not on the first visible line, snap to it.
if offset := m.viewport.YOffset() - m.cursorLineNumber(); offset < 0 {
m.setCursorLineRelative(offset)
return
}

    // Already on first visible line, move up by a full page.
    m.setCursorLineRelative(-m.height)

}

// PageDown moves the cursor down by one page. First call snaps to the last
// visible line, subsequent calls move down by a full page.
func (m \*Model) PageDown() {
// If not on the last visible line, snap to it.
if offset := m.cursorLineNumber() - m.viewport.YOffset(); offset < m.height-1 {
m.setCursorLineRelative(m.height - 1 - offset)
return
}

    // Already on last visible line, move down by a full page.
    m.setCursorLineRelative(m.height)

}

// SetWidth sets the width of the textarea to fit exactly within the given width.
// This means that the textarea will account for the width of the prompt and
// whether or not line numbers are being shown.
//
// Ensure that SetWidth is called after setting the Prompt and ShowLineNumbers,
// It is important that the width of the textarea be exactly the given width
// and no more.
func (m \*Model) SetWidth(w int) {
// Update prompt width only if there is no prompt function as
// [SetPromptFunc] updates the prompt width when it is called.
if m.promptFunc == nil {
// XXX: Do we even need this or can we calculate the prompt width
// at render time?
m.promptWidth = uniseg.StringWidth(m.Prompt)
}

    // Add base style borders and padding to reserved outer width.
    reservedOuter := m.activeStyle().Base.GetHorizontalFrameSize()

    // Add prompt width to reserved inner width.
    reservedInner := m.promptWidth

    // Add line number width to reserved inner width.
    if m.ShowLineNumbers {
    	// XXX: this was originally documented as needing "1 cell" but was,
    	// in practice, effectively hardcoded to 2 cells. We can, and should,
    	// reduce this to one gap and update the tests accordingly.
    	const gap = 2

    	// Number of digits plus 1 cell for the margin.
    	reservedInner += numDigits(m.MaxHeight) + gap
    }

    // Input width must be at least one more than the reserved inner and outer
    // width. This gives us a minimum input width of 1.
    minWidth := reservedInner + reservedOuter + 1
    inputWidth := max(w, minWidth)

    // Input width must be no more than maximum width.
    if m.MaxWidth > 0 {
    	inputWidth = min(inputWidth, m.MaxWidth)
    }

    // Since the width of the viewport and input area is dependent on the width of
    // borders, prompt and line numbers, we need to calculate it by subtracting
    // the reserved width from them.

    m.viewport.SetWidth(inputWidth - reservedOuter)
    m.width = inputWidth - reservedOuter - reservedInner

    // Update statusbar width to match
    if m.Statusbar != nil {
    	m.Statusbar.SetWidth(inputWidth)
    }

}

// SetPromptFunc supersedes the Prompt field and sets a dynamic prompt instead.
//
// If the function returns a prompt that is shorter than the specified
// promptWidth, it will be padded to the left. If it returns a prompt that is
// longer, display artifacts may occur; the caller is responsible for computing
// an adequate promptWidth.
func (m \*Model) SetPromptFunc(promptWidth int, fn func(PromptInfo) string) {
m.promptFunc = fn
m.promptWidth = promptWidth
}

// Height returns the current height of the textarea.
func (m Model) Height() int {
return m.height
}

// SetHeight sets the height of the textarea.
func (m \*Model) SetHeight(h int) {
if m.MaxHeight > 0 {
m.height = clamp(h, minHeight, m.MaxHeight)
m.viewport.SetHeight(clamp(h, minHeight, m.MaxHeight))
} else {
m.height = max(h, minHeight)
m.viewport.SetHeight(max(h, minHeight))
}

    m.repositionView()

}

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
if !m.focus {
m.virtualCursor.Blur()
return m, nil
}

    // Used to determine if the cursor should blink.
    oldRow, oldCol := m.cursorLineNumber(), m.col

    var cmds []tea.Cmd

    if m.value[m.row] == nil {
    	m.value[m.row] = make([]rune, 0)
    }

    if m.MaxHeight > 0 && m.MaxHeight != m.cache.Capacity() {
    	m.cache = memoization.NewMemoCache[line, [][]rune](m.MaxHeight)
    }

    switch msg := msg.(type) {
    case tea.PasteMsg:
    	m.insertRunesFromUserInput([]rune(msg.Content))
    case tea.KeyPressMsg:
    	// Handle mode switching first
    	switch {
    	case key.Matches(msg, m.KeyMap.EnterViewMode):
    		m.SetViewMode()
    		m.updateWordCount()
    		return m, nil
    	case key.Matches(msg, m.KeyMap.EnterInsertMode):
    		m.SetInsertMode()
    		m.updateWordCount()
    		return m, nil
    	}

    	// Only allow editing in INSERT mode
    	if !m.InsertMode {
    		// VIEW mode - navigation only
    		switch {
    		case key.Matches(msg, m.KeyMap.LineEnd):
    			m.CursorEnd()
    		case key.Matches(msg, m.KeyMap.LineStart):
    			m.CursorStart()
    		case key.Matches(msg, m.KeyMap.CharacterForward):
    			m.characterRight()
    		case key.Matches(msg, m.KeyMap.LineNext):
    			m.setCursorLineRelative(+1)
    		case key.Matches(msg, m.KeyMap.WordForward):
    			m.wordRight()
    		case key.Matches(msg, m.KeyMap.CharacterBackward):
    			m.characterLeft(false)
    		case key.Matches(msg, m.KeyMap.LinePrevious):
    			m.setCursorLineRelative(-1)
    		case key.Matches(msg, m.KeyMap.WordBackward):
    			m.wordLeft()
    		case key.Matches(msg, m.KeyMap.InputBegin):
    			m.MoveToBegin()
    		case key.Matches(msg, m.KeyMap.InputEnd):
    			m.MoveToEnd()
    		case key.Matches(msg, m.KeyMap.PageUp):
    			m.PageUp()
    		case key.Matches(msg, m.KeyMap.PageDown):
    			m.PageDown()
    		}
    		// Skip all editing keys in VIEW mode
    		break
    	}

    	// INSERT mode - full editing capabilities
    	switch {
    	case key.Matches(msg, m.KeyMap.DeleteAfterCursor):
    		m.col = clamp(m.col, 0, len(m.value[m.row]))
    		if m.col >= len(m.value[m.row]) {
    			m.mergeLineBelow(m.row)
    			m.updateWordCount()
    			break
    		}
    		m.deleteAfterCursor()
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.DeleteBeforeCursor):
    		m.col = clamp(m.col, 0, len(m.value[m.row]))
    		if m.col <= 0 {
    			m.mergeLineAbove(m.row)
    			m.updateWordCount()
    			break
    		}
    		m.deleteBeforeCursor()
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.DeleteCharacterBackward):
    		m.col = clamp(m.col, 0, len(m.value[m.row]))
    		if m.col <= 0 {
    			m.mergeLineAbove(m.row)
    			m.updateWordCount()
    			break
    		}
    		if len(m.value[m.row]) > 0 {
    			m.value[m.row] = append(m.value[m.row][:max(0, m.col-1)], m.value[m.row][m.col:]...)
    			if m.col > 0 {
    				m.SetCursorColumn(m.col - 1)
    			}
    		}
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.DeleteCharacterForward):
    		if len(m.value[m.row]) > 0 && m.col < len(m.value[m.row]) {
    			m.value[m.row] = slices.Delete(m.value[m.row], m.col, m.col+1)
    		}
    		if m.col >= len(m.value[m.row]) {
    			m.mergeLineBelow(m.row)
    		}
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.DeleteWordBackward):
    		if m.col <= 0 {
    			m.mergeLineAbove(m.row)
    			m.updateWordCount()
    			break
    		}
    		m.deleteWordLeft()
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.DeleteWordForward):
    		m.col = clamp(m.col, 0, len(m.value[m.row]))
    		if m.col >= len(m.value[m.row]) {
    			m.mergeLineBelow(m.row)
    			m.updateWordCount()
    			break
    		}
    		m.deleteWordRight()
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.InsertNewline):
    		if m.MaxHeight > 0 && len(m.value) >= m.MaxHeight {
    			return m, nil
    		}
    		m.col = clamp(m.col, 0, len(m.value[m.row]))
    		m.splitLine(m.row, m.col)
    		m.updateWordCount()
    	case key.Matches(msg, m.KeyMap.LineEnd):
    		m.CursorEnd()
    	case key.Matches(msg, m.KeyMap.LineStart):
    		m.CursorStart()
    	case key.Matches(msg, m.KeyMap.CharacterForward):
    		m.characterRight()
    	case key.Matches(msg, m.KeyMap.LineNext):
    		m.CursorDown()
    	case key.Matches(msg, m.KeyMap.WordForward):
    		m.wordRight()
    	case key.Matches(msg, m.KeyMap.Paste):
    		return m, Paste
    	case key.Matches(msg, m.KeyMap.CharacterBackward):
    		m.characterLeft(false /* insideLine */)
    	case key.Matches(msg, m.KeyMap.LinePrevious):
    		m.CursorUp()
    	case key.Matches(msg, m.KeyMap.WordBackward):
    		m.wordLeft()
    	case key.Matches(msg, m.KeyMap.InputBegin):
    		m.MoveToBegin()
    	case key.Matches(msg, m.KeyMap.InputEnd):
    		m.MoveToEnd()
    	case key.Matches(msg, m.KeyMap.PageUp):
    		m.PageUp()
    	case key.Matches(msg, m.KeyMap.PageDown):
    		m.PageDown()
    	case key.Matches(msg, m.KeyMap.LowercaseWordForward):
    		m.lowercaseRight()
    	case key.Matches(msg, m.KeyMap.UppercaseWordForward):
    		m.uppercaseRight()
    	case key.Matches(msg, m.KeyMap.CapitalizeWordForward):
    		m.capitalizeRight()
    	case key.Matches(msg, m.KeyMap.TransposeCharacterBackward):
    		m.transposeLeft()

    	default:
    		m.insertRunesFromUserInput([]rune(msg.Text))
    		m.updateWordCount()
    	}

    case pasteMsg:
    	m.insertRunesFromUserInput([]rune(msg))
    	m.updateWordCount()

    case pasteErrMsg:
    	m.Err = msg
    }

    // Make sure we set the content of the viewport before updating it.
    view := m.view()
    m.viewport.SetContent(view)
    vp, cmd := m.viewport.Update(msg)
    m.viewport = &vp
    cmds = append(cmds, cmd)

    if m.useVirtualCursor {
    	m.virtualCursor, cmd = m.virtualCursor.Update(msg)

    	// If the cursor has moved, reset the blink state. This is a small UX
    	// nuance that makes cursor movement obvious and feel snappy.
    	newRow, newCol := m.cursorLineNumber(), m.col
    	if (newRow != oldRow || newCol != oldCol) && m.virtualCursor.Mode() == cursor.CursorBlink {
    		m.virtualCursor.IsBlinked = false
    		cmd = m.virtualCursor.Blink()
    	}
    	cmds = append(cmds, cmd)
    }

    m.repositionView()

    return m, tea.Batch(cmds...)

}

func (m \*Model) view() string {
if len(m.Value()) == 0 && m.row == 0 && m.col == 0 && m.Placeholder != "" {
return m.placeholderView()
}
m.virtualCursor.TextStyle = m.activeStyle().computedCursorLine()

    var (
    	s                strings.Builder
    	style            lipgloss.Style
    	newLines         int
    	widestLineNumber int
    	lineInfo         = m.LineInfo()
    	styles           = m.activeStyle()
    )

    displayLine := 0
    for l, line := range m.value {
    	wrappedLines := m.memoizedWrap(line, m.width)

    	if m.row == l {
    		style = styles.computedCursorLine()
    	} else {
    		style = styles.computedText()
    	}

    	for wl, wrappedLine := range wrappedLines {
    		prompt := m.promptView(displayLine)
    		prompt = styles.computedPrompt().Render(prompt)
    		s.WriteString(style.Render(prompt))
    		displayLine++

    		var ln string
    		if m.ShowLineNumbers {
    			if wl == 0 { // normal line
    				isCursorLine := m.row == l
    				s.WriteString(m.lineNumberView(l+1, isCursorLine))
    			} else { // soft wrapped line
    				isCursorLine := m.row == l
    				s.WriteString(m.lineNumberView(-1, isCursorLine))
    			}
    		}

    		// Note the widest line number for padding purposes later.
    		lnw := uniseg.StringWidth(ln)
    		if lnw > widestLineNumber {
    			widestLineNumber = lnw
    		}

    		strwidth := uniseg.StringWidth(string(wrappedLine))
    		padding := m.width - strwidth
    		// If the trailing space causes the line to be wider than the
    		// width, we should not draw it to the screen since it will result
    		// in an extra space at the end of the line which can look off when
    		// the cursor line is showing.
    		if strwidth > m.width {
    			// The character causing the line to be wider than the width is
    			// guaranteed to be a space since any other character would
    			// have been wrapped.
    			wrappedLine = []rune(strings.TrimSuffix(string(wrappedLine), " "))
    			padding -= m.width - strwidth
    		}
    		if m.row == l && lineInfo.RowOffset == wl {
    			s.WriteString(style.Render(string(wrappedLine[:lineInfo.ColumnOffset])))
    			if m.col >= len(line) && lineInfo.CharOffset >= m.width {
    				m.virtualCursor.SetChar(" ")
    				s.WriteString(m.virtualCursor.View())
    			} else {
    				m.virtualCursor.SetChar(string(wrappedLine[lineInfo.ColumnOffset]))
    				s.WriteString(style.Render(m.virtualCursor.View()))
    				s.WriteString(style.Render(string(wrappedLine[lineInfo.ColumnOffset+1:])))
    			}
    		} else {
    			s.WriteString(style.Render(string(wrappedLine)))
    		}
    		s.WriteString(style.Render(strings.Repeat(" ", max(0, padding))))
    		s.WriteRune('\n')
    		newLines++
    	}
    }

    // Always show at least m.Height lines at all times.
    // To do this we can simply pad out a few extra new lines in the view.
    for range m.height {
    	s.WriteString(m.promptView(displayLine))
    	displayLine++

    	// Write end of buffer content
    	leftGutter := string(m.EndOfBufferCharacter)
    	rightGapWidth := m.Width() - uniseg.StringWidth(leftGutter) + widestLineNumber
    	rightGap := strings.Repeat(" ", max(0, rightGapWidth))
    	s.WriteString(styles.computedEndOfBuffer().Render(leftGutter + rightGap))
    	s.WriteRune('\n')
    }

    return s.String()

}

// View renders the text area in its current state.
func (m Model) View() string {
// XXX: This is a workaround for the case where the viewport hasn't
// been initialized yet like during the initial render. In that case,
// we need to render the view again because Update hasn't been called
// yet to set the content of the viewport.
m.viewport.SetContent(m.view())
view := m.viewport.View()
styles := m.activeStyle()
textareaView := styles.Base.Render(view)

    // Append statusbar if it exists
    if m.Statusbar != nil {
    	return textareaView + "\n" + m.Statusbar.Render()
    }
    return textareaView

}

// promptView renders a single line of the prompt.
func (m Model) promptView(displayLine int) (prompt string) {
prompt = m.Prompt
if m.promptFunc == nil {
return prompt
}
prompt = m.promptFunc(PromptInfo{
LineNumber: displayLine,
Focused: m.focus,
})
width := lipgloss.Width(prompt)
if width < m.promptWidth {
prompt = fmt.Sprintf("%\*s%s", m.promptWidth-width, "", prompt)
}

    return m.activeStyle().computedPrompt().Render(prompt)

}

// lineNumberView renders the line number.
//
// If the argument is less than 0, a space styled as a line number is returned
// instead. Such cases are used for soft-wrapped lines.
//
// The second argument indicates whether this line number is for a 'cursorline'
// line number.
func (m Model) lineNumberView(n int, isCursorLine bool) (str string) {
if !m.ShowLineNumbers {
return ""
}

    if n <= 0 {
    	str = " "
    } else {
    	str = strconv.Itoa(n)
    }

    // XXX: is textStyle really necessary here?
    textStyle := m.activeStyle().computedText()
    lineNumberStyle := m.activeStyle().computedLineNumber()
    if isCursorLine {
    	textStyle = m.activeStyle().computedCursorLine()
    	lineNumberStyle = m.activeStyle().computedCursorLineNumber()
    }

    // Format line number dynamically based on the maximum number of lines.
    digits := len(strconv.Itoa(m.MaxHeight))
    str = fmt.Sprintf(" %*v ", digits, str)

    return textStyle.Render(lineNumberStyle.Render(str))

}

// placeholderView returns the prompt and placeholder, if any.
func (m Model) placeholderView() string {
var (
s strings.Builder
p = m.Placeholder
styles = m.activeStyle()
)
// word wrap lines
pwordwrap := ansi.Wordwrap(p, m.width, "")
// hard wrap lines (handles lines that could not be word wrapped)
pwrap := ansi.Hardwrap(pwordwrap, m.width, true)
// split string by new lines
plines := strings.Split(strings.TrimSpace(pwrap), "\n")

    for i := range m.height {
    	isLineNumber := len(plines) > i

    	lineStyle := styles.computedPlaceholder()
    	if len(plines) > i {
    		lineStyle = styles.computedCursorLine()
    	}

    	// render prompt
    	prompt := m.promptView(i)
    	prompt = styles.computedPrompt().Render(prompt)
    	s.WriteString(lineStyle.Render(prompt))

    	// when show line numbers enabled:
    	// - render line number for only the cursor line
    	// - indent other placeholder lines
    	// this is consistent with vim with line numbers enabled
    	if m.ShowLineNumbers {
    		var ln int

    		switch {
    		case i == 0:
    			ln = i + 1
    			fallthrough
    		case len(plines) > i:
    			s.WriteString(m.lineNumberView(ln, isLineNumber))
    		default:
    		}
    	}

    	switch {
    	// first line
    	case i == 0:
    		// first character of first line as cursor with character
    		m.virtualCursor.TextStyle = styles.computedPlaceholder()

    		ch, rest, _, _ := uniseg.FirstGraphemeClusterInString(plines[0], 0)
    		m.virtualCursor.SetChar(ch)
    		s.WriteString(lineStyle.Render(m.virtualCursor.View()))

    		// the rest of the first line
    		s.WriteString(lineStyle.Render(styles.computedPlaceholder().Render(rest)))

    		// extend the first line with spaces to fill the width, so that
    		// the entire line is filled when cursorline is enabled.
    		gap := strings.Repeat(" ", max(0, m.width-lipgloss.Width(plines[0])))
    		s.WriteString(lineStyle.Render(gap))
    	// remaining lines
    	case len(plines) > i:
    		// current line placeholder text
    		if len(plines) > i {
    			placeholderLine := plines[i]
    			gap := strings.Repeat(" ", max(0, m.width-uniseg.StringWidth(plines[i])))
    			s.WriteString(lineStyle.Render(placeholderLine + gap))
    		}
    	default:
    		// end of line buffer character
    		eob := styles.computedEndOfBuffer().Render(string(m.EndOfBufferCharacter))
    		s.WriteString(eob)
    	}

    	// terminate with new line
    	s.WriteRune('\n')
    }

    m.viewport.SetContent(s.String())
    return styles.Base.Render(m.viewport.View())

}

// Blink returns the blink command for the virtual cursor.
func Blink() tea.Msg {
return cursor.Blink()
}

// Cursor returns a [tea.Cursor] for rendering a real cursor in a Bubble Tea
// program. This requires that [Model.VirtualCursor] is set to false.
//
// Note that you will almost certainly also need to adjust the offset cursor
// position per the textarea's per the textarea's position in the terminal.
//
// Example:
//
// // In your top-level View function:
// f := tea.NewFrame(m.textarea.View())
// f.Cursor = m.textarea.Cursor()
// f.Cursor.Position.X += offsetX
// f.Cursor.Position.Y += offsetY
func (m Model) Cursor() \*tea.Cursor {
if m.useVirtualCursor || !m.Focused() {
return nil
}

    lineInfo := m.LineInfo()
    w := lipgloss.Width
    baseStyle := m.activeStyle().Base

    xOffset := lineInfo.CharOffset +
    	w(m.promptView(0)) +
    	w(m.lineNumberView(0, false)) +
    	baseStyle.GetMarginLeft() +
    	baseStyle.GetPaddingLeft() +
    	baseStyle.GetBorderLeftSize()

    yOffset := m.cursorLineNumber() -
    	m.viewport.YOffset() +
    	baseStyle.GetMarginTop() +
    	baseStyle.GetPaddingTop() +
    	baseStyle.GetBorderTopSize()

    c := tea.NewCursor(xOffset, yOffset)
    c.Blink = m.styles.Cursor.Blink
    c.Color = m.styles.Cursor.Color
    c.Shape = m.styles.Cursor.Shape
    return c

}

func (m Model) memoizedWrap(runes []rune, width int) [][]rune {
input := line{runes: runes, width: width}
if v, ok := m.cache.Get(input); ok {
return v
}
v := wrap(runes, width)
m.cache.Set(input, v)
return v
}

// cursorLineNumber returns the line number that the cursor is on.
// This accounts for soft wrapped lines.
func (m Model) cursorLineNumber() int {
line := 0
for i := range m.row {
// Calculate the number of lines that the current line will be split
// into.
line += len(m.memoizedWrap(m.value[i], m.width))
}
line += m.LineInfo().RowOffset
return line
}

// mergeLineBelow merges the current line the cursor is on with the line below.
func (m \*Model) mergeLineBelow(row int) {
if row >= len(m.value)-1 {
return
}

    // To perform a merge, we will need to combine the two lines and then
    m.value[row] = append(m.value[row], m.value[row+1]...)

    // Shift all lines up by one
    for i := row + 1; i < len(m.value)-1; i++ {
    	m.value[i] = m.value[i+1]
    }

    // And, remove the last line
    if len(m.value) > 0 {
    	m.value = m.value[:len(m.value)-1]
    }

}

// mergeLineAbove merges the current line the cursor is on with the line above.
func (m \*Model) mergeLineAbove(row int) {
if row <= 0 {
return
}

    m.col = len(m.value[row-1])
    m.row = m.row - 1

    // To perform a merge, we will need to combine the two lines and then
    m.value[row-1] = append(m.value[row-1], m.value[row]...)

    // Shift all lines up by one
    for i := row; i < len(m.value)-1; i++ {
    	m.value[i] = m.value[i+1]
    }

    // And, remove the last line
    if len(m.value) > 0 {
    	m.value = m.value[:len(m.value)-1]
    }

}

func (m \*Model) splitLine(row, col int) {
// To perform a split, take the current line and keep the content before
// the cursor, take the content after the cursor and make it the content of
// the line underneath, and shift the remaining lines down by one
head, tailSrc := m.value[row][:col], m.value[row][col:]
tail := make([]rune, len(tailSrc))
copy(tail, tailSrc)

    m.value = append(m.value[:row+1], m.value[row:]...)

    m.value[row] = head
    m.value[row+1] = tail

    m.col = 0
    m.row++

}

// Paste is a command for pasting from the clipboard into the text input.
func Paste() tea.Msg {
str, err := clipboard.ReadAll()
if err != nil {
return pasteErrMsg{err}
}
return pasteMsg(str)
}

func wrap(runes []rune, width int) [][]rune {
var (
lines = [][]rune{{}}
word = []rune{}
row int
spaces int
)

    // Word wrap the runes
    for _, r := range runes {
    	if unicode.IsSpace(r) {
    		spaces++
    	} else {
    		word = append(word, r)
    	}

    	if spaces > 0 { //nolint:nestif
    		if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces > width {
    			row++
    			lines = append(lines, []rune{})
    			lines[row] = append(lines[row], word...)
    			lines[row] = append(lines[row], repeatSpaces(spaces)...)
    			spaces = 0
    			word = nil
    		} else {
    			lines[row] = append(lines[row], word...)
    			lines[row] = append(lines[row], repeatSpaces(spaces)...)
    			spaces = 0
    			word = nil
    		}
    	} else {
    		// If the last character is a double-width rune, then we may not be able to add it to this line
    		// as it might cause us to go past the width.
    		lastCharLen := rw.RuneWidth(word[len(word)-1])
    		if uniseg.StringWidth(string(word))+lastCharLen > width {
    			// If the current line has any content, let's move to the next
    			// line because the current word fills up the entire line.
    			if len(lines[row]) > 0 {
    				row++
    				lines = append(lines, []rune{})
    			}
    			lines[row] = append(lines[row], word...)
    			word = nil
    		}
    	}
    }

    if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces >= width {
    	lines = append(lines, []rune{})
    	lines[row+1] = append(lines[row+1], word...)
    	// We add an extra space at the end of the line to account for the
    	// trailing space at the end of the previous soft-wrapped lines so that
    	// behaviour when navigating is consistent and so that we don't need to
    	// continually add edges to handle the last line of the wrapped input.
    	spaces++
    	lines[row+1] = append(lines[row+1], repeatSpaces(spaces)...)
    } else {
    	lines[row] = append(lines[row], word...)
    	spaces++
    	lines[row] = append(lines[row], repeatSpaces(spaces)...)
    }

    return lines

}

func repeatSpaces(n int) []rune {
return []rune(strings.Repeat(string(' '), n))
}

// numDigits returns the number of digits in an integer.
func numDigits(n int) int {
if n == 0 {
return 1
}
count := 0
num := abs(n)
for num > 0 {
count++
num /= 10
}
return count
}

func clamp(v, low, high int) int {
if high < low {
low, high = high, low
}
return min(high, max(low, v))
}

func abs(n int) int {
if n < 0 {
return -n
}
return n
}

this is now the textarea

I can now be much more concrete, because this shows exactly where edits happen. The main takeaway is that your textarea already has the right mutation points — you do not need to sprinkle event logic everywhere, just wrap the true buffer-changing operations.

Thought for a few seconds
Yes — and this is actually very helpful, because now we can see the real mutation surface.

Your textarea is already close to a good architecture for event tracking, because almost all document changes funnel through a small set of functions.

The key realization is:

your event system should track buffer mutations, not keypresses.

So you do not want “user pressed backspace” as the primary event.
You want:

what range changed

what text was there before

what text is there after

cursor before / after

time

That means your textarea should expose a few canonical edit operations, and those operations should emit events.

What your textarea already gives you
Your mutable document state is basically:

value [][]rune
row int
col int
This means the actual text is stored as lines, but for event tracking you should think of it as one logical document:

line 0

newline

line 1

newline

...

That is important because your event system should probably use global rune offsets, not (row, col) alone.

Your current mutation points are mainly these:

insertRunesFromUserInput

deleteBeforeCursor

deleteAfterCursor

deleteWordLeft

deleteWordRight

mergeLineBelow

mergeLineAbove

splitLine

direct edits in Update, like backspace/delete-char

case transforms:

lowercaseRight

uppercaseRight

capitalizeRight

transposeLeft

These are the places where content actually changes.

The most important thing you need next
You need a single edit emission layer between textarea mutation and event manager.

Not this:

every keybinding manually creates events

But this:

every document mutation calls one canonical function like:

applyReplace(start, end, after)
which:

computes before

mutates m.value

updates cursor

emits an event to your EventMgr

That will keep the system sane.

What you probably need to add

1. A canonical “replace span” operation
   Right now your mutations are distributed across many methods and directly edit m.value.

That is the first thing I would change.

Add one internal primitive such as:

func (m \*Model) replaceRange(start, end DocPos, after []rune)
or even better, if you can convert to flat offsets:

func (m \*Model) replaceRangeByOffset(start, end int, after []rune)
This function should be the only place that actually changes text.

Then:

insert becomes replace empty span with text

backspace becomes replace one-rune span with empty

newline becomes replace empty span with \n

merge lines becomes replace newline span with empty

word delete becomes replace span with empty

paste becomes replace span with pasted text

capitalization becomes replace span with transformed text

This is a huge simplification.

2. A way to convert (row, col) to global offset
   This is probably the next concrete thing you need.

Something like:

func (m \*Model) Offset(row, col int) int
Meaning: how many runes from the start of the document is this cursor location?

You will use that for event recording.

For example:

row 0, col 3 in "abc\ndef" => offset 3

row 1, col 2 => offset 6 (abc\n is 4, then 2 more)

This is essential because your event manager wants document-relative positions, not line-relative ones.

You may also want the inverse:

func (m \*Model) PosFromOffset(offset int) (row, col int)
That helps replay and debugging.

3. A document position type
   Instead of using naked ints everywhere, a small struct helps:

type DocPos struct {
Row int
Col int
}
Even if your event log stores flat offsets, the textarea implementation can use DocPos internally.

4. An event sink in the textarea model
   Your textarea should probably hold a reference or callback for event recording.

For example:

type EditRecorder interface {
Record(e events.Event)
}
and in Model:

Recorder EditRecorder
or:

OnEdit func(e events.Event)
That keeps textarea independent from the concrete event manager package.

This is much better than importing your event manager logic directly into every textarea method.

5. Cursor before / after capture
   For reflection and replay, you likely want:

cursor before edit

cursor after edit

So before a mutation, capture:

beforeRow, beforeCol := m.row, m.col
beforeOffset := m.Offset(beforeRow, beforeCol)
After the mutation:

afterRow, afterCol := m.row, m.col
afterOffset := m.Offset(afterRow, afterCol)
Store both in the event.

6. Stop modeling events as just Write/Delete
   Now that we see the textarea code, this becomes even clearer.

You have many mutations that are neither purely write nor delete:

split line

merge lines

transpose

uppercase/lowercase/capitalize

paste-over-selection if you add selection later

replace due to autocomplete or programmatic edits

So I would move even more strongly toward:

type Event struct {
ID int64
CreatedAt time.Time
UpdatedAt time.Time
StartOffset int
EndOffset int
Before []rune
After []rune
CursorBefore int
CursorAfter int
Source EventSource
}
This will fit your actual editor much better.

Where your current code should emit events
You do not need to track every helper function directly.
You only need to track the functions that semantically edit the document.

A. insertRunesFromUserInput
This is a major one.

This should emit one replace event:

start = cursor offset before insert

end = same as start

before = empty

after = inserted runes

Even if the inserted text has newlines, it is still one replace event.

That is good.

B. Backspace / delete-char branches in Update
Right now these directly mutate slices.

Those are very important to route through your canonical replace function.

For example, backspace should become:

replace [offset-1 : offset] with ""

Delete forward:

replace [offset : offset+1] with ""

If at line boundary and merging lines, then delete the implicit newline between lines.

C. splitLine
This should really just be:

replace current cursor position with "\n"

That is conceptually what Enter does.

So instead of a special split primitive for event tracking, model it as newline insertion.

D. mergeLineAbove / mergeLineBelow
These are really deletion of a newline.

So event-wise:

replace [newlineOffset : newlineOffset+1] with ""

This is much cleaner than thinking of it as “merge lines.”

E. deleteBeforeCursor, deleteAfterCursor, deleteWordLeft, deleteWordRight
All of these are span deletions.

They should compute the span, then call replace with empty after-text.

F. transposeLeft, uppercaseRight, lowercaseRight, capitalizeRight
These are replacements, not writes/deletes.

For example transpose of "ab" to "ba" is:

replace [i-1 : i+1] from "ab" to "ba"

Case changes similarly replace a word span with transformed text.

If you want your edit history to reflect real thought evolution, these may or may not matter semantically later, but structurally they are still edits.

You can also mark them as a different Source so higher layers can ignore them if desired.

What is still missing architecturally

1. Selection support
   I do not see selection support here.

That is fine for now, but once you add selection, many operations become natural replaces over a range.

So designing around replace now will save you pain later.

2. Undo/redo
   Your event model is going to be much stronger if it can eventually support undo/redo.

That becomes much easier if every edit is a replace event with before and after.

Because undo is just the inverse replace.

3. Snapshots
   If you keep an event log, eventually store snapshots every N edits or every few seconds.

Otherwise reconstruction gets annoying.

4. Separation between raw editor events and semantic events
   Do not confuse:

raw edit events

later semantic segments / reflection units

Your textarea should only produce raw edit events.

The concrete problem in your current textarea
Right now your code mutates m.value in many separate places.

That means if you try to bolt event tracking on top naively, you will end up with:

duplicated event creation logic

inconsistent spans

bugs around line merges/splits

some edits forgotten entirely

So the biggest thing you need is:

centralize all document mutations through one edit primitive

That is the real next step.

What I would add to Model
Something like:

type EditRecorder interface {
RecordEdit(EditEvent)
}

type EditEvent struct {
Time time.Time
StartOffset int
EndOffset int
Before []rune
After []rune
CursorBefore int
CursorAfter int
Source string
}

type Model struct {
...
Recorder EditRecorder
}
Then add helpers:

func (m *Model) docOffset(row, col int) int
func (m *Model) selectedTextByOffset(start, end int) []rune
func (m \*Model) replaceRangeByOffset(start, end int, after []rune, source string)
And then gradually rewrite mutations to go through that.

Important design decision: keep [][]rune or flatten?
You currently store text as [][]rune, which is good for rendering and line ops.

That is okay. You do not need to switch internal storage yet.

But for event tracking, I would still define all spans in terms of the logical flat document.

So:

storage can remain [][]rune

event coordinates should be flat offsets

That is a good compromise.

What you probably need immediately, in priority order
docOffset(row, col)

one canonical replace/edit function

an edit recorder callback/interface

convert insert/delete/newline/merge to use replace

event type based on Before/After, not Write/Delete

later: snapshots, undo/redo, selection

The simplest conceptual mapping for your current editor
Here is the clean mental model:

typing x → replace empty span with "x"

enter → replace empty span with "\n"

backspace → replace previous rune span with ""

delete → replace next rune span with ""

ctrl-u / ctrl-k → replace prefix/suffix span with ""

word delete → replace word span with ""

paste → replace empty span with pasted text

transpose → replace two-rune span with swapped runes

lowercase/uppercase/capitalize → replace word span with transformed version

Everything becomes replace.

That is why I keep pushing that model.

Bottom line
From this textarea code, the main thing you still need is not more data structures first.

You need to refactor the editor so every content mutation flows through a single edit primitive.

That primitive should:

compute flat offsets

know before/after text

mutate buffer

emit event

Once you have that, your event tracking system becomes solid.

Without that, your event tracking will become scattered and fragile.

I can next help you design the exact replaceRangeByOffset / docOffset layer for this textarea, or sketch how to retrofit your current methods one by one.

the main concern i have it, after these changes, willl typing-showing on textarea still be fast ?

Thought for a couple of seconds
Yes, it can still stay fast — if you keep the event layer thin and avoid doing expensive work on every keystroke.

The key is:

recording an edit event is cheap; recomputing meaning is expensive.

So the event system itself usually will not be what makes typing lag. The danger is when you accidentally attach heavyweight work to each edit.

What is cheap enough per keystroke
These are usually fine on every keystroke:

capture cursor before/after

compute a small span

copy a few runes

append one event to a slice

maybe merge into a pending event

increment a revision counter

That is tiny compared to everything your textarea is already doing for:

wrapping

rendering

viewport updates

string width calculations

blinking cursor

status bar updates

So if your event layer is just “record what changed,” typing should remain fast.

What will make it slow
These are the things that would hurt responsiveness if done inline on each keypress:

generating embeddings

LLM calls

full semantic segmentation

writing to SQLite synchronously every key

recomputing reflections immediately

flattening the whole document repeatedly

taking full snapshots every keypress

large allocations or full-buffer copies on every edit

That is what you want to avoid.

The safe rule
Your textarea update path should only do:

mutate buffer

record small edit event

maybe merge with pending event

return

Everything else should happen later.

The biggest performance risk in your case
Given your current textarea, the risk is probably not the event object itself.

The bigger risks are:

1. Repeated full-document conversions
   If every keystroke does something like:

turn [][]rune into one giant []rune

or call Value() to rebuild the entire string

or scan the whole document to compute offsets

then that can become noticeable as the document gets large.

2. Full-buffer copying
   If your replace/event layer copies the entire textarea content for every edit, that will be bad.

3. Synchronous disk writes
   If every keypress writes event logs to SQLite or disk immediately, that can introduce latency.

How to keep it fast

1. Keep event recording in memory first
   Do not write every event to disk immediately.

Instead:

append to in-memory pending event log

flush periodically:

after idle

every N events

on blur

on save

on exit

That preserves responsiveness.

2. Coalesce keystrokes
   Do not store one event per character if you do not need to.

For example, typing hello can become one pending insert event.

This helps both performance and later analysis.

3. Do not trigger AI from the hot path
   Never do this inside the core textarea mutation path:

embedding

entity extraction

summarization

reflection retrieval

Instead, fire a deferred signal like:

“document changed”

“dirty since revision X”

and let another layer handle it after debounce/idle.

4. Avoid full flattening on every keypress
   If you need global offsets, try not to rebuild the whole document every time.

You have options:

Option A: compute offset by summing preceding line lengths
That is often good enough at first.

If your cursor is at (row, col), then:

sum lengths of lines 0..row-1

add one newline per preceding line

add col

This is O(number of lines before cursor).

For many diary cases, that is acceptable initially.

Option B: maintain a line prefix-sum structure
If performance becomes an issue, maintain something like:

per-line rune counts

prefix sums

Then offset lookup becomes much cheaper.

You probably do not need this immediately.

5. Only copy changed text, not whole text
   For an event, store only:

replaced span

before text

after text

That is small.

Do not store the full document snapshot per keystroke.

6. Snapshot rarely
   Snapshots are useful, but do them:

every 100 edits

every few seconds while idle

on save

not every keypress.

Practical answer for your architecture
If you implement it like this, typing should remain fast:

Hot path
keypress arrives

textarea mutates local buffer

event manager records lightweight replace event

maybe merge with pending event

redraw

Warm path, debounced
after 500ms–2s idle:

flush pending event group

maybe persist to disk

maybe update derived segments

Cold path, background
generate embeddings

update search index

recompute summaries

reflection candidates

That separation is what keeps the UI snappy.

What I would avoid in version 1
I would not start with:

immediate SQLite write per key

full CRDT-style bookkeeping

full graph updates live while typing

resegmenting whole document every key

storing every cursor move as an event

That is where complexity and lag creep in.

A useful mental budget
Per keystroke, you want work to be roughly:

local

bounded

proportional to the edit size, not the whole document size

Good:

insert 1 rune

delete 1 rune

replace 5 runes

append/merge one event

Bad:

rescan 20,000 characters

rebuild multiple full strings

sync database write

call model

A subtle point: your current textarea already does nontrivial work
Even now, your textarea already performs work that is more expensive than event logging in many cases:

updateWordCount() calls Value()

Value() rebuilds the full string

strings.Fields(text) scans the whole text

wrapping/view rendering scans lots of text

So honestly, your event system may not be the first bottleneck.

For example, this:

text := m.Value()
chars := len([]rune(text))
words := len(strings.Fields(text))
on every edit is already a whole-document operation.

That is more suspicious performance-wise than appending a tiny edit event.

So if you notice lag later, I would inspect things like:

full text reconstruction in updateWordCount

repeated wrapping

repeated width computations

before blaming the event log.

My honest estimate
For a diary-sized TUI, a lightweight edit-event layer will almost certainly be fine.

You are much more likely to hit lag from:

rendering

repeated full string conversion

status/count recomputation

future AI hooks

than from storing a replace event.

Best compromise
Start with the clean architecture even if it is slightly less optimal than the absolute fastest possible code.

Then profile.

If something becomes slow, optimize the real bottleneck.

Most likely progression:

naive event recording is fine

large-document word count becomes noticeable

full-document offset scans become noticeable

then you add counters / prefix sums / buffered persistence

That is a good path.

If you want, I can point out the specific places in this textarea that are more likely than event tracking to cause typing slowdown.
