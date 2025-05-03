# Todo

- [x] Use a markdown file as input (INPUT.md)
- [x] Advanced model process input and create a markdown checklist of tasks to complete (TASKS.md)
- [x] Make the AI iterate the tasks without extra user input
- [x] AI executes the tasks step by step and updates the TASKS.md with the changes per step
- [x] High level documentation in `wiki` folder
  - [x] Data structures and relationships
  - [x] Methods and functions
  - [x] Files and folders
- [x] Each step must have 3 clear phases:
  - [x] Write tests
  - [x] Write code
  - [x] Run tests and fix bugs (iterate until passing)
- [x] Each step must be atomic and self-contained. If a step fails, the AI must be able to continue from where it left off.
- [x] Independent AI chat with Nano model checks for loopholes.
- [x] Run tests after each file modification.
- [x] Make a commit per step.
- [x] Add a new tool to fetch the documentation, it should return the content of all the markdown files in the `wiki` folder as a string for the LLM to use.
- [x] Update the code in main.go to erase the content of the INPUT.md file when all the tasks are completed, before finishing the main() function.
