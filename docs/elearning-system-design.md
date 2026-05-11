# E-Learning Platform — System Design Document

---

## 1. System overview

The e-learning platform is a multi-tenant SaaS product that enables organisations to create, distribute, and track online training for their workforce. Each organisation gets a branded **Portal** at `<code>.e-learning.com` (public course catalogue and player) and an **LMS** at `<code>.e-learning.com/lms` (back-office management interface).

The platform is built as a set of domain-driven microservices organised around six bounded contexts: **Identity & Access**, **Catalog**, **Content Creation**, **Learning**, **Subscription & Licensing**, and **Reporting**. An event-driven backbone (Kafka) connects all services; a ClickHouse analytics store receives every domain event for reporting.

### Key actors

| Actor | Description |
|---|---|
| **Supervisor** | Company admin. Creates and configures portals, manages accounts, issues subscriptions and licenses. |
| **Admin** | LMS operator. Authors courses, manages students and groups, assigns courses, views reports. |
| **Student** | End user. Browses the portal catalogue, plays courses, tracks own progress, earns certificates. |

---

## 2. Functional requirements

### 2.1 Identity & Access

- The system must support multiple independent portals, each identified by a unique short code. Each portal gets a derived domain (`<code>.e-learning.com`) and a corresponding LMS path (`<code>.e-learning.com/lms`).
- A Supervisor must be able to create, configure, and archive portals and accounts.
- A portal can have multiple accounts; an account belongs to one portal.
- An account can have multiple Admins. Both Supervisors and Admins can create Admin users. Admins have no role hierarchy — all Admins within an account have equal permissions.
- Accounts and portals must be archivable (soft delete; data retained).
- Users authenticate via a token-based auth service. The API gateway validates every token before forwarding requests downstream. Optional SSO (SAML / OIDC) is supported per portal.

### 2.2 Course catalogue

- The catalogue is public-facing and scoped per portal.
- Courses are organised into a tree of categories. Each category may have a parent category; root categories have no parent.
- Category operations: create, update (rename, re-parent), delete.
- Course listings in the catalogue are denormalised read models, cached per portal domain, updated when a course is published or unpublished by an Admin.

### 2.3 Content creation

- Admins can create courses containing an ordered set of chapters, each containing an ordered set of lessons.
- Supported lesson types: video, PDF, image, quiz.
- A course may optionally include a certification test.
- Uploaded video assets are processed asynchronously via an external transcoder; the course editor reflects processing status.
- Publishing a course is a separate catalogue action from authoring it. A course can be edited without affecting its published state.

### 2.4 Student management

- Admins can create, deactivate, and manage students within their account.
- Students can be organised into named groups. Groups support full lifecycle management: create, update, delete.
- A student must hold a valid license seat (allocated from the account's subscription) to access assigned courses.

### 2.5 Assignments

- Admins can assign courses to individual students or to entire groups.
- Group assignment fans out asynchronously to create per-student assignment records.
- An assignment can have an optional due date.
- Assignments can be revoked.
- The system must support a high-performance course access check: given a student and a course, determine whether the student has an active assignment and a valid license. Results are cached in Redis with a 5-minute TTL and invalidated on assignment revocation or subscription expiry.

### 2.6 Progress tracking

- Students can start a course, complete lessons, submit quiz attempts, and complete a course.
- The system tracks progress at lesson granularity.
- On completion of all lessons and a passing quiz score, the system automatically marks the course as completed and issues a certificate if a certification test is attached.
- Progress writes are append-only; read projections are built asynchronously.

### 2.7 Subscriptions & licensing

- Supervisors create subscriptions for an account. A subscription contains a pool of N licenses.
- Licenses can be added to an active subscription.
- Subscriptions can be renewed or allowed to expire.
- On expiry, all license seats for the account are revoked and the corresponding course-access cache entries are invalidated.

### 2.8 Reporting

- Admins and Supervisors can create and update report templates that define a report type and its configurable parameters (e.g. date range).
- Any user with access can order a report by selecting a template and supplying parameters.
- Report generation is asynchronous: the system enqueues a job, queries ClickHouse, formats the result, uploads it to object storage, and notifies the requester when ready.
- Example built-in report type: **Course Completion** — shows course completion counts and rates for a given date range, derived from domain events stored in ClickHouse.

### 2.9 Notifications

- The notification service subscribes to domain events from all other contexts and sends email and in-app notifications for key lifecycle moments: course assigned, due date set, course completed, certificate earned, subscription expiring, licenses added, report ready.
- Notification preferences are configurable per user.

---

## 3. Non-functional requirements

### 3.1 Scale targets

| Metric | Target |
|---|---|
| Monthly active users (MAU) | 5,000,000 |
| Daily active users (DAU) | 200,000 |
| Catalogue service read RPS | 50 |
| Assignment service read RPS | 100,000 |
| Progress service write RPS | 120,000 |

### 3.2 Performance

- Course access checks (assignment service) must respond in **< 10 ms** at p99 under normal load, achieved via Redis caching with a 5-minute TTL.
- Catalogue responses must be served from per-portal CDN/cache; origin latency target **< 50 ms** p95.
- Progress writes must sustain **120,000 writes/sec** using an append-only event log with asynchronous read-model projection. Standard relational writes alone are insufficient at this volume; a sharded or time-series-optimised storage strategy is required.
- Report generation is fully asynchronous; no latency SLA on generation time, but the user must be notified on completion.

### 3.3 Availability & reliability

- Core learning path (auth, catalogue, course player, progress) targets **99.9% uptime** (≤ 8.7 hours downtime/year).
- LMS back-office and reporting target **99.5% uptime**.
- All services publish domain events to Kafka with at-least-once delivery semantics. Consumers must be idempotent.
- The ClickHouse event sink must tolerate consumer lag without data loss; Kafka retention must cover at least 7 days.

### 3.4 Data architecture

| Store | Purpose |
|---|---|
| **PostgreSQL** (per service) | Each microservice owns its own database; no cross-service SQL joins. |
| **Redis** | Short-lived course access view cache. Key: `courseaccess:{studentId}:{courseId}`, TTL 5 min. |
| **Apache Kafka** | Domain event bus. All services publish events; ClickHouse sink and notification service are consumers. |
| **ClickHouse** | Append-only analytics store. Receives all domain events via Kafka sink. Source of truth for all reports. |
| **Object storage (S3)** | Raw and transcoded video, PDFs, images, and generated report files. |

### 3.5 Security

- All inter-service communication inside the cluster is internal; only the API gateway is internet-facing.
- The API gateway validates auth tokens on every request before routing downstream.
- Tenant isolation is enforced at the application layer: every query is scoped to a `portalId` / `accountId` derived from the validated token.
- Media assets are served via pre-signed URLs or CDN-signed tokens; direct S3 bucket access is not permitted.
- SSO per portal is supported via SAML 2.0 / OIDC; credentials are never stored by the platform when SSO is active.

### 3.6 Observability

- All services emit structured logs (JSON) with a `correlationId` propagated from the API gateway.
- Distributed tracing is required across the full request path (gateway → service → database).
- Kafka consumer lag is monitored per consumer group; alerting triggers if lag exceeds a configurable threshold.
- ClickHouse query performance is monitored; slow queries on the report worker are logged and alerted.

### 3.7 Scalability & deployment

- All services are independently deployable and horizontally scalable.
- The progress service and assignment service are the primary scaling targets given their RPS profiles; they should be scaled and deployed independently of all other services.
- Group assignment fan-out (potentially 10,000+ records per operation) is handled by a background worker to avoid blocking the API response.
- The report generator worker is stateless and can be scaled horizontally to process concurrent report jobs.
