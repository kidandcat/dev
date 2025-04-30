# Todo

- [x] Use a markdown file as input (INPUT.md)
- [x] Advanced model process input and create a markdown checklist of tasks to complete (TASKS.md)
- [x] Make the AI iterate the tasks without extra user input
- [x] AI executes the tasks step by step and updates the TASKS.md with the changes per step
- [ ] Each step must have 3 clear phases:
  - [ ] Write tests
  - [ ] Write code
  - [ ] Run tests and fix bugs (iterate until passing)
- [ ] Each step must be atomic and self-contained. If a step fails, the AI must be able to continue from where it left off.
- [ ] High level documentation in `wiki` folder (maybe graph structure using json? find best way to structure relational data for a LLM to understand)
- [ ] Independent AI chat with Nano model checks for loopholes.
- [ ] Run tests after each file modification.
- [ ] Make a commit per step.

