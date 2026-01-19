grextor/
├── README.md
├── docker-compose.yml        # Qdrant + Neo4j
├── cmd/
│   ├── ingest/               # index docs into vector + graph
│   └── query/                # semantic + graph-constrained search
├── internal/
│   ├── embed/                # embedding interface
│   ├── vector/               # Qdrant client
│   ├── graph/                # Neo4j client
│   └── engine/               # Grextor core logic
├── data/
│   └── sample_docs/
└── examples/
