# MangaHub Report - Step-by-Step Content Filling Guide

This guide provides the exact content to fill in each placeholder in your template document.

## 📋 How to Use This Guide

1. Open `MangaHub_Report_TEMPLATE.docx` in Microsoft Word
2. Find each `[Section Name - Content to be added]` placeholder
3. Replace it with the content provided below for that section
4. Keep the formatting consistent (justified text, 12pt font)

---

## ABSTRACT

### [Abstract Paragraph 1]
```
MangaHub is a comprehensive network programming project developed as part of the IT096IU Network Programming course at International University, Vietnam National University - Ho Chi Minh City. The system demonstrates the practical implementation of five major network protocols: HTTP, TCP, UDP, WebSocket, and gRPC, integrated into a unified manga and comic tracking application.
```

### [Abstract Paragraph 2]
```
The project addresses the real-world challenge of tracking manga reading progress across multiple devices and platforms. Through a client-server architecture built with Go programming language, the system enables users to manage their manga library, synchronize reading progress in real-time, receive notifications about new chapter releases, and participate in community discussions.
```

### [Abstract Paragraph 3]
```
The HTTP REST API server handles user authentication, manga catalog management, and library operations. The TCP synchronization server ensures real-time progress updates across all connected devices using JSON-based messaging. UDP broadcasts deliver instant chapter release notifications to subscribed clients. WebSocket technology powers the real-time chat system for manga discussions, while gRPC facilitates efficient internal service-to-service communication.
```

### [Abstract Paragraph 4]
```
The system employs modern software engineering practices including clean architecture, concurrent programming with goroutines, JWT-based authentication, and comprehensive testing strategies. A Next.js web frontend provides an intuitive user interface, while a command-line tool offers advanced functionality for power users. This report presents a detailed examination of the project's design, implementation, testing, and deployment, demonstrating how multiple network protocols can be integrated cohesively to create a practical, scalable, and user-friendly application.
```

---

## CHAPTER 1: INTRODUCTION

### [Background Paragraph 1]
```
The digital consumption of manga and comics has experienced exponential growth in recent years, with millions of readers worldwide accessing content across multiple platforms and devices. However, this proliferation has created a significant challenge: readers often lose track of their progress when switching between devices, struggle to stay updated on new chapter releases, and lack effective tools for organizing their extensive reading lists.
```

### [Background Paragraph 2]
```
Network programming forms the backbone of modern distributed systems, enabling seamless communication between clients and servers across the internet. Understanding how different network protocols work together to create cohesive applications is essential for computer science students preparing for careers in software development, cloud computing, and systems engineering.
```

### [Background Paragraph 3]
```
MangaHub was conceived to address both challenges simultaneously. By building a practical manga tracking system, the project demonstrates how multiple network protocols can work together to solve real-world problems. The application requires HTTP for RESTful API operations, TCP for reliable real-time synchronization, UDP for broadcast notifications, WebSocket for bidirectional communication, and gRPC for efficient inter-service communication. The choice of Go (Golang) as the primary programming language was motivated by its excellent support for concurrent programming, built-in networking capabilities, and growing adoption in cloud-native applications.
```

### [Objectives Introduction]
```
The primary objective of this project is to demonstrate comprehensive understanding of network programming concepts through practical implementation of five major network protocols. Each protocol was carefully selected to address specific requirements of the manga tracking system while providing educational value in understanding their strengths, limitations, and appropriate use cases. The specific objectives include:
```

### [Objective 1]
```
Implement HTTP REST API: Design and develop a complete RESTful API server supporting user authentication, CRUD operations for manga management, library management, and reading progress tracking. The API follows REST architectural principles, implements proper HTTP status codes, and provides comprehensive error handling with standardized response formats.
```

### [Objective 2]
```
Build TCP Synchronization Server: Create a TCP server that maintains persistent connections with clients to enable real-time synchronization of reading progress across multiple devices. The implementation handles concurrent connections efficiently using goroutines and implements a robust JSON-based messaging protocol with acknowledgment mechanisms.
```

### [Objective 3]
```
Develop UDP Notification Service: Implement a UDP broadcast system for sending chapter release notifications to subscribed clients. The system handles client registration, maintains subscription lists, and efficiently broadcasts updates without guaranteeing delivery order, demonstrating UDP's best-effort delivery model.
```

### [Objective 4]
```
Create WebSocket Chat System: Build a real-time chat application using WebSocket protocol, supporting multiple chat rooms, user presence detection, and message broadcasting. The system handles connection upgrades from HTTP and manages room-based message routing efficiently.
```

### [Objective 5]
```
Implement gRPC Internal Service: Design and implement gRPC services for internal communication between system components, utilizing Protocol Buffers for efficient data serialization. The implementation demonstrates both unary and streaming RPC patterns, showcasing gRPC's type safety and performance benefits.
```

### [Scope Introduction]
```
The MangaHub project encompasses a comprehensive set of features designed to demonstrate network programming concepts while providing practical functionality. The scope has been carefully defined to balance educational objectives with project feasibility within the course timeline. The following sections detail what is included and explicitly excluded from the project.
```

### [In Scope Items]
```
The project includes five independent server processes implementing HTTP (port 8080), TCP (port 9090), UDP (port 9091), WebSocket (via HTTP upgrade), and gRPC (port 9092) protocols. User authentication and authorization is implemented using JWT tokens with bcrypt password hashing. The system provides a complete manga catalog with search and filtering capabilities, personal library management for organizing collections, and real-time progress synchronization across multiple client connections. Chapter release notifications are broadcast via UDP to subscribed clients. Multi-room chat functionality enables manga-specific discussions and a general community chat. The SQLite database includes a migration system for schema versioning. A web frontend built with Next.js, TypeScript, and Tailwind CSS provides an intuitive interface, while a command-line tool offers advanced functionality. Comprehensive API documentation, unit and integration testing suites, and Docker containerization complete the deliverables.
```

### [Out of Scope Items]
```
The project explicitly excludes actual manga reading functionality (viewing pages), payment processing or subscription management, social features beyond basic chat (friend systems, user profiles), mobile native applications, content recommendation algorithms, advanced analytics or reporting features, multi-language internationalization support, and production-grade scalability optimizations. These exclusions allow the team to focus on core network programming demonstrations while maintaining a realistic project scope.
```

### [Tools Introduction]
```
The MangaHub project leverages modern development tools and technologies selected for their suitability to network programming education and real-world applicability. Each tool was chosen based on its strengths in demonstrating specific concepts, community support, documentation quality, and industry relevance.
```

### [Go Language Description]
```
Go 1.19+ serves as the primary backend language, chosen for its exceptional concurrency model with goroutines and channels. Goroutines are lightweight threads managed by the Go runtime, enabling efficient concurrent handling of thousands of connections with minimal overhead. The language's built-in network programming support through the net package, fast compilation times, static typing for safety, and comprehensive standard library make it ideal for building network services. Go's garbage collection provides memory safety without sacrificing performance, while its simplicity and clear syntax enhance code readability and maintainability.
```

### [TypeScript Description]
```
TypeScript 5.0+ is used for frontend development, providing type safety and improved developer experience over JavaScript. The language's static typing catches errors at compile time, its interface system enables clear API contracts, and enum support improves code organization. TypeScript enhances large codebases' maintainability while remaining fully interoperable with the JavaScript ecosystem.
```

### [Gin Framework]
```
Gin is a high-performance HTTP web framework featuring middleware support for cross-cutting concerns, JSON validation and binding, routing groups for API versioning, and excellent documentation. Chosen for its speed (40x faster than Martini) and simplicity in building RESTful APIs, Gin provides a familiar Express.js-like API for developers while maintaining Go's performance characteristics.
```

### [Gorilla WebSocket]
```
Gorilla WebSocket is a battle-tested library providing full RFC 6455 compliance and extensive control over connection handling. It offers both high-level and low-level APIs, comprehensive test coverage, and active community maintenance. The library handles complex WebSocket frame parsing, masking requirements, and connection management while exposing simple interfaces for application development.
```

### [gRPC-Go]
```
gRPC-Go is the official Go implementation of gRPC, providing high-performance RPC framework with Protocol Buffer support. It enables both unary and streaming RPC patterns, automatic code generation from proto files, built-in load balancing and retry mechanisms, and strong typing across language boundaries. The framework's HTTP/2 transport provides multiplexing and header compression benefits.
```

### [Next.js]
```
Next.js 14 is a React framework providing server-side rendering for improved SEO and initial load performance, file-based routing for intuitive page organization, API routes for backend functionality within the same codebase, and optimized production builds with automatic code splitting. The framework's hybrid rendering approach enables both static and dynamic content strategies.
```

### [React]
```
React 18 is a component-based UI library featuring hooks for state management and lifecycle handling. useState enables local component state, useEffect handles side effects and subscriptions, useContext provides global state management, and custom hooks enable logic reusability. React's virtual DOM optimization ensures efficient updates while maintaining declarative programming model.
```

### [Tailwind CSS]
```
Tailwind CSS 3 is a utility-first framework enabling rapid UI development with consistent design tokens. It provides responsive design utilities, customizable theme system, JIT compiler for minimal CSS output, and component extraction through @apply directive. The framework's approach reduces CSS bloat while maintaining design consistency across the application.
```

### [SQLite3 Description]
```
SQLite3 is an embedded relational database chosen for its simplicity, zero-configuration setup, and suitability for educational projects. It provides full ACID compliance (Atomicity, Consistency, Isolation, Durability), supports complex SQL queries, requires no separate server process, and offers cross-platform compatibility. The file-based storage makes deployment and backup straightforward while maintaining full SQL feature support.
```

### [Git/GitHub]
```
Git provides distributed version control with branching for parallel development, commit history for change tracking, and merge capabilities for collaboration. GitHub adds pull request workflow for code review, issue tracking for bug management, GitHub Actions for CI/CD automation, and project boards for task organization. The team uses feature branches, pull requests for all changes, and maintains a linear commit history.
```

### [VS Code]
```
Visual Studio Code serves as the primary IDE with Go extension for language support including autocomplete, debugging, and test runner integration. ESLint and Prettier extensions ensure code quality and consistent formatting. The integrated terminal and Git support streamline development workflow.
```

### [Docker]
```
Docker provides containerization for packaging applications with their dependencies, ensuring consistent development and production environments. Multi-stage builds optimize image sizes, docker-compose orchestrates multiple services, and volume mounting enables live development. The containerization approach simplifies deployment and scaling strategies.
```

---

## CHAPTER 2: LITERATURE REVIEW

### [Network Programming Introduction]
```
Network programming forms the foundation of distributed systems, enabling communication between processes running on different machines across local or wide area networks. At its core, network programming involves creating applications that can send and receive data using various protocols and communication patterns. Understanding these fundamentals is essential for building robust, scalable network services.
```

### [OSI Model]
```
The OSI (Open Systems Interconnection) model provides a conceptual framework for understanding network communications through seven distinct layers. Layer 1 (Physical) handles the physical transmission of bits over media. Layer 2 (Data Link) manages error detection and correction within a physical link. Layer 3 (Network) routes packets across networks using IP addresses. Layer 4 (Transport) provides end-to-end communication with TCP or UDP. Layer 5 (Session) establishes and manages connections. Layer 6 (Presentation) handles data format translation and encryption. Layer 7 (Application) interfaces directly with applications through protocols like HTTP, FTP, and SMTP.
```

### [TCP/IP Model]
```
While the OSI model serves as a theoretical reference, practical implementations typically follow the simpler four-layer TCP/IP model. The Network Access layer combines OSI layers 1-2, handling physical transmission. The Internet layer (OSI layer 3) routes packets using IP. The Transport layer (OSI layer 4) provides reliable (TCP) or unreliable (UDP) data transfer. The Application layer (OSI layers 5-7) contains application-specific protocols. This model more accurately reflects real-world network implementations.
```

### [Socket Programming]
```
Socket programming represents the primary mechanism for network communication in modern systems. A socket acts as an endpoint for sending or receiving data across a network, identified by an IP address and port number. The Berkeley sockets API, introduced in BSD Unix, has become the de facto standard for network programming across operating systems. Go's net package provides a high-level abstraction over these low-level socket operations while maintaining efficiency and ease of use.
```

### [Concurrency in Network Programming]
```
Concurrent programming is essential for building scalable network services that handle multiple simultaneous connections. Traditional approaches using threads face scalability challenges due to context switching overhead and memory consumption (typically 1-2MB per thread). Go addresses these limitations through goroutines, lightweight threads managed by the Go runtime with just 2KB initial stack size. Goroutines enable efficient concurrent handling of thousands of connections through cooperative scheduling and efficient channel-based communication.
```

### [TCP/IP Introduction]
```
The TCP/IP protocol suite provides the foundation for internet communication, offering both reliable (TCP) and unreliable (UDP) transport mechanisms. Understanding when to use each protocol is crucial for building efficient network applications that balance reliability, latency, and resource usage.
```

### [TCP Protocol Details]
```
The Transmission Control Protocol (TCP) provides reliable, ordered, and error-checked delivery of data between applications. TCP establishes connections through a three-way handshake: client sends SYN, server responds with SYN-ACK, client confirms with ACK. This ensures both parties are ready before data transmission begins. TCP implements flow control using a sliding window mechanism, preventing faster senders from overwhelming slower receivers. The protocol includes congestion control algorithms (Slow Start, Congestion Avoidance, Fast Retransmit, Fast Recovery) to optimize network utilization while preventing collapse. Each byte is sequenced and acknowledged, enabling retransmission of lost packets. Connection termination uses a four-way handshake to gracefully close both communication channels.
```

### [UDP Protocol Details]
```
The User Datagram Protocol (UDP) offers a connectionless, best-effort delivery service. Unlike TCP, UDP does not establish connections or guarantee delivery, making it faster but less reliable. UDP packets (datagrams) may arrive out of order, be duplicated, or be lost entirely without notification. Each datagram is independent, containing source and destination ports, length, and checksum. These characteristics make UDP suitable for applications that can tolerate some data loss but require low latency, such as real-time video streaming, online gaming, DNS queries, and broadcast notifications.
```

### [IP Layer]
```
The Internet Protocol (IP) layer handles routing packets across networks using IP addresses. IPv4 uses 32-bit addresses (e.g., 192.168.1.1) supporting approximately 4.3 billion unique addresses. IPv6 uses 128-bit addresses providing virtually unlimited address space while improving routing efficiency and security. Port numbers (16-bit, range 0-65535) distinguish multiple services on a single host: well-known ports (0-1023) for standard services, registered ports (1024-49151) for specific applications, and dynamic ports (49152-65535) for client-side communication.
```

### [Protocol Comparison]
```
TCP provides reliability guarantees at the cost of additional overhead and latency. It is ideal for applications requiring guaranteed delivery such as file transfers, email, web browsing, and database synchronization. UDP offers minimal overhead and lower latency but no delivery guarantees. It suits real-time applications, broadcast/multicast scenarios, and stateless request-response protocols. In MangaHub, TCP serves the synchronization server where reliable, ordered delivery is critical for maintaining consistent reading progress across devices. UDP powers the notification system where occasional message loss is acceptable but rapid distribution is valuable.
```

### [HTTP Introduction]
```
The Hypertext Transfer Protocol (HTTP) serves as the foundation of data communication on the World Wide Web. HTTP operates as a request-response protocol in the client-server computing model, where clients send requests to servers, which then return responses. Understanding HTTP's evolution and REST architectural principles is essential for building modern web applications.
```

### [HTTP Evolution]
```
HTTP/1.0 introduced basic request-response mechanisms but required a new connection for each request. HTTP/1.1 added persistent connections allowing multiple requests over a single TCP connection, pipelining for sending multiple requests without waiting for responses, chunked transfer encoding for streaming, and host header enabling virtual hosting. HTTP/2 introduced binary framing for efficiency, multiplexing multiple streams over one connection, header compression with HPACK, and server push for proactive resource delivery. HTTP/3 uses QUIC protocol over UDP, providing connection migration, improved congestion control, and reduced latency.
```

### [REST Principles]
```
Representational State Transfer (REST) is an architectural style for distributed hypermedia systems. RESTful services treat resources as first-class entities identified by URIs, with operations performed using standard HTTP methods. Key principles include: Client-Server separation enabling independent evolution, Statelessness where each request contains all needed information improving scalability and caching, Cacheability through HTTP cache headers reducing server load, Uniform Interface simplifying architecture, Layered System allowing intermediary servers, and Code-On-Demand optionally delivering executable code.
```

### [HTTP Methods]
```
HTTP defines standard methods with specific semantics. GET retrieves resources and is both safe (no side effects) and idempotent (multiple identical requests produce same result). POST creates new resources, is neither safe nor idempotent. PUT replaces entire resources or creates if nonexistent, is idempotent. PATCH applies partial modifications, may or may not be idempotent. DELETE removes resources, is idempotent. HEAD retrieves headers only, useful for checking existence. OPTIONS discovers supported methods.
```

### [Status Codes]
```
HTTP status codes indicate request outcomes. 2xx series signals success: 200 OK for successful GET, 201 Created for successful POST, 204 No Content for successful DELETE. 3xx series handles redirection. 4xx series indicates client errors: 400 Bad Request for malformed requests, 401 Unauthorized requiring authentication, 403 Forbidden denying access despite authentication, 404 Not Found for missing resources, 422 Unprocessable Entity for validation errors. 5xx series signals server errors: 500 Internal Server Error for unexpected conditions, 502 Bad Gateway for invalid upstream response, 503 Service Unavailable during maintenance.
```

### [API Design Best Practices]
```
RESTful API design follows several best practices. Resource-oriented URLs use nouns not verbs (/api/v1/manga not /api/v1/getManga). Versioning strategies include URL versioning (/v1/), header versioning (Accept: application/vnd.api+json; version=1), or accept header versioning. Pagination handles large datasets with cursor or offset-based approaches. Filtering, sorting, and searching use query parameters. Error responses include machine-readable codes and human-readable messages. Rate limiting prevents abuse through request quotas per user or IP. CORS headers enable cross-origin requests. Comprehensive documentation with examples aids integration.
```

### [WebSocket Introduction]
```
WebSocket is a protocol providing full-duplex communication channels over a single TCP connection. Unlike HTTP's request-response model, WebSocket enables bi-directional data flow where both client and server can send messages independently. This makes WebSocket ideal for applications requiring real-time updates with minimal latency.
```

### [WebSocket Protocol]
```
WebSocket connections begin with an HTTP upgrade request. The client sends standard HTTP request with Upgrade: websocket, Connection: Upgrade headers, and Sec-WebSocket-Key for validation. If the server supports WebSocket, it responds with HTTP 101 Switching Protocols, after which the connection upgrades to WebSocket protocol. Data transmission uses frames containing FIN bit (final fragment), RSV bits (reserved), opcode (text, binary, close, ping, pong), mask bit (client-to-server masking required), payload length, and actual data. The protocol supports message fragmentation for large messages and ping/pong frames for keepalive. Security considerations include origin checking to prevent CSRF attacks and WSS (WebSocket Secure) for encryption.
```

### [Comparison with Alternatives]
```
Long polling simulates real-time by repeatedly requesting updates, creating latency and server overhead. Server-Sent Events (SSE) enable server-to-client streaming over HTTP but lack client-to-server communication and are limited to text data. WebSocket provides true bidirectional communication with minimal overhead after initial handshake, binary data support, and lower latency through persistent connection. For MangaHub's chat system, WebSocket's bidirectional capability and low latency make it the optimal choice.
```

### [Use Cases]
```
WebSocket excels in real-time chat applications enabling instant message delivery, live notifications providing immediate updates without polling, collaborative editing synchronizing changes across users, online gaming transmitting player actions with minimal delay, live dashboards streaming metrics, and financial trading platforms delivering market data. The protocol's efficiency and low latency make it ideal for any application requiring real-time bidirectional communication.
```

### [gRPC Introduction]
```
gRPC is a high-performance, open-source Remote Procedure Call (RPC) framework developed by Google. It enables client applications to directly call methods on server applications as if they were local objects, abstracting network communication complexities. gRPC uses HTTP/2 for transport, Protocol Buffers for serialization, and provides authentication, load balancing, and bi-directional streaming.
```

### [gRPC Architecture]
```
gRPC services are defined in .proto files using Protocol Buffers Interface Definition Language (IDL). The protoc compiler generates client and server code in multiple languages from these definitions, ensuring type safety and API consistency. gRPC supports four service methods: unary RPCs (single request, single response like traditional functions), server streaming (single request, stream of responses for large datasets or real-time updates), client streaming (stream of requests, single response for data aggregation), and bidirectional streaming (independent streams in both directions for chat-like applications). HTTP/2 transport provides multiplexing multiple calls over one connection, header compression reducing overhead, and flow control preventing resource exhaustion.
```

### [Protocol Buffers]
```
Protocol Buffers (protobuf) is a language-neutral, platform-neutral mechanism for serializing structured data. Unlike JSON or XML, Protocol Buffers uses binary format that is smaller, faster, and more efficient. Schema is defined in .proto files with strong typing, compiled into language-specific code by protoc compiler. The format provides backward and forward compatibility through field numbering, default values for missing fields, and optional/required markers. Compared to JSON, protobuf offers 3-10x smaller size, 20-100x faster parsing, and strong typing catching errors at compile time.
```

### [gRPC vs REST]
```
gRPC offers superior performance through binary serialization, HTTP/2 multiplexing, and smaller message sizes. Strong typing with Protocol Buffers catches errors at compile time. Streaming support enables real-time bidirectional communication. Code generation reduces boilerplate. However, gRPC has limited browser support requiring gRPC-Web gateway, less human-readable binary format, and steeper learning curve. REST provides universal browser support, human-readable JSON, simpler debugging, and established tooling. In MangaHub, gRPC serves internal service-to-service communication where performance and type safety are paramount, while REST API serves external clients requiring broad compatibility.
```

---

## Page Count Management

To reach exactly 30 pages:

**Current Structure:**
- Title Page: 1 page
- TOC: 1 page  
- Contribution: 1 page
- Abstract: 1 page
- Chapter 1: ~4 pages
- Chapter 2: ~6 pages
- Chapter 3: ~8 pages (add diagrams)
- Chapter 4: ~10 pages (add code blocks)
- Chapter 5: ~5 pages (add tables)
- Chapter 6: ~3 pages
- References: 1 page
- Appendices: ~2 pages

**Total: ~33 pages**

If you need fewer pages:
- Reduce code examples
- Condense implementation details
- Combine similar sections

If you need more pages:
- Add more code examples
- Include more diagrams
- Expand testing details
- Add more appendices

---

## Next Steps

1. I'll provide content for Chapter 3, 4, 5, and 6 in the next response
2. Each section will have detailed, ready-to-paste content
3. You simply replace placeholders with the provided text
4. Adjust formatting as needed to maintain consistency

Would you like me to continue with Chapters 3-6 content now?
