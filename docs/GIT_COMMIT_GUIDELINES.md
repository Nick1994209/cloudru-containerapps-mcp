# Git Commit Guidelines for Kilo Code

This document provides guidelines for writing git commits when working with the Kilo Code project in VS Code.

## General Rules

1. **Language**: All commit messages must be written in English
2. **Length**: Keep commit messages concise and to the point
3. **Format**: Use Conventional Commits format

## Conventional Commits Format

Follow the Conventional Commits specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries such as documentation generation

### Examples

#### Feature Addition
```
feat(containerapps): add scaling and environment variables support

- Add MinInstanceCount and MaxInstanceCount fields with defaults
- Add environment variables parsing from string format
- Add command and args support for container configuration
```

#### Bug Fix
```
fix(auth): resolve token expiration handling

Update token refresh logic to handle expired tokens gracefully
```

#### Refactoring
```
refactor(mcp): simplify request parameter handling

Consolidate parameter parsing into helper functions to reduce code duplication
```

## Commit Message Best Practices

1. **Use the imperative mood**: "Add feature" not "Added feature" or "Adds feature"
2. **Keep the first line under 50 characters**: This is the subject line
3. **Separate subject from body with a blank line**
4. **Wrap the body at 72 characters**
5. **Use the body to explain what and why vs. how**

## Before Committing

1. **Review all changes**: Check git diff to ensure all intended changes are included
2. **Verify no unintended changes**: Make sure only relevant files are modified
3. **Test your changes**: Ensure the code builds and tests pass
4. **Check for sensitive data**: Never commit secrets, API keys, or passwords

## VS Code Integration

To make following these guidelines easier in VS Code:

1. Install the "Conventional Commits" extension
2. Configure it to enforce these guidelines
3. Use the integrated source control view to review changes before committing

## Example Workflow

```bash
# Check what files have changed
git status

# Review the changes
git diff

# Stage the relevant files
git add internal/domain/types.go internal/application/cloudru/containerapps.go internal/presentation/mcp_handlers.go

# Commit with conventional commit format
git commit -m "feat(containerapps): add scaling and env vars support"
```

Remember: Good commit messages make it easier to understand project history and collaborate effectively.