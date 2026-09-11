# Architecture Guide for Contributors and Agents

This repository is a **modular monolith**: it is deployed as one application and
kept in one repository, but its code is organized into independently understandable
business modules. Every contributor and agent must preserve those boundaries when
adding, changing, reviewing, or refactoring code.

When a change does not clearly fit an existing module, first identify the business
capability it serves. Create a new module only when that capability has a durable,
cohesive responsibility; do not create modules merely to mirror technical layers.

## Core Rules

1. **Organize modules around business capabilities.**

   A module represents a cohesive domain concern, such as billing, identity,
   messaging, or subscriptions—not a generic technical category such as
   `controllers`, `services`, or `helpers`. Keep the feature's domain logic, use
   cases, and internal implementation near the module that owns the capability.

2. **Give every module an explicit public API.**

   Other modules may depend only on the owner module's documented, intentional
   entry points: exported use cases, contracts, events, or query interfaces. Do
   not import another module's internal files, reach into its private state, or
   rely on incidental implementation details. Treat public APIs as contracts:
   make changes deliberately and update their consumers as needed.

3. **Do not create a catch-all shared `utils` area.**

   Keep code inside the module that owns the business meaning. A helper should be
   private to a module unless it is genuinely generic, stable, and needed by
   multiple modules. In that exceptional case, place it in a small, purpose-named
   shared package with a narrow API—not in an ever-growing `utils` folder.

4. **Keep dependency direction acyclic.**

   Module dependencies must form a directed acyclic graph (DAG). A module may
   depend on another module's public API, but dependencies must never circle back.
   If two modules need each other, revisit the boundary: extract a neutral contract,
   introduce an event, move the responsibility to its true owner, or merge modules
   if they are not independently cohesive.

5. **Let modules own their documents and models.**

   The module that owns a business concept also owns its domain models, persistence
   documents/schema, validation rules, and migrations for that concept. Other
   modules must not alter those records directly or duplicate the model as an
   unofficial source of truth; they should request work through the owner's API.

6. **Place interfaces with the consumer.**

   Define an interface/port where it is used—the consuming module specifies the
   capability it requires. The provider implements that contract or is adapted to
   it. This avoids leaking provider details into consumers and keeps dependencies
   pointing toward stable business needs. Do not introduce an interface unless it
   has a real consumer-driven purpose.

7. **Avoid premature abstraction.**

   Prefer simple, local code until repeated requirements demonstrate a stable
   common abstraction. Do not build generic frameworks, base classes, or shared
   layers solely because future reuse is possible. Extract only after identifying
   a concrete common behavior and a clear owner for it.

8. **Make database access module-oriented.**

   Database reads and writes belong behind the owning module's repository/data
   access boundary. A module must not query, mutate, or join another module's
   tables/collections directly. Cross-module workflows should use public APIs,
   explicit read models, events, or owner-provided queries instead of bypassing
   the boundary.

9. **Design modules to be replaceable.**

   A module should be removable or replaceable with limited, well-defined changes
   outside it. Keep coupling to its public contracts, avoid shared mutable state
   and cross-module persistence access, and make integrations explicit. If removing
   a module requires widespread edits across unrelated areas, improve its boundary
   before adding more dependencies.

## Working Checklist

Before finishing a change, verify the following:

- The change has a clear business-module owner.
- Cross-module calls use only explicit public APIs.
- No new circular dependency or internal import was introduced.
- Models and database access remain with their owning module.
- Any shared code is narrowly named, truly shared, and has a justified owner.
- New abstractions solve an observed current need rather than a hypothetical one.
- The affected module could still be replaced without broad changes elsewhere.
