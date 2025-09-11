# Instructions for AI Agent

Follow all instructions CAREFULLY to make sure we succeed with our tasks.

## Core Philosophy

1. **Simplicity:** Prioritize simple, clear, and maintainable solutions. Avoid unnecessary complexity or over-engineering.
2. **Iterate:** Prefer iterating on existing, working code rather than building entirely new solutions from scratch, unless fundamentally necessary or explicitly requested.
3. **Focus:** Concentrate efforts on the specific task assigned. Avoid unrelated changes or scope creep.
4. **Quality:** Strive for a clean, organized, well-tested, and secure codebase.


<project-structure>

## Project structure

This is a repository hosting a json parser written in Go.

- **Languages**: Go
- **Language Version**: 1.25
- **Key roots**:
  - `main/` main applications
  - `internal/` private application and library code
  - `pkg/` shared libraries
  - `vendor/` application dependencies
  - `configs/` configuration file templates or default configs
  - `init/` system init (systemd, upstart, sysv) and process manager/supervisor (runit, supervisord) configs
  - `scripts/` scripts to perform various build, install, analysis, etc operations
  - `build/` packaging and Continuous Integration
  - `deployments/` IaaS, PaaS, system and container orchestration deployment configurations and templates (docker-compose, kubernetes/helm, terraform)
  - `test/` additional external test apps and test data
  - `docs/` design and user documents
  - `tools/` supporting tools for this project
  - `examples/`examples for your applications and/or public libraries
  - `third_party/` external helper tools
  - `githooks/` Git hooks
  - `assets/` other assets to go along with your repository (images, logos, etc)
  - `website/` this is the place to put your project's website data
- **Build/Run**: `Makefile`

</project-structure>


<documentation>

## Documentation

### Documentation and Context Files for AI assistant's use

- `ai-docs/ai_tasks.md`: Task list (✅ done, ⏳ in progress, ❌ not started). Update status, don't remove tasks.
- `ai-docs/ai_changelog.md`: Log of changes made by AI. Add concise summaries here.
- `ai-docs/ai_learnings.md`: Technical learnings, best practices, error solutions. Add new findings here.

### Design Documents

- `docs/architecture.md`: The architecture of the JSON parser.

### important files at project root level  

- `AGENTS.md`: AI agent instructions with core philosophy, behavioral rules, coding standards, and structured workflows.
- `README.md`: The top-level developer guide.
- `CONTRIBUTING.md`: Contributing guidelines with pull request rules, commit message conventions, and trunk-based development workflow.
- `RELEASE.md`: Production release process including deployment workflows, database migration procedures, and rollback strategies.
- `go.mod`: List of used Go libraries and their versions.
- `Makefile`: Lists all the various actions for running the application, tests, linting etc.
- `.golangci.yml`: The Go linter rules for the project

</documentation>


<behavior>

## Behavioral Rules

- **Existing Project:** NEVER create a new project. Analyze existing structure.
- **Up-to-date Info:** Do NOT rely solely on internal knowledge. ALWAYS use web search to verify external API/library documentation, versions, and best practices.
- **Latest Libraries:** Prefer latest stable versions for dependencies, unless it would introduce significant breaking changes to the existing codebase or conflict with versions defined in `go.mod`.
- **Cautious Edits:** Search the codebase to understand impact before modifying existing code. Avoid deletions unless clearly required.
- **Follow Instructions:** When following a workflow defined in `<workflows>`, adhere precisely to it's steps and these rules.
- **Persistence:** Be persistent, do not give up. Don't stop when you have not completed the task. 
- **Simplicity:** Prefer the simplest, most direct solution meeting all requirements.
- **Tool Usage:** Use provided tools appropriately. Explain *why* a tool is used and *how* the prompt is constructed.
- **Clarity:** Communicate plans, workflows, results clearly. Ask clarifying questions if ambiguous.
- **Exact Request Fulfillment:** Implement *only* what is explicitly asked. No extras.
    - Confirm all parts of the request are addressed, nothing more.
    - Ask: "Am I adding anything not explicitly requested?" If yes, remove it.

</behavior>


<coding-rules>

## Coding Rules

### Rules for Fixing Issues
When encountering errors, failing tests or other failure states, where fixing the issue is needed, then always follow the `<fix>` workflow (in `<workflows>`) to fix the issue. After the issue is fixed, return to where you left off. I recommend taking a mental note of the task at hand before fixing the issue, to make sure you know where to return after the fix is done.

### Rules for Writing Useful Comments in Code
- Comments should answer the question "why?" instead of "what?"
- Add comments to give the reader the necessary context for understanding the code
- Avoid commenting obvious and self-documenting code. Strive to make the code self-documenting instead by using clear naming conventions

### Rules for Writing Git Branch Names
- The name should summarize what the changes in the branch achieve (the intent)
- Try to stay between two and eight words, the less the better.
- If this is a JIRA issue. Include the JIRA issue key as the first entry in the branch name
  - The JIRA issue should be separated with a `/` from the rest of the branch name
- If this is not and JIRA Issue. Just use a descriptive name

## Example Format
```
ABC-1/support_multiple_json_objects
```

### Rules for Writing Git Commit Messages
- The title should summarize what the changes in the commit achieve (the intent)
- Include all relevant information about what has changed and why the change has been made
- Include the type of the change with the target component separated by a forward slash `/` as a prefix for the commit message title inside square brackets. The allowed change types are:
  - feat
  - fix
  - test
	- chore
	- docs
	- refactor
	- style
	- perf
  - ci
  - tool

#### Example Format
```
[feat/json-parser] add support for multiple json objects

- Feature now supports multiple JSON objects
- Added unit tests
- Updated documentation to include examples
```

### Go Specific Coding Rules
- You are an expert Go programmer with lots of experience
- You follow the industry standard best practices when programming in Go
- You are mindful about code style and follow the linter rules as described in `#### Go Linting Rules`
- Use `any` instead of `ìnterface{}`
- `ioutil.ReadFile` is deprecated: As of Go 1.16, this function simply calls `os.ReadFile`.

#### Go Version
- Use Go version 1.25

#### Go Linting Rules
- Use `golangci-lint` for linting files.
- Make sure that `golangci-lint` uses the top level `.golangci.yml` configuration.

</coding-rules>


## Request Processing Steps

This is very important, so DO NOT ignore this!

Follow these instructions carefully for **every** user request and follow behavioral rules above (see `<behavior>`):

1. **Initial Analysis:**
  * Read the user request carefully.
  * Consult `<documentation>` and instructions for understanding and planning the request
  * Identify the primary goal of the request (like research, implement code, fix a problem, etc.)
2. Workflow Planning
  * Read and understand the set of pre-defined workflows in `<workflows>`.
  * Based on the primary goal of the request, choose the appropriate `<workflows>` to include. If unclear, `<research>` is a sensible default.
  * Plan the order of the `<workflows>` to execute based on the requirements of the request and the dependency hierarchy of the chosen workflows.
  * If the user's request contains a very specific workflow which is in great conflict with `<workflows>`, then adhere to the workflow that the user provided. Refer to this workflow as `<user-defined-workflow>` in the upcoming steps. 
3. **Present Analysis and Workflow Plan**
  * **Present Analysis:**
    ```
    Primary goal: [primary goal of the user's request]
    Workflow plan: [the workflow plan with transitions, eg. "<research> -> <plan> -> <implement> -> <validate> -> <record>"]
    ```
4. **Workflow Execution:**
  * Execute the workflow plan.
  * Make sure to follow every step within the workflows to complete them successfully.
  * Adhere strictly to `<behavior_rules>`.
  * Review outputs and iterate with follow-up prompts if necessary.
5.  **Final Recording (on task completion):**
  * Once the *entire* request is fulfilled:
  * Summarize significant changes.
  * Ensure that `<record>` workflow has been performed if relevant for the request.
  * Ensure `ai-docs/ai_learnings.md` was updated during `<fix>` workflow if applicable.


<workflows>

## Workflows for AI agent

This defines the workflows available for the AI agent to use when performing a task.

Each workflow has the following sections:
- **Goal:** The goal of the task.
- **Dependencies:** If the workflow requires other workflows to be performed first, then they are listed here.
- **Steps:** The concrete steps to follow then using the workflow. All steps are mandatory to follow.


<research>

**Goal:** Research a topic related to the code base.
**Dependencies:** None
**Steps:**
  1. Inform the user that we are starting this workflow in the format "_Initiated workflow: **research**_"
  2. Familiarize yourself with the task at hand and make sure you understand it.
  3. Look at `<documentation>` and read the documents that are relevant for the task.
  4. Gather further resources that are necessary for the task. If needed, perform actual web searches.
  5. If needed, perform actual web searches for more relevant information.
  6. Look for reference implementations for similar features in the code base.
  7. Make sure you understand the task in the light of all the information you have. If the task is not clear, ask for clarification at this point
  8. Inform the user that we have completed this workflow in the format "_Completed workflow: **research**_"

</research>


<plan>

**Goal:** Plan the implementation of a coding task.
**Dependencies:** `<research>`
**Steps:**
  1. Inform the user that we are starting this workflow in the format "_Initiated workflow: **plan**_"
  2. Based on the research, plan the implementation of the coding task.
  3. Break down the plan into logical, safe, and manageable items.
  4. Verify that the plan and it's items fulfill the user`s requirements.
  5. **Verify `ai-docs/ai_tasks.md` exists.** First, check if the file `ai-docs/ai_tasks.md` is present.
  6. **Create the file if needed.** If the file does not exist, create it with an empty state.
  7. **Write the plan.** Always write the plan with it's items into the file ON DISK `ai-docs/ai_tasks.md` and mark the individual items as ❌ (not started).
  8. Inform the user that you have completed this workflow in the format "_Completed workflow: **plan**_"

</plan>


<implement>

**Goal:** Implement a coding task.
**Dependencies:** `<plan>`
**Steps:**
  1. Inform the user that we are starting this workflow in the format "_Initiated workflow: **implement**_"
  2. Refer to `ai-docs/ai_tasks.md` to understand the items required to implement the coding task.
  3. Mark the items that are going to be implemented as ⏳ (in progress) in `ai-docs/ai_tasks.md`.
  4. Implement the task. Make sure you follow the coding rules defined in `<coding-rules>`. 
  5. **Run the linter on the modified files and fix any new issues.**
  6. Inform the user that we have completed this workflow in the format "_Completed workflow: **implement**_"

</implement>


<validate>

**Goal:** Validate the implementation of a coding task
**Dependencies:** `<implement>`
**Steps:**
  1. Inform the user that we are starting this workflow in the format "_Initiated workflow: **validate**_"
  2. **Understand Requirements:** Consider whether the task requires any of these:
      - Manual testing
      - Unit tests
      - Integration tests
  3. If the change resulted from the implementation is so trivial that no validation is required, then inform the user and jump to step 9.
  4. **Manual Testing:** If the task only requires manual testing:
    - Perform manual testing. In case of errors adhere to the `<fix>` workflow to fix the issue, then return here.
  5. **Write, Run, and Fix Unit Tests:** If the task requires unit tests:
    - Write the unit tests
    - Run the unit tests. In case of errors adhere to the `<fix>` workflow to fix the issue, then return here.
  6. **Write, Run, and Fix Integration tests:** If the task requires integration tests:
    - Write the integration tests
    - Run the integration tests. In case of errors adhere to the `<fix>` workflow to fix the issue, then return here.
  7. **Document:** Ensure any non-trivial fixes are documented in `ai-docs/ai_learnings.md`
  8. **Repeat:** Continue until all relevant tests pass.
  9. Inform the user that we have completed this workflow in the format "_Completed workflow: **validate**_"

</validate>


<document>

**Goal:** Write project documentation. Should be done when documentation is created or updated.
**Dependencies:** None
**Steps:**
  1. Inform the user that you are starting this workflow in the format "_Initiated workflow: **document**_"
  2. Identify whether we need to update or create documentation in `<documentation>`.
  3. Perform the documentation task.
  4. If you created new documentation, then add it to the document index in `<documentation>` in the `AGENTS.md` file.
  5. Inform the user that you have completed this workflow in the format "_Completed workflow: **document**_"

</document>


<record>

**Goal:** Document completed work and update the task backlog.
**Dependencies:** `<implement>` AND/OR `<validate>`
**Steps:**
  1. Inform the user that you are starting this workflow in the format "_Initiated workflow: **record**_"
  2. Add a summary of what has been done to the file `ai-docs/ai_changelog.md`.
  3. Update the status (✅, ⏳, ❌) of the item(s) by writing them into the file ON DISK `ai-docs/ai_tasks.md`. **Do not remove tasks.**
  4. Inform the user that you have completed this workflow in the format "_Completed workflow: **record**_"

</record>


<fix>

**Goal:** Diagnose and fix a problem
**Dependencies:** None
**General Guidelines:**
  - Concentrate on fixing the root cause, not the symptom.
  - Commenting out code that causes an error is NOT a valid fix.
  - Deleting failing test cases is NOT a valid fix. 
**Steps:**
  1. Inform the user that we have started fixing an issue. Use the format "_Started fixing issue: **ISSUE**_", replacing ISSUE
  with a short summary of the issue.
  2. **Gather Context:**
    - Read `ai-docs/ai_learnings.md` for previous solutions.
    - If error messages contain URLs, use web search to understand them.
    - Fetch relevant agent instructions and documentation.
    - Use web search for general error resolution information.
  3. **Iterative Fixing Loop:**
    a. **Hypothesize:** Based on context, identify 1-2 likely root causes.
    b. **Validate Hypothesis (Optional but Recommended):** Add temporary logging, then use the terminal to run relevant tests/code and observe logs.
    c. **Implement Fix:** Apply the proposed code change.
    d. **Validate Fix:** Use the terminal to run tests or execute the relevant code path.
    e. **Record Outcome:**
      - If fixed: Update `ai-docs/ai_learnings.md` with the solution. Delete `/temp/ai_fix_backlog.md` if it exists.
      - If not fixed: Create or update `ai-docs/temp/ai_fix_backlog.md` with what was tried. Return to step (a). *Do not loop more than 3 times on the same core issue without asking the user.*
  4. Inform the user that we have completed fixing the issue by saying "_Completed fixing the issue: **SUMMARY**_", replacing
  SUMMARY with a one sentence summary of the fix.

</fix>

</workflows>
