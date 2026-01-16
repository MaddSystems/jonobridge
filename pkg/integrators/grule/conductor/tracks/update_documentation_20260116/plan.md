# Implementation Plan: Update Documentation and Programming Guides

**Track ID:** `update_documentation_20260116`
**Date:** January 16, 2026
**Priority:** Medium
**Status:** Proposed
**Scope:** `backend/` documentation

## Goal
Update `backend/README.md` and create comprehensive programming guides (`HOW2PROGRAM.md`, `LEEME.md`, `COMOPROGRAMAR.md`) to assist developers in understanding the architecture and adding new features.

## Objectives
1.  **Update `backend/README.md`:** Ensure it reflects the latest architectural changes, specifically the recent audit system refactor ("Pure Explicit Capture", declarative audit) and the `IncomingPacket` flags structure.
2.  **Create `backend/HOW2PROGRAM.md`:** A step-by-step guide for developers.
    *   How to add new logic flags (`IncomingPacket`).
    *   How to add new capabilities (Strategy pattern).
    *   How to update the audit system to capture new data (`SnapshotProvider`).
    *   How to create rules using these new features.
    *   Explanation of the project structure (Adapters -> Capabilities -> Persistence -> Schema -> Audit).
3.  **Create `backend/LEEME.md`:** A Spanish version of `README.md`.
4.  **Create `backend/COMOPROGRAMAR.md`:** A Spanish version of `HOW2PROGRAM.md`.

## Strategy
*   **Analyze Recent Changes:** Focus on the audit system refactoring (explicit snapshots, removing listener/worker auto-capture) and the `IncomingPacket` flags.
*   **Structure Guides:** Follow the requested order: Adapters -> Capabilities -> Persistence -> Schema -> Audit.
*   **Strategy Pattern:** Explicitly explain how capabilities implementation follows the Strategy pattern.
*   **Practical Examples:** Include code snippets for adding a flag, registering a capability, and updating a rule.

## Implementation Steps

### 1. Update `backend/README.md`
*   Review current content.
*   Add section on "Audit System Refactor" (Pure Explicit Capture).
*   Ensure `IncomingPacket` structure is mentioned.

### 2. Create `backend/HOW2PROGRAM.md`
*   **Introduction:** Overview of the "Lego Brick" architecture.
*   **Step 1: Modifying Logic Flags:** How to add fields to `IncomingPacket` in `backend/grule/packet.go`.
*   **Step 2: Adding Capabilities:**
    *   Defining the interface.
    *   Implementing the strategy.
    *   Registering in `main.go` and `context_builder.go`.
*   **Step 3: Persistence & Schema:** Brief mention of updating storage and manifests.
*   **Step 4: Modifying Audit:**
    *   Implementing `SnapshotProvider` interface for new capabilities.
    *   How to ensure data appears in `rule_execution_state`.
*   **Step 5: Creating Rules:** Syntax for using new flags and capabilities.

### 3. Create `backend/LEEME.md`
*   Translate `backend/README.md` to Spanish.

### 4. Create `backend/COMOPROGRAMAR.md`
*   Translate `backend/HOW2PROGRAM.md` to Spanish.

## Verification
*   Review all created markdown files for clarity, accuracy, and formatting.
*   Ensure code examples are correct and follow project conventions.
