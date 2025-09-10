Create Jira issues based on an provided implementation plan.

Please follow the workflow below carefully to create the Jira issues successfully.

## Workflow

1. Read the provided implementation plan and architecture. If no plan or architecture file is provided, ask for it and continue only when both of the files are provided.
2. Identify the key steps in the plan that can be translated into Jira issues
3. For each step, create a Jira issue with the following details:
   - **Summary**: A concise title for the issue
   - **Description**: A detailed description of the task, including the goal and deliverables. See [Jira Ticket Description](#jira-ticket-description) for details
   - **Type**: Set the issue type to "Task"
4. Provide a summary of the Jira issues created to the human developer.

DO NOT DO ANYTHING ELSE THAN CREATING THE JIRA ISSUES! E.g. do not implement any code or initialize anything.

### Jira Ticket Description

**Why**
Explain the rationale behind the task or feature. Describe the business need or user problem being addressed. Explain the benefits of implementing this change. Connect the task to overall project goals.

**What**
Clearly define the task or feature. Provide a summary and a detailed explanation of what needs to be done. Include specific functionalities, changes, or improvements.

**Acceptance Criteria**
List specific, measurable, achievable, relevant, and time-bound (SMART) criteria that must be met for the task to be considered complete. These criteria should be testable and verifiable. Include positive and negative scenarios. Cover functional and non-functional requirements.