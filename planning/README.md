# GoQueue - Project Planning Documentation

**Product:** GoQueue - A Golang Job Processor
**Version:** MVP 0.1.0
**Team:** 2-3 Engineers
**Timeline:** 8 weeks (4 sprints)

---

## Quick Navigation

### 📋 Strategy & Vision
Start here to understand **what we're building and why**.

- **[PRD - Product Requirements](strategy/PRD.md)** - Problem statement, vision, scope decisions
- **[Roadmap](strategy/ROADMAP.md)** - Timeline, milestones, risk register

### 🏗️ Architecture Decisions
Key technical choices and their rationale.

- **[ADR-001: Storage Backend](adrs/ADR-001-storage-backend.md)** - PostgreSQL selection
- _More ADRs to be added as decisions are made_

### 🎯 Execution & Implementation
Detailed breakdown for sprint planning and development.

- **[Epics](execution/epics.md)** - 14 major feature areas with story points
- **[User Stories](execution/user-stories.md)** - 48 detailed, Jira-ready stories
- **[Dependencies](execution/dependencies.md)** - Tech stack, setup, testing strategy

### 📅 Milestone Details
Week-by-week implementation plans.

- **[Milestones Overview](milestones/milestones.md)** - 10 milestones with success criteria
- **[M1: Foundation](milestones/M1-foundation.md)** - Weeks 1-2: Database layer
- **[M2: Processing Engine](milestones/M2-processing.md)** - Weeks 3-4: Workers & retries
- _M3 & M4 to be detailed as we progress_

---

## Project Summary

### What We're Building
A lightweight, self-hosted job processing system for Go teams:
- Single binary deployment
- PostgreSQL backend (use what you have)
- Simple HTTP API
- Go client library and CLI
- Prometheus metrics

### What We're NOT Building (MVP)
- Web UI dashboard
- Job scheduling (cron)
- Multi-tenancy
- Job dependencies/workflows
- Redis backend (post-MVP)

### Success Criteria
**Functional:**
- Submit jobs via HTTP API ✓
- Process with configurable concurrency ✓
- Automatic retries with backoff ✓
- Query job status ✓
- Cancel and retry jobs ✓

**Performance:**
- < 10ms p99 submission latency
- \> 1,000 jobs/sec throughput
- Zero job loss on graceful restart

---

## How to Use This Planning

### 🎯 For Product Owners
1. Start with **[PRD](strategy/PRD.md)** to understand scope and vision
2. Review **[Roadmap](strategy/ROADMAP.md)** for timeline and milestones
3. Use **[Epics](execution/epics.md)** for release planning
4. Track progress using milestone success criteria

### 🛠️ For Tech Leads
1. Review **[ADRs](adrs/)** for architectural context
2. Check **[Dependencies](execution/dependencies.md)** for tech stack decisions
3. Use **[Milestone docs](milestones/)** for sprint planning
4. Break down epics into tasks using **[User Stories](execution/user-stories.md)**

### 👨‍💻 For Developers
1. Read **[Dependencies](execution/dependencies.md)** to set up environment
2. Pick stories from current **[Milestone](milestones/)** during sprint planning
3. Check story dependencies before starting work
4. Refer to acceptance criteria for definition of done

### 🏃 For Scrum Masters
1. Use **[Epics](execution/epics.md)** to create Jira epics
2. Import **[User Stories](execution/user-stories.md)** into Jira
3. Track velocity using story points
4. Monitor milestone completion

---

## Milestone Timeline

| Milestone | Weeks | Focus | Key Deliverables |
|-----------|-------|-------|------------------|
| **M1** | 1-2 | Foundation | DB schema, storage layer, config |
| **M2** | 3-4 | Processing | Workers, retries, graceful shutdown |
| **M3** | 5-6 | Interface | HTTP API, Go client, CLI |
| **M4** | 7-8 | Polish | Metrics, Docker, docs, examples |

---

## Scope Overview

### Total Project Metrics
- **Milestones:** 10
- **Epics:** 14
- **User Stories:** 48
- **Total Story Points:** ~172

### MVP Scope (Target)
- **Milestones:** 1-6
- **Epics:** 1-8
- **Story Points:** ~102
- **Duration:** 8 weeks (4 x 2-week sprints)

---

## Key Technical Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Storage | PostgreSQL | Teams have it; SKIP LOCKED is perfect for queues |
| API | REST/HTTP | Universal, debuggable, fast to build |
| Workers | Embedded Go lib | Best DX for Go teams |
| Retries | Exponential backoff | Industry standard |

See **[ADRs](adrs/)** for detailed decision documentation.

---

## Creating Jira Tickets

### Creating Epics
From **[epics.md](execution/epics.md)**:
1. Epic Name: Use epic title
2. Description: Copy epic description
3. Story Points: Use estimate from epic
4. Priority: Use P0/P1/P2 from epic

### Creating Stories
From **[user-stories.md](execution/user-stories.md)**:
1. Story ID: Use STORY-XXX identifier
2. Title: Use story title
3. Acceptance Criteria: Copy checklist
4. Story Points: Use estimate
5. Epic Link: Link to parent epic
6. Dependencies: Note blocking stories

---

## Sprint Planning Recommendations

### Sprint Structure
- **Sprint Length:** 2 weeks
- **Velocity Target:** 13-21 story points per sprint (1 developer)
- **Total Sprints for MVP:** 4 sprints

### Suggested Sprint Allocation

**Sprint 1 (M1):** Foundation
- Epic 1: Core Domain Model (13 points)
- Epic 2: Repository Layer (8 points)
- **Total:** 21 points

**Sprint 2 (M2):** Concurrency
- Epic 3: Worker Pool (21 points)
- Epic 4: Job Execution (13 points)
- **Total:** 34 points → Split across 2 sprints

**Sprint 3 (M3-M4):** Retry & Cancellation
- Epic 5: Retry Mechanism (13 points)
- Epic 6: Job Cancellation (8 points)
- **Total:** 21 points

**Sprint 4 (M5-M6):** Handlers & API
- Epic 7: Job Handler System (13 points)
- Epic 8: HTTP API (21 points)
- **Total:** 34 points → Split across 2 sprints

---

## Definition of Done

### Story-Level DoD
- [ ] All acceptance criteria met
- [ ] Unit tests written and passing
- [ ] Code reviewed
- [ ] No new linter warnings
- [ ] Documentation updated
- [ ] Merged to main branch

### Epic-Level DoD
- [ ] All stories in epic completed
- [ ] Integration tests passing
- [ ] Feature documented
- [ ] Demo-able to stakeholders

### Milestone-Level DoD
- [ ] All epics in milestone completed
- [ ] Milestone success criteria met
- [ ] No critical bugs
- [ ] Retrospective completed
- [ ] Plan updated for next milestone

---

## Getting Started

### For First-Time Readers
1. Read **[PRD](strategy/PRD.md)** (10 min) - Understand the "why"
2. Review **[Roadmap](strategy/ROADMAP.md)** (5 min) - See the timeline
3. Skim **[Epics](execution/epics.md)** (15 min) - Understand feature breakdown
4. Deep dive **[M1 Milestone](milestones/M1-foundation.md)** (20 min) - See detailed execution plan

### For Developers Starting Work
1. Review **[Dependencies](execution/dependencies.md)** - Set up environment
2. Check **[Current Milestone](milestones/)** - Understand sprint context
3. Pick story from **[User Stories](execution/user-stories.md)** - Choose task
4. Check dependencies before starting

### For Stakeholders
1. Read **[PRD](strategy/PRD.md)** - Vision and scope
2. Track **[Roadmap](strategy/ROADMAP.md)** - Timeline and risks
3. Review milestone completion via success criteria

---

## Document Ownership

| Document Type | Owner | Update Frequency |
|---------------|-------|------------------|
| PRD, Roadmap | Product Owner | After each milestone |
| ADRs | Tech Lead | As decisions are made |
| Epics, Stories | Scrum Master / PO | Weekly backlog grooming |
| Milestones | Tech Lead | Sprint planning |
| Dependencies | Tech Lead | As stack evolves |

---

## Questions & Support

### Common Questions

**Q: Can we change story order?**
A: Yes, but respect dependencies. Check [dependencies.md](execution/dependencies.md) for critical path.

**Q: Can we skip stories?**
A: Some stories are critical (P0), others can be deferred (P1, P2). Consult with Tech Lead and PO.

**Q: What if estimates are wrong?**
A: Re-estimate after first sprint and adjust plan accordingly.

**Q: Do we need all 14 epics?**
A: No, MVP only needs Epics 1-8. Epics 9-14 are for production readiness.

---

## Next Steps

### Before Development Starts
1. [ ] Review PRD with stakeholders for approval
2. [ ] Review and approve all ADRs with tech team
3. [ ] Create repository and CI/CD pipeline
4. [ ] Set up project management tool (Jira)
5. [ ] Import epics and stories into Jira

### Sprint 0 (Optional)
1. [ ] Team environment setup
2. [ ] Development workflow agreed
3. [ ] Definition of Done finalized
4. [ ] First sprint planning session

### Sprint 1 - M1 Foundation
1. [ ] Begin with Epic 1: Core Domain Model
2. [ ] Set up test infrastructure
3. [ ] Daily standups to track progress
4. [ ] Sprint review and retrospective at end

---

## Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | 2024-01-XX | Initial combined planning docs | Team |

---

**Ready to start?** → Begin with **[PRD](strategy/PRD.md)** to understand the vision, then dive into **[M1: Foundation](milestones/M1-foundation.md)** for immediate next steps.
