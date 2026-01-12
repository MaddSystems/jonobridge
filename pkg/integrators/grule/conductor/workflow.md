# Project Workflow

## Guiding Principles

1. **The Plan is the Source of Truth:** All work must be tracked in `plan.md`
2. **The Tech Stack is Deliberate:** Changes to the tech stack must be documented in `tech-stack.md` *before* implementation
3. **Tests are Optional but Recommended:** Write unit tests when critical or complex functionality is implemented, as decided by the human.
4. **Human Control:** The user controls all commands, git operations, and deployment. No automatic prompts or CI/CD pipelines.
5. **Consult the Manual:** Always consult the documentation in the `Manual/` folder to understand the engine's rules and behavior before writing code.
6. **User Experience First:** Every decision should prioritize user experience
7. **Frontend & Browser Testing:** Frontend verification and browser testing are strictly manual processes handled by the human.
8. **Non-Interactive:** Prefer non-interactive commands where possible, but always await user confirmation for critical actions.

## Task Workflow

All tasks follow a strict lifecycle:

### Standard Task Workflow

1. **Select Task:** Choose the next available task from `plan.md` in sequential order

2. **Mark In Progress:** Before beginning work, edit `plan.md` and change the task from `[ ]` to `[~]`

3. **Consult Documentation:**
   - Read relevant sections in the `Manual/` folder to ensure understanding of the engine's rules and architecture.

4. **Write Tests (Optional):**
   - If the task is critical or complex, or if the user requests it:
     - Create a new test file for the feature or bug fix.
     - Write one or more unit tests that clearly define the expected behavior.
     - Run the tests and confirm they fail (Red Phase).

5. **Implement Feature:**
   - Write the application code necessary to implement the task.
   - If tests were written, run them to ensure they pass (Green Phase).

6. **Refactor (If Recommended):**
   - Only refactor if explicitly recommended by the human.
   - Refactor code for clarity and performance.
   - Rerun tests if applicable.

7. **Verify Coverage (Optional):** Run coverage reports if tests were written.

8. **Document Deviations:** If implementation differs from tech stack:
   - **STOP** implementation
   - Update `tech-stack.md` with new design
   - Add dated note explaining the change
   - Resume implementation

9. **Stage Code Changes:**
   - Request the user to stage all code changes related to the task (`git add ...`).
   - *Note: Commits are performed at the end of the Phase.*

10. **Mark Task Complete:**
    - Update `plan.md`: change the task from `[~]` to `[x]`.

### Phase Completion Verification and Checkpointing Protocol

**Trigger:** This protocol is executed immediately after all tasks in a phase are completed.

1.  **Announce Protocol Start:** Inform the user that the phase is complete and the verification and checkpointing protocol has begun.

2.  **Ensure Test Coverage (If Applicable):**
    -   **Step 2.1: Determine Phase Scope:** Identify all files changed in this phase (staged changes).
    -   **Step 2.2: Verify Tests:** For critical/complex code files, check if tests exist.
        -   Ask the user if additional tests are required for the changed files.
        -   If yes, create and run them.

3.  **Execute Automated Tests (If Applicable):**
    -   If tests exist, announce the command and run them.
    -   If tests fail, debugging is required.

4.  **Propose a Detailed, Actionable Manual Verification Plan:**
    -   **CRITICAL:** To generate the plan, first analyze `product.md`, `product-guidelines.md`, and `plan.md` to determine the user-facing goals of the completed phase.
    -   You **must** generate a step-by-step plan that walks the user through the verification process.
    -   **Frontend verification is manual.** Do not attempt to script browser interactions.
    -   The plan you present to the user **must** follow this format:

        **For a Frontend Change:**
        ```
        The automated tests (if any) have passed. For manual verification, please follow these steps:

        **Manual Verification Steps:**
        1.  **Start the development server with the command:** `[Insert Command]`
        2.  **Open your browser to:** `[Insert URL]`
        3.  **Confirm that you see:** The new user profile page...
        ```

5.  **Await Explicit User Feedback:**
    -   After presenting the detailed plan, ask the user for confirmation: "**Does this meet your expectations? Please confirm with yes or provide feedback on what needs to be changed.**"
    -   **PAUSE** and await the user's response. Do not proceed without an explicit yes or confirmation.

6.  **Request Phase Commit:**
    -   Ask the user to perform the commit with a suggested message (e.g., `feat(phase): Complete Phase X - <Description>`).

7.  **Attach Phase Summary & Verification Report using Git Notes:**
    -   **Step 7.1: Get Commit Hash:** Obtain the hash of the *just-completed commit*.
    -   **Step 7.2: Draft Note Content:** Create a detailed summary including:
        - Phase Name
        - Summary of tasks completed
        - List of created/modified files
        - Automated test command used (if any)
        - Manual verification steps and user confirmation
    -   **Step 7.3: Request Attach Note:** Ask the user to run the `git notes` command to attach the summary to the commit.
     ```bash
     git notes add -m "<note content>" <commit_hash>
     ```

8.  **Get and Record Phase Checkpoint SHA:**
    -   **Step 8.1: Get Commit Hash:** Obtain the hash of the phase commit.
    -   **Step 8.2: Update Plan:** Read `plan.md`, find the heading for the completed phase, and append the first 7 characters of the commit hash in the format `[checkpoint: <sha>]`.
    -   **Step 8.3: Write Plan:** Write the updated content back to `plan.md`.

9. **Commit Plan Update:**
    - **Action:** Request the user to stage the modified `plan.md` file.
    - **Action:** Request the user to commit this change with a descriptive message following the format `conductor(plan): Mark phase '<PHASE NAME>' as complete`.

10.  **Announce Completion:** Inform the user that the phase is complete and the checkpoint has been created.

### Quality Gates

Before marking any phase complete, verify:

- [ ] All tests pass (if any)
- [ ] Code follows project's code style guidelines
- [ ] Works correctly on mobile (if applicable)
- [ ] Documentation updated if needed
- [ ] No security vulnerabilities introduced

## Development Commands

**AI AGENT INSTRUCTION: This section should be adapted to the project's specific language, framework, and build tools.**

### Setup
```bash
# Example: Commands to set up the development environment
# e.g., for a Node.js project: npm install
# e.g., for a Go project: go mod tidy
```

### Daily Development
```bash
# Example: Commands for common daily tasks
# e.g., for a Node.js project: npm run dev, npm test
# e.g., for a Go project: go run main.go, go test ./...
```

### Before Committing
```bash
# Example: Commands to run all pre-commit checks
# e.g., for a Node.js project: npm run check
# e.g., for a Go project: make check
```

## Testing Requirements

### Unit Testing (Optional)
- Tests are recommended for critical business logic.
- Use appropriate test setup/teardown mechanisms.

## Code Review Process

### Self-Review Checklist
Before requesting review:

1. **Functionality**
   - Feature works as specified
   - Edge cases handled

2. **Code Quality**
   - Follows style guide
   - DRY principle applied
   - Clear variable/function names

3. **Testing**
   - Tests written if required by user

4. **Security**
   - No hardcoded secrets
   - Input validation present

## Commit Guidelines

### Message Format
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests
- `chore`: Maintenance tasks