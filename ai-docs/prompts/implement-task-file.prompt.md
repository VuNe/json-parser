Implement a task file.

If the task file is not provided, ask for it.

Follow the workflow below carefully to implement the task file successfully.

Always be aware of which step in the workflow we currently are in, so that you can proceed to the next step after for example iterating on a step.

## Workflow

1. Fetch the details using the provided task file. If one is not provided ask for it.
2. Read the task description to understand the requirements of the task.
3. Check `README.md`, `docs/architecture.md` and `ai_docs/ai_learnings.md` (if exists) for any relevant information for successfully completing the task
4. Check reference implementations in the codebase, if available. ALWAYS follow the existing patterns and conventions in the codebase, unless they contradict with `docs/architecture.md` or given intructions
5. If the issue and it's requirements are not clear at this point, ask for clarification
6. Outline the steps required to complete the task
7. Make sure you are on the "main" or "master" branch (whichever name is used in this repo) and pull the latest changes
8. Checkout a new branch, named based on the `<coding-rules>`
9. Implement the task
10. Make sure you have written unit tests for the implementation
11. Update (or create) `README.md` and `docs/architecture.md` if necessary to reflect the changes made
12. Ask the human developer to review the implementation. If applicable, provide the human developer with steps to test the implementation themselves.
13. Once the human developer approves the implementation, suggest a commit message based on the `<coding-rules>` for the human developer
14. Once the human developer approves the commit message, commit the changes
15. The task is now complete. Good job! Remember to provide the link to the pull request to the human developer.
