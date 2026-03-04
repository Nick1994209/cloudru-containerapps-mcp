---
name: git-commit-best-practices
description: Write effective git commit messages following conventional commits format
---

# Git Commit Best Practices

## When to use this skill
Use this skill when you need to write git commit messages, review commit history, or establish commit message standards for a project.

## Language guidelines
- **Always write commit messages in English**
- This ensures consistency and international collaboration

## How to format commit messages

### Use Conventional Commits format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Choose the right type
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring without functional changes
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates, etc.

### Write effective descriptions
- Use present tense ("add feature" not "added feature")
- Use imperative mood ("fix bug" not "fixes bug")
- Keep subject line concise (50-72 characters)
- Separate subject from body with blank line
- Limit body lines to 72 characters

## How to write good commit messages

### Before committing
1. **Review all changes** in the files being committed
2. **Check for unintended changes** or debug code
3. **Ensure the commit message accurately describes** what changed and why
4. **Keep messages concise but descriptive**

### Example commit messages
```
feat(dns): add DNS monitoring service

Implement DNS monitoring with automatic health checks
and alerting for failed queries.

Closes #123
```

```
fix(handler): resolve null pointer exception in DNS check

The handler was not properly validating input parameters,
causing crashes when receiving malformed requests.
```

```
docs: update API documentation

Add examples for new endpoints and clarify authentication
requirements.
```

```
refactor(database): optimize query performance

Replace N+1 queries with JOIN operations to reduce
database load from 500ms to 50ms.
```

```
test: add unit tests for DNS checker

Cover edge cases including timeout handling and
invalid response parsing.
```

## How to reference issues
Add issue numbers in the footer when applicable:
```
feat(auth): implement OAuth2 login

Add support for Google and GitHub OAuth providers.

Closes #456
Related #789
```

## Best practices checklist
- [ ] Message is in English
- [ ] Type is appropriate for the change
- [ ] Subject line is under 72 characters
- [ ] Subject uses imperative mood
- [ ] Body explains "why" not just "what"
- [ ] No trailing whitespace
- [ ] References related issues if applicable
