# fastjob

Fastjob is a fast and robust job queue using Google Cloud PubSub 🛰

**Work In Progress**

Design objectives:
- Robustness: never lose a job.
- Reliability: never let the main queue be blocked by failing jobs.

Strategies:
- Robustness: only one external dependencies: PubSub.
- Robustness: the durability is garanteed by PubSub.
- Reliability: the core features are mostly only the PubSub semantics and features (but extensible).
- Reliability: route failing jobs to a dead letter queue.

